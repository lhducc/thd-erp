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

# SQL files in order of execution
SCHEMA_FILE="$SCRIPT_DIR/create_database_schema.sql"
SIMPLE_DATA_FILE="$SCRIPT_DIR/generate_simple_test_data.sql"
BASIC_DATA_FILE="$SCRIPT_DIR/generate_test_data.sql"
COMPREHENSIVE_DATA_FILE="$SCRIPT_DIR/generate_comprehensive_test_data.sql"

# PostgreSQL connection details for Docker
POSTGRES_CONTAINER="erp-postgres-1"
POSTGRES_DB="erp_db"
POSTGRES_USER="erp_user"

# Function to check if container exists and is running
check_postgres_container() {
    if ! docker ps | grep -q "$POSTGRES_CONTAINER\|postgres"; then
        print_error "PostgreSQL container is not running!"
        print_info "Please start your PostgreSQL container first:"
        print_info "cd ../deploys && docker-compose -f docker-compose.uat.yml up postgres -d"
        
        # Try to find the correct container name
        local containers=$(docker ps --format "table {{.Names}}" | grep -i postgres || true)
        if [ ! -z "$containers" ]; then
            print_info "Found these PostgreSQL containers:"
            echo "$containers"
            print_info "You may need to update POSTGRES_CONTAINER in this script"
        fi
        exit 1
    fi
    
    # Try to detect the actual container name
    local actual_container=$(docker ps --format "{{.Names}}" | grep -E "(postgres|erp.*postgres)" | head -1)
    if [ ! -z "$actual_container" ]; then
        POSTGRES_CONTAINER="$actual_container"
        print_info "Using PostgreSQL container: $POSTGRES_CONTAINER"
    fi
}

# Function to check if SQL file exists
check_sql_file() {
    local file=$1
    local name=$2
    if [ ! -f "$file" ]; then
        print_error "$name file not found: $file"
        exit 1
    fi
}

# Function to run SQL file
run_sql_file() {
    local sql_file=$1
    local description=$2
    
    print_info "$description"
    
    # Copy SQL file to container and execute it
    local container_path="/tmp/$(basename "$sql_file")"
    docker cp "$sql_file" "$POSTGRES_CONTAINER:$container_path"
    
    # Execute the SQL file
    if docker exec "$POSTGRES_CONTAINER" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$container_path"; then
        print_success "$description completed successfully"
    else
        print_error "$description failed"
        return 1
    fi
    
    # Clean up
    docker exec "$POSTGRES_CONTAINER" rm -f "$container_path"
    echo ""
}

# Function to show database status
show_database_status() {
    print_info "Checking database status..."
    docker exec "$POSTGRES_CONTAINER" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "
        SELECT 
            schemaname,
            tablename,
            n_tup_ins as total_rows
        FROM pg_stat_user_tables 
        WHERE schemaname = 'public' 
          AND n_tup_ins > 0
        ORDER BY tablename;
    " 2>/dev/null || print_warning "Could not fetch database statistics"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Generate test data for ERP PostgreSQL database"
    echo ""
    echo "Options:"
    echo "  -h, --help         Show this help message"
    echo "  --check            Check prerequisites only (don't run)"
    echo "  --simple           Run only simple data generation (basic lookup tables)"
    echo "  --basic            Run simple + basic data (employees, accounts)"
    echo "  --comprehensive    Run all data generation (full ERP data)"
    echo "  --schema-only      Create/update database schema only"
    echo ""
    echo "Prerequisites:"
    echo "  - PostgreSQL container must be running"
    echo "  - Database 'erp_db' must exist"
    echo "  - User 'erp_user' must have access"
    echo ""
    echo "Examples:"
    echo "  $0                     Run comprehensive data generation"
    echo "  $0 --simple           Generate only lookup tables"
    echo "  $0 --basic            Generate lookup tables + employees"
    echo "  $0 --schema-only      Update database schema only"
}

