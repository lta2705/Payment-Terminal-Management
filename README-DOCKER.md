# Docker Deployment Guide

## Yêu cầu
- Docker
- Docker Compose

## Tổng quan

Hệ thống bao gồm các service:
- **App**: Ứng dụng Go kết nối Kafka và PostgreSQL
- **Kafka**: Message broker
- **Zookeeper**: Quản lý Kafka cluster
- **AKHQ**: Web UI để monitor Kafka
- **PostgreSQL**: Database

## Các chế độ triển khai

Project hỗ trợ 2 chế độ:

### 1. **Chế độ Full Stack** (Recommended) - `docker-compose.full.yml`
Chạy tất cả services: App + Kafka + Zookeeper + AKHQ + PostgreSQL

### 2. **Chế độ App-only** - `docker-compose.yml`
Chỉ chạy App. Kết nối đến Kafka, AKHQ, PostgreSQL đã chạy sẵn.

---

## 🚀 Triển khai Full Stack (Recommended)

### Bước 1: Build và start tất cả services

```bash
docker-compose -f docker-compose.full.yml up -d --build
```

### Bước 2: Kiểm tra logs

```bash
# Xem logs của app
docker-compose -f docker-compose.full.yml logs -f app

# Xem logs của kafka
docker-compose -f docker-compose.full.yml logs -f kafka

# Xem logs của tất cả services
docker-compose -f docker-compose.full.yml logs -f
```

### Bước 3: Truy cập các services

- **App API**: http://localhost:8089
- **AKHQ (Kafka UI)**: http://localhost:8080
- **PostgreSQL**: localhost:5432
  - Database: `payment_terminal`
  - User: `postgres`
  - Password: `postgres`
- **Kafka**: 
  - External (from host): localhost:9092
  - Internal (docker network): kafka:29092

### Bước 4: Stop services

```bash
# Stop nhưng giữ data
docker-compose -f docker-compose.full.yml down

# Stop và xóa volumes (clean data)
docker-compose -f docker-compose.full.yml down -v
```

---

## Triển khai App-only (kết nối services có sẵn)

Nếu bạn đã có Kafka, AKHQ, và PostgreSQL chạy sẵn trong Docker network khác:

### Bước 1: Cấu hình môi trường

Copy file `.env.docker.example` thành `.env`:
```bash
cp .env.docker.example .env
```

Chỉnh sửa file `.env` để kết nối đến services:

**Ví dụ - Services trong cùng Docker network:**
```env
PORT=8089
KAFKA_BOOTSTRAP_SERVERS=kafka:29092
DB_HOST=postgres
DB_PORT=5432
DB_NAME=payment_terminal
DB_USER=postgres
DB_PASSWORD=postgres
```

**Ví dụ - Services trên host machine:**
```env
PORT=8089
KAFKA_BOOTSTRAP_SERVERS=host.docker.internal:9092
DB_HOST=host.docker.internal
```

**Ví dụ - Services trên server khác:**
```env
KAFKA_BOOTSTRAP_SERVERS=192.168.1.100:9092
DB_HOST=192.168.1.100
```

### Bước 2: Cấu hình Docker network

Nếu services chạy ở network khác, cập nhật `docker-compose.yml`:
```yaml
networks:
  payment-network:
    external: true
    name: <tên_network_của_bạn>  # ví dụ: kafka-network
```

### Bước 3: Build và chạy
```bash
docker-compose up -d --build
```

### Bước 4: Xem logs
```bash
docker-compose logs -f app
```

### Bước 5: Dừng app
```bash
docker-compose down
```

---

## 🔧 Network Configuration

### Internal Docker Network (Trong docker network)
Các services giao tiếp qua hostname:
- `kafka:29092` - Kafka broker
- `postgres:5432` - PostgreSQL
- `akhq:8080` - AKHQ UI
- `zookeeper:2181` - Zookeeper

### External Access (Từ host machine)
- `localhost:9092` - Kafka broker
- `localhost:5432` - PostgreSQL
- `localhost:8080` - AKHQ UI
- `localhost:8089` - App API
- `localhost:2181` - Zookeeper

