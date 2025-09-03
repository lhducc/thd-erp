# UAT Deployment Instructions (Docker Compose)

## Prerequisites

- Docker and Docker Compose installed on the server
- Access to the private container registry (login if required)
- SSL certificates (`server.crt` and `server.key`) placed in both `frontend/ssl/` and `backend/ssl/` directories
- `.env` file configured with UAT-safe values (never use production credentials)

## Steps

1. **Clone or update the repository on the UAT server:**
   ```bash
   git clone <your-repo-url> /opt/erp
   cd /opt/erp
   git pull
   ```

2. **Prepare environment variables:**
   - Copy and edit the `.env` file for UAT:
     ```bash
     cp example.env .env
     vim .env
     ```
   - Ensure all database and service credentials are UAT-specific.

3. **Place SSL certificates:**
   - Copy your SSL certificate and key to:
     - `erp/frontend/ssl/server.crt`
     - `erp/frontend/ssl/server.key`
     - `erp/backend/ssl/server.crt`
     - `erp/backend/ssl/server.key`

4. **Build and deploy the stack:**
   ```bash
   cd deploys
   docker compose -f docker-compose.uat.yml up -d --build
   ```

5. **Verify services:**
   ```bash
   docker compose -f docker-compose.uat.yml ps
   docker compose -f docker-compose.uat.yml logs -f backend
   docker compose -f docker-compose.uat.yml logs -f frontend
   ```

6. **Access the application:**
   Frontend: [https://192.168.1.58](https://192.168.1.58)

7. **Stop the stack (when needed):**
   ```bash
   docker compose -f docker-compose.uat.yml down
   ```

## Notes

- **Database data is persisted** via Docker named volumes. Data will not be lost unless you run `docker compose down -v` or manually remove volumes.
- **Never use production credentials** in UAT `.env` or config files.
- For troubleshooting, check logs with `docker compose logs -f <service>`.