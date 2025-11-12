# WhatsApp Bot Setup Guide

## Prerequisites

### System Requirements
- **Go**: Version 1.21 or higher
- **PostgreSQL**: Version 12 or higher
- **Redis**: Version 6 or higher
- **Docker & Docker Compose**: For containerized deployment
- **Git**: For version control

### WhatsApp Business API Requirements
- **Facebook Business Account**
- **WhatsApp Business Account**
- **Phone Number** (must not be registered with WhatsApp)
- **Meta Developer Account**

## Installation

### 1. Clone the Repository
```bash
git clone https://github.com/your-repo/whatsapp-bot.git
cd whatsapp-bot
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Environment Variables
Copy the example environment file:
```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=whatsapp_bot

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=7d

# WhatsApp Business API Configuration
WHATSAPP_API_KEY=your-whatsapp-api-key
WHATSAPP_PHONE_NUMBER_ID=your-phone-number-id
WHATSAPP_BUSINESS_ACCOUNT_ID=your-business-account-id
WHATSAPP_WEBHOOK_VERIFY_TOKEN=your-webhook-verify-token

# Application Configuration
APP_ENV=development
APP_PORT=8080
APP_NAME=WhatsApp Bot
APP_URL=http://localhost:8080

# External Services
WEATHER_API_KEY=your-weather-api-key
TRANSLATION_API_KEY=your-translation-api-key
CURRENCY_API_KEY=your-currency-api-key

# Security
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Monitoring
PROMETHEUS_ENABLED=true
PROMETHEUS_PORT=9090
GRAFANA_PORT=3000
```

## Database Setup

### 1. Create PostgreSQL Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE whatsapp_bot;

# Create user (optional)
CREATE USER whatsapp_bot WITH PASSWORD 'your-password';
GRANT ALL PRIVILEGES ON DATABASE whatsapp_bot TO whatsapp_bot;
```

### 2. Run Database Migration
```bash
go run cmd/migrate/main.go
```

Or using the application:
```bash
go run main.go migrate
```

## WhatsApp Business API Setup

