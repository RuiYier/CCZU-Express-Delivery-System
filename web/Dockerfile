FROM node:18-alpine AS build

WORKDIR /app

# Install dependencies
COPY web/package*.json ./
RUN npm ci --silent

# Copy source and build
COPY web/ .
RUN npm run build --silent

FROM nginx:stable-alpine

# Remove default nginx static files
RUN rm -rf /usr/share/nginx/html/*

# Copy build output
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
