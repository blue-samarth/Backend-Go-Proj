.PHONY: install-local uninstall-local help dev test server clean setup status

# Local bin directory inside project
LOCAL_BIN := $(CURDIR)/bin
SCRIPT := $(CURDIR)/scripts/run_local.sh
LINK := $(LOCAL_BIN)/run_local

# Add local bin to PATH for this make session
export PATH := $(LOCAL_BIN):$(PATH)

# Default target
all: help

install-local: ## Symlink scripts/run_local.sh to bin/run_local
	@echo "🔧 Installing run_local locally..."
	@mkdir -p $(LOCAL_BIN)
	@if [ ! -f "$(SCRIPT)" ]; then \
		echo "❌ Error: $(SCRIPT) not found"; \
		exit 1; \
	fi
	@ln -sf $(SCRIPT) $(LINK)
	@chmod +x $(SCRIPT) $(LINK)
	@echo "✅ Installed: run_local -> $(SCRIPT)"
	@echo "💡 To use in your shell, run:"
	@echo "   export PATH=\"$(LOCAL_BIN):\\$$PATH\""
	@echo "   run_local --help"

uninstall-local: ## Remove the local symlink
	@rm -f $(LINK)
	@echo "✅ Uninstalled local run_local"

setup: install-local ## Full project setup (alias for install-local)

dev: install-local ## Run in development mode
	@echo "🚀 Starting development server..."
	@run_local --debug

test: install-local ## Run tests
	@echo "🧪 Running tests..."
	@run_local --test

server: install-local ## Run server in background
	@echo "⚡ Starting background server..."
	@run_local --server

clean: uninstall-local ## Clean up local installation
	@if [ -d "$(LOCAL_BIN)" ] && [ -z "$(ls -A $(LOCAL_BIN))" ]; then \
		rmdir $(LOCAL_BIN); \
		echo "✅ Removed empty bin directory"; \
	fi

status: ## Show installation status
	@echo "📋 Installation Status:"
	@echo "   Script: $(SCRIPT)"
	@if [ -f "$(SCRIPT)" ]; then echo "   ✅ Script exists"; else echo "   ❌ Script missing"; fi
	@echo "   Link: $(LINK)"
	@if [ -L "$(LINK)" ]; then echo "   ✅ Symlink exists"; else echo "   ❌ Symlink missing"; fi
	@if command -v run_local >/dev/null 2>&1; then \
		echo "   ✅ run_local command available"; \
	else \
		echo "   ⚠️  run_local not in PATH"; \
	fi

help: ## Show Makefile command usage
	@echo ""
	@echo "🛠️  Makefile Commands for Go Backend Project"
	@echo "==========================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .+$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "📦 Project Structure:"
	@echo "  scripts/run_local.sh  → Main script"
	@echo "  bin/run_local         → Symlink created by 'make install-local'"
	@echo ""
	@echo "💡 Tip: To see runtime options, run: \033[33mrun_local --help\033[0m"
