#!/bin/sh

docker network create range-allocator-network

docker run --name postgres \
  --network range-allocator-network \
  -e POSTGRES_USER=username \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=alloc \
  -p 5432:5432 \
  -v /home/raj/projects/golang/gotiny-range-allocator/db_data:/var/lib/postgresql/data \
  -d postgres:latest

# docker exec -it postgres psql -U username -d alloc
