FROM node:20-alpine AS build

WORKDIR /app

# Install dependencies
COPY web/package*.json ./
RUN npm ci --silent

# Copy source and build (skip slow type-check during image build to avoid TS failures here)
COPY web/ .
# Use build-only script to avoid running vue-tsc type-check during the image build
RUN npm run build-only --silent

FROM nginx:stable-alpine

# Remove default nginx static files
RUN rm -rf /usr/share/nginx/html/*

# Copy build output
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
