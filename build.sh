#!/bin/bash
# Exit on any error
set -e
#echo "Building the Go application..."
#go build -o journeymaster
echo "Building the Docker image..."
docker build -t journeymaster:1.0 .
echo "Build completed successfully!"
