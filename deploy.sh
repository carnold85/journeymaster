#!/bin/bash
# Exit on any error
set -e
echo "Deploying the container..."
docker run -d -p 8080:8080 --name journeymaster journeymaster:1.0
echo "Deployment completed successfully!"
