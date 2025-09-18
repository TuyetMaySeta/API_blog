Blog API - Hệ thống API Blog hiệu năng cao
Đây là một hệ thống API Blog được xây dựng với Golang, tích hợp PostgreSQL, Redis, và Elasticsearch để đảm bảo hiệu năng và khả năng mở rộng cao.

Cấu trúc project
blog-api/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── config/config.go         # Configuration
│   ├── database/               # Database connections
│   ├── models/post.go          # Data models
│   ├── handlers/post_handler.go # HTTP handlers
│   ├── services/               # Business logic
│   └── middleware/             # Middleware
├── migrations/                 # SQL migrations
├── docker-compose.yml          # Services definition
├── Dockerfile                  # API service build
└── README.md                   # This file

Cách chạy dự án
1. Clone repository
bashgit clone <repository-url>
cd blog-api
2. Khởi chạy với Docker Compose
bash# Khởi chạy tất cả services
docker-compose up -d

# Xem logs
docker-compose logs -f blog_api

# Dừng services
docker-compose down
3. Kiểm tra health
bashcurl http://localhost:8080/health

API Documentation

Base URL: http://localhost:8080/api/v1
1. Tạo bài viết mới
bashcurl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Bài viết đầu tiên",
    "content": "Nội dung bài viết về Golang và microservices",
    "tags": ["golang", "microservices", "api"]
  }'
2. Lấy chi tiết bài viết (với cache và related posts)
bashcurl http://localhost:8080/api/v1/posts/1
3. Cập nhật bài viết (với cache invalidation)
bashcurl -X PUT http://localhost:8080/api/v1/posts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Bài viết được cập nhật",
    "content": "Nội dung mới",
    "tags": ["golang", "redis", "elasticsearch"]
  }'
4. Tìm kiếm theo tag (sử dụng GIN index)
bashcurl "http://localhost:8080/api/v1/posts/search-by-tag?tag=golang"
5. Tìm kiếm full-text (sử dụng Elasticsearch)
bashcurl "http://localhost:8080/api/v1/posts/search?q=microservices"
6. Lấy danh sách bài viết (với pagination)
bashcurl "http://localhost:8080/api/v1/posts?limit=10&offset=0"
7. Xóa bài viết
bashcurl -X DELETE http://localhost:8080/api/v1/posts/1

Testing Commands

Tạo dữ liệu test
bash# Tạo vài bài viết để test
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Học Golang từ cơ bản",
    "content": "Golang là ngôn ngữ lập trình hiệu năng cao, được phát triển bởi Google",
    "tags": ["golang", "programming", "tutorial"]
  }'

curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Redis Caching Strategies",
    "content": "Cache-aside pattern là một trong những chiến lược caching phổ biến nhất",
    "tags": ["redis", "caching", "performance"]
  }'

curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Elasticsearch cho Full-text Search",
    "content": "Elasticsearch cung cấp khả năng tìm kiếm mạnh mẽ cho ứng dụng web",
    "tags": ["elasticsearch", "search", "indexing"]
  }'
Test caching (chạy 2 lần để thấy cache hit)
bash# Lần 1: Cache miss
time curl http://localhost:8080/api/v1/posts/1

# Lần 2: Cache hit (nhanh hơn)
time curl http://localhost:8080/api/v1/posts/1
Test search functionality
bash# Tìm kiếm theo tag
curl "http://localhost:8080/api/v1/posts/search-by-tag?tag=golang"

# Full-text search
curl "http://localhost:8080/api/v1/posts/search?q=performance"

 Production Deployment
 
1. Build cho production
bash# Build binary
go build -o blog-api ./cmd/server

# Hoặc build Docker image
docker build -t blog-api:latest .
2. Environment variables cho production
bashexport DB_HOST=prod-postgres-host
export DB_PASSWORD=secure-password
export REDIS_HOST=prod-redis-host
export ELASTICSEARCH_HOST=prod-es-host

 Monitoring & Troubleshooting
 
Kiểm tra PostgreSQL
bash# Kết nối vào PostgreSQL container
docker exec -it blog_postgres psql -U blog_user -d blog_db

# Kiểm tra GIN index
\d+ posts
SELECT * FROM pg_indexes WHERE tablename = 'posts';
Kiểm tra Redis
bash# Kết nối vào Redis
docker exec -it blog_redis redis-cli

# Xem các keys
KEYS post:*
TTL post:1
Kiểm tra Elasticsearch
bash# Xem index
curl http://localhost:9200/posts/_mapping

# Tìm kiếm trực tiếp
curl -X POST http://localhost:9200/posts/_search \
  -H "Content-Type: application/json" \
  -d '{"query": {"match": {"title": "golang"}}}'
Logs
bash# Xem logs của từng service
docker-compose logs postgres
docker-compose logs redis  
docker-compose logs elasticsearch
docker-compose logs blog_api
