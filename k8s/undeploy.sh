#!/bin/bash

# Script to undeploy the magic-video application from Kubernetes
# Usage: ./undeploy.sh

set -e

echo "Undeploying magic-video application from Kubernetes..."

# Remove Ingress
echo "Removing Ingress..."
kubectl delete -f ingress.yaml --ignore-not-found=true

# Remove API deployment
echo "Removing API deployment..."
kubectl delete -f api-deployment.yaml --ignore-not-found=true

# Remove migration job
echo "Removing migration job..."
kubectl delete -f migration-job.yaml --ignore-not-found=true

# Remove MySQL deployment
echo "Removing MySQL deployment..."
kubectl delete -f mysql-deployment.yaml --ignore-not-found=true

# Remove storage (optional - comment out if you want to keep data)
echo "Removing storage..."
kubectl delete -f mysql-pv.yaml --ignore-not-found=true

# Remove ConfigMaps
echo "Removing ConfigMaps..."
kubectl delete -f api-configmap.yaml --ignore-not-found=true
kubectl delete -f mysql-configmap.yaml --ignore-not-found=true

# Remove namespace (this will remove everything in the namespace)
echo "Removing namespace..."
kubectl delete -f namespace.yaml --ignore-not-found=true

echo "Undeployment completed successfully!"
echo ""
echo "Note: If you want to keep the database data, comment out the storage removal section in this script."