### 1. Create Facebook Business Account
1. Go to [Facebook Business Manager](https://business.facebook.com/)
2. Create a new business account
3. Verify your business (if required)

### 2. Setup WhatsApp Business Account
1. Go to [WhatsApp Business Platform](https://developers.facebook.com/apps/)
2. Create a new app
3. Add WhatsApp product to your app
4. Setup WhatsApp Business account

### 3. Configure Phone Number
1. Add a phone number to your WhatsApp Business account
2. Verify the phone number (you'll receive a verification code)
3. Wait for approval from Meta

### 4. Get API Credentials
1. In your Facebook app, go to WhatsApp > Getting Started
2. Copy your:
   - Phone Number ID
   - WhatsApp Business Account ID
   - Access Token (API Key)
3. Update these in your `.env` file

### 5. Setup Webhook
1. In your Facebook app, go to WhatsApp > Configuration
2. Add webhook URL: `https://your-domain.com/webhooks/whatsapp`
3. Set verify token (same as in `.env`)
4. Subscribe to events:
   - messages
   - message_deliveries
   - message_reads

## Running the Application

### Development Mode
```bash
# Run with hot reload
go run main.go

# Or using air (if installed)
air
```

### Production Mode
```bash
# Build the application
go build -o whatsapp-bot

# Run the binary
./whatsapp-bot
```

### Using Docker
```bash
# Build and run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Configuration

### Application Settings
Edit `internal/config/config.go` to modify:
- Database connection settings
- Redis configuration
- JWT settings
- Rate limiting rules
- CORS policies

### Middleware Configuration
Edit files in `internal/middleware/` to customize:
- Authentication behavior
- Rate limiting rules
- CORS settings
- Logging format

### Service Configuration
Edit files in `internal/services/` to modify:
- WhatsApp API integration
- Business logic rules
- External API integrations
- Data processing logic

## Testing

### Run Unit Tests
```bash
go test ./test -v
```

### Run Integration Tests
```bash
go test ./test -v -tags=integration
```

### Run Specific Test
```bash
go test ./test -v -run TestUserRegistration
```

### Test Coverage
```bash
go test ./test -cover
go test ./test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Deployment

### Docker Deployment
```bash
# Build Docker image
docker build -t whatsapp-bot:latest .

# Run with Docker Compose
docker-compose up -d

# Scale services
docker-compose up -d --scale app=3
```

### Kubernetes Deployment
```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods
kubectl get services
```

### Manual Deployment
1. Build the application: `go build -o whatsapp-bot`
2. Copy binary to server
3. Setup PostgreSQL and Redis
4. Configure environment variables
5. Setup reverse proxy (Nginx/Apache)
6. Setup SSL certificate
7. Run the application

## Monitoring

### Prometheus Metrics
Access Prometheus dashboard at: `http://localhost:9090`

### Grafana Dashboard
Access Grafana dashboard at: `http://localhost:3000`
- Default credentials: admin/admin123

### Application Logs
Logs are stored in:
- Console output
- File: `logs/app.log`
- JSON format for structured logging

### Health Check
Check application health:
```bash
curl http://localhost:8080/health
```

## Security Best Practices

### 1. Environment Variables
- Never commit `.env` file to version control
- Use different secrets for different environments
- Rotate secrets regularly

### 2. Database Security
- Use strong passwords
- Enable SSL/TLS for database connections
- Regular backups
- Limit database access

### 3. API Security
- Use HTTPS in production
- Implement proper rate limiting
- Validate all input data
- Use CORS appropriately

### 4. WhatsApp API Security
- Keep API keys secure
- Monitor webhook requests
- Implement message validation
- Use proper error handling

## Troubleshooting

### Common Issues

#### Database Connection Failed
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check connection
psql -h localhost -U postgres -d whatsapp_bot

# Check logs
tail -f /var/log/postgresql/postgresql.log
```

#### Redis Connection Failed
```bash
# Check Redis status
sudo systemctl status redis

# Test connection
redis-cli ping

# Check logs
tail -f /var/log/redis/redis.log
```

#### WhatsApp API Issues
```bash
# Check API credentials
curl -X GET "https://graph.facebook.com/v18.0/YOUR_PHONE_NUMBER_ID" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Check webhook verification
curl -X GET "http://localhost:8080/webhooks/whatsapp?hub.verify_token=YOUR_VERIFY_TOKEN&hub.challenge=TEST"
```

#### Application Won't Start
```bash
# Check logs
tail -f logs/app.log

# Check port availability
netstat -tlnp | grep :8080

# Check environment variables
go run main.go config
```

### Debug Mode
Enable debug logging:
```env
LOG_LEVEL=debug
```

### Performance Issues
1. Check database query performance
2. Monitor Redis memory usage
3. Review application metrics
4. Check system resources

## Maintenance

### Regular Tasks
1. **Database Maintenance**
   - Vacuum and analyze tables
   - Update statistics
   - Check for slow queries

2. **Log Rotation**
   - Setup logrotate for application logs
   - Monitor disk space usage

3. **Security Updates**
   - Keep dependencies updated
   - Apply security patches
   - Review access logs

4. **Backup Strategy**
   - Database backups
   - Configuration backups
   - User data backups

### Monitoring Alerts
Setup alerts for:
- High error rates
- Database connection issues
- Memory usage spikes
- API rate limit exceeded
- Webhook failures

## Support

### Getting Help
- Check documentation: `API_DOCUMENTATION.md`
- Review logs: `logs/app.log`
- Check GitHub issues
- Contact support team

### Contributing
1. Fork the repository
2. Create feature branch
3. Make changes
4. Add tests
5. Submit pull request

## License
See LICENSE file for details.

## Changelog
See CHANGELOG.md for version history.