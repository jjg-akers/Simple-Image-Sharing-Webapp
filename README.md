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

## TODOs ##
1. Upgrade UI
2. Handle View All
3. Implement Search
4. Delete an image/image(s)
5. View an individual image
6. “Permalink” to an individual image
7. User Input Validation
8. resize /format - standardize uploaded photos

