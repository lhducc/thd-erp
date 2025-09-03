#!/bin/bash

# ERP UAT Environment Runner Script
# This script manages the UAT environment using docker-compose

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$PROJECT_ROOT/deploys/docker-compose.uat.yml"

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if docker-compose file exists
check_compose_file() {
    if [ ! -f "$COMPOSE_FILE" ]; then
        print_error "Docker compose file not found: $COMPOSE_FILE"
        exit 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  up, start         Start all services"
    echo "  down, stop        Stop all services"
    echo "  restart           Restart all services"
    echo "  rebuild           Rebuild and start all services"
    echo "  logs              Show logs for all services"
    echo "  status, ps        Show status of all services"
    echo "  clean             Stop services and remove volumes"
    echo ""
    echo "Options:"
    echo "  -d, --detach      Run in detached mode (background)"
    echo "  -f, --follow      Follow logs output"
    echo "  --build           Force rebuild when starting"
    echo ""
    echo "Examples:"
    echo "  $0 up              Start services in foreground"
    echo "  $0 up -d           Start services in background"
    echo "  $0 logs -f         Follow logs output"
    echo "  $0 restart         Restart all services"
    echo "  $0 clean           Stop and clean everything"
}

# Function to run docker-compose command
run_compose() {
    cd "$PROJECT_ROOT/deploys"
    docker compose -f docker-compose.uat.yml --env-file "$PROJECT_ROOT/.env" "$@"
}

# Main script logic
main() {
    check_compose_file

    case "${1:-}" in
        "up"|"start")
            shift
            print_info "Starting UAT environment..."
            if [[ "$*" == *"-d"* ]] || [[ "$*" == *"--detach"* ]]; then
                run_compose up -d "$@"
                print_success "UAT environment started in background"
                print_info "Use '$0 logs' to view logs"
                print_info "Use '$0 status' to check service status"
            else
                run_compose up "$@"
            fi
            ;;
        
        "down"|"stop")
            shift
            print_info "Stopping UAT environment..."
            run_compose down "$@"
            print_success "UAT environment stopped"
            ;;
        
        "restart")
            shift
            print_info "Restarting UAT environment..."
            run_compose restart "$@"
            print_success "UAT environment restarted"
            ;;
        
        "rebuild")
            shift
            print_info "Rebuilding and starting UAT environment..."
            run_compose up --build -d "$@"
            print_success "UAT environment rebuilt and started"
            ;;
        
        "logs")
            shift
            if [[ "$*" == *"-f"* ]] || [[ "$*" == *"--follow"* ]]; then
                print_info "Following logs... (Press Ctrl+C to stop)"
                run_compose logs -f "$@"
            else
                run_compose logs "$@"
            fi
            ;;
        
        "status"|"ps")
            shift
            print_info "UAT environment status:"
            run_compose ps "$@"
            ;;
        
        "clean")
            shift
            print_warning "This will stop all services and remove volumes!"
            read -p "Are you sure? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                print_info "Cleaning UAT environment..."
                run_compose down -v "$@"
                print_success "UAT environment cleaned"
            else
                print_info "Operation cancelled"
            fi
            ;;
        
        "help"|"-h"|"--help")
            show_usage
            ;;
        
        "")
            print_error "No command specified"
            show_usage
            exit 1
            ;;
        
        *)
            print_error "Unknown command: $1"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
