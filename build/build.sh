#!/bin/bash

# Function to handle the build process with optimizations
build() {
    export GOOS=$1
    export GOARCH=amd64
    FILETYPE=$2
    echo "Building for $GOOS..."

    # Disable CGO and set garbage collection threshold
    export CGO_ENABLED=0
    export GOGC=100

    go build -o "../bin/$GOOS/tiny-proxy$FILETYPE" -ldflags "-s -w" ../src/cmd/main.go
    if [ $? -ne 0 ]; then
        echo "Build failed for $GOOS!"
        exit 1
    fi
    echo "Build succeeded for $GOOS!"
}

# Menu to select OS
echo "Choose the OS you want to build for:"
echo "1. Windows"
echo "2. Linux"
echo "3. Darwin (macOS)"
echo "4. All"
read -p "Enter your choice (1-4): " choice

case $choice in
    1)
        build windows .exe
        ;;
    2)
        build linux ""
        ;;
    3)
        build darwin ""
        ;;
    4)
        build windows .exe
        build linux ""
        build darwin ""
        ;;
    *)
        echo "Invalid choice."
        ;;
esac

echo "Build process completed."
