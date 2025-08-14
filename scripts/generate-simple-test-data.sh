#!/bin/bash

# Simple script to generate limited test data for existing tables only
# This works with the current database schema

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
SQL_FILE="$SCRIPT_DIR/generate_simple_test_data.sql"
POSTGRES_CONTAINER="erp-postgres"
POSTGRES_DB="erp_db"
POSTGRES_USER="erp_user"

print_info "Running simple test data generation..."
print_warning "This only creates data for existing tables (office, position, role)"

# Check if container is running
if ! docker ps | grep -q "$POSTGRES_CONTAINER"; then
    print_error "PostgreSQL container '$POSTGRES_CONTAINER' is not running!"
    exit 1
fi

# Copy and execute SQL
docker cp "$SQL_FILE" "$POSTGRES_CONTAINER:/tmp/generate_simple_test_data.sql"
docker exec -it "$POSTGRES_CONTAINER" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /tmp/generate_simple_test_data.sql
docker exec "$POSTGRES_CONTAINER" rm -f /tmp/generate_simple_test_data.sql

print_success "Simple test data created!"
print_warning "To create full test data, the backend needs to run migrations first."
print_info "Check that all models in backend/config/db.go AutoMigrate() are uncommented."
