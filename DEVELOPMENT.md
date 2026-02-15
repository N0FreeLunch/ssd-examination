# Development Guide

This guide describes how to set up the local development environment and run the application.

## Prerequisites

*   [Docker Desktop](https://www.docker.com/products/docker-desktop/)
*   [AWS CLI v2](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
*   [AWS Session Manager Plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)
*   `make` (GNU Make)

## 1. AWS Configuration (SSO)

This project uses AWS SSM Parameter Store for secrets management. You must configure AWS SSO first.

```bash
aws configure sso
# SSO Session Name: [your-session-name]
# Start URL: [your-start-url] (e.g., https://d-xxxxxxxxxx.awsapps.com/start)
# SSO Region: [your-region] (e.g., ap-northeast-2 or ap-northeast-1)
# Registration Scopes: sso:account:access
# Profile Name: sdd-exam (⚠️ IMPORTANT: Must match the profile used in Makefile)
```

**Verify login:**
```bash
aws sso login --profile sdd-exam
```

## 2. Running the Application (Local)

Use `make` commands to inject secrets automatically. **Do not run `docker-compose up` directly.**

```bash
# Start the application (local environment)
make up

# Rebuild containers (if dependencies changed)
make rebuild

# Stop the application
make down

# View logs
make logs
```

## 3. Infrastructure & Deployment

To deploy to the development server (AWS EC2):

```bash
# Deploy to dev environment
make deploy
```
(This will inject `dev` secrets and execute `deploy.sh`.)

## Troubleshooting

*   **`AccessDeniedException`**: Your SSO session might have expired. Run `aws sso login --profile sdd-exam`.
*   **`The config profile (sdd-exam) could not be found`**: Check `~/.aws/config` if the profile name matches.
