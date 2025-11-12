# WhatsApp Bot API Documentation

## Overview
This API provides comprehensive functionality for managing a WhatsApp Business bot with features including auto-replies, broadcasts, games, business tools, and more.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Rate Limiting
- General API endpoints: 10 requests per second
- Authentication endpoints: 5 requests per second

## Response Format
All responses follow this structure:
```json
{
  "success": true,
  "data": {},
  "error": null
}
```

## Endpoints

### Authentication

#### Register User
**POST** `/auth/register`
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "password": "password123"
}
```

#### Login
**POST** `/auth/login`
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

#### Refresh Token
**POST** `/auth/refresh`
```json
{
  "refresh_token": "your-refresh-token"
}
```

#### Forgot Password
**POST** `/auth/forgot-password`
```json
{
  "email": "john@example.com"
}
```

#### Reset Password
**POST** `/auth/reset-password`
```json
{
  "token": "reset-token",
  "new_password": "newpassword123"
}
```

### User Management

#### Get Profile
**GET** `/users/profile`
**Headers:** Authorization: Bearer <token>

#### Update Profile
**PUT** `/users/profile`
```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com",
  "phone": "+1234567891",
  "avatar": "https://example.com/avatar.jpg",
  "bio": "Updated bio",
  "language": "en"
}
```

#### Change Password
**POST** `/users/change-password`
```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword123"
}
```

#### Get Settings
**GET** `/users/settings`

#### Update Settings
**PUT** `/users/settings`
```json
{
  "notifications_enabled": true,
  "auto_reply_enabled": true,
  "game_notifications": true,
  "business_alerts": true,
  "language": "en",
  "timezone": "Asia/Jakarta"
}
```

### WhatsApp Management

#### Send Message
**POST** `/whatsapp/send-message`
```json
{
  "to": "+1234567890",
  "message": "Hello from WhatsApp Bot!"
}
```

#### Send Media
**POST** `/whatsapp/send-media`
```json
{
  "to": "+1234567890",
  "media_url": "https://example.com/image.jpg",
  "media_type": "image",
  "caption": "Check out this image"
}
```

#### Send Template
**POST** `/whatsapp/send-template`
```json
{
  "to": "+1234567890",
  "template_name": "welcome_message",
  "language": "en",
  "parameters": ["John", "Welcome"]
}
```

#### Get Contacts
**GET** `/whatsapp/contacts`

#### Get Chats
**GET** `/whatsapp/chats`

#### Get Chat Messages
**GET** `/whatsapp/chat/{chat_id}`

#### Mark Messages as Read
**POST** `/whatsapp/mark-read`
```json
{
  "message_ids": ["msg1", "msg2"]
}
```

#### Get WhatsApp Status
**GET** `/whatsapp/status`

### Auto-Reply Management

#### Get Auto-Replies
**GET** `/auto-replies?user_id={user_id}&status={status}`

#### Create Auto-Reply
**POST** `/auto-replies`
```json
{
  "name": "Greeting Reply",
  "trigger": "hello",
  "response": "Hello! How can I help you?",
  "match_type": "exact",
  "keywords": ["hello", "hi"],
  "is_active": true,
  "user_id": "user-uuid"
}
```

#### Get Auto-Reply
**GET** `/auto-replies/{reply_id}`

#### Update Auto-Reply
**PUT** `/auto-replies/{reply_id}`
```json
{
  "name": "Updated Greeting",
  "response": "Hello! Welcome to our service!",
  "is_active": true
}
```

#### Delete Auto-Reply
**DELETE** `/auto-replies/{reply_id}`

#### Toggle Auto-Reply
**POST** `/auto-replies/{reply_id}/toggle`

### Broadcast Management

#### Get Broadcasts
**GET** `/broadcasts?user_id={user_id}&status={status}`

#### Create Broadcast
**POST** `/broadcasts`
```json
{
  "name": "Holiday Greetings",
  "message": "Happy holidays! Wishing you joy and prosperity.",
  "recipients": ["+1234567890", "+1234567891"],
  "schedule_at": "2024-12-25 10:00:00",
  "user_id": "user-uuid"
}
```

#### Get Broadcast
**GET** `/broadcasts/{broadcast_id}`

#### Update Broadcast
**PUT** `/broadcasts/{broadcast_id}`
```json
{
  "name": "Updated Holiday Message",
  "message": "Updated holiday greetings!"
}
```

#### Delete Broadcast
**DELETE** `/broadcasts/{broadcast_id}`

#### Send Broadcast
**POST** `/broadcasts/{broadcast_id}/send`

#### Get Broadcast Stats
**GET** `/broadcasts/{broadcast_id}/stats`

### Game Management

#### Get Available Games
**GET** `/games`

#### Start Game
**POST** `/games/start`
```json
{
  "game_type": "trivia",
  "user_id": "user-uuid"
}
```

#### Play Game
**POST** `/games/{game_id}/play`
```json
{
  "user_id": "user-uuid",
  "answer": "correct answer"
}
```

#### Get Leaderboard
**GET** `/games/{game_id}/leaderboard`

#### Get Game Stats
**GET** `/games/{game_id}/stats`

#### Start Trivia
**POST** `/games/trivia/start`
```json
{
  "category": "general",
  "difficulty": "medium",
  "user_id": "user-uuid"
}
```

#### Answer Trivia
**POST** `/games/trivia/answer`
```json
{
  "game_id": "game-uuid",
  "question_id": "question-uuid",
  "answer": "correct answer",
  "user_id": "user-uuid"
}
```

### Utility Services

#### Create QR Code
**POST** `/utils/qr-code`
```json
{
  "data": "https://example.com",
  "size": 200,
  "format": "png"
}
```

#### Create Short Link
**POST** `/utils/short-link`
```json
{
  "original_url": "https://example.com/very/long/url",
  "custom_alias": "short",
  "expiry_days": 7
}
```

#### Get Short Link
**GET** `/utils/short-link/{alias}`

#### Convert Currency
**POST** `/utils/currency-convert`
```json
{
  "amount": 100,
  "from": "USD",
  "to": "EUR"
}
```

#### Get Weather
**GET** `/utils/weather?city=Jakarta&country=Indonesia`

#### Translate Text
**POST** `/utils/translate`
```json
{
  "text": "Hello",
  "target_lang": "id",
  "source_lang": "en"
}
```

#### Get Location Info
**POST** `/utils/location-info`
```json
{
  "latitude": -6.2088,
  "longitude": 106.8456
}
```

#### Create Poll
**POST** `/utils/polls`
```json
{
  "question": "What's your favorite color?",
  "options": ["Red", "Blue", "Green"],
  "user_id": "user-uuid"
}
```

#### Vote on Poll
**POST** `/utils/polls/{poll_id}/vote`
```json
{
  "option_id": 0,
  "user_id": "user-uuid"
}
```

#### Get Poll Results
**GET** `/utils/polls/{poll_id}/results`

#### Create Reminder
**POST** `/utils/reminders`
```json
{
  "title": "Meeting Reminder",
  "description": "Team meeting at 2 PM",
  "remind_at": "2024-01-15 14:00:00",
  "user_id": "user-uuid",
  "repeat": "none"
}
```

#### Get Reminders
**GET** `/utils/reminders?user_id={user_id}`

#### Delete Reminder
**DELETE** `/utils/reminders/{reminder_id}`

#### Create Note
**POST** `/utils/notes`
```json
{
  "title": "Meeting Notes",
  "content": "Discussed project timeline and deliverables",
  "user_id": "user-uuid",
  "category": "Work",
  "tags": ["meeting", "project"]
}
```

#### Get Notes
**GET** `/utils/notes?user_id={user_id}&category={category}&tag={tag}`

#### Update Note
**PUT** `/utils/notes/{note_id}`
```json
{
  "title": "Updated Meeting Notes",
  "content": "Updated content",
  "category": "Work",
  "tags": ["meeting", "updated"]
}
```

#### Delete Note
**DELETE** `/utils/notes/{note_id}`

#### Search Notes
**GET** `/utils/notes/search?user_id={user_id}&query={query}`

#### Create Timer
**POST** `/utils/timers`
```json
{
  "name": "Pomodoro Timer",
  "duration": 1500,
  "user_id": "user-uuid"
}
```

#### Get Timer
**GET** `/utils/timers/{timer_id}`

#### Stop Timer
**DELETE** `/utils/timers/{timer_id}`

#### Create File Upload
**POST** `/utils/file-upload`
```json
{
  "filename": "document.pdf",
  "file_size": 1024000,
  "file_type": "application/pdf",
  "user_id": "user-uuid"
}
```

#### Get File Upload
**GET** `/utils/file-upload/{upload_id}`

#### Delete File Upload
**DELETE** `/utils/file-upload/{upload_id}`

### Business Management

#### Create Product
**POST** `/business/products`
```json
{
  "name": "Smartphone",
  "description": "Latest model smartphone",
  "price": 999.99,
  "category": "Electronics",
  "image_url": "https://example.com/phone.jpg",
  "stock": 50,
  "user_id": "user-uuid"
}
```

#### Get Products
**GET** `/business/products?category={category}&search={search}`

#### Get Product
**GET** `/business/products/{product_id}`

#### Update Product
**PUT** `/business/products/{product_id}`
```json
{
  "name": "Updated Smartphone",
  "price": 899.99,
  "stock": 75
}
```

#### Delete Product
**DELETE** `/business/products/{product_id}`

#### Create Order
**POST** `/business/orders`
```json
{
  "customer_name": "John Doe",
  "customer_phone": "+1234567890",
  "items": [
    {
      "product_id": "product-uuid",
      "quantity": 2
    }
  ],
  "user_id": "user-uuid"
}
```

#### Get Orders
**GET** `/business/orders?status={status}&user_id={user_id}`

#### Get Order
**GET** `/business/orders/{order_id}`

#### Update Order Status
**PUT** `/business/orders/{order_id}/status`
```json
{
  "status": "shipped"
}
```

#### Create Invoice
**POST** `/business/invoices`
```json
{
  "order_id": "order-uuid",
  "amount": 1999.98,
  "description": "Invoice for smartphone order",
  "due_date": "2024-02-15",
  "user_id": "user-uuid"
}
```

#### Get Invoices
**GET** `/business/invoices?status={status}&user_id={user_id}`

#### Get Invoice
**GET** `/business/invoices/{invoice_id}`

#### Update Invoice Status
**PUT** `/business/invoices/{invoice_id}/status`
```json
{
  "status": "paid"
}
```

#### Create Customer
**POST** `/business/customers`
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "address": "123 Main St, City",
  "user_id": "user-uuid"
}
```

