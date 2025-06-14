#!/bin/bash

set -e

# === COLORS ===
GREEN="\033[1;32m"
RED="\033[1;31m"
YELLOW="\033[1;33m"
BLUE="\033[1;34m"
NC="\033[0m"

timestamp() {
  date '+%Y-%m-%d %H:%M:%S'
}

log()     { echo -e "${BLUE}[*]${NC} [$(timestamp)] $*"; }
success() { echo -e "${GREEN}[✓]${NC} [$(timestamp)] $*"; }
warn()    { echo -e "${YELLOW}[!]${NC} [$(timestamp)] $*"; }
error()   { echo -e "${RED}[✗]${NC} [$(timestamp)] $*" >&2; }



# === CONFIG ===
TARGET="${TARGET:-Backend-Go-Project}"
BINARY_NAME="backend"
PORT=8080

MODE="server" # default

for arg in "$@"; do
  case $arg in
    --debug) MODE="debug" ;;
    --test)  MODE="test" ;;
    --server) MODE="server" ;;
    *)
      echo -e "\033[1;31m[✗] Unknown flag: $arg\033[0m" >&2
      exit 1
      ;;
  esac
done

log "Running script in directory: $(pwd)"
log "Target directory: $TARGET"
log "Binary name: $BINARY_NAME"
log "Port: $PORT"





# === Clean up old binary if exists ===
if [ -f "$BINARY_NAME" ]; then
  warn "Removing existing binary '$BINARY_NAME'..."
  rm "$BINARY_NAME"
fi


# === Ensure Go is installed ===
GO_REQUIRED_VERSION="1.24"
GO_LOCAL_DIR=".go-bin"
GO_LOCAL_BIN="$GO_LOCAL_DIR/go/bin/go"


# Try to find Go or install locally
if ! command -v go &> /dev/null; then
  warn "Go not found. Attempting to install Go $GO_REQUIRED_VERSION locally..."

  mkdir -p "$GO_LOCAL_DIR"
  ARCH=$(uname -m)
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')

  if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
  elif [[ "$ARCH" =~ arm64|aarch64 ]]; then
    ARCH="arm64"
  else
    error "Unsupported architecture: $ARCH"
    exit 1
  fi

  GO_TAR="go${GO_REQUIRED_VERSION}.${OS}-${ARCH}.tar.gz"
  GO_URL="https://go.dev/dl/${GO_TAR}"

  if ! command -v curl &> /dev/null; then
    warn "'curl' not found. Attempting to install..."
    
    if [ -f /etc/alpine-release ]; then
        log "Detected Alpine Linux"
        apk add --no-cache curl || { error "Failed to install curl via apk"; exit 1; }
    elif [ -f /etc/debian_version ]; then
        log "Detected Debian-based OS"
        if ! sudo -n true 2>/dev/null; then
            error "sudo access required to install curl on Debian-based systems"
            exit 1
        fi
        sudo apt-get update && sudo apt-get install -y curl || { error "Failed to install curl via apt"; exit 1; }
    elif command -v brew &> /dev/null; then
        log "Detected macOS with Homebrew"
        brew install curl || { error "Failed to install curl via brew"; exit 1; }
    else
        error "'curl' is not installed and automatic installation is not supported on this OS."
        exit 1
    fi
    
    success "'curl' installed successfully."
  else
    log "'curl' is already installed, proceeding with download."
  fi 
    
  for i in {1..3}; do
    if curl -fsSL "$GO_URL" -o "$GO_TAR"; then
        break
    elif [ $i -eq 3 ]; then
        error "Failed to download Go tarball after 3 attempts"
        exit 1
    fi
        warn "Download attempt $i failed, retrying in 2s..."
        sleep 2
    done
  tar -C "$GO_LOCAL_DIR" -xzf "$GO_TAR" || { error "Failed to extract Go tarball"; exit 1; }
  rm "$GO_TAR"

  export PATH="$PWD/$GO_LOCAL_DIR/go/bin:$PATH"
  success "Go $GO_REQUIRED_VERSION installed locally at $GO_LOCAL_DIR/go"
else
  log "Go already installed: $(go version)"
fi


# Finalize GO_CMD after any install
GO_CMD=$(command -v go || echo "$GO_LOCAL_BIN")


# === Add .go-bin to .gitignore if not already ===
if [ ! -f .gitignore ]; then
  log "Creating .gitignore file..."
  touch .gitignore
fi
grep -qxF '.go-bin/' .gitignore || echo '.go-bin/' >> .gitignore


