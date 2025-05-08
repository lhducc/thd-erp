#!/bin/bash
set -e

MODE=$1

if [[ "$MODE" != "up" && "$MODE" != "down" ]]; then
  echo "Usage: $0 [up|down]"
  exit 1
fi

if [[ "$MODE" == "up" ]]; then
  docker compose -f ./deploys/docker-compose.dev.infras.yml up -d

  echo "Waiting for infrastructure to be ready..."
  sleep 2

  echo "Running backend locally..."
  cd backend
  make run
else
  docker compose -f ./deploys/docker-compose.dev.infras.yml down
fi

