# Stage 1: Build React
FROM node:alpine as builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm config set fetch-retries 5 && \
    npm config set fetch-retry-factor 2 && \
    npm config set fetch-retry-mintimeout 10000 && \
    npm config set fetch-retry-maxtimeout 60000 && \
    npm install
COPY . .
RUN npm run build

# Stage 2: Prepare Nginx
FROM nginx:alpine
# เอาไฟล์ที่ Build เสร็จจาก Stage 1 มาวาง
COPY --from=builder /app/dist /usr/share/nginx/html
# (สังเกตว่าเรา ไม่ต้อง COPY nginx.conf ตรงนี้แล้ว)
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
