#!/bin/bash
set -e

MODE=$1

if [[ "$MODE" != "up" && "$MODE" != "down" ]]; then
  echo "Usage: $0 [up|down]"
  exit 1
fi

if [[ "$MODE" == "up" ]]; then
  docker compose -f ./deploys/docker-compose.dev.be.yml up -d

  echo "Waiting for backend to be ready..."
  sleep 2

  echo "Running frontend locally..."
  cd frontend
  npm install
  npm run dev
else
  docker compose -f ./deploys/docker-compose.dev.be.yml down
fi
