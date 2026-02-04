#!/bin/bash

# Payment Terminal Management - Docker Helper Script
# Sử dụng: ./docker-helper.sh [command]

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

COMPOSE_FULL="docker-compose -f docker-compose.full.yml"
COMPOSE_APP="docker-compose"

# Functions
print_header() {
    echo -e "${BLUE}===================================================${NC}"
    echo -e "${BLUE}  Payment Terminal Management - Docker Helper${NC}"
    echo -e "${BLUE}===================================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Commands
start_full() {
    print_header
    print_info "Starting full stack (App + Kafka + AKHQ + PostgreSQL)..."
    $COMPOSE_FULL up -d --build
    print_success "Full stack started successfully!"
    echo ""
    print_info "Access points:"
    echo "  - App API: http://localhost:8089"
    echo "  - AKHQ: http://localhost:8080"
    echo "  - PostgreSQL: localhost:5432"
    echo "  - Kafka: localhost:9092"
}

start_app() {
    print_header
    print_info "Starting app only..."
    $COMPOSE_APP up -d --build
    print_success "App started successfully!"
    echo ""
    print_info "App API: http://localhost:8089"
}

stop_full() {
    print_info "Stopping full stack..."
    $COMPOSE_FULL down
    print_success "Full stack stopped!"
}

stop_app() {
    print_info "Stopping app..."
    $COMPOSE_APP down
    print_success "App stopped!"
}

clean_full() {
    print_info "Stopping and removing volumes (full stack)..."
    $COMPOSE_FULL down -v
    print_success "Full stack cleaned!"
}

clean_app() {
    print_info "Stopping app..."
    $COMPOSE_APP down
    print_success "App cleaned!"
}

logs_full() {
    print_info "Showing logs (full stack)..."
    $COMPOSE_FULL logs -f
}

logs_app() {
    print_info "Showing app logs..."
    $COMPOSE_FULL logs -f app
}

logs_kafka() {
    print_info "Showing Kafka logs..."
    $COMPOSE_FULL logs -f kafka
}

status() {
    print_header
    print_info "Container status:"
    docker ps --filter "name=payment-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
}

health() {
    print_header
    print_info "Health checks:"
    echo ""
    
    # App health
    if docker ps --filter "name=payment-terminal-app" --format "{{.Names}}" | grep -q payment-terminal-app; then
        APP_HEALTH=$(docker inspect payment-terminal-app --format='{{.State.Health.Status}}' 2>/dev/null || echo "no healthcheck")
        echo -e "App: ${GREEN}Running${NC} (Health: $APP_HEALTH)"
    else
        echo -e "App: ${RED}Not Running${NC}"
    fi
    
    # Kafka health
    if docker ps --filter "name=payment-kafka" --format "{{.Names}}" | grep -q payment-kafka; then
        KAFKA_HEALTH=$(docker inspect payment-kafka --format='{{.State.Health.Status}}' 2>/dev/null || echo "no healthcheck")
        echo -e "Kafka: ${GREEN}Running${NC} (Health: $KAFKA_HEALTH)"
    else
        echo -e "Kafka: ${RED}Not Running${NC}"
    fi
    
    # PostgreSQL health
    if docker ps --filter "name=payment-postgres" --format "{{.Names}}" | grep -q payment-postgres; then
        POSTGRES_HEALTH=$(docker inspect payment-postgres --format='{{.State.Health.Status}}' 2>/dev/null || echo "no healthcheck")
        echo -e "PostgreSQL: ${GREEN}Running${NC} (Health: $POSTGRES_HEALTH)"
    else
        echo -e "PostgreSQL: ${RED}Not Running${NC}"
    fi
}

rebuild_app() {
    print_info "Rebuilding app (no cache)..."
    $COMPOSE_FULL build --no-cache app
    $COMPOSE_FULL up -d app
    print_success "App rebuilt successfully!"
}

shell_app() {
    print_info "Accessing app container shell..."
    docker exec -it payment-terminal-app sh
}

shell_kafka() {
    print_info "Accessing Kafka container shell..."
    docker exec -it payment-kafka sh
}

shell_postgres() {
    print_info "Accessing PostgreSQL shell..."
    docker exec -it payment-postgres psql -U postgres payment_terminal
}

kafka_topics() {
    print_info "Listing Kafka topics..."
    docker exec payment-kafka kafka-topics --list --bootstrap-server localhost:9092
}

kafka_create_topics() {
    print_info "Creating default topics..."
    docker exec payment-kafka kafka-topics --create \
        --bootstrap-server localhost:9092 \
        --topic transaction_request \
        --partitions 3 \
        --replication-factor 1 \
        --if-not-exists
    
    docker exec payment-kafka kafka-topics --create \
        --bootstrap-server localhost:9092 \
        --topic transaction_response \
        --partitions 3 \
        --replication-factor 1 \
        --if-not-exists
    
    print_success "Topics created successfully!"
}

backup_postgres() {
    BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
    print_info "Backing up PostgreSQL to $BACKUP_FILE..."
    docker exec payment-postgres pg_dump -U postgres payment_terminal > "$BACKUP_FILE"
    print_success "Backup completed: $BACKUP_FILE"
}

show_usage() {
    print_header
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  start-full          Start full stack (app + kafka + akhq + postgres)"
    echo "  start-app           Start app only"
    echo "  stop-full           Stop full stack"
    echo "  stop-app            Stop app only"
    echo "  clean-full          Stop and remove volumes (full stack)"
    echo "  clean-app           Stop app"
    echo "  logs-full           Show logs (all services)"
    echo "  logs-app            Show app logs"
    echo "  logs-kafka          Show Kafka logs"
    echo "  status              Show container status"
    echo "  health              Show health status"
    echo "  rebuild-app         Rebuild app (no cache)"
    echo "  shell-app           Access app container shell"
    echo "  shell-kafka         Access Kafka container shell"
    echo "  shell-postgres      Access PostgreSQL shell"
    echo "  kafka-topics        List Kafka topics"
    echo "  kafka-create        Create default topics"
    echo "  backup-postgres     Backup PostgreSQL database"
    echo ""
    echo "Examples:"
    echo "  $0 start-full       # Start all services"
    echo "  $0 logs-app         # View app logs"
    echo "  $0 health           # Check service health"
}

# Main
case "${1:-}" in
    start-full)
        start_full
        ;;
    start-app)
        start_app
        ;;
    stop-full)
        stop_full
        ;;
    stop-app)
        stop_app
        ;;
    clean-full)
        clean_full
        ;;
    clean-app)
        clean_app
        ;;
    logs-full)
        logs_full
        ;;
    logs-app)
        logs_app
        ;;
    logs-kafka)
        logs_kafka
        ;;
    status)
        status
        ;;
    health)
        health
        ;;
    rebuild-app)
        rebuild_app
        ;;
    shell-app)
        shell_app
        ;;
    shell-kafka)
        shell_kafka
        ;;
    shell-postgres)
        shell_postgres
        ;;
    kafka-topics)
        kafka_topics
        ;;
    kafka-create)
        kafka_create_topics
        ;;
    backup-postgres)
        backup_postgres
        ;;
    *)
        show_usage
        exit 1
        ;;
esac
