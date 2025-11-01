# Stage 1: Builder
# Install dependencies and build the app
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

# Define a build-time argument for the secret
ARG NEXT_PUBLIC_LIVEKIT_URL

# Make it available as an environment variable during the build
ENV NEXT_PUBLIC_LIVEKIT_URL=$NEXT_PUBLIC_LIVEKIT_URL

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