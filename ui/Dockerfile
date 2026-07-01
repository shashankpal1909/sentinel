# Build stage
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci || npm install

COPY . .

ARG VITE_ADMIN_API_URL=http://localhost:9901
ENV VITE_ADMIN_API_URL=${VITE_ADMIN_API_URL}

RUN npm run build

# Production serve stage
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
