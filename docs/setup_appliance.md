# Setup Guide: Appliance Mode (On-Premise / Local)

This guide describes how to deploy the Log Management System on a single server or local machine for testing/development.

## Prerequisites
- Docker & Docker Compose
- PowerShell (Windows) or Bash (Linux/Mac)

## Installation Steps

1.  **Clone the Repository**
    ```bash
    git clone https://github.com/balliolon2/log-management-system.git
    cd log-management-system
    ```

2.  **Certificate Setup (Self-Signed)**
    Required for HTTPS. Run this command in `deploy/`:
    ```powershell
    cd deploy
    docker run --rm -v ${PWD}/certs:/certs alpine/openssl req -x509 -newkey rsa:4096 -keyout /certs/key.pem -out /certs/cert.pem -days 365 -nodes -subj "/CN=localhost"
    ```

3.  **Environment Setup**
    ```powershell
    cp .env.example .env
    ```
    *Note: For appliance mode, the default values in `.env.example` are sufficient.*

4.  **Start the System**
    ```powershell
    docker-compose up --build -d
    ```

5.  **Access the System**
    - Go to `https://localhost`
    - **Login**:
        - Admin: `admin` / `password` (See `init.sql` for seeds)
        - Viewer: `viewer1` / `password`

## Troubleshooting
- If **Connection Refused**: Check if containers are running `docker ps`.
- If **Certificate Warning**: This is expected for self-signed certs. Click "Proceed".
