# syntax=docker/dockerfile:1

# ---- build stage ----
FROM node:22-alpine AS build
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# The browser calls the API on the same origin; the Ingress routes /api to the
# backend. That avoids CORS in the cluster and keeps one public entry point.
ARG VITE_API_URL=/api/v1
ARG VITE_APP_VERSION=dev
ENV VITE_API_URL=${VITE_API_URL} VITE_APP_VERSION=${VITE_APP_VERSION}
RUN npm run build

# ---- runtime stage: unprivileged nginx listens on 8080 as UID 101 ----
FROM nginxinc/nginx-unprivileged:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 8080
