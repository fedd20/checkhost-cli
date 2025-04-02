#!/bin/bash

PROJECT_NAME="checkhost-cli"
COMPILED_FILE="./checkhost"

echo "Compiling $PROJECT_NAME..."
go build -o "$COMPILED_FILE" $PROJECT_NAME
if [ $? -ne 0 ]; then
    echo "Error: Compilation failed!"
    exit 1
fi

move_to_bin() {
    local TARGET_DIR=$1
    if mv "$COMPILED_FILE" "$TARGET_DIR/"; then
        echo "$PROJECT_NAME installed successfully to $TARGET_DIR!"
        echo "Now you can run $PROJECT_NAME from anywhere in your terminal."
        exit 0
    fi
}

move_to_bin "/usr/local/bin"
move_to_bin "/usr/bin"

echo "Error: Failed to move $PROJECT_NAME to any bin folder!"
echo "Please ensure you have the necessary permissions to write to these directories."
echo "...Or try to install it manually."
exit 1