---

## 📊 Useful Commands

### Check service health
```bash
# Xem tất cả containers
docker ps

# Xem health status của app
docker inspect payment-terminal-app | grep -A 5 Health

# Check Kafka connectivity
docker exec payment-kafka kafka-broker-api-versions --bootstrap-server localhost:9092

# Check PostgreSQL
docker exec payment-postgres pg_isready -U postgres
```

### Debug

```bash
# Vào app container
docker exec -it payment-terminal-app sh

# Vào kafka container
docker exec -it payment-kafka sh

# Vào postgres container
docker exec -it payment-postgres psql -U postgres payment_terminal

# Xem Kafka topics
docker exec payment-kafka kafka-topics --list --bootstrap-server localhost:9092

# Tạo topic manually
docker exec payment-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic transaction_request \
  --partitions 3 \
  --replication-factor 1

# Consume messages từ topic
docker exec payment-kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic transaction_request \
  --from-beginning

# Produce message vào topic
docker exec -it payment-kafka kafka-console-producer \
  --bootstrap-server localhost:9092 \
  --topic transaction_request
```

### Rebuild app only

```bash
# Full rebuild
docker-compose -f docker-compose.full.yml up -d --build app

# Rebuild without cache
docker-compose -f docker-compose.full.yml build --no-cache app
docker-compose -f docker-compose.full.yml up -d app
```

### Monitor resource usage

```bash
docker stats payment-terminal-app payment-kafka payment-postgres payment-akhq
```

---

## ⚙️ Environment Variables

### App Configuration
| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8089 | App listening port |
| KAFKA_BOOTSTRAP_SERVERS | kafka:29092 | Kafka broker addresses |
| DB_HOST | postgres | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_NAME | payment_terminal | Database name |
| DB_USER | postgres | Database user |
| DB_PASSWORD | postgres | Database password |

### Kafka Producer Config
| Variable | Default | Description |
|----------|---------|-------------|
| KAFKA_PRODUCER_ACKS | all | Acknowledgment mode |
| KAFKA_PRODUCER_RETRIES | 10 | Number of retries |
| KAFKA_PRODUCER_ENABLE_IDEMPOTENCE | true | Enable idempotence |
| KAFKA_PRODUCER_TOPIC | transaction_response | Producer topic |

### Kafka Consumer Config
| Variable | Default | Description |
|----------|---------|-------------|
| KAFKA_CONSUMER_GROUP_ID | transaction_forwarder_group | Consumer group ID |
| KAFKA_CONSUMER_TOPIC | transaction_request | Consumer topic |
| KAFKA_CONSUMER_AUTO_OFFSET_RESET | latest | Offset reset strategy |
| KAFKA_CONSUMER_ENABLE_AUTO_COMMIT | false | Auto commit mode |

---

## 🔍 Troubleshooting

### App không kết nối được Kafka
1. Kiểm tra Kafka đã healthy chưa:
   ```bash
   docker-compose -f docker-compose.full.yml ps
   ```

2. Kiểm tra network:
   ```bash
   docker network inspect payment-network
   ```

3. Kiểm tra logs:
   ```bash
   docker-compose -f docker-compose.full.yml logs kafka
   ```

4. Test connectivity từ app container:
   ```bash
   docker exec payment-terminal-app ping kafka
   docker exec payment-terminal-app nc -zv kafka 29092
   ```

### App không kết nối được PostgreSQL
1. Kiểm tra PostgreSQL đã ready:
   ```bash
   docker exec payment-postgres pg_isready -U postgres
   ```

2. Test connection từ app container:
   ```bash
   docker exec payment-terminal-app ping postgres
   ```

3. Test login:
   ```bash
   docker exec payment-postgres psql -U postgres -c "SELECT 1"
   ```

