# Deployment Guide

This guide covers deploying the Neko-Sync backend to various environments.

## Prerequisites

- Docker and Docker Compose
- PostgreSQL 15+
- Go 1.24+ (for building from source)
- Access to target deployment environment

## Environment Configuration

### Environment Variables

Create appropriate `.env` files for each environment:

#### Development (.env.development)
```env
# Server Configuration
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=debug

# Database
DATABASE_URL=postgres://nekosync_user:password@localhost:5432/nekosync_dev?sslmode=disable

# Authentication
JWT_SECRET=your-development-jwt-secret-key
JWT_EXPIRY=24h

# External APIs
TMDB_API_KEY=your-tmdb-api-key
MAL_CLIENT_ID=your-mal-client-id

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:8080

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_HOUR=1000

# File Upload
MAX_UPLOAD_SIZE=10MB
UPLOAD_PATH=./uploads

# Email (if implemented)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

#### Production (.env.production)
```env
# Server Configuration
PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info

# Database (use connection pooling in production)
DATABASE_URL=postgres://nekosync_user:secure_password@db-host:5432/nekosync_prod?sslmode=require&pool_max_conns=20

# Authentication (use strong secrets)
JWT_SECRET=your-very-secure-production-jwt-secret-key-at-least-32-chars
JWT_EXPIRY=24h

# External APIs
TMDB_API_KEY=your-production-tmdb-api-key
MAL_CLIENT_ID=your-production-mal-client-id

# CORS (restrict to your domains)
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_HOUR=500

# File Upload
MAX_UPLOAD_SIZE=5MB
UPLOAD_PATH=/app/uploads

# Monitoring
SENTRY_DSN=your-sentry-dsn
PROMETHEUS_ENABLED=true
```

## Docker Deployment

### Single Container Deployment

#### 1. Build the image
```bash
# Development build
docker build -t nekosync-backend:dev .

# Production build
docker build -t nekosync-backend:prod --target production .
```

#### 2. Run with Docker
```bash
# Development
docker run -d \
  --name nekosync-backend \
  --env-file .env.development \
  -p 8080:8080 \
  nekosync-backend:dev

# Production
docker run -d \
  --name nekosync-backend \
  --env-file .env.production \
  -p 8080:8080 \
  --restart unless-stopped \
  nekosync-backend:prod
```

### Docker Compose Deployment

#### Development Stack
```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  backend:
    build:
      context: .
      target: development
    ports:
      - "8080:8080"
    env_file:
      - .env.development
    volumes:
      - .:/app
      - /app/bin
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: nekosync_dev
      POSTGRES_USER: nekosync_user
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

#### Production Stack
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  backend:
    build:
      context: .
      target: production
    ports:
      - "8080:8080"
    env_file:
      - .env.production
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: nekosync_prod
      POSTGRES_USER: nekosync_user
      POSTGRES_PASSWORD_FILE: /run/secrets/postgres_password
    secrets:
      - postgres_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    depends_on:
      - backend

secrets:
  postgres_password:
    external: true

volumes:
  postgres_data:
  redis_data:
```

#### Deploy with Docker Compose
```bash
# Development
docker-compose -f docker-compose.dev.yml up -d

