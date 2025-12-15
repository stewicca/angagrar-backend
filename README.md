# Angagrar Backend - AI-Powered Personal Budget Assistant

Backend API untuk Angagrar, aplikasi budget management dengan AI conversation-based budget generation menggunakan OpenAI.

## 🌟 Features

- **AI-Powered Conversation**: Chat natural dengan Aira (AI assistant) untuk budget planning
- **Personalized Budget**: LLM analyze conversation context untuk generate truly personal budget
- **No Rigid Forms**: Free-form conversation, tidak kaku
- **Smart Analysis**: Consider salary, location, lifestyle, habits, dan goals simultaneously
- **Budget Tracking**: Track dan adjust budget yang sudah di-generate
- **Transaction Management**: Record actual spending

## 🏗️ Architecture

**LLM-First Approach**:
- Conversation natural dengan OpenAI GPT-4o-mini
- Context-aware budget generation
- No hardcoded formulas, fully AI-driven
- Conversation history untuk continuity & audit trail

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL
- OpenAI API Key

### Installation

1. **Clone repository**
```bash
git clone <repo-url>
cd angagrar-backend
```

2. **Setup environment**
```bash
cp .env.example .env
# Edit .env dengan credentials Anda (terutama OPENAI_API_KEY)
```

3. **Install dependencies**
```bash
go mod download
```

4. **Setup database**
```bash
# Create database
createdb angagrar_db

# Migrations run automatically on startup
```

5. **Run server**
```bash
go run cmd/server/main.go
```

Server berjalan di `http://localhost:8080`

## 📚 API Documentation

Lihat [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) untuk detail lengkap semua endpoints.

### Quick Example

```bash
# 1. Create guest user
curl -X POST http://localhost:8080/api/v1/auth/guest

# 2. Start conversation (gunakan token dari step 1)
curl -X POST http://localhost:8080/api/v1/conversations/start \
  -H "Authorization: Bearer <token>"

# 3. Chat dengan Aira
curl -X POST http://localhost:8080/api/v1/conversations/<session_id>/messages \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"message": "hai! gaji gue 8 juta, tinggal di bandung, lifestyle moderate. pengeluaran rutin gue sewa 2jt, makan sekitar 1.5jt"}'

# 4. Request budget generation
curl -X POST http://localhost:8080/api/v1/conversations/<session_id>/messages \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"message": "oke buatin budget gue dong!"}'

# 5. Get generated budgets
curl -X GET http://localhost:8080/api/v1/budgets \
  -H "Authorization: Bearer <token>"
```

## 🎯 How It Works

### Conversation Flow:
1. User start conversation → Aira greets
2. User cerita tentang keuangan mereka (free-form, natural)
3. Aira bertanya follow-up untuk gather lebih banyak context
4. User request "buatin budget"
5. LLM analyze seluruh conversation
6. Generate personalized budget (6 categories)
7. Save ke database

### Budget Categories:
- 💸 **Kewajiban**: Sewa, utilities, cicilan
- 🍜 **Makan**: Makanan sehari-hari
- 🚗 **Transport**: Transportasi
- 🎮 **Healing**: Hiburan, hobi, self-care
- 💰 **Tabungan**: Savings & investasi
- 📦 **Lain-lain**: Misc expenses

## 🛠️ Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL + GORM
- **AI/LLM**: OpenAI GPT-4o-mini
- **Auth**: JWT (guest-based)
- **Deployment**: Docker + Docker Compose

## 📁 Project Structure

```
angagrar-backend/
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── internal/
│   ├── database/        # DB connection & migrations
│   ├── handlers/        # HTTP handlers (controllers)
│   ├── middleware/      # Auth, logging, error handling
│   ├── models/          # Domain models (GORM)
│   ├── repositories/    # Data access layer
│   └── services/        # Business logic (incl. AI)
├── pkg/utils/           # Utility helpers
├── .env.example         # Environment template
└── docker-compose.yml   # Docker setup
```

## 🔧 Configuration

Edit `.env`:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=angagrar_db

# App
APP_PORT=8080

# JWT Secret
JWT_SECRET=your-secret-key

# OpenAI (REQUIRED)
OPENAI_API_KEY=sk-your-api-key-here
OPENAI_MODEL=gpt-4o-mini
OPENAI_MAX_TOKENS=1000
OPENAI_TEMPERATURE=0.7
```

## 🐳 Docker

```bash
# Run dengan Docker Compose
docker-compose up -d

# Stop
docker-compose down
```

## 💰 Cost Estimation

**OpenAI API (GPT-4o-mini)**:
- Input: $0.60 per 1M tokens
- Output: $2.40 per 1M tokens
- Average conversation: ~3000 tokens
- **Cost per budget generation**: ~$0.002 (2 cents)
- **10,000 users/month**: ~$20/month

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run specific package
go test ./internal/services/...

# With coverage
go test -cover ./...
```

## 📝 Development

### Adding New Features

1. **Model**: Define di `internal/models/`
2. **Repository**: CRUD operations di `internal/repositories/`
3. **Service**: Business logic di `internal/services/`
4. **Handler**: HTTP endpoints di `internal/handlers/`
5. **Routes**: Wire di `cmd/server/main.go`

### Code Style
- Follow Go conventions
- Use `gofmt` untuk formatting
- Comment exported functions
- Keep handlers thin, logic in services

## 🚧 Roadmap (Future Enhancements)

- [ ] Multiple budget plans per user
- [ ] Budget vs Actual tracking & alerts
- [ ] Voice input untuk conversation
- [ ] Export budget to PDF/Excel
- [ ] Budget sharing dengan family/friends
- [ ] Admin panel untuk analytics
- [ ] Multi-language support (English)
- [ ] Self-hosted LLM option (Llama)

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👥 Authors

- **Wicca** - Initial work

## 🙏 Acknowledgments

- OpenAI for GPT-4o-mini API
- Gin framework
- GORM
- All open source contributors

---

**Need help?** Open an issue atau hubungi tim development.