### Port conflicts
Nếu port đã được sử dụng, cập nhật trong `docker-compose.full.yml`:
```yaml
ports:
  - "8090:8089"  # Map to different host port
  - "9093:9092"  # Kafka
  - "5433:5432"  # PostgreSQL
  - "8081:8080"  # AKHQ
```

### Kafka topics không tự tạo
```bash
# Tạo topics manually
docker exec payment-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic transaction_request \
  --partitions 3 \
  --replication-factor 1

docker exec payment-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic transaction_response \
  --partitions 3 \
  --replication-factor 1
```

---

## 🔐 Production Considerations

### Security
- ✅ App chạy với non-root user (đã config trong Dockerfile)
- ⚠️ Thay đổi default passwords trong production
- ⚠️ Sử dụng Docker secrets cho sensitive data
- ⚠️ Enable SSL/TLS cho Kafka và PostgreSQL
- ⚠️ Sử dụng environment-specific configs

### Performance
- Tăng Kafka partitions cho high throughput
- Configure PostgreSQL connection pooling
- Monitor với AKHQ và PostgreSQL monitoring tools
- Scale app horizontally: `docker-compose -f docker-compose.full.yml up --scale app=3`

### Backup & Recovery

**Backup PostgreSQL:**
```bash
docker exec payment-postgres pg_dump -U postgres payment_terminal > backup_$(date +%Y%m%d).sql
```

**Restore PostgreSQL:**
```bash
docker exec -i payment-postgres psql -U postgres payment_terminal < backup_20260114.sql
```

**Backup Kafka data (volumes):**
```bash
docker run --rm -v payment-kafka-data:/data -v $(pwd):/backup alpine tar czf /backup/kafka-backup.tar.gz /data
```

---

## 🎯 Quick Reference

### Start full stack
```bash
docker-compose -f docker-compose.full.yml up -d --build
```

### Start app only
```bash
docker-compose up -d --build
```

### View logs
```bash
docker-compose -f docker-compose.full.yml logs -f app
```

### Stop and clean
```bash
docker-compose -f docker-compose.full.yml down -v
```

### Rebuild app
```bash
docker-compose -f docker-compose.full.yml up -d --build app
```

### Access services
- App: http://localhost:8089
- AKHQ: http://localhost:8080
- PostgreSQL: localhost:5432
```bash
# Tạo topics manually
docker exec payment-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic transaction_request \
  --partitions 3 \
  --replication-factor 1

docker exec payment-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic transaction_response \
  --partitions 3 \
  --replication-factor 1
```

---

## 🔐 Production Considerations

### Security
- ✅ App chạy với non-root user (đã config trong Dockerfile)
- ⚠️ Thay đổi default passwords trong production
- ⚠️ Sử dụng Docker secrets cho sensitive data
- ⚠️ Enable SSL/TLS cho Kafka và PostgreSQL
- ⚠️ Sử dụng environment-specific configs

### Performance
- Tăng Kafka partitions cho high throughput
- Configure PostgreSQL connection pooling
- Monitor với AKHQ và PostgreSQL monitoring tools
- Scale app horizontally: `docker-compose -f docker-compose.full.yml up --scale app=3`

### Backup & Recovery

**Backup PostgreSQL:**
```bash
docker exec payment-postgres pg_dump -U postgres payment_terminal > backup_$(date +%Y%m%d).sql
```

**Restore PostgreSQL:**
```bash
docker exec -i payment-postgres psql -U postgres payment_terminal < backup_20260114.sql
```

**Backup Kafka data (volumes):**
```bash
docker run --rm -v payment-kafka-data:/data -v $(pwd):/backup alpine tar czf /backup/kafka-backup.tar.gz /data
```

---

## 🎯 Quick Reference

### Start full stack
```bash
docker-compose -f docker-compose.full.yml up -d --build
```

### Start app only
```bash
docker-compose up -d --build
```

### View logs
```bash
docker-compose -f docker-compose.full.yml logs -f app
```