#### Get Customers
**GET** `/business/customers?user_id={user_id}&search={search}`

#### Get Customer
**GET** `/business/customers/{customer_id}`

#### Update Customer
**PUT** `/business/customers/{customer_id}`
```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com",
  "phone": "+1234567891",
  "address": "456 New St, City"
}
```

#### Delete Customer
**DELETE** `/business/customers/{customer_id}`

#### Create Payment
**POST** `/business/payments`
```json
{
  "invoice_id": "invoice-uuid",
  "amount": 1999.98,
  "payment_method": "credit_card",
  "reference": "PAY123",
  "user_id": "user-uuid"
}
```

#### Get Payments
**GET** `/business/payments?invoice_id={invoice_id}&user_id={user_id}`

#### Get Payment
**GET** `/business/payments/{payment_id}`

#### Get Business Stats
**GET** `/business/stats?user_id={user_id}`

#### Get Sales Report
**GET** `/business/sales-report?user_id={user_id}&start_date=2024-01-01&end_date=2024-01-31`

### Moderation

#### Get Blocked Users
**GET** `/moderation/blocked-users?user_id={user_id}`

#### Block User
**POST** `/moderation/block-user`
```json
{
  "user_id": "user-uuid",
  "blocked_user_id": "blocked-user-uuid",
  "reason": "Spam behavior"
}
```

