# Stage 1: Build React
FROM node:alpine as builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Prepare Nginx
FROM nginx:alpine
# เอาไฟล์ที่ Build เสร็จจาก Stage 1 มาวาง
COPY --from=builder /app/dist /usr/share/nginx/html
# (สังเกตว่าเรา ไม่ต้อง COPY nginx.conf ตรงนี้แล้ว)
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
