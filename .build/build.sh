#!/bin/sh

set -e

CUR_DIR="./.build"
TARGET_DIR="target"
BUILD_DIR="${CUR_DIR}/${TARGET_DIR}"

COMMAND_PATH="$1"
COMMAND_NAME="$2"
TAGS="$3"

TARGET_NAME="${BUILD_DIR}/${COMMAND_NAME}"

mkdir -p "${BUILD_DIR}"
echo "Building: ${COMMAND_PATH}..."
# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags="${TAGS}" -o "${TARGET_NAME}" "${COMMAND_PATH}"
CGO_ENABLED=0 go build -tags="${TAGS}" -o "${TARGET_NAME}" "${COMMAND_PATH}"
