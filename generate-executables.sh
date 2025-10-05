# macOS Apple Silicon (M1/M2/M3/M4) - ARM64
GOOS=darwin GOARCH=arm64 go build -o ClassServer-mac-arm64

# macOS Intel - AMD64
GOOS=darwin GOARCH=amd64 go build -o ClassServer-mac-amd64

# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o ClassServer-windows.exe

# Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o ClassServer-linux-amd64

# Linux ARM64 (Raspberry Pi, etc.)
GOOS=linux GOARCH=arm64 go build -o ClassServer-linux-arm64