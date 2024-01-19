# tasks-go

## Prerequisites

- Prepare your own `.env` file

  ```bash
  cp .env.example .env
  ```

## Run with Go

1. Install dependencies

   ```bash
   go mod download
   ```

2. Run the server

   ```bash
   go run cmd/tasks/main.go
   ```

## Run with Docker

1. Build docker image

   ```bash
   docker build -t tasks-go .
   ```

2. Run docker container

   ```bash
   docker run -d -p 8080:8080 --env-file ./.env tasks-go:latest
   ```

3. Happy adding tasks

   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"name":"Add more tasks"}' http://localhost:8080/tasks
   ```

4. Stop the container

   ```bash
   $ docker ps
   CONTAINER ID   IMAGE             COMMAND      CREATED         STATUS         PORTS                    NAMES
   b17486cbd255   tasks-go:latest   "./server"   3 minutes ago   Up 3 minutes   0.0.0.0:8080->8080/tcp   condescending_sammet

   # docker stop using CONTAINER ID
   $ docker stop b17486cbd255
   ```

## Authentication

All API endpoints are protected by JWT authentication. To get the token, you need to login first.

   ```bash
   echo -n "username:password" | base64
   # encodeString: "dXNlcm5hbWU6cGFzc3dvcmQ="

   curl --location --request POST 'http://localhost:8080/auth' \
   --header 'Authorization: Basic <encodeString>'
   ```

## API Endpoints

1. Create Task

    ```bash
    curl --location 'http://localhost:8080/task' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Bearer <jwtToken>' \
    --data '{
      "name": "eat dinner"
    }'
    ```

2. Get Tasks

   ```bash
   curl --location 'http://localhost:8080/tasks' \
   --header 'Authorization: Bearer <jwtToken>'
   ```

3. Update Task

   ```bash
   curl --location --request PUT 'http://localhost:8080/task/<taskId>' \
   --header 'Content-Type: application/json' \
   --header 'Authorization: Bearer <jwtToken>'
   --data '{
      "name": "go climbing", // optional
      "status": 1 // optional
   }'
   ```

4. Delete Task

   ```bash
   curl --location --request DELETE 'http://localhost:8080/task/2' \
   --header 'Authorization: Bearer <jwtToken>'
   ```
