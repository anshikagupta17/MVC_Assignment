#!/bin/sh
set -e

echo "Waiting for database to be ready..."
until nc -z "$DB_HOST" "$DB_PORT"; do
  sleep 1
done
echo "Database is ready."

echo "Running migrations..."
migrate -path ./db/migrations -database "$DB_URI" up

echo "Starting server..."
exec ./server
