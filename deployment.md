Plan: Production AWS Deployment with Auto-Scaling
TL;DR: Deploy Neko-Sync to AWS ECS Fargate with production-grade infrastructure: multi-AZ RDS PostgreSQL, ElastiCache Redis, Application Load Balancers, S3 file storage, Route53 DNS with ACM SSL certificates, CloudWatch monitoring, and GitHub Actions CI/CD pipeline. This setup provides high availability, auto-scaling, and automated deployments with an estimated cost of $200-400/month.

Steps
Phase 1: Prepare Application for Production

Update backend config in apps/backend/internal/config/config.go:

Add AWS-specific environment variables (AWS_REGION, S3_BUCKET_NAME, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
Add production database SSL mode configuration
Add comprehensive logging configuration for CloudWatch
Add health check configuration with database and Redis connectivity checks
Create S3 file upload service in apps/backend/internal/infrastructure/storage/:

Implement S3 client for file uploads (replace local filesystem in UPLOAD_DIR)
Add file service interface in domain layer
Update upload handlers in apps/backend/internal/interfaces/http/handlers/ to use S3
Add signed URLs for private file access
Update frontend configuration in apps/web/next.config.js:

Configure for production builds with optimizations
Set up image optimization domains for S3 bucket
Configure environment variable validation
Create production Docker images:

Update Dockerfile with multi-stage builds for minimal image size
Add health check commands to Docker configurations
Create apps/web/Dockerfile optimized for Next.js production
Add database migration management:

Create apps/backend/migrations/ directory
Convert scripts/init-db.sql to versioned migrations
Add golang-migrate integration in cmd/nekosync/main.go
Create migration runner for ECS task definition
Phase 2: AWS Infrastructure Setup (Infrastructure as Code)

Create Terraform configuration in new infrastructure/terraform/ directory:

VPC Module: Create VPC with 3 availability zones, public subnets (ALB), private subnets (ECS, RDS, ElastiCache), NAT gateways, route tables
RDS Module: PostgreSQL 15 Multi-AZ instance (db.t3.medium or larger), automated backups, parameter groups with production settings, subnet groups
ElastiCache Module: Redis cluster with Multi-AZ replica, subnet groups, parameter groups with AOF persistence
ECS Module: ECS cluster, task definitions for backend and frontend, Fargate compute, IAM roles for task execution (S3, CloudWatch, Secrets Manager access)
ALB Module: Application Load Balancers with HTTPS listeners, target groups for backend (8080) and frontend (3000), health checks, WAF integration
S3 Module: Private bucket for uploads with versioning, lifecycle policies, CORS configuration for web access
Route53 Module: Hosted zone configuration, A records for API and web subdomains pointing to ALBs
ACM Module: SSL certificates for domain and wildcard subdomain with Route53 DNS validation
CloudWatch Module: Log groups for services, metric alarms (CPU, memory, error rates), SNS topics for alerts
Auto Scaling Module: Target tracking policies for ECS services (CPU > 70%, memory > 80%), scale out/in configurations
Secrets Manager: Store sensitive configs (DATABASE_URL, REDIS_URL, JWT_SECRET)
ECR Repositories: Create repositories for backend and frontend images
Create Terraform variables file (infrastructure/terraform/terraform.tfvars):

Define AWS region, domain name, database credentials, environment tags, instance sizes
Initialize and apply infrastructure:

Run terraform init and terraform plan to preview changes
Apply with terraform apply to provision all resources (takes 15-30 minutes)
Save outputs (ALB DNS, RDS endpoint, Redis endpoint, ECR URLs) to .env.production
Phase 3: Secrets and Configuration Management

Store secrets in AWS Secrets Manager:

Create secret for DATABASE_URL with RDS connection string (include SSL mode)
Create secret for REDIS_URL with ElastiCache endpoint
Create secret for JWT_SECRET (generate strong random key)
Store third-party API keys (ANIME_API_URL, MUSIC_API_URL tokens if applicable)
Configure ECS task definitions with environment variables:

Backend: PORT=8080, ENVIRONMENT=production, AWS_REGION, S3_BUCKET_NAME, CORS_ALLOWED_ORIGINS (your domain)
Frontend: NEXT_PUBLIC_API_URL (your API domain with HTTPS)
Both: Reference Secrets Manager ARNs for sensitive values
Phase 4: CI/CD Pipeline Setup

Create GitHub Actions workflow in .github/workflows/deploy-production.yml:

Trigger: Push to main branch, manual workflow dispatch
Backend job: Build Go binary, run tests, build Docker image, push to ECR, update ECS service
Frontend job: Run Next.js build, build Docker image, push to ECR, update ECS service
Database migration job: Run migrations as ECS one-off task before deployment
Use AWS credentials from GitHub Secrets (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION)
Create GitHub repository secrets:

Add AWS IAM user credentials with ECR and ECS permissions
Add production environment variables
Configure deployment approvals for production environment
Add deployment scripts in scripts/:

scripts/deploy.sh: Helper script for manual deployments
scripts/rollback.sh: Quick rollback to previous ECS task definition
scripts/run-migrations.sh: Run database migrations on AWS
Phase 5: Initial Deployment

Build and push initial Docker images:

Authenticate to ECR: aws ecr get-login-password | docker login
Build backend: docker build -t nekosync-backend -f apps/backend/Dockerfile .
Tag and push: docker tag nekosync-backend:latest [ECR_URL]/nekosync-backend:latest
Repeat for frontend
Run database migrations:

Execute ECS task with migration command
Verify tables created in RDS instance
Update ECS services to run tasks:

Backend service: Desired count 2 (for high availability across AZs)
Frontend service: Desired count 2
Verify tasks reach healthy status in ECS console
Configure DNS records in Route53:

Point api.yourdomain.com → Backend ALB
Point www.yourdomain.com and yourdomain.com → Frontend ALB
Wait for DNS propagation (5-10 minutes)
Verify SSL certificates are attached to ALB listeners and test HTTPS access

Phase 6: Production Hardening

Configure CloudWatch monitoring:

Set up dashboards for ECS metrics, RDS performance, Redis hit rates
Create alarms for: High CPU/memory, elevated error rates (5xx responses), database connections, disk space
Configure SNS email notifications for critical alerts
Set up automated backups:

Enable RDS automated backups (7-day retention minimum)
Configure RDS snapshot schedule
Enable S3 bucket versioning and lifecycle policies (archive old uploads to Glacier after 90 days)
Implement security hardening:

Review security groups (ensure database/Redis only accessible from ECS security group)
Enable VPC Flow Logs for network monitoring
Configure WAF rules on ALB (rate limiting, SQL injection protection, XSS)
Enable AWS GuardDuty for threat detection
Implement least-privilege IAM policies
Configure auto-scaling policies:

Backend: Scale out when CPU > 70% or memory > 80%, scale in when < 30%
Frontend: Scale out when CPU > 60% or memory > 75%
Set min capacity: 2, max capacity: 10 (adjust based on expected load)
Set up centralized logging:

Configure CloudWatch Logs agent in ECS tasks
Create log insights queries for error tracking
Enable CloudWatch Container Insights for ECS cluster
Phase 7: Post-Deployment Testing

Verify all services are running:

Check ECS console for healthy task count
Test backend health endpoint: https://api.yourdomain.com/health
Test frontend loads: https://www.yourdomain.com
Verify WebSocket connections work through ALB
Test file upload functionality:

Upload test file through web interface
Verify file appears in S3 bucket
Test file download/access with signed URLs
Load testing (optional but recommended):

Use tools like k6, Artillery, or Locust
Simulate expected traffic patterns
Verify auto-scaling triggers correctly
Identify bottlenecks before real traffic
Create runbook documentation in docs/RUNBOOK.md:

Deployment procedures
Rollback procedures
Troubleshooting common issues
Monitoring dashboard links
On-call escalation procedures
Verification
Automated Checks:

GitHub Actions pipeline succeeds on push to main
ECS health checks pass (backend and frontend tasks healthy)
RDS connection successful from backend tasks
Redis connectivity verified
Manual Testing:

Access https://www.yourdomain.com - frontend loads correctly
Access https://api.yourdomain.com/health - returns 200 OK status
User registration/login works (tests JWT, database, Redis)
File upload works (tests S3 integration)
WebSocket connection establishes (tests ALB WebSocket support)
Check CloudWatch dashboards show metrics
Trigger manual scaling test (increase load) and verify auto-scaling
Security Verification:

Run AWS Trusted Advisor checks
Verify SSL Labs rating (A or higher for HTTPS endpoints)
Confirm database not publicly accessible
Test WAF rules block malicious requests
Decisions
ECS Fargate over EKS: Fargate provides simpler management without Kubernetes complexity while still offering container orchestration and auto-scaling needed for production
RDS Multi-AZ over Aurora: More cost-effective for medium traffic while providing high availability; can migrate to Aurora later if needed
S3 over EFS: Better scalability and durability for user uploads; no shared filesystem complexity
Terraform over CloudFormation: More readable, better state management, easier to extend/modify
Application Load Balancer over Network Load Balancer: Supports HTTP/HTTPS routing, WebSocket, and integrates with ACM for SSL
Secrets Manager over Parameter Store: Better secret rotation capabilities and audit logging for production credentials
CloudWatch over third-party: Native AWS integration, no additional vendors, simpler setup
Estimated Timeline
Application preparation: 4-6 hours
Infrastructure setup: 3-4 hours
CI/CD configuration: 2-3 hours
Initial deployment: 1-2 hours
Production hardening: 3-4 hours
Testing and documentation: 2-3 hours
Total: 15-22 hours (can be spread over several days)

Cost Breakdown (Monthly)
ECS Fargate (4 tasks average): ~$50-80
RDS PostgreSQL Multi-AZ (db.t3.medium): ~$80-120
ElastiCache Redis (cache.t3.small): ~$30-40
ALB (2 load balancers): ~$35-45
NAT Gateway: ~$35-45
S3 storage and requests: ~$5-20
Data transfer: ~$10-30
CloudWatch logs/metrics: ~$10-20
Total: ~$255-400/month (varies with traffic and storage)
