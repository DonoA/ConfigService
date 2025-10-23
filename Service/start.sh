#!/bin/bash

# Stop the script if any command fails
set -e

# --- Configuration ---
# The main package for your project
MAIN_PACKAGE="."

# The desired name for the output binary
BINARY_NAME="config-service"
# ---------------------

echo "Building Go project..."

# Build the project
# -o specifies the output file
# ${MAIN_PACKAGE} tells go build which package to build
go build -o ${BINARY_NAME} ${MAIN_PACKAGE}

# Check if the build was successful (set -e handles this, but explicit check is good)
if [ ! -f "${BINARY_NAME}" ]; then
    echo "Build failed! Binary not found."
    exit 1
fi

echo "Build successful. Binary created: ${BINARY_NAME}"

# Make the binary executable (optional, go build usually does this)
chmod +x ${BINARY_NAME}

echo "Running the project..."

# Run the compiled binary
# ./ indicates the file is in the current directory
./${BINARY_NAME}
