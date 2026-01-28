# Setup Guide: SaaS / Cloud Deployment

This guide describes how to deploy the system for Production use on a Public Cloud (AWS, DigitalOcean, etc.).

## Security Features in SaaS Mode
- **Isolated Database**: Port 5432 is NOT exposed to the internet.
- **Secure Secrets**: Uses `.env` for passwords and keys.
- **TLS/HTTPS**: Enforced via Nginx.

## Installation Steps

1.  **Provision Server**: Ubuntu 22.04 LTS (Recommended 1 CPU / 2GB RAM).
2.  **Install Docker**:
    ```bash
    curl -fsSL https://get.docker.com | sh
    ```
3.  **Upload Code**:
    ```bash
    scp -r log-management-system root@<SERVER_IP>:~/
    ```
4.  **Configure Production Secrets**:
    ```bash
    cd ~/log-management-system/deploy
    cp .env.example .env
    nano .env
    ```
    *Change `DB_PASSWORD` and `JWT_SECRET` to strong random values.*

5.  **Generate/Import Certificates**:
    - **Option A: Let's Encrypt (Real CA)**: Install Certbot and map paths to `./certs`.
    - **Option B: Self-Signed (Testing)**:
      ```bash
      mkdir certs
      docker run --rm -v ${PWD}/certs:/certs alpine/openssl req -x509 -newkey rsa:4096 -keyout /certs/key.pem -out /certs/cert.pem -days 365 -nodes -subj "/CN=<YOUR_DOMAIN_OR_IP>"
      ```

6.  **Run in Production Mode**:
    The main `docker-compose.yml` is already optimized.
    ```bash
    docker compose up -d --build
    ```

## Verify Deployment
1.  Access `https://<YOUR_SERVER_IP>`.
2.  Verify Database Port is closed: `telnet <SERVER_IP> 5432` (Should fail).
