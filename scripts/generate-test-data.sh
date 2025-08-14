#!/bin/bash

# Script to generate test data for ERP PostgreSQL database
# This script connects to the PostgreSQL container and runs the test data generation SQL

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SQL_FILE="$SCRIPT_DIR/generate_test_data.sql"

# PostgreSQL connection details
POSTGRES_CONTAINER="erp-postgres"
POSTGRES_DB="erp_db"
POSTGRES_USER="erp_user"

# Function to check if container exists and is running
check_postgres_container() {
    if ! docker ps | grep -q "$POSTGRES_CONTAINER"; then
        print_error "PostgreSQL container '$POSTGRES_CONTAINER' is not running!"
        print_info "Please start your PostgreSQL container first:"
        print_info "cd ../deploys && docker-compose -f docker-compose.uat.yml up postgres -d"
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
    print_info "Connecting to PostgreSQL container..."
    print_info "Running test data generation script..."
    
    # Copy SQL file to container and execute it
    docker cp "$SQL_FILE" "$POSTGRES_CONTAINER:/tmp/generate_test_data.sql"
    
    # Execute the SQL file
    docker exec -it "$POSTGRES_CONTAINER" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /tmp/generate_test_data.sql
    
    # Clean up
    docker exec "$POSTGRES_CONTAINER" rm -f /tmp/generate_test_data.sql
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Generate test data for ERP PostgreSQL database"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  --check        Check prerequisites only (don't run)"
    echo ""
    echo "Prerequisites:"
    echo "  - PostgreSQL container 'erp-postgres' must be running"
    echo "  - Database 'erp_db' must exist"
    echo "  - User 'erp_user' must have access"
    echo ""
    echo "Example:"
    echo "  $0              Generate test data"
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
        "")
            print_info "Starting test data generation..."
            check_postgres_container
            check_sql_file
            
            print_warning "This will add test data to your database!"
            read -p "Continue? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                run_sql_file
                print_success "Test data generation completed!"
                print_info "Generated data includes:"
                print_info "  - 3 offices (Hanoi, Ho Chi Minh City, Da Nang)"
                print_info "  - 6 departments"
                print_info "  - Sample positions and job titles"
                print_info "  - 50 random employees"
                print_info "  - 22 user accounts (with default password: 'password123')"
                print_info ""
                print_info "Login credentials:"
                print_info "  - admin@company.vn / password123 (admin role)"
                print_info "  - manager@company.vn / password123 (manager role)"
                print_info "  - Various employee accounts (employee role)"
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