# === Ensure 'bc' is installed ===
if ! command -v bc &> /dev/null; then
  warn "'bc' not found. Attempting to install..."

  if [ -f /etc/alpine-release ]; then
    log "Detected Alpine Linux"
    apk add --no-cache bc || { error "Failed to install bc via apk"; exit 1; }

  elif [ -f /etc/debian_version ]; then
    log "Detected Debian-based OS"
    if ! sudo -n true 2>/dev/null; then
      error "sudo access required to install bc on Debian-based systems"
      exit 1
    fi
    sudo apt-get update && sudo apt-get install -y bc || { error "Failed to install bc via apt"; exit 1; }

  elif command -v brew &> /dev/null; then
    log "Detected macOS with Homebrew"
    brew install bc || { error "Failed to install bc via brew"; exit 1; }

  else
    error "'bc' is not installed and automatic installation is not supported on this OS."
    exit 1
  fi

  success "'bc' installed successfully."
fi


# === Check if we are in the correct project root ===
CURRENT_DIR_NAME="$(basename "$PWD")"

if [ "$CURRENT_DIR_NAME" == "$TARGET" ]; then
  log "Already in the target directory: '$TARGET'"
elif [ -d "$TARGET" ]; then
  log "Not in target directory. Changing into '$TARGET'..."
  cd "$TARGET" || { error "Failed to enter directory '$TARGET'"; exit 1; }
else
  error "Target directory '$TARGET' not found in current path."
  exit 1
fi

# === Ensure this is really the project root ===
if [ ! -f "main.go" ] || [ ! -d ".git" ]; then
  error "This does not appear to be the project root (missing main.go or .git)."
  exit 1
fi


# === Kill any previous process ===
PID=$(pgrep -f "./$BINARY_NAME" || true)
if [ -n "$PID" ]; then
  warn "Killing existing $BINARY_NAME process (PID $PID)..."
  kill -9 "$PID"
fi


# === Check if port is in use ===
if lsof -i:$PORT &>/dev/null; then
  warn "Port $PORT is already in use"
  lsof -i:$PORT
  read -p "Do you want to kill the process using port $PORT? [y/N]: " choice
  if [[ "$choice" == "y" || "$choice" == "Y" ]]; then
    PID_TO_KILL=$(lsof -ti:$PORT)
    kill -9 "$PID_TO_KILL"
    success "Killed process on port $PORT (PID $PID_TO_KILL)"
  else
    error "Aborting due to port conflict."
    exit 1
  fi
fi


# === Download dependencies ===
log "Downloading Go modules..."
"$GO_CMD" mod download


# === Check Go version ===
GO_VERSION=$("$GO_CMD" version | grep -oE 'go[0-9]+\.[0-9]+(\.[0-9]+)?' | sed 's/go//')

if [ -z "$GO_VERSION" ]; then
  error "Failed to detect installed Go version"
  exit 1
fi

REQUIRED="$GO_REQUIRED_VERSION"
INSTALLED="$GO_VERSION"

if [ "$(printf '%s\n' "$REQUIRED" "$INSTALLED" | sort -V | head -n1)" != "$REQUIRED" ]; then
  error "Go $GO_REQUIRED_VERSION+ required, but found $GO_VERSION"
  exit 1
fi

# === Run tests if --test flag is provided ===
if [ "$MODE" = "test" ]; then
  log "Running tests..."
  if ! "$GO_CMD" test ./... -v -count=1; then
    error "Tests failed. Aborting."
    exit 1
  fi
  success "All tests passed"
  exit 0
fi

success "The Test passed successfully"


# === Build ===
log "Building the project..."
"$GO_CMD" build -o "$BINARY_NAME" main.go

success "Build completed successfully"


# === Check if build was successful ===
if [ ! -f "$BINARY_NAME" ]; then
  error "Build failed - binary not found"
  exit 1
fi

success "Build successful - binary '$BINARY_NAME' created"


# === Cleanup on exit ===
cleanup() {
  if [ -n "$SERVER_PID" ]; then
    warn "Killing background server (PID $SERVER_PID)..."
    kill "$SERVER_PID" 2>/dev/null || true
  fi

  if [ -f "$BINARY_NAME" ]; then
    warn "Removing binary '$BINARY_NAME'..."
    rm -f "$BINARY_NAME"
  fi

  log "Cleanup complete. Exiting script."
}


# === Run the server ===
if [ "$DEBUG" = true ]; then
  trap cleanup EXIT
  ./"$BINARY_NAME"  # foreground mode; script stays attached
else
  ./"$BINARY_NAME" &> /dev/null &
  SERVER_PID=$!
  echo "$SERVER_PID" > .server.pid
  success "Server is running in background (PID $SERVER_PID). Use '--debug' to see logs."
  trap cleanup EXIT
fi
