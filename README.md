# tasks-go

## How to run

1. Build docker image
    ```bash
    docker build -t tasks-go .
    ```

2. Run docker container 
    
    ```bash
    docker run -d -p 8080:8080 tasks-go:latest
    ```

3. Happy adding tasks

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"name":"Add more tasks"}' http://localhost:8080/tasks
    ```
