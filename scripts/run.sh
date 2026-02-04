#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored message
print_message() {
    echo -e "${2}${1}${NC}"
}

# Print header
print_header() {
    echo ""
    print_message "========================================" "$BLUE"
    print_message "$1" "$BLUE"
    print_message "========================================" "$BLUE"
    echo ""
}

# Check if docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_message "Error: Docker is not running. Please start Docker first." "$RED"
        exit 1
    fi
}

# Development commands
dev_up() {
    print_header "Starting Development Environment"
    check_docker
    docker compose -f docker-compose.dev.yml up --build -d
    print_message "Development environment started!" "$GREEN"
    print_message "App: http://localhost:3000" "$YELLOW"
    print_message "pgAdmin: http://localhost:5050" "$YELLOW"
}

dev_down() {
    print_header "Stopping Development Environment"
    docker compose -f docker-compose.dev.yml down
    print_message "Development environment stopped!" "$GREEN"
}

dev_logs() {
    print_header "Development Logs"
    docker compose -f docker-compose.dev.yml logs -f
}

dev_logs_app() {
    print_header "App Logs (Development)"
    docker compose -f docker-compose.dev.yml logs -f app
}

dev_restart() {
    print_header "Restarting Development Environment"
    docker compose -f docker-compose.dev.yml restart
    print_message "Development environment restarted!" "$GREEN"
}

dev_rebuild() {
    print_header "Rebuilding Development Environment"
    docker compose -f docker-compose.dev.yml down
    docker compose -f docker-compose.dev.yml up --build -d
    print_message "Development environment rebuilt!" "$GREEN"
}

# Production commands
prod_up() {
    print_header "Starting Production Environment"
    check_docker
    docker compose -f docker-compose.prod.yml up --build -d
    print_message "Production environment started!" "$GREEN"
    print_message "App: http://localhost:3000" "$YELLOW"
    print_message "pgAdmin: http://localhost:5050" "$YELLOW"
}

prod_down() {
    print_header "Stopping Production Environment"
    docker compose -f docker-compose.prod.yml down
    print_message "Production environment stopped!" "$GREEN"
}

prod_logs() {
    print_header "Production Logs"
    docker compose -f docker-compose.prod.yml logs -f
}

prod_logs_app() {
    print_header "App Logs (Production)"
    docker compose -f docker-compose.prod.yml logs -f app
}

prod_restart() {
    print_header "Restarting Production Environment"
    docker compose -f docker-compose.prod.yml restart
    print_message "Production environment restarted!" "$GREEN"
}

prod_rebuild() {
    print_header "Rebuilding Production Environment"
    docker compose -f docker-compose.prod.yml down
    docker compose -f docker-compose.prod.yml up --build -d
    print_message "Production environment rebuilt!" "$GREEN"
}

# Database commands
db_shell() {
    print_header "Connecting to PostgreSQL"
    if [ "$1" == "prod" ]; then
        docker exec -it fiber-postgres-prod psql -U postgres -d fiber_gorm_db
    else
        docker exec -it fiber-postgres-dev psql -U postgres -d fiber_gorm_db
    fi
}

db_backup() {
    print_header "Creating Database Backup"
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    BACKUP_FILE="backup_${TIMESTAMP}.sql"
    
    if [ "$1" == "prod" ]; then
        docker exec fiber-postgres-prod pg_dump -U postgres fiber_gorm_db > "$BACKUP_FILE"
    else
        docker exec fiber-postgres-dev pg_dump -U postgres fiber_gorm_db > "$BACKUP_FILE"
    fi
    
    print_message "Backup created: $BACKUP_FILE" "$GREEN"
}

db_restore() {
    if [ -z "$2" ]; then
        print_message "Usage: ./scripts/run.sh db:restore [dev|prod] <backup_file>" "$RED"
        exit 1
    fi
    
    print_header "Restoring Database from $2"
    
    if [ "$1" == "prod" ]; then
        cat "$2" | docker exec -i fiber-postgres-prod psql -U postgres -d fiber_gorm_db
    else
        cat "$2" | docker exec -i fiber-postgres-dev psql -U postgres -d fiber_gorm_db
    fi
    
    print_message "Database restored from $2" "$GREEN"
}

# Clean commands
clean_all() {
    print_header "Cleaning All Docker Resources"
    docker compose -f docker-compose.dev.yml down -v --rmi all --remove-orphans 2>/dev/null
    docker compose -f docker-compose.prod.yml down -v --rmi all --remove-orphans 2>/dev/null
    print_message "All Docker resources cleaned!" "$GREEN"
}

clean_volumes() {
    print_header "Cleaning Volumes"
    docker compose -f docker-compose.dev.yml down -v 2>/dev/null
    docker compose -f docker-compose.prod.yml down -v 2>/dev/null
    print_message "Volumes cleaned!" "$GREEN"
}

# Status command
status() {
    print_header "Docker Container Status"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
}

# Help command
show_help() {
    print_header "Available Commands"
    echo "Development:"
    echo "  dev:up        - Start development environment with hot reload"
    echo "  dev:down      - Stop development environment"
    echo "  dev:logs      - View all development logs"
    echo "  dev:logs:app  - View app logs only (development)"
    echo "  dev:restart   - Restart development environment"
    echo "  dev:rebuild   - Rebuild and restart development environment"
    echo ""
    echo "Production:"
    echo "  prod:up       - Start production environment"
    echo "  prod:down     - Stop production environment"
    echo "  prod:logs     - View all production logs"
    echo "  prod:logs:app - View app logs only (production)"
    echo "  prod:restart  - Restart production environment"
    echo "  prod:rebuild  - Rebuild and restart production environment"
    echo ""
    echo "Database:"
    echo "  db:shell [dev|prod]              - Connect to PostgreSQL shell"
    echo "  db:backup [dev|prod]             - Create database backup"
    echo "  db:restore [dev|prod] <file>     - Restore database from backup"
    echo ""
    echo "Utilities:"
    echo "  status        - Show container status"
    echo "  clean:all     - Remove all Docker resources"
    echo "  clean:volumes - Remove volumes only"
    echo "  help          - Show this help message"
}

# Main command handler
case "$1" in
    # Development
    "dev:up")
        dev_up
        ;;
    "dev:down")
        dev_down
        ;;
    "dev:logs")
        dev_logs
        ;;
    "dev:logs:app")
        dev_logs_app
        ;;
    "dev:restart")
        dev_restart
        ;;
    "dev:rebuild")
        dev_rebuild
        ;;
    
    # Production
    "prod:up")
        prod_up
        ;;
    "prod:down")
        prod_down
        ;;
    "prod:logs")
        prod_logs
        ;;
    "prod:logs:app")
        prod_logs_app
        ;;
    "prod:restart")
        prod_restart
        ;;
    "prod:rebuild")
        prod_rebuild
        ;;
    
    # Database
    "db:shell")
        db_shell "$2"
        ;;
    "db:backup")
        db_backup "$2"
        ;;
    "db:restore")
        db_restore "$2" "$3"
        ;;
    
    # Utilities
    "status")
        status
        ;;
    "clean:all")
        clean_all
        ;;
    "clean:volumes")
        clean_volumes
        ;;
    "help"|"")
        show_help
        ;;
    *)
        print_message "Unknown command: $1" "$RED"
        show_help
        exit 1
        ;;
esac
