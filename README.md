# tasks-go

## Prerequisites

   - Prepare your own `.env` file

       ```bash
       cp .env.example .env
       ```
   
## How to run

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