# Production
docker-compose -f docker-compose.prod.yml up -d
```

## Cloud Deployment

### AWS ECS

#### 1. Create Task Definition
```json
{
  "family": "nekosync-backend",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::account:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "nekosync-backend",
      "image": "your-account.dkr.ecr.region.amazonaws.com/nekosync-backend:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "ENVIRONMENT",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "DATABASE_URL",
          "valueFrom": "arn:aws:ssm:region:account:parameter/nekosync/database-url"
        },
        {
          "name": "JWT_SECRET",
          "valueFrom": "arn:aws:ssm:region:account:parameter/nekosync/jwt-secret"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/nekosync-backend",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

#### 2. Deploy to ECS
```bash
# Register task definition
aws ecs register-task-definition --cli-input-json file://task-definition.json

# Create or update service
aws ecs create-service \
  --cluster nekosync-cluster \
  --service-name nekosync-backend \
  --task-definition nekosync-backend:1 \
  --desired-count 2 \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-12345,subnet-67890],securityGroups=[sg-12345],assignPublicIp=ENABLED}"
```

### Google Cloud Run

#### 1. Build and push image
```bash
# Build for Cloud Run
docker build -t gcr.io/your-project/nekosync-backend:latest .

# Push to Container Registry
docker push gcr.io/your-project/nekosync-backend:latest
```

#### 2. Deploy to Cloud Run
```bash
gcloud run deploy nekosync-backend \
  --image gcr.io/your-project/nekosync-backend:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10 \
  --set-env-vars ENVIRONMENT=production \
  --set-secrets DATABASE_URL=database-url:latest,JWT_SECRET=jwt-secret:latest
```

### Kubernetes

#### 1. Create ConfigMap
```yaml
# configmap.yml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nekosync-config
data:
  ENVIRONMENT: "production"
  PORT: "8080"
  LOG_LEVEL: "info"
```

#### 2. Create Secret
```yaml
# secret.yml
apiVersion: v1
kind: Secret
metadata:
  name: nekosync-secrets
type: Opaque
data:
  DATABASE_URL: <base64-encoded-database-url>
  JWT_SECRET: <base64-encoded-jwt-secret>
```

#### 3. Create Deployment
```yaml
# deployment.yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nekosync-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nekosync-backend
  template:
    metadata:
      labels:
        app: nekosync-backend
    spec:
      containers:
      - name: nekosync-backend
        image: your-registry/nekosync-backend:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: nekosync-config
        - secretRef:
            name: nekosync-secrets
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

#### 4. Create Service
```yaml
# service.yml
apiVersion: v1
kind: Service
metadata:
  name: nekosync-backend-service
spec:
  selector:
    app: nekosync-backend
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

#### 5. Deploy to Kubernetes
```bash
kubectl apply -f configmap.yml
kubectl apply -f secret.yml
kubectl apply -f deployment.yml
kubectl apply -f service.yml
```

## Database Setup

### PostgreSQL Setup

#### 1. Create Database
```sql
-- Connect as superuser
CREATE DATABASE nekosync_prod;
CREATE USER nekosync_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE nekosync_prod TO nekosync_user;

-- Connect to nekosync_prod database
\c nekosync_prod;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
GRANT ALL ON SCHEMA public TO nekosync_user;
```

#### 2. Run Migrations
```bash
# Using the application (if migration system is implemented)
./bin/nekosync migrate

# Or manually run SQL files
psql -U nekosync_user -d nekosync_prod -f migrations/001_initial_schema.sql
```

### Database Backup and Restore

#### Backup
```bash
# Full backup
pg_dump -U nekosync_user -h localhost nekosync_prod > backup_$(date +%Y%m%d_%H%M%S).sql

# Compressed backup
pg_dump -U nekosync_user -h localhost nekosync_prod | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz

# Docker backup
docker exec postgres_container pg_dump -U nekosync_user nekosync_prod > backup.sql
```

#### Restore
```bash
# Restore from backup
psql -U nekosync_user -d nekosync_prod < backup.sql

# Restore compressed backup
gunzip -c backup.sql.gz | psql -U nekosync_user -d nekosync_prod
```

## Load Balancing and High Availability

### Nginx Configuration
```nginx
# nginx.conf
upstream backend {
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}

server {
    listen 80;
    listen 443 ssl;
    server_name yourdomain.com;

    # SSL configuration
    ssl_certificate /etc/ssl/certs/yourdomain.crt;
    ssl_certificate_key /etc/ssl/private/yourdomain.key;

    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";

    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    location /health {
        proxy_pass http://backend/health;
        access_log off;
    }
}
```

## Monitoring and Logging

### Health Checks
The application provides a health check endpoint at `/health`:

```bash
# Check application health
curl http://localhost:8080/health

# Expected response
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0",
  "services": {
    "database": "healthy"
  }
}
```

### Logging Configuration
Configure structured logging for production:

```go
// Example logging configuration
logger := log.New().WithFields(log.Fields{
    "service": "nekosync-backend",
    "version": version,
    "environment": environment,
})
```

### Prometheus Metrics
If Prometheus is enabled, metrics are available at `/metrics`:

```bash
curl http://localhost:8080/metrics
```

## SSL/TLS Configuration

### Let's Encrypt with Certbot
```bash
# Install certbot
sudo apt-get install certbot

# Get certificate
sudo certbot certonly --standalone -d yourdomain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### Custom SSL Certificate
```bash
# Generate self-signed certificate (development only)
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

## Deployment Checklist

### Pre-deployment
- [ ] Environment variables configured
- [ ] Database setup and migrations run
- [ ] SSL certificates configured
- [ ] Health checks working
- [ ] Load testing completed
- [ ] Backup strategy in place

### Deployment
- [ ] Build and test Docker image
- [ ] Deploy to staging environment
- [ ] Run integration tests
- [ ] Deploy to production
- [ ] Verify health checks
- [ ] Test critical API endpoints

### Post-deployment
- [ ] Monitor application logs
- [ ] Check database performance
- [ ] Verify SSL certificate
- [ ] Test load balancer
- [ ] Set up monitoring alerts
- [ ] Document any issues

## Rollback Strategy

### Quick Rollback
```bash
# Docker Compose
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d

# Kubernetes
kubectl rollout undo deployment/nekosync-backend

# ECS
aws ecs update-service --cluster nekosync-cluster --service nekosync-backend --task-definition nekosync-backend:previous-version
```

### Database Rollback
```bash
# Restore from backup
psql -U nekosync_user -d nekosync_prod < backup_before_deployment.sql
```

## Security Considerations

### Production Security
- Use strong, randomly generated secrets
- Enable HTTPS/TLS encryption
- Configure proper CORS origins
- Set up rate limiting
- Use secure database connections
- Regularly update dependencies
- Monitor for security vulnerabilities
- Implement proper logging and monitoring

### Environment Isolation
- Use separate databases for each environment
- Use different API keys for external services
- Implement proper access controls
- Use secrets management systems
- Regular security audits
