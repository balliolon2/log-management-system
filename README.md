# Log Management System

A centralized log management system with multi-tenant support, real-time alerting, and secure SaaS capabilities. Built for the Cyber Defense Internship.

## Features

- **Multi-Protocol Ingestion**: Supports Syslog (UDP/TCP), HTTP API (JSON), and Batch File Ingestion.
- **Normalization**: Automatically parses and normalizes logs into a common schema.
- **Storage**: High-performance storage using PostgreSQL with JSONB indexing.
- **Alerting**: Real-time alerting engine with configurable rules (e.g., Brute Force Detection).
- **Retention Policy**: Automatic cleanup of logs older than 7 days.
- **Multi-Tenancy**: Data isolation per tenant using Row Level Security (RLS).
- **Secure Deployment**: TLS/HTTPS support and production-ready Docker configuration.

## Prerequisites

- **Docker** and **Docker Compose** installed on your machine.

## Quick Start (Appliance/Dev Mode)

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/balliolon2/log-management-system.git
    cd log-management-system
    ```

2.  **Generate Self-Signed Certificates** (Required for HTTPS):
    ```powershell
    # Windows PowerShell
    cd deploy
    docker run --rm -v ${PWD}/certs:/certs alpine/openssl req -x509 -newkey rsa:4096 -keyout /certs/key.pem -out /certs/cert.pem -days 365 -nodes -subj "/CN=localhost"
    ```

3.  **Setup Environment**:
    ```powershell
    cp .env.example .env
    ```

4.  **Start the System**:
    ```powershell
    docker-compose up --build
    ```

5.  **Access the Dashboard**:
    - Open `https://localhost` in your browser.
    - Accept the self-signed certificate warning.

## Production / SaaS Deployment

This system is designed to be deployed as a SaaS platform.

### Deploying to Cloud (AWS, DigitalOcean, Azure)

1.  **Provision a Linux Server** (Ubuntu 22.04 LTS Recommended).
2.  **Install Docker**.
3.  **Upload Code** or `git clone` the repository.
4.  **Configure `.env`**: Set secure passwords and a strong `JWT_SECRET`.
5.  **Run with Production Config**:
    The default `docker-compose.yml` is already optimized for production (hiding DB ports, etc.).
    ```bash
    cd deploy
    docker compose up -d --build
    ```

### Public Access for Demo (ngrok)

If you need to demonstrate the SaaS capability without a cloud server, use **ngrok**:

1.  Start the project locally (`docker-compose up`).
2.  Run ngrok to tunnel HTTPS:
    ```bash
    ngrok http https://localhost:443
    ```
3.  Share the generated `https://xxxx.ngrok-free.app` URL.

## API Usage

### Ingest Logs (Batch)

**Endpoint**: `POST /ingest`

```json
[
  {
    "tenant": "tenant_A",
    "event_type": "login_failed",
    "source": "app-01",
    "body": { "user": "admin", "src_ip": "10.0.0.1" }
  },
  {
    "tenant": "tenant_A",
    "event_type": "login_success",
    "source": "app-01",
    "body": { "user": "alice" }
  }
]
```

### File Ingestion (Raw Text)

To ingest a raw text log file (e.g., specific application logs):

```bash
cd samples
# Install requests if needed: pip install requests
python file_ingestor.py dummy.log
```

## Project Structure

```
/backend       # Go Backend
  /cmd         # Main entry point
  /internal    # Application Logic
    /alert     # Alert Engine & Rules
    /auth      # JWT Authentication & Middleware
    /handler   # HTTP API Handlers
    /ingest    # Syslog (UDP/TCP) & Ingestion Logic
    /job       # Background Jobs (Retention)
    /models    # Database Structs
    /normalize # Log Parsing & Normalization
    /repository # Database Access Layer
/frontend      # React Frontend (Vite + Tailwind)
/deploy        # Deployment Config
  /certs       # SSL Certificates (Self-signed)
  .env         # Environment Variables (Secrets)
  docker-compose.yml # Production-ready Compose
  nginx.conf   # Reverse Proxy & TLS Config
  init.sql     # DB Schema & Seed Data
/tests         # Python Simulation Scripts
```