### Stop and clean
```bash
docker-compose -f docker-compose.full.yml down -v
```bash
docker-compose -f docker-compose.full.yml down -v
```

---

## Truy cập Services

### Chế độ App-only:
- **Payment Terminal Server**: `0.0.0.0:8089` (TCP)

### Chế độ Full Stack:
- **Payment Terminal Server**: `0.0.0.0:8089` (TCP)
- **Kafka**: `localhost:9092` (từ host machine)
- **Kafka UI**: `http://localhost:8080` (Web interface để quản lý Kafka)
- **Zookeeper**: `localhost:2181`

---

## Kết nối từ máy khác

Server lắng nghe trên port `8089` và có thể nhận kết nối từ mọi nơi:

```bash
# Từ máy khác trong cùng network
nc <IP_CỦA_SERVER> 8089

# hoặc
telnet <IP_CỦA_SERVER> 8089
```

## Kafka Topics

Các topics được tự động tạo khi có message đầu tiên:
- `transaction_request` - Consumer topic
- `transaction_response` - Producer topic

Bạn có thể xem và quản lý topics qua Kafka UI tại `http://localhost:8080`

## Rebuild Application

### Chế độ App-only:
```bash
docker-compose up -d --build
```

### Chế độ Full Stack:
```bash
# Rebuild chỉ app
docker-compose -f docker-compose.full.yml up -d --build app

# Rebuild tất cả
docker-compose -f docker-compose.full.yml up -d --build
```

# App-only
docker-compose ps

# Full stack
docker-compose -f docker-compose.full.yml ps
```

### Kiểm tra kết nối Kafka từ app container
```bash
docker exec -it payment-terminal-app sh
# Trong container, kiểm tra biến môi trường
env | grep KAFKA
```

### Kiểm tra health của Kafka (Full Stack)
```bash
docker exec -it payment-kafka kafka-topics --list --bootstrap-server localhost:9092
```

### Truy cập vào container
```bash
# Vào container app
docker exec -it payment-terminal-app sh

# Vào container kafka (Full Stack)
docker exec -it payment-kafka bash
```

### Reset toàn bộ
### App-only mode:
Tạo file `.env` từ template:
```bash
cp .env.docker.example .env
```

Chỉnh sửa các giá trị trong file `.env` theo nhu cầu.

### Full Stack mode:
Các biến môi trường được cấu hình trực tiếp trong `docker-compose.full.ymlcker-compose -f docker-compose.full.yml down -v
docker-compose -f docker-compose.full.yml
### Reset toàn bộ
```bash
docker-compose down -v
docker-compose up -d --build
```

## Cấu hình môi trường

Các biến môi trường có thể được thay đổi trong `docker-compose.yml` hoặc tạo file `.env` trong thư mục gốc:

```env
PORT=8089
KAFKA_PRODUCER_TOPIC=transaction_response
KAFKA_CONSUMER_TOPIC=transaction_request
```
### Sử dụng Kafka có sẵn (khuyến nghị):

1. Copy project lên server
2. Tạo file `.env` và cấu hình Kafka connection:
   ```bash
   cp .env.docker.example .env
   nano .env  # Chỉnh sửa KAFKA_BOOTSTRAP_SERVERS
   ```
3. Đảm bảo firewall mở port 8089
4. Deploy:
   ```bash
   docker-compose up -d --build
   ```
5. Monitoring:
   ```bash
   docker-compose logs -f
   ```

### Deploy Full Stack (dev/test):

1. Copy project lên server
2. Đảm bảo firewall mở các ports: 8089, 9092, 8080, 2181
3. Deploy:
   ```bash
   docker-compose -f docker-compose.full.yml up -d
   ```
4. Monitoring:
   ```bash
   docker-compose -f docker-compose.full.yml logs -f
   ``
2. Đảm bảo firewall mở port 8089
3. Chạy: `docker-compose up -d`
4. Monitoring: `docker-compose logs -f`

### Mở port cho external connections (nếu cần)

```bash
# Ubuntu/Debian
sudo ufw allow 8089/tcp

# CentOS/RHEL
sudo firewall-cmd --permanent --add-port=8089/tcp
sudo firewall-cmd --reload
```
