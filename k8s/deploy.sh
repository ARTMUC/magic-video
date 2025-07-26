#!/bin/bash

# Script to deploy the magic-video application to Kubernetes
# Usage: ./deploy.sh

set -e

echo "Deploying magic-video application to Kubernetes..."

# Apply namespace first
echo "Creating namespace..."
kubectl apply -f namespace.yaml

# Apply ConfigMaps
echo "Applying ConfigMaps..."
kubectl apply -f mysql-configmap.yaml
kubectl apply -f api-configmap.yaml

# Apply PersistentVolume and PersistentVolumeClaim
echo "Creating storage..."
kubectl apply -f mysql-pv.yaml

# Deploy MySQL
echo "Deploying MySQL..."
kubectl apply -f mysql-deployment.yaml

# Wait for MySQL to be ready
echo "Waiting for MySQL to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/mysql -n magic-video

# Run database migrations
echo "Running database migrations..."
kubectl apply -f migration-job.yaml

# Wait for migration job to complete
echo "Waiting for migrations to complete..."
kubectl wait --for=condition=complete --timeout=300s job/db-migration -n magic-video

# Deploy API
echo "Deploying API..."
kubectl apply -f api-deployment.yaml

# Wait for API to be ready
echo "Waiting for API to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/magic-video-api -n magic-video

# Apply Ingress
echo "Applying Ingress..."
kubectl apply -f ingress.yaml

echo "Deployment completed successfully!"
echo ""
echo "To access the application:"
echo "1. Add '127.0.0.1 magic-video.local' to your /etc/hosts file"
echo "2. Access the API at: http://magic-video.local/api"
echo ""
echo "To check the status:"
echo "kubectl get pods -n magic-video"
echo "kubectl get services -n magic-video"
echo "kubectl logs -f deployment/magic-video-api -n magic-video"
