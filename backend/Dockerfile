# ใช้ Go Image ตัวเล็ก
FROM golang:alpine

# ติดตั้งของจำเป็น
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy โค้ดทั้งหมด
COPY . .

# Build เป็นไฟล์ binary ชื่อ server
RUN go build -o server cmd/main.go

# เปิด Port 8080 (API) และ 514 (Syslog UDP)
EXPOSE 8080
EXPOSE 514/udp

# รันโปรแกรม
CMD ["./server"]