#### Unblock User
**DELETE** `/moderation/unblock-user/{user_id}`

#### Get Reported Messages
**GET** `/moderation/reported-messages?status={status}`

#### Report Message
**POST** `/moderation/report-message`
```json
{
  "message_id": "message-uuid",
  "reporter_id": "reporter-uuid",
  "reason": "inappropriate_content",
  "details": "Contains offensive language"
}
```

#### Update Report Status
**PUT** `/moderation/reported-messages/{report_id}`
```json
{
  "status": "resolved",
  "notes": "Issue resolved after warning"
}
```

#### Get Spam Detections
**GET** `/moderation/spam-detections?user_id={user_id}`

#### Mark as Spam
**POST** `/moderation/mark-spam`
```json
{
  "message_id": "message-uuid",
  "user_id": "user-uuid",
  "reason": "repetitive content"
}
```

#### Get Content Filters
**GET** `/moderation/content-filters?user_id={user_id}`

#### Create Content Filter
**POST** `/moderation/content-filters`
```json
{
  "name": "Profanity Filter",
  "type": "keyword",
  "pattern": "badword",
  "action": "block",
  "replacement": "",
  "is_active": true,
  "user_id": "user-uuid"
}
```

#### Update Content Filter
**PUT** `/moderation/content-filters/{filter_id}`
```json
{
  "name": "Updated Filter",
  "pattern": "newbadword",
  "action": "replace",
  "replacement": "***"
}
```

#### Delete Content Filter
**DELETE** `/moderation/content-filters/{filter_id}`

### Analytics

#### Get Dashboard
**GET** `/analytics/dashboard?user_id={user_id}&time_range={time_range}`

#### Get Message Analytics
**GET** `/analytics/messages?user_id={user_id}&time_range={time_range}&group_by={group_by}`

