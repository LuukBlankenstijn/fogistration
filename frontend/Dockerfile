# syntax=docker/dockerfile:1.7

# --- Build stage ---
FROM node:20-alpine AS build
WORKDIR /app

# Enable corepack for yarn 4
RUN corepack enable

# Copy manifests first for better caching
COPY package.json yarn.lock .yarnrc.yml ./
COPY .yarn ./.yarn

# Install dependencies
RUN yarn install --immutable


# Copy the rest of the app
COPY . .

# generate client
RUN yarn gen-client

# Build (vite outputs to dist/)
RUN yarn build


# --- Serve stage ---
FROM nginx:alpine AS serve

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy build output
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
