#!/bin/bash

PROJECT_NAME="checkhost-cli"
COMPILED_FILE="./checkhost"

echo "Compiling $PROJECT_NAME..."
go build -o "$COMPILED_FILE" $PROJECT_NAME
if [ $? -ne 0 ]; then
    echo "Error: Compilation failed!"
    exit 1
fi

echo "Compilation successful!"
exit 0
