## Dockerfile - frontend

`Stage 1: Builder
# Install dependencies and build the app
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

# Define a build-time argument for the secret
ARG NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY

# Make it available as an environment variable during the build
ENV NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=$NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY

RUN npm run build

# Stage 2: Production
# Create a smaller image for running the app
FROM node:20-alpine AS runner

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create the application directory
RUN mkdir -p /app

# Copy built assets from the builder stage to explicit paths
COPY --from=builder /app/public /app/public
COPY --from=builder /app/.next/standalone /app
COPY --from=builder /app/.next/static /app/.next/static

# Set ownership of the app directory
RUN chown -R appuser:appgroup /app

# Set the user to the non-root user
USER appuser

# Set the working directory
WORKDIR /app

# Expose the port the app runs on
EXPOSE 3000

# Set the default environment variables.
ENV NODE_ENV=production
ENV PORT=3000

# Start the server by sourcing the env file and then running the node server
CMD ["/bin/sh", "-c", "set -a && . /secrets/.env && set +a && node server.js"]


## Dockerfile - backend
# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server -ldflags="-w -s" ./main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Add a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create the application directory
RUN mkdir -p /app

# Explicitly copy the binary from the builder to the absolute path /app/server
COPY --from=builder /app/server /app/server

# Set ownership of the entire app directory
RUN chown -R appuser:appgroup /app

# Set the user
USER appuser

# Set the working directory
WORKDIR /app

# Expose the port
EXPOSE 8080

# Start the server using a relative path from the working directory
CMD ["/bin/sh", "-c", "set -a && . /secrets/.env && set +a && /app/server"]

## Cloudbuild.yaml

steps:
  # =====================================================================================
  # Step 1: Build Docker Images (Parallel)
  # =====================================================================================

  # Build the backend service
  - name: 'gcr.io/cloud-builders/docker'
    id: 'Build Backend'
    args:
      - 'build'
      - '-t'
      - 'us-central1-docker.pkg.dev/$PROJECT_ID/meetmind-webapp-docker-images/meetmind-backend:$COMMIT_SHA'
      - '.'
      - '-f'
      - 'Dockerfile'
    dir: 'backend'

  # Build the frontend service
  - name: 'gcr.io/cloud-builders/docker'
    id: 'Build Frontend'
    entrypoint: 'bash'
    dir: 'frontend_web'
    #### here the frontend need to be built with these env variable so that the webapp on the browser have access to these environment variables:
    NEXT_PUBLIC_LIVEKIT_URL="wss://con19.livekit.cloud"

######
    secretEnv: ['NEXT_PUBLIC_LIVEKIT_URL_SECRET']
    args:
    - '-c'
    - |
      docker build \
        --build-arg "NEXT_PUBLIC_LIVEKIT_URL=$$NEXT_PUBLIC_LIVEKIT_URL_SECRET" \
        -t "us-central1-docker.pkg.dev/$PROJECT_ID/meetmind-webapp-docker-images/meetmind-frontend:$COMMIT_SHA" \
        -f Dockerfile \
        .

  # =====================================================================================
  # Step 2: Push Docker Images (Parallel)
  # =====================================================================================

  # Push the backend service image
  - name: 'gcr.io/cloud-builders/docker'
    id: 'Push Backend'
    args: ['push', 'us-central1-docker.pkg.dev/$PROJECT_ID/meetmind-webapp-docker-images/meetmind-backend:$COMMIT_SHA']
    wait_for: ['Build Backend']

  # Push the frontend service image
  - name: 'gcr.io/cloud-builders/docker'
    id: 'Push Frontend'
    args: ['push', 'us-central1-docker.pkg.dev/$PROJECT_ID/meetmind-webapp-docker-images/meetmind-frontend:$COMMIT_SHA']
    wait_for: ['Build Frontend']

  # =====================================================================================
  # Step 3: Deploy to Cloud Run (Parallel)
  # =====================================================================================

  # Deploy the backend service
  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    id: 'Deploy Backend'
    entrypoint: 'gcloud'
    args:
      - 'run'
      - 'deploy'
      - 'meetmind-backend'
      - '--image'
      - 'us-central1-docker.pkg.dev/$PROJECT_ID/meetmind-webapp-docker-images/meetmind-backend:$COMMIT_SHA'
      - '--region'
      - 'us-central1'
      - '--platform'
      - 'managed'
      - '--allow-unauthenticated'
      - '--port'
      - '8080'
      - '--service-account=dial2ai-webapp-cloudrun@dial2ai.iam.gserviceaccount.com'
      - '--set-secrets=/secrets/.env=meetmind_webapp_golang_backend:latest'
    wait_for: ['Push Backend']

  # Deploy the frontend service
  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    id: 'Deploy Frontend'
    entrypoint: 'gcloud'
    args:
      - 'run'
      - 'deploy'
      - 'meetmind-frontend'
      - '--image'
      - 'us-central1-docker.pkg.dev/$PROJECT_ID/meetmind-webapp-docker-images/meetmind-frontend:$COMMIT_SHA'
      - '--region'
      - 'us-central1'
      - '--platform'
      - 'managed'
      - '--allow-unauthenticated'
      - '--port'
      - '3000'
      - '--service-account=dial2ai-webapp-cloudrun@dial2ai.iam.gserviceaccount.com'
      - '--set-secrets=/secrets/.env=meetmind_frontend_secret:latest'
    wait_for: ['Push Frontend']

# =====================================================================================
# Final Images
# =====================================================================================
images:
  - 'us-central1-docker.pkg.dev/$PROJECT_ID/meetmind-webapp-docker-images/meetmind-backend:$COMMIT_SHA'
  - 'us-central1-docker.pkg.dev/$PROJECT_ID/meetmind-webapp-docker-images/meetmind-frontend:$COMMIT_SHA'

# =====================================================================================
# Build Options
# =====================================================================================
options:
  logging: CLOUD_LOGGING_ONLY

# cloudbuild.yaml

# Define the secrets to be fetched from Secret Manager at the top level
availableSecrets:      
  secretManager:
    - versionName: projects/$PROJECT_ID/secrets/NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY/versions/latest
      env: STRIPE_PUBLISHABLE_KEY             # This path will be available to all steps
