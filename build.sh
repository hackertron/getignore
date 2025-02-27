#!/bin/bash
echo "Building gitignore CLI tool..."
go build -o gitignore main.go
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi
chmod +x gitignore
echo "Build successful! You can now use ./gitignore"
