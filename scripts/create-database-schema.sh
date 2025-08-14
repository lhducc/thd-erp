#!/bin/bash

# Script to create database schema from scratch
# This replaces GORM AutoMigrate with pure SQL table creation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SQL_FILE="$SCRIPT_DIR/create_database_schema.sql"
POSTGRES_CONTAINER="erp-postgres"
POSTGRES_DB="erp_db"
POSTGRES_USER="erp_user"

# Function to check if container exists and is running
check_postgres_container() {
    if ! docker ps | grep -q "$POSTGRES_CONTAINER"; then
        print_error "PostgreSQL container '$POSTGRES_CONTAINER' is not running!"
        print_info "Please start your PostgreSQL container first:"
        print_info "cd ../deploys && ./run-uat.sh up postgres"
        exit 1
    fi
}

# Function to check if SQL file exists
check_sql_file() {
    if [ ! -f "$SQL_FILE" ]; then
        print_error "SQL file not found: $SQL_FILE"
        exit 1
    fi
}

# Function to run SQL file
run_sql_file() {
    print_info "Creating database schema from scratch..."
    print_warning "This will drop existing tables and recreate them!"
    
    # Copy SQL file to container and execute it
    docker cp "$SQL_FILE" "$POSTGRES_CONTAINER:/tmp/create_database_schema.sql"
    
    # Execute the SQL file
    docker exec -i "$POSTGRES_CONTAINER" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /tmp/create_database_schema.sql
    
    # Clean up
    docker exec "$POSTGRES_CONTAINER" rm -f /tmp/create_database_schema.sql
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Create database schema from scratch (replaces GORM AutoMigrate)"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  --check        Check prerequisites only (don't run)"
    echo "  --force        Skip confirmation prompt"
    echo ""
    echo "Prerequisites:"
    echo "  - PostgreSQL container 'erp-postgres' must be running"
    echo "  - Database 'erp_db' must exist"
    echo "  - User 'erp_user' must have access"
    echo ""
    echo "Example:"
    echo "  $0              Create schema (with confirmation)"
    echo "  $0 --force      Create schema (no confirmation)"
    echo "  $0 --check     Check if prerequisites are met"
}

# Main function
main() {
    case "${1:-}" in
        "--help"|"-h")
            show_usage
            exit 0
            ;;
        "--check")
            print_info "Checking prerequisites..."
            check_postgres_container
            check_sql_file
            print_success "All prerequisites met!"
            exit 0
            ;;
        "--force")
            print_info "Starting database schema creation (forced)..."
            check_postgres_container
            check_sql_file
            run_sql_file
            print_success "Database schema created successfully!"
            ;;
        "")
            print_info "Starting database schema creation..."
            check_postgres_container
            check_sql_file
            
            print_warning "This will DROP existing tables and recreate them!"
            print_warning "All existing data will be LOST!"
            read -p "Continue? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                run_sql_file
                print_success "Database schema created successfully!"
                print_info ""
                print_info "Next steps:"
                print_info "1. Update backend config to disable AutoMigrate"
                print_info "2. Run './generate-test-data.sh' to populate with sample data"
                print_info "3. Restart backend: cd ../deploys && ./run-uat.sh restart backend"
            else
                print_info "Operation cancelled"
                exit 1
            fi
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
