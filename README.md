# GO CLASS SERVER 2025 
## Development branch
This is the branch that generates the executable on 10.2.2025

## 10.3.2025
1. commented out replace
1. issues with blackfriday version 1, currently on version 2
1. go run main.go waiting.....
1. cd chakraFrontend and npm install
1. package-lock.json deleted
1. npm install
1. it rendered the home page.
1. offline website OK, course content nothing, about nothing


The version v1.1.0 doesn't exist for that package. Try these solutions:

## Fix 1: Let Go Find the Correct Version

```bash
# Remove the specific version requirement
go get github.com/russross/blackfriday

# Or update to latest
go get -u github.com/russross/blackfriday
```

## Fix 2: Use blackfriday v2

The package was replaced by v2. Update your import:

```bash
# Install v2
go get github.com/russross/blackfriday/v2
```

Then in your code, change:
```go
// Old
import "github.com/russross/blackfriday"

// New
import "github.com/russross/blackfriday/v2"
```

## Fix 3: Check Available Versions

```bash
# See what versions exist
go list -m -versions github.com/russross/blackfriday
```

## Fix 4: Clean and Retry

```bash
# Clean module cache
go clean -modcache

# Tidy dependencies
go mod tidy

# Download again
go mod download
```

## Fix 5: Check go.mod

Look at your `go.mod` file and verify the version:

```bash
cat go.mod
```

If it shows an invalid version, edit it or run:
```bash
go mod edit -replace github.com/russross/blackfriday@v1.1.0=github.com/russross/blackfriday/v2@latest
go mod tidy
```

The blackfriday v1 line is likely outdated. The current maintained version is v2.
Here are the essential Go commands for running and building:

## Run Go Programs

```bash
# Run directly (doesn't create executable)
go run main.go

# Run with multiple files
go run main.go utils.go

# Run entire package
go run .
```

## Build Go Programs

```bash
# Build executable in current directory
go build

# Build with specific output name
go build -o myapp

# Build with full path
go build -o bin/myapp

# Build for current platform
go build main.go
```

## Cross-Platform Building

```bash
# Build for Windows (from any OS)
GOOS=windows GOARCH=amd64 go build -o app.exe

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o app

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o app

# Build for ARM (Raspberry Pi, etc.)
GOOS=linux GOARCH=arm go build -o app
```

## Common Build Flags

```bash
# Strip debug info (smaller binary)
go build -ldflags="-s -w" -o app

# Enable race detector
go build -race -o app

# Static binary (no external dependencies)
CGO_ENABLED=0 go build -o app

# Verbose output
go build -v

# Show all commands
go build -x
```

## Install (builds and moves to $GOPATH/bin)

```bash
# Install current package
go install

# Install specific package
go install github.com/user/package@latest
```

## Complete Example

```bash
# Development: quick run
go run main.go

# Test build
go build -o test-app

# Production build (optimized)
CGO_ENABLED=0 go build -ldflags="-s -w" -o myapp

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o myapp.exe
```

## Module Commands (if using Go modules)

```bash
# Initialize module
go mod init github.com/user/project

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

The most common workflow is `go run .` for development and `go build` for creating executables.

## Code Assessment

### Critical Cross-Platform Issues

**Hard-coded Windows paths throughout:**
```go
"C:/Users/Public/ClassServer/static"
"C:/websites"
"C:/Users/Public/CANVAS_FILE_CACHES/"
userHome+"/AppData/Local/Google/Chrome/Application/chrome.exe"
```

This code will fail on Linux/macOS. These systems don't have `C:` drives or `AppData` directories.

**Problems:**
1. Application won't start on non-Windows systems (path errors)
2. Chrome launch will fail (hardcoded Windows path)
3. Static file serving breaks completely
4. Symlink creation uses Windows-specific paths

### Recommended Fixes

**1. Use `os.UserConfigDir()` and `os.UserCacheDir()`:**
```go
configDir, _ := os.UserConfigDir()  // Returns platform-specific config dir
dataDir := filepath.Join(configDir, "ClassServer")
```

**2. Make Chrome path cross-platform:**
```go
func openBrowser(url string) error {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("cmd", "/c", "start", url)
    case "darwin":
        cmd = exec.Command("open", url)
    case "linux":
        cmd = exec.Command("xdg-open", url)
    default:
        return fmt.Errorf("unsupported platform")
    }
    return cmd.Start()
}
```

**3. Use `filepath.Join()` everywhere:**
```go
// Instead of:
"C:/websites/" + item

// Use:
filepath.Join(websitesDir, item)
```

**4. Make base paths configurable:**
```go
var (
    websitesDir = flag.String("websites", getDefaultWebsitesDir(), "Websites directory")
    staticDir   = flag.String("static", getDefaultStaticDir(), "Static files directory")
)

func getDefaultWebsitesDir() string {
    if runtime.GOOS == "windows" {
        return "C:/websites"
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, "websites")
}
```

### Other Issues

**Security concerns:**
- File permissions `0777` is overly permissive
- No input validation on `ctx.Param()` values
- Direct file serving without sanitization could allow directory traversal

**Error handling:**
- Many errors logged but execution continues
- `serverLog.Panicln()` will crash the server on missing files

**Code quality:**
- Mixed case in directory names (`classServer` vs `ClassServer`)
- Unused imports commented out (should remove)
- Hardcoded IP address `192.168.1.28` won't work for others

### Minimal Fix for Cross-Platform

At minimum, add this at the top:
```go
import (
    "path/filepath"
    "runtime"
)

var websitesRoot = func() string {
    if runtime.GOOS == "windows" {
        return "C:/websites"
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, "websites")
}()
```

Then replace all hardcoded paths with variables. The current code is Windows-only and will fail immediately on other operating systems.