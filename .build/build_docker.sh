#!/bin/sh

set -e

COMMAND_PATH="$1"
COMMAND_NAME="$2"
IMAGE_NAME="$3"
IMAGE_TAG="$4"
TAGS="$5"

docker buildx build --build-arg COMMAND_PATH="${COMMAND_PATH}" --build-arg COMMAND_NAME="${COMMAND_NAME}" --build-arg TAGS="${TAGS}" --tag "${IMAGE_NAME}:${IMAGE_TAG}" .
