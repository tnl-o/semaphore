# Deployment Guide

This guide covers deployment of Semaphore UI in various environments and configurations.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Deployment Methods](#deployment-methods)
- [Configuration](#configuration)
- [High Availability](#high-availability)
- [Production Considerations](#production-considerations)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **CPU**: 2+ cores recommended
- **Memory**: 2GB+ RAM (4GB+ for production)
- **Disk**: 10GB+ free space
- **Network**: Internet access for downloading dependencies

### Software Requirements

- **Operating System**: Linux, macOS, or Windows
- **Database**: MySQL 8.0+, PostgreSQL 12+, or SQLite 3
- **Optional**: Redis (for High Availability mode)

## Quick Start

### Docker (Recommended)

The fastest way to get started:

```bash
docker run -d -p 3000:3000 --name semaphore \
  -e SEMAPHORE_DB_DIALECT=bolt \
  -e SEMAPHORE_ADMIN=admin \
  -e SEMAPHORE_ADMIN_PASSWORD=changeme \
  -e SEMAPHORE_ADMIN_NAME=Admin \
  -e SEMAPHORE_ADMIN_EMAIL=admin@localhost \
  semaphoreui/semaphore:latest
```

Access at: http://localhost:3000

### Binary Installation

1. Download the binary for your platform from [releases](https://github.com/semaphoreui/semaphore/releases)
2. Make it executable: `chmod +x semaphore`
3. Run setup: `./semaphore setup`
4. Start server: `./semaphore server --config ./config.json`

## Deployment Methods

### 1. Docker

#### Using Docker Compose

**Basic Setup:**

```yaml
version: '3.8'
services:
  semaphore:
    image: semaphoreui/semaphore:latest
    ports:
      - "3000:3000"
    environment:
      - SEMAPHORE_DB_DIALECT=bolt
      - SEMAPHORE_ADMIN=admin
      - SEMAPHORE_ADMIN_PASSWORD=changeme
      - SEMAPHORE_ADMIN_NAME=Admin
      - SEMAPHORE_ADMIN_EMAIL=admin@localhost
    volumes:
      - semaphore_data:/etc/semaphore
volumes:
  semaphore_data:
```

**With MySQL:**

```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: semaphore
      MYSQL_USER: semaphore
      MYSQL_PASSWORD: semaphorepassword
    volumes:
      - mysql_data:/var/lib/mysql

  semaphore:
    image: semaphoreui/semaphore:latest
    ports:
      - "3000:3000"
    environment:
      - SEMAPHORE_DB_DIALECT=mysql
      - SEMAPHORE_DB_HOST=mysql
      - SEMAPHORE_DB_USER=semaphore
      - SEMAPHORE_DB_PASS=semaphorepassword
      - SEMAPHORE_DB_NAME=semaphore
      - SEMAPHORE_ADMIN=admin
      - SEMAPHORE_ADMIN_PASSWORD=changeme
    depends_on:
      - mysql
    volumes:
      - semaphore_data:/etc/semaphore

volumes:
  mysql_data:
  semaphore_data:
```

See [deployment/compose/README.md](../deployment/compose/README.md) for more examples.

#### Using Pre-built Images

Images are available on [Docker Hub](https://hub.docker.com/r/semaphoreui/semaphore):

```bash
docker pull semaphoreui/semaphore:latest
```

### 2. Binary Installation

#### Linux/macOS

1. Download the latest release:
```bash
wget https://github.com/semaphoreui/semaphore/releases/latest/download/semaphore_linux_amd64.tar.gz
tar -xzf semaphore_linux_amd64.tar.gz
```

2. Run setup:
```bash
./semaphore setup
```

3. Start server:
```bash
./semaphore server --config ./config.json
```

#### Windows

1. Download the Windows binary
2. Run setup in PowerShell:
```powershell
.\semaphore.exe setup
```

3. Start server:
```powershell
.\semaphore.exe server --config .\config.json
```

### 3. Package Installation

#### Debian/Ubuntu

```bash
wget https://github.com/semaphoreui/semaphore/releases/latest/download/semaphore_amd64.deb
sudo dpkg -i semaphore_amd64.deb
sudo systemctl start semaphore
```

#### RHEL/CentOS

```bash
wget https://github.com/semaphoreui/semaphore/releases/latest/download/semaphore_amd64.rpm
sudo rpm -i semaphore_amd64.rpm
sudo systemctl start semaphore
```

### 4. Systemd Service

Create `/etc/systemd/system/semaphore.service`:

```ini
[Unit]
Description=Semaphore UI
After=network.target

[Service]
Type=simple
User=semaphore
WorkingDirectory=/opt/semaphore
ExecStart=/usr/local/bin/semaphore server --config /etc/semaphore/config.json
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable semaphore
sudo systemctl start semaphore
```

See [deployment/systemd/](../deployment/systemd/) for more details.

### 5. Kubernetes

**Basic Deployment:**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: semaphore
spec:
  replicas: 1
  selector:
    matchLabels:
      app: semaphore
  template:
    metadata:
      labels:
        app: semaphore
    spec:
      containers:
      - name: semaphore
        image: semaphoreui/semaphore:latest
        ports:
        - containerPort: 3000
        env:
        - name: SEMAPHORE_DB_DIALECT
          value: "mysql"
        - name: SEMAPHORE_DB_HOST
          value: "mysql-service"
        volumeMounts:
        - name: config
          mountPath: /etc/semaphore
      volumes:
      - name: config
        configMap:
          name: semaphore-config
---
apiVersion: v1
kind: Service
metadata:
  name: semaphore-service
spec:
  selector:
    app: semaphore
  ports:
  - port: 80
    targetPort: 3000
  type: LoadBalancer
```

## Configuration

### Environment Variables

Key environment variables:

- `SEMAPHORE_DB_DIALECT`: Database type (`mysql`, `postgres`, `sqlite`, `bolt`)
- `SEMAPHORE_DB_HOST`: Database host
- `SEMAPHORE_DB_USER`: Database user
- `SEMAPHORE_DB_PASS`: Database password
- `SEMAPHORE_DB_NAME`: Database name
- `SEMAPHORE_ADMIN`: Admin username
- `SEMAPHORE_ADMIN_PASSWORD`: Admin password
- `SEMAPHORE_ADMIN_EMAIL`: Admin email
- `SEMAPHORE_ADMIN_NAME`: Admin display name

See [util/config.go](../util/config.go) for all available options.

### Configuration File

Create `config.json`:

```json
{
  "mysql": {
    "host": "localhost:3306",
    "user": "semaphore",
    "pass": "password",
    "name": "semaphore"
  },
  "dialect": "mysql",
  "port": ":3000",
  "interface": "",
  "tmp_path": "/tmp/semaphore"
}
```

### Database Setup

#### MySQL

```sql
CREATE DATABASE semaphore CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'semaphore'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON semaphore.* TO 'semaphore'@'localhost';
FLUSH PRIVILEGES;
```

#### PostgreSQL

```sql
CREATE DATABASE semaphore;
CREATE USER semaphore WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE semaphore TO semaphore;
```

## High Availability

### Requirements

- Multiple Semaphore instances
- Shared database (MySQL/PostgreSQL)
- Redis for shared state

### Configuration

Enable HA mode in `config.json`:

```json
{
  "ha": {
    "enabled": true,
    "redis": {
      "addr": "redis:6379",
      "db": 0,
      "pass": "redispassword"
    }
  }
}
```

### Deployment

1. Deploy multiple Semaphore instances
2. Configure shared database
3. Configure Redis
4. Use load balancer (nginx, HAProxy, etc.)

**Example with nginx:**

```nginx
upstream semaphore {
    least_conn;
    server semaphore1:3000;
    server semaphore2:3000;
}

server {
    listen 80;
    server_name semaphore.example.com;

    location / {
        proxy_pass http://semaphore;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Production Considerations

### Security

1. **Change default passwords**: Always change admin password
2. **Use HTTPS**: Configure reverse proxy with SSL/TLS
3. **Firewall**: Restrict access to necessary ports only
4. **Database security**: Use strong passwords, limit network access
5. **Regular updates**: Keep Semaphore and dependencies updated

### Performance

1. **Database**: Use MySQL or PostgreSQL for better performance
2. **Connection pooling**: Configure appropriate pool sizes
3. **Caching**: Enable Redis caching (when available)
4. **Resource limits**: Set appropriate CPU/memory limits

### Monitoring

1. **Health checks**: Monitor `/api/ping` endpoint
2. **Logs**: Set up log aggregation
3. **Metrics**: Configure Prometheus metrics (when available)
4. **Alerts**: Set up alerts for critical issues

### Backup

1. **Database backups**: Regular database backups
2. **Configuration backups**: Backup `config.json`
3. **Disaster recovery**: Test restore procedures

## Troubleshooting

### Common Issues

**Issue: Cannot connect to database**

- Check database is running
- Verify connection credentials
- Check network connectivity
- Review firewall rules

**Issue: Port already in use**

- Change port in configuration
- Check for other processes using the port
- Use `netstat` or `lsof` to find process

**Issue: Permission denied**

- Check file permissions
- Verify user has access to directories
- Check SELinux/AppArmor settings

**Issue: High memory usage**

- Review task concurrency settings
- Check for memory leaks
- Increase available memory
- Review database connection pool

### Logs

Check logs for errors:

```bash
# Systemd
journalctl -u semaphore -f

# Docker
docker logs -f semaphore

# Direct
tail -f /var/log/semaphore.log
```

### Getting Help

- [GitHub Issues](https://github.com/semaphoreui/semaphore/issues)
- [Discord Community](https://discord.gg/5R6k7hNGcH)
- [Documentation](https://docs.semaphoreui.com)

## Resources

- [Architecture Documentation](./ARCHITECTURE.md)
- [Development Guide](./DEVELOPMENT.md)
- [Docker Deployment](../deployment/docker/README.md)
- [Docker Compose Examples](../deployment/compose/README.md)
- [Systemd Service](../deployment/systemd/README.md)

