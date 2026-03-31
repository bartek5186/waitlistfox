#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

COMPOSE_FILE="${COMPOSE_FILE:-compose.prod.yaml}"
ENV_FILE="${ENV_FILE:-.env.production}"
SERVICE="${SERVICE:-api}"
RUN_GIT_PULL=1
RUN_MIGRATIONS=1
ACTION="deploy"

usage() {
  cat <<'EOF'
Usage:
  ./deploy.sh [deploy|build|up|migrate|logs|status] [--skip-git-pull] [--skip-migrate]

Examples:
  ./deploy.sh
  ./deploy.sh deploy --skip-git-pull
  ./deploy.sh migrate
  ./deploy.sh logs
EOF
}

compose_cmd() {
  if [[ -f "$ENV_FILE" ]]; then
    docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" "$@"
  else
    docker compose -f "$COMPOSE_FILE" "$@"
  fi
}

ensure_runtime_files() {
  mkdir -p config log

  if [[ ! -f config/config.json ]]; then
    echo "Missing config/config.json. Create it first, for example from config/config.example.json." >&2
    exit 1
  fi
}

run_git_pull() {
  if (( RUN_GIT_PULL )); then
    git pull --ff-only
  fi
}

run_migrations() {
  compose_cmd run --rm "$SERVICE" ./waitlistfox -config /app/config/config.json -migrate=true
}

build_service() {
  compose_cmd build "$SERVICE"
}

start_service() {
  compose_cmd up -d --remove-orphans "$SERVICE"
}

show_status() {
  compose_cmd ps
}

follow_logs() {
  compose_cmd logs -f "$SERVICE"
}

if [[ $# -gt 0 ]]; then
  case "$1" in
    deploy|build|up|migrate|logs|status)
      ACTION="$1"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
  esac
fi

while [[ $# -gt 0 ]]; do
  case "$1" in
    --skip-git-pull)
      RUN_GIT_PULL=0
      ;;
    --skip-migrate)
      RUN_MIGRATIONS=0
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
  shift
done

ensure_runtime_files

case "$ACTION" in
  deploy)
    run_git_pull
    build_service
    if (( RUN_MIGRATIONS )); then
      run_migrations
    fi
    start_service
    show_status
    ;;
  build)
    build_service
    ;;
  up)
    start_service
    show_status
    ;;
  migrate)
    run_migrations
    ;;
  logs)
    follow_logs
    ;;
  status)
    show_status
    ;;
esac
