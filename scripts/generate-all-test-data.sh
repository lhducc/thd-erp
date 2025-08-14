#!/bin/bash

# Script to generate comprehensive test data for the ERP system
# This script runs all data generation scripts in the correct order

set -e  # Exit on any error

echo "🚀 Starting comprehensive test data generation..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-erp_db}
DB_USER=${DB_USER:-erp_user}
DB_PASSWORD=${DB_PASSWORD:-erp_password}

# Function to run SQL script with error handling
run_sql_script() {
    local script_path=$1
    local script_name=$(basename "$script_path")
    
    echo -e "${BLUE}📄 Running $script_name...${NC}"
    
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$script_path"; then
        echo -e "${GREEN}✅ $script_name completed successfully${NC}"
    else
        echo -e "${RED}❌ Error running $script_name${NC}"
        exit 1
    fi
    echo ""
}

# Function to check database connection
check_db_connection() {
    echo -e "${BLUE}🔍 Checking database connection...${NC}"
    
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Database connection successful${NC}"
    else
        echo -e "${RED}❌ Cannot connect to database. Please check connection parameters:${NC}"
        echo "  Host: $DB_HOST"
        echo "  Port: $DB_PORT"
        echo "  Database: $DB_NAME"
        echo "  User: $DB_USER"
        exit 1
    fi
    echo ""
}

# Function to show database status
show_db_status() {
    echo -e "${BLUE}📊 Database Status:${NC}"
    
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        SELECT 
            schemaname,
            tablename,
            n_tup_ins as rows_inserted,
            n_tup_upd as rows_updated,
            n_tup_del as rows_deleted
        FROM pg_stat_user_tables 
        WHERE schemaname = 'public' 
        ORDER BY tablename;
    "
    echo ""
}

# Main execution
main() {
    echo -e "${YELLOW}=====================================${NC}"
    echo -e "${YELLOW}  ERP System Data Generation Script  ${NC}"
    echo -e "${YELLOW}=====================================${NC}"
    echo ""
    
    # Check if we're in the right directory
    if [ ! -f "scripts/create_database_schema.sql" ]; then
        echo -e "${RED}❌ Please run this script from the ERP project root directory${NC}"
        exit 1
    fi
    
    # Check database connection
    check_db_connection
    
    # Step 1: Create/Update database schema
    echo -e "${YELLOW}Step 1: Creating/Updating database schema...${NC}"
    run_sql_script "scripts/create_database_schema.sql"
    
    # Step 2: Generate basic test data (employees, offices, etc.)
    echo -e "${YELLOW}Step 2: Generating basic test data...${NC}"
    run_sql_script "scripts/generate_test_data.sql"
    
    # Step 3: Generate comprehensive test data (contracts, attendance, etc.)
    echo -e "${YELLOW}Step 3: Generating comprehensive test data...${NC}"
    run_sql_script "scripts/generate_comprehensive_test_data.sql"
    
    # Show final status
    echo -e "${GREEN}🎉 All data generation completed successfully!${NC}"
    echo ""
    
    # Show database statistics
    show_db_status
    
    # Show summary
    echo -e "${BLUE}📋 Summary:${NC}"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        SELECT 
            'Total Tables' as metric,
            COUNT(*) as value 
        FROM information_schema.tables 
        WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
        
        UNION ALL
        
        SELECT 
            'Total Employees' as metric,
            COUNT(*)::TEXT as value 
        FROM employee
        
        UNION ALL
        
        SELECT 
            'Total Accounts' as metric,
            COUNT(*)::TEXT as value 
        FROM account
        
        UNION ALL
        
        SELECT 
            'Total Attendance Records' as metric,
            COUNT(*)::TEXT as value 
        FROM attendance_record
        
        UNION ALL
        
        SELECT 
            'Total Contracts' as metric,
            COUNT(*)::TEXT as value 
        FROM contract
        
        ORDER BY metric;
    "
    
    echo ""
    echo -e "${GREEN}✨ Ready to use! You can now start the application with all test data.${NC}"
    echo -e "${BLUE}💡 Default admin login: admin@company.vn / password123${NC}"
    echo -e "${BLUE}💡 Default manager login: manager@company.vn / password123${NC}"
    echo -e "${BLUE}💡 Employee logins: [employee-name]@company.vn / password123${NC}"
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [options]"
        echo ""
        echo "Environment variables:"
        echo "  DB_HOST      Database host (default: localhost)"
        echo "  DB_PORT      Database port (default: 5432)"
        echo "  DB_NAME      Database name (default: erp_db)"
        echo "  DB_USER      Database user (default: erp_user)"
        echo "  DB_PASSWORD  Database password (default: erp_password)"
        echo ""
        echo "Example:"
        echo "  DB_HOST=192.168.1.100 DB_PASSWORD=mypass $0"
        exit 0
        ;;
    *)
        main
        ;;
esac
