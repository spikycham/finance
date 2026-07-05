#!/usr/bin/env bash
set -e

SERVER="cham@devcham.xyz"
REMOTE_DB="~/projects/go/finance/data.db"
LOCAL_DIR="./backup"

mkdir -p "$LOCAL_DIR"

scp "$SERVER:$REMOTE_DB" \
    "$LOCAL_DIR/$(date +%F).db"

echo "Backup saved to $LOCAL_DIR/$(date +%F).db"
