#### Photo Sharing Toy Project ####

## Objective ##
Build a simple photo sharing webapp to demonstrate use to DB, Remote Storage, UI Framework

## Components ##
1. Sql DB
2. MinIO
3. API

## Usage ##
# Run Locally #
This project is built on docker containers and is intended to be ran locally through docker-compose.

NOTE: To allow minio signed urls to resolve to the containerized service, add `minio` to your hosts file:
```
127.0.0.1       minio
```

Run the application:
```
docker-compose up --build
```

## Interaction ##

The wepapp can be viewed at `http://localhost:8080`
