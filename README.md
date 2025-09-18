# 📝 Blog API - Hệ thống API Blog hiệu năng cao

[![Go](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)](https://postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7-red.svg)](https://redis.io)
[![Elasticsearch](https://img.shields.io/badge/Elasticsearch-8.11-orange.svg)](https://elastic.co)
[![Docker](https://img.shields.io/badge/Docker-Compose-blue.svg)](https://docker.com)

Đây là một hệ thống API Blog được xây dựng với **Golang**, tích hợp **PostgreSQL**, **Redis**, và **Elasticsearch** để đảm bảo hiệu năng và khả năng mở rộng cao.

## 🚀 Tính năng chính

- ✅ **PostgreSQL** với GIN index cho tìm kiếm tag nhanh chóng
- ✅ **Redis Cache-Aside Pattern** với TTL và invalidation tự động  
- ✅ **Elasticsearch** full-text search với match query
- ✅ **Transaction** đảm bảo toàn vẹn dữ liệu
- ✅ **Related Posts** thông minh dựa trên tags
- ✅ **Docker Compose** setup đầy đủ
- ✅ **Clean Architecture** với middleware và error handling

## 🛠️ Tech Stack

| Công nghệ | Version | Mục đích |
|-----------|---------|----------|
| **Golang** | 1.21 | Backend API với Gin framework |
| **PostgreSQL** | 15 | Primary database với GIN index |
| **Redis** | 7 | Caching layer cho performance |
| **Elasticsearch** | 8.11 | Full-text search engine |
| **Docker** | Latest | Containerization |
| **GORM** | Latest | ORM cho Golang |

## 📁 Cấu trúc Project

```
blog-api/
├── 📁 cmd/
│   └── server/main.go              # 🚀 Entry point
├── 📁 internal/
│   ├── config/config.go            # ⚙️ Configuration
│   ├── 📁 database/                # 🗄️ Database connections
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   └── elasticsearch.go
│   ├── models/post.go              # 📊 Data models
│   ├── handlers/post_handler.go    # 🌐 HTTP handlers
│   ├── 📁 services/                # 💼 Business logic
│   │   ├── post_service.go
│   │   ├── cache_service.go
│   │   └── search_service.go
│   └── 📁 middleware/              # 🔧 Middleware
├── 📁 migrations/                  # 📝 SQL migrations
├── docker-compose.yml              # 🐳 Services definition
├── Dockerfile                      # 📦 API service build
├── Makefile                        # 🛠️ Development commands
└── README.md                       # 📖 Documentation
```

## 🚦 Quick Start

### Prerequisites
- Docker & Docker Compose
- Git

### 1️⃣ Clone Repository
```bash
git clone https://github.com/your-username/blog-api.git
cd blog-api
```

### 2️⃣ Khởi chạy Services
```bash
# Khởi chạy tất cả services
docker-compose up -d

# Hoặc sử dụng Makefile
make run
```

### 3️⃣ Kiểm tra Health
```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "message": "Blog API is running"
}
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/posts` | Tạo bài viết mới |
| `GET` | `/posts/:id` | Lấy chi tiết bài viết |
| `PUT` | `/posts/:id` | Cập nhật bài viết |
| `DELETE` | `/posts/:id` | Xóa bài viết |
| `GET` | `/posts` | Danh sách bài viết (pagination) |
| `GET` | `/posts/search-by-tag?tag=<name>` | Tìm kiếm theo tag |
| `GET` | `/posts/search?q=<query>` | Full-text search |

### 📝 Examples

#### Tạo bài viết mới
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Học Golang từ cơ bản",
    "content": "Golang là ngôn ngữ lập trình hiệu năng cao, được phát triển bởi Google",
    "tags": ["golang", "programming", "tutorial"]
  }'
```

#### Lấy chi tiết bài viết (với related posts)
```bash
curl http://localhost:8080/api/v1/posts/1
```

#### Tìm kiếm theo tag
```bash
curl "http://localhost:8080/api/v1/posts/search-by-tag?tag=golang"
```

#### Full-text search
```bash
curl "http://localhost:8080/api/v1/posts/search?q=performance"
```

## 🧪 Testing

### Seed dữ liệu test
```bash
make seed
```

### Test caching performance
```bash
# Lần 1: Cache miss
time curl http://localhost:8080/api/v1/posts/1

# Lần 2: Cache hit (nhanh hơn)
time curl http://localhost:8080/api/v1/posts/1
```

### Health check tất cả services
```bash
make health
```

## 🏗️ Architecture

### Database Schema
```sql
-- Posts với GIN index cho tags
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_tags ON posts USING GIN (tags);

-- Activity logs cho data integrity
CREATE TABLE activity_logs (
    id SERIAL PRIMARY KEY,
    action VARCHAR(100) NOT NULL,
    post_id INTEGER NOT NULL,
    logged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Cache Strategy
- 🔄 **Cache-Aside Pattern**
- ⏱️ **TTL**: 5 phút
- 🔑 **Key Pattern**: `post:<id>`
- 🗑️ **Auto Invalidation** khi update/delete

### Search Architecture
- 📊 **Elasticsearch Index**: `posts`
- 🔍 **Multi-match Query**: title + content
- 🔗 **Related Posts**: Bool query với tags matching

## 🔧 Development

### Available Make Commands
```bash
make help          # Show available commands
make run           # Start all services
make stop          # Stop services  
make clean         # Clean up containers
make logs          # Show API logs
make db            # Connect to PostgreSQL
make redis         # Connect to Redis CLI
make es            # Show Elasticsearch info
make seed          # Create test data
make health        # Check services health
```

### Local Development
```bash
# Start dependencies only
docker-compose up -d postgres redis elasticsearch

# Run API locally
make dev
```

## 📊 Performance Features

### ⚡ Database Optimizations
- **GIN Index** cho array tags search
- **Connection Pooling** với GORM
- **Transaction Rollback** cho data integrity

### 🚀 Caching Optimizations  
- **Cache-Aside Pattern** giảm DB load
- **Async Cache Operations** không block request
- **Smart Invalidation** đảm bảo data consistency

### 🔍 Search Optimizations
- **Elasticsearch Scoring** cho relevance
- **Multi-field Search** title + content
- **Related Posts Algorithm** cho engagement

## 🛠️ Monitoring & Troubleshooting

### PostgreSQL
```bash
# Connect to database
make db

# Check GIN index
\d+ posts
SELECT * FROM pg_indexes WHERE tablename = 'posts';
```

### Redis
```bash
# Connect to Redis
make redis

# Check cache keys
KEYS post:*
TTL post:1
```

### Elasticsearch
```bash
# Check index mapping
curl http://localhost:9200/posts/_mapping

# Direct search
curl -X POST http://localhost:9200/posts/_search \
  -H "Content-Type: application/json" \
  -d '{"query": {"match": {"title": "golang"}}}'
```

### Logs
```bash
# API logs
make logs

# All services
docker-compose logs -f
```


## 🤝 Contributing

1. Fork repository
2. Tạo feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push branch: `git push origin feature/amazing-feature`
5. Tạo Pull Request



</div>
