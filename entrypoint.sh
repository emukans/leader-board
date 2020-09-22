#!/usr/bin/env sh

set -e

db/migrate.sh

exec "$@"
