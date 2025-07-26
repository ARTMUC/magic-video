#!/bin/bash

# Script to build the Docker image for magic-video API
# Usage: ./build.sh [tag]

set -e

# Default tag
TAG=${1:-latest}
IMAGE_NAME="magic-video-api"

echo "Building Docker image: $IMAGE_NAME:$TAG"

# Build the Docker image
docker build -t $IMAGE_NAME:$TAG -f ../Dockerfile ..

echo "Docker image built successfully: $IMAGE_NAME:$TAG"
echo ""
echo "To push to a registry:"
echo "docker tag $IMAGE_NAME:$TAG your-registry/$IMAGE_NAME:$TAG"
echo "docker push your-registry/$IMAGE_NAME:$TAG"
echo ""
echo "To load into minikube (if using minikube):"
echo "minikube image load $IMAGE_NAME:$TAG"
