# Multi-Platform Bot (WhatsApp & Telegram) - Golang Implementation

Bot WhatsApp dan Telegram lengkap dengan 60+ fitur yang dibangun menggunakan Golang. Bot ini mendukung WhatsApp Business API dan Telegram Bot API dengan berbagai fitur interaktif, bisnis, produktivitas, dan keamanan.

## Fitur Utama

### 1. Fitur Dasar WhatsApp & Telegram
- Auto-reply dengan berbagai pola (exact, contains, regex)
- Broadcast message ke banyak kontak
- Manajemen grup WhatsApp & Telegram
- Media handling (gambar, audio, video, dokumen)
- Welcome & away messages
- Quick reply buttons
- Interactive menu dengan inline keyboards (Telegram)
- Poll creation dan voting (Telegram)
- Webhook dan polling support (Telegram)

### 2. Fitur Interaktif & Hiburan
- Game dan kuis dengan leaderboard
- Cek khodam / zodiak generator
- Love calculator
- Math challenge
- Jokes generator
- Story telling
- Tebak gambar

### 3. Fitur Produktivitas
- Reminder system
- Schedule message
- To-do list
- Weather forecast
- Currency converter
- Translator
- QR Code generator
- Password generator

### 4. Fitur Bisnis & E-commerce
- Product catalog
- Order management
- Payment integration
- Invoice generator
- Promo & discount system
- Loyalty program
- Business analytics

### 5. Fitur Keamanan & Moderasi
- Anti-spam protection
- Word filter & content moderation
- Anti-link protection
- Flood control
- Admin tools
- Report system

### 6. Fitur Kustomisasi
- Custom commands
- Bot personality
- Theme support
- Multi-language support
- Response templates
- API integration

## Teknologi yang Digunakan

- **Backend**: Golang 1.21+
- **Database**: PostgreSQL dengan GORM
- **Cache**: Redis
- **API Framework**: Gin
- **WhatsApp**: WhatsApp Business API
- **Telegram**: Telegram Bot API
- **Authentication**: JWT
- **Logging**: Logrus
- **Task Scheduling**: Cron
- **PDF Generation**: gofpdf
- **QR Code**: go-qrcode
- **Monitoring**: Prometheus & Grafana

## Persyaratan Sistem

- Go 1.21 atau lebih baru
- PostgreSQL 12+
- Redis 6+
- WhatsApp Business API Account
- Telegram Bot Token

## Instalasi

### 1. Clone Repository
```bash
git clone https://github.com/0xHadiRamdhani/chatbotcate
cd chatbotcate
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Environment
```bash
cp .env.example .env
# Edit .env dengan konfigurasi Anda
```

### 4. Setup Database
```bash
# Buat database PostgreSQL
createdb whatsapp_bot

# Jalankan migrasi otomatis saat aplikasi dijalankan
```

### 5. Jalankan Aplikasi
```bash
go run main.go
```

## Konfigurasi

### Environment Variables

```bash
# Server
SERVER_PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=whatsapp_bot

# WhatsApp Business API
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
WHATSAPP_ACCESS_TOKEN=your_access_token
WHATSAPP_WEBHOOK_SECRET=your_webhook_secret

# Telegram Bot API
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_WEBHOOK_URL=https://your-domain.com/api/telegram/webhook

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
```

## Dokumentasi API

### Authentication Endpoints

- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/refresh` - Refresh access token

### WhatsApp Endpoints

- `POST /api/v1/whatsapp/send` - Send message
- `POST /api/v1/whatsapp/broadcast` - Broadcast message
- `GET /api/v1/whatsapp/contacts` - Get contacts
- `POST /api/v1/whatsapp/groups` - Create group
- `POST /api/v1/whatsapp/webhook` - Webhook endpoint

### Telegram Endpoints

- `POST /api/v1/telegram/send` - Send message
- `POST /api/v1/telegram/send-poll` - Create poll
- `POST /api/v1/telegram/broadcast` - Broadcast message
- `GET /api/v1/telegram/messages` - Get message history
- `POST /api/v1/telegram/webhook` - Webhook endpoint

### Bot Feature Endpoints

- `GET /api/v1/bot/features` - Get available features
- `POST /api/v1/bot/features/:feature/enable` - Enable feature
- `POST /api/v1/bot/features/:feature/disable` - Disable feature
- `GET /api/v1/bot/analytics` - Get bot analytics

### Game Endpoints

- `GET /api/v1/game/leaderboard` - Get game leaderboard
- `POST /api/v1/game/quiz/start` - Start quiz
- `POST /api/v1/game/quiz/answer` - Submit quiz answer
- `GET /api/v1/game/khodam/:name` - Check khodam

### Business Endpoints

- `GET /api/v1/business/products` - Get products
- `POST /api/v1/business/products` - Create product
- `GET /api/v1/business/orders` - Get orders
- `POST /api/v1/business/orders` - Create order

### Utility Endpoints

- `GET /api/v1/utils/weather/:city` - Get weather
- `GET /api/v1/utils/currency` - Convert currency
- `POST /api/v1/utils/translate` - Translate text
- `GET /api/v1/utils/qrcode` - Generate QR code

## Cara Penggunaan

### Auto-Reply (WhatsApp & Telegram)
```
User: "halo"
Bot: "Halo! Ada yang bisa saya bantu?"
```

### Interactive Menu (Telegram)
```
User: Klik tombol "Games"
Bot: Menampilkan menu games dengan inline keyboard
```

### Poll Creation (Telegram)
```
User: "Create poll: Apa makanan favoritmu? Options: Nasi Goreng, Mie Ayam, Sate"
Bot: Membuat poll dengan opsi yang diberikan
```

### Game - Cek Khodam
```
User: "cek khodam"
Bot: "KHODAM ANDA
Nama: Khodam Macan Putih
Kekuatan: Memberikan kekuatan dan keberanian"
```

### Weather
```
User: "cuaca jakarta"
Bot: "CUACA JAKARTA
Suhu: 28.5°C
Kelembapan: 75%
Kondisi: Berawan"
```

### Order Product
```
User: "katalog"
Bot: "KATALOG PRODUK
1. T-Shirt - Rp 150.000
2. Hoodie - Rp 250.000"

User: "pesan 1 2"
Bot: "PESANAN BERHASIL
Nomor: ORD-123456
Total: Rp 300.000"
```

## Pengembangan

### Menambah Fitur Baru

1. Buat service baru di `internal/services/`
2. Tambahkan handler di `internal/handlers/`
3. Daftarkan route di `main.go`
4. Update dokumentasi

### Testing
```bash
go test ./...
```

### Docker Support
```bash
docker build -t whatsapp-bot .
docker run -p 8080:8080 --env-file .env whatsapp-bot
```

## Monitoring

Aplikasi ini mendukung:
- Health check endpoint
- Metrics collection dengan Prometheus
- Error logging dengan Logrus
- Performance monitoring dengan Grafana
- Real-time dashboard untuk WhatsApp & Telegram

## Keamanan

- JWT authentication
- Rate limiting
- Input validation
- SQL injection protection
- XSS protection
- CORS configuration

## Lisensi

MIT License - lihat file [LICENSE](LICENSE) untuk detail

## Kontribusi

1. Fork repository
2. Buat branch fitur (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## Support

Untuk pertanyaan dan support:
- Email: hadsxdev@gmail.com
- Discord: [Join Server](https://discord.gg/imphnen)

## 🙏 Acknowledgments

- WhatsApp Business API
- Telegram Bot API
- Golang community
- All contributors

---

**Jika Anda menyukai project ini, jangan lupa untuk memberikan bintang!**