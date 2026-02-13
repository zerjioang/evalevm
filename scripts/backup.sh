#!/bin/bash

echo "Backing up Docker images..."
docker images --format '{{.Repository}}:{{.Tag}} {{.ID}}' \
| while read name id; do
  safe_name=$(echo "$name" | tr '/:' '_')
  docker save "$id" | gzip > "${safe_name}.tar.gz"
done