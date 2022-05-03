## How to run ##
- Creation of a server (Execute from the ``interact/server/`` folder)
    ```
    go run server_exec.go --addr=":8080" -v=0
    ```
- Run tests of the server handlers (Execute from ``interact/server/test/`` folder)
    ```
    go test -v
    ```