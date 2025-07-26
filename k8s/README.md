# Kubernetes Deployment for Magic Video

This directory contains Kubernetes manifests and scripts to deploy the magic-video application.

## Prerequisites

- Kubernetes cluster (minikube, kind, or cloud provider)
- kubectl configured to access your cluster
- Docker for building images
- NGINX Ingress Controller installed in your cluster

### Installing NGINX Ingress Controller

For minikube:
```bash
minikube addons enable ingress
```

For other clusters:
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml
```

## Quick Start

1. **Build the Docker image:**
   ```bash
   cd k8s
   ./build.sh
   ```

   For minikube, load the image:
   ```bash
   minikube image load magic-video-api:latest
   ```

2. **Deploy the application:**
   ```bash
   ./deploy.sh
   ```

3. **Add to your hosts file:**
   ```bash
   echo "127.0.0.1 magic-video.local" | sudo tee -a /etc/hosts
   ```

4. **Access the application:**
   - API: http://magic-video.local/api
   - Frontend (when deployed): http://magic-video.local/

## Components

### Database
- **MySQL 8.0** with persistent storage
- Configured with database `magic_video`
- Automatic schema migration on deployment

### API
- **Go application** with 2 replicas for high availability
- Health checks and resource limits configured
- Automatic database connection waiting

### Storage
- **PersistentVolume** for MySQL data persistence
- 10GB storage allocation

### Networking
- **Ingress** with NGINX for external access
- CORS configuration for frontend integration
- Ready for TLS/SSL configuration

## Scripts

### `build.sh`
Builds the Docker image for the API application.
```bash
./build.sh [tag]  # Default tag is 'latest'
```

### `deploy.sh`
Deploys the entire application stack to Kubernetes.
```bash
./deploy.sh
```

### `undeploy.sh`
Removes the application from Kubernetes.
```bash
./undeploy.sh
```

## Configuration

### Database Configuration
Edit `mysql-configmap.yaml` to change database settings:
- Database name
- Username/password
- Root password

### API Configuration
Edit `api-configmap.yaml` to change application settings:
- Database connection parameters
- Server port and host
- Other environment variables

## Monitoring and Debugging

### Check deployment status:
```bash
kubectl get pods -n magic-video
kubectl get services -n magic-video
kubectl get ingress -n magic-video
```

### View logs:
```bash
# API logs
kubectl logs -f deployment/magic-video-api -n magic-video

# MySQL logs
kubectl logs -f deployment/mysql -n magic-video

# Migration job logs
kubectl logs job/db-migration -n magic-video
```

### Access MySQL directly:
```bash
kubectl exec -it deployment/mysql -n magic-video -- mysql -u magic_user -p magic_video
```

## Scaling

### Scale API replicas:
```bash
kubectl scale deployment magic-video-api --replicas=3 -n magic-video
```

### Update application:
```bash
# Build new image
./build.sh v1.1.0

# Update deployment
kubectl set image deployment/magic-video-api api=magic-video-api:v1.1.0 -n magic-video
```

## Production Considerations

1. **Security:**
   - Use Kubernetes Secrets for sensitive data
   - Enable TLS/SSL certificates
   - Configure network policies
   - Use non-root containers

2. **Storage:**
   - Use cloud provider storage classes
   - Configure backup strategies
   - Monitor disk usage

3. **Monitoring:**
   - Add Prometheus metrics
   - Configure alerting
   - Set up log aggregation

4. **High Availability:**
   - Deploy across multiple nodes
   - Configure pod disruption budgets
   - Use multiple replicas for all components

## Future Enhancements

- Frontend deployment configuration
- Redis for caching
- Message queue for background jobs
- Horizontal Pod Autoscaler (HPA)
- Vertical Pod Autoscaler (VPA)
