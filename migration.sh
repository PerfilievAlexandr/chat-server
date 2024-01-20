#!/bin/bash
source .env

export MIGRATION_DSN="host=$MIGRATION_HOST port=$MIGRATION_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD sslmode=disable"

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v