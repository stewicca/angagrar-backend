# Angagrar Backend API Documentation

## Architecture Overview

**LLM-First Approach**: Menggunakan OpenAI untuk personalized budget generation melalui free-form conversation dengan Aira (AI assistant).

### Tech Stack
- **Framework**: Go + Gin
- **Database**: PostgreSQL + GORM
- **AI**: OpenAI GPT-4o-mini
- **Auth**: JWT (Guest users)

---

## Authentication

### Create Guest User
```http
POST /api/v1/auth/guest
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "guest_id": "uuid-here"
    },
    "token": "jwt-token-here"
  }
}
```

**Usage:**
- Save token untuk requests berikutnya
- Header: `Authorization: Bearer <token>`

---

## AI Conversation Flow

### 1. Start Conversation
```http
POST /api/v1/conversations/start
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "session_id": "uuid-session-id",
    "message": "hai! 👋 gue aira, siap bantu kamu atur budget..."
  }
}
```

### 2. Send Message to Aira
```http
POST /api/v1/conversations/:sessionId/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "message": "hai! gaji gue 8 juta, tinggal di bandung, suka healing"
}
```

**Response (Normal Chat):**
```json
{
  "success": true,
  "data": {
    "assistant_message": "nice! 8 juta di bandung oke tuh. kira-kira pengeluaran rutin kamu apa aja?",
    "completed": false
  }
}
```

**Response (Budget Generated):**
```json
{
  "success": true,
  "data": {
    "assistant_message": "done! ✨ ini budget recommendation yang gue bikinin...",
    "completed": true,
    "budget_generated": true,
    "budgets": [
      {
        "id": 1,
        "category": "Kewajiban",
        "amount": 3500000,
        "period": "monthly",
        "description": "sewa, utilities, cicilan"
      },
      {
        "id": 2,
        "category": "Makan",
        "amount": 1500000,
        "period": "monthly",
        "description": "makanan sehari-hari"
      }
      // ... 4 more categories
    ]
  }
}
```

### 3. Get Conversation History
```http
GET /api/v1/conversations/:sessionId/history
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "messages": [
      {
        "id": 1,
        "role": "assistant",
        "content": "hai! 👋 gue aira...",
        "created_at": "2025-01-01T10:00:00Z"
      },
      {
        "id": 2,
        "role": "user",
        "content": "hai! gaji gue 8 juta...",
        "created_at": "2025-01-01T10:01:00Z"
      }
    ]
  }
}
```

### 4. Reset Conversation
```http
POST /api/v1/conversations/:sessionId/reset
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Conversation reset. Starting new interview.",
    "new_session_id": "new-uuid-session-id",
    "greeting": "hai! 👋 gue aira..."
  }
}
```

---

## Budget Management

### Get User Budgets
```http
GET /api/v1/budgets
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "budgets": [
      {
        "id": 1,
        "user_id": 1,
        "category": "Kewajiban",
        "amount": 3500000,
        "period": "monthly",
        "start_date": "2025-01-01T00:00:00Z",
        "end_date": "2025-01-31T23:59:59Z",
        "description": "sewa, utilities, cicilan",
        "created_at": "2025-01-01T10:00:00Z"
      }
      // ... more budgets
    ]
  }
}
```

### Update Budget
```http
PATCH /api/v1/budgets/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "amount": 4000000
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "budget": {
      "id": 1,
      "amount": 4000000,
      "category": "Kewajiban"
      // ... other fields
    }
  }
}
```

---

## User Profile

### Get User Profile
```http
GET /api/v1/users/profile
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "guest_id": "uuid-here",
      "created_at": "2025-01-01T10:00:00Z"
    }
  }
}
```

---

## Transactions

### Create Transaction
```http
POST /api/v1/transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "type": "expense",
  "category": "Makan",
  "amount": 50000,
  "description": "lunch",
  "date": "2025-01-01T12:00:00Z"
}
```

### Get Transactions
```http
GET /api/v1/transactions
Authorization: Bearer <token>
```

---

## Example Flow

### Complete User Journey:

1. **Create guest user**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/guest
   ```

2. **Start conversation** (save session_id)
   ```bash
   curl -X POST http://localhost:8080/api/v1/conversations/start \
     -H "Authorization: Bearer <token>"
   ```

3. **Chat with Aira** (natural conversation)
   ```bash
   curl -X POST http://localhost:8080/api/v1/conversations/<session_id>/messages \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"message": "hai! gaji gue 8 juta, tinggal di bandung, lifestyle moderate"}'
   ```

4. **Continue chatting** until Aira has enough info

5. **Request budget generation**
   ```bash
   curl -X POST http://localhost:8080/api/v1/conversations/<session_id>/messages \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"message": "oke buatin budget gue dong!"}'
   ```

6. **Check generated budgets**
   ```bash
   curl -X GET http://localhost:8080/api/v1/budgets \
     -H "Authorization: Bearer <token>"
   ```

---

## How It Works

### LLM-Powered Budget Generation:

1. **Free-form Conversation**: User chats naturally dengan Aira tentang keuangan mereka
2. **Context Collection**: Aira (OpenAI) mengumpulkan info: salary, location, lifestyle, habits, goals
3. **Smart Analysis**: Ketika user siap, LLM analyze seluruh conversation context
4. **Personalized Budget**: LLM generate budget allocation yang truly personal, bukan hardcoded formula
5. **Database Storage**: Budget results disimpan untuk tracking & adjustment

### Why LLM Approach?

- **Personal**: Budget disesuaikan dengan context unik user
- **Flexible**: No rigid state machine, conversation natural
- **Smart**: LLM consider cost of living, lifestyle, goals simultaneously  
- **Scalable**: Easy to add more context/features tanpa rewrite logic

---

## Environment Variables

Required `.env`:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=angagrar_db

# Application
APP_PORT=8080

# JWT
JWT_SECRET=your-secret-key-here

# OpenAI
OPENAI_API_KEY=sk-your-api-key-here
OPENAI_MODEL=gpt-4o-mini
OPENAI_MAX_TOKENS=1000
OPENAI_TEMPERATURE=0.7
```

---

## Running the Server

```bash
# Install dependencies
go mod download

# Run migrations (auto on startup)
go run cmd/server/main.go

# Server runs on http://localhost:8080
```

---

## Database Models

### User
- Guest-based authentication
- One-to-many: Conversations, Budgets, Transactions

### Conversation
- Tracks AI chat sessions
- Stores conversation completion status
- Links to Messages

### Message
- Chat history (user & assistant)
- Provides context for LLM

### Budget
- 6 categories: Kewajiban, Makan, Transport, Healing, Tabungan, Lain-lain
- User can manually adjust amounts
- Monthly period

### Transaction
- Track actual spending
- Optional link to Budget

---

## API Response Format

### Success Response:
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response:
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error"
}
```

---

## Notes

- **MVP**: 1 user = 1 active conversation (must complete or reset before starting new)
- **OpenAI Costs**: ~$0.002 per conversation (GPT-4o-mini)
- **Budget Categories**: Fixed 6 categories untuk MVP
- **Conversation**: Disimpan untuk history & audit trail
- **Manual Adjustment**: User bisa edit budget amounts setelah generated