# Main function
main() {
    local mode="comprehensive"
    
    case "${1:-}" in
        "--help"|"-h")
            show_usage
            exit 0
            ;;
        "--check")
            print_info "Checking prerequisites..."
            check_postgres_container
            check_sql_file "$SCHEMA_FILE" "Database schema"
            check_sql_file "$SIMPLE_DATA_FILE" "Simple test data"
            check_sql_file "$BASIC_DATA_FILE" "Basic test data"
            check_sql_file "$COMPREHENSIVE_DATA_FILE" "Comprehensive test data"
            print_success "All prerequisites met!"
            exit 0
            ;;
        "--simple")
            mode="simple"
            ;;
        "--basic")
            mode="basic"
            ;;
        "--comprehensive"|"")
            mode="comprehensive"
            ;;
        "--schema-only")
            mode="schema"
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
    
    print_info "Starting test data generation (mode: $mode)..."
    check_postgres_container
    
    # Always ensure schema is up to date
    check_sql_file "$SCHEMA_FILE" "Database schema"
    run_sql_file "$SCHEMA_FILE" "Creating/updating database schema"
    
    if [ "$mode" = "schema" ]; then
        print_success "Database schema update completed!"
        exit 0
    fi
    
    print_warning "This will add test data to your database!"
    read -p "Continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Operation cancelled"
        exit 1
    fi
    
    # Run data generation based on mode
    case "$mode" in
        "simple")
            check_sql_file "$SIMPLE_DATA_FILE" "Simple test data"
            run_sql_file "$SIMPLE_DATA_FILE" "Generating simple test data (lookup tables)"
            ;;
        "basic")
            check_sql_file "$SIMPLE_DATA_FILE" "Simple test data"
            check_sql_file "$BASIC_DATA_FILE" "Basic test data"
            run_sql_file "$SIMPLE_DATA_FILE" "Generating simple test data (lookup tables)"
            run_sql_file "$BASIC_DATA_FILE" "Generating basic test data (employees, accounts)"
            ;;
        "comprehensive")
            check_sql_file "$SIMPLE_DATA_FILE" "Simple test data"
            check_sql_file "$BASIC_DATA_FILE" "Basic test data"
            check_sql_file "$COMPREHENSIVE_DATA_FILE" "Comprehensive test data"
            run_sql_file "$SIMPLE_DATA_FILE" "Generating simple test data (lookup tables)"
            run_sql_file "$BASIC_DATA_FILE" "Generating basic test data (employees, accounts)"
            run_sql_file "$COMPREHENSIVE_DATA_FILE" "Generating comprehensive test data (contracts, attendance, etc.)"
            ;;
    esac
    
    print_success "Test data generation completed!"
    
    # Show what was created
    case "$mode" in
        "simple")
            print_info "Generated data includes:"
            print_info "  - Offices, departments, positions, job titles"
            print_info "  - Work shifts and attendance categories"
            print_info "  - Role definitions"
            ;;
        "basic")
            print_info "Generated data includes:"
            print_info "  - All simple data (lookup tables)"
            print_info "  - 50 random employees with Vietnamese names"
            print_info "  - 22 user accounts with default password: 'password123'"
            ;;
        "comprehensive")
            print_info "Generated data includes:"
            print_info "  - All basic data (employees, accounts)"
            print_info "  - Employee documents and contracts"
            print_info "  - Decisions, allowances, and insurance records"
            print_info "  - Work schedules and attendance records (30 days)"
            print_info "  - Complete ERP system data"
            ;;
    esac
    
    print_info ""
    print_info "Login credentials:"
    print_info "  - admin@company.vn / password123 (admin role)"
    print_info "  - manager@company.vn / password123 (manager role)"
    print_info "  - Various employee accounts: [name]@company.vn / password123"
    
    # Show database status
    show_database_status
}

# Run main function
main "$@"