#### Get User Analytics
**GET** `/analytics/users?time_range={time_range}&group_by={group_by}`

#### Get Game Analytics
**GET** `/analytics/games?time_range={time_range}&game_type={game_type}`

#### Get Business Analytics
**GET** `/analytics/business?user_id={user_id}&time_range={time_range}`

#### Export Analytics
**GET** `/analytics/export?user_id={user_id}&time_range={time_range}&format={format}`

### Admin Endpoints

#### Get Admin Dashboard
**GET** `/admin/dashboard`

#### Get Users
**GET** `/admin/users`

#### Get User
**GET** `/admin/users/{user_id}`

#### Update User
**PUT** `/admin/users/{user_id}`

#### Delete User
**DELETE** `/admin/users/{user_id}`

#### Ban User
**POST** `/admin/users/{user_id}/ban`

#### Unban User
**POST** `/admin/users/{user_id}/unban`

#### Get System Stats
**GET** `/admin/system-stats`

#### Get Logs
**GET** `/admin/logs`

#### Clear Logs
**DELETE** `/admin/logs`

#### Get Settings
**GET** `/admin/settings`

#### Update Settings
**PUT** `/admin/settings`

#### Create Backup
**POST** `/admin/backup`

#### Get Backups
**GET** `/admin/backups`

#### Restore Backup
**POST** `/admin/backups/{backup_id}/restore`

#### Delete Backup
**DELETE** `/admin/backups/{backup_id}`

#### Get System Health
**GET** `/admin/system-health`

#### Set Maintenance Mode
**POST** `/admin/system-maintenance`

#### Get Broadcast Messages
**GET** `/admin/broadcast-messages`

#### Create Broadcast Message
**POST** `/admin/broadcast-messages`

#### Delete Broadcast Message
**DELETE** `/admin/broadcast-messages/{message_id}`

#### Get Spam Reports
**GET** `/admin/spam-reports`

#### Update Spam Report
**PUT** `/admin/spam-reports/{report_id}`

#### Get Content Reports
**GET** `/admin/content-reports`

#### Update Content Report
**PUT** `/admin/content-reports/{report_id}`

### Webhooks

#### WhatsApp Webhook
**POST** `/webhooks/whatsapp`

#### WhatsApp Webhook Verification
**GET** `/webhooks/whatsapp`

### Health Check

#### Health Status
**GET** `/health`

## Error Codes

| Code | Description |
|------|-------------|
| 400  | Bad Request - Invalid input data |
| 401  | Unauthorized - Invalid or missing token |
| 403  | Forbidden - Insufficient permissions |
| 404  | Not Found - Resource not found |
| 429  | Too Many Requests - Rate limit exceeded |
| 500  | Internal Server Error - Server error |

## Common Parameters

### Pagination
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)

### Filtering
- `status`: Filter by status (pending, active, completed, etc.)
- `user_id`: Filter by user ID
- `search`: Search query

### Time Range
- `time_range`: today, week, month
- `start_date`: Start date (YYYY-MM-DD)
- `end_date`: End date (YYYY-MM-DD)

## WebSocket Events
The API supports WebSocket connections for real-time updates:
- Message notifications
- Game updates
- Broadcast status
- System alerts

## Rate Limiting Headers
Response includes rate limiting information:
```
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 8
X-RateLimit-Reset: 1640995200
```

## CORS
The API supports Cross-Origin Resource Sharing (CORS) for web applications.

## Security
- All endpoints use HTTPS in production
- JWT tokens expire after 24 hours
- Refresh tokens expire after 7 days
- Passwords are hashed using bcrypt
- Rate limiting prevents abuse
- Input validation and sanitization

## Examples

### Complete User Registration and Login Flow
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'

# Use token for protected endpoints
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <your-token>"
```

### Send WhatsApp Message
```bash
curl -X POST http://localhost:8080/api/v1/whatsapp/send-message \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "+1234567890",
    "message": "Hello from API!"
  }'
```

### Create Auto-Reply
```bash
curl -X POST http://localhost:8080/api/v1/auto-replies \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Greeting Reply",
    "trigger": "hello",
    "response": "Hello! How can I help?",
    "match_type": "exact",
    "keywords": ["hello", "hi"],
    "user_id": "user-uuid"
  }'
```

## SDKs and Libraries
- **Go**: Native implementation
- **JavaScript/Node.js**: Available on npm
- **Python**: Available on PyPI
- **PHP**: Available on Packagist

## Support
For API support, please contact:
- Email: support@whatsapp-bot.com
- Documentation: https://docs.whatsapp-bot.com
- Issues: https://github.com/whatsapp-bot/api/issues