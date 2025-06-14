## 🛠️ Makefile Usage

This project includes a [Makefile](Makefile) to simplify common development tasks such as running the server, executing tests, and managing local scripts.

### 📋 Available Commands

- **Setup & Install**
  - `make setup` or `make install-local`  
    Installs the `run_local` script as a symlink in the `bin/` directory for easy access.

- **Run the Server**
  - `make server`  
    Builds and runs the server in the background.

  - `make dev`  
    Runs the server in the foreground with debug output.

- **Run Tests**
  - `make test`  
    Runs all Go tests in the project.

- **Clean Up**
  - `make clean`  
    Removes the local symlink and cleans up the `bin/` directory if empty.

- **Status**
  - `make status`  
    Shows the installation status of the `run_local` script and symlink.

- **Help**
  - `make help`  
    Lists all available Makefile commands and their descriptions.

### 🚀 Quick Start

1. **Install the local script:**
   ```sh
   make setup
   ```

2. **Run the server in development mode:**
   ```sh
   make dev
   ```

3. **Run tests:**
   ```sh
   make test
   ```

4. **Check status or get help:**
   ```sh
   make status
   make help
   ```

### 💡 Notes

- The `run_local` script is symlinked to `bin/run_local` and can be added to your `PATH` for convenience:
  ```sh
  export PATH="$PWD/bin:$PATH"
  ```
- All commands are designed to work cross-platform and will ensure dependencies are installed as needed.
- For advanced usage and options, run:
  ```sh
  run_local --help
  ```

See the [Makefile](Makefile) for full details and customization