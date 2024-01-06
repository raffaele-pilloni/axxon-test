# raffaele-pilloni/axxon-test
Simple application to run asynchronous HTTP request to 3rd-party.

POST and PUT requests only handle JSON type payloads.

## Development environment
To ease local development you have to install these tools:

* [Docker CE](https://www.docker.com/)
* [Docker-Compose](https://docs.docker.com/compose/)

In the project's root folder there is the main file where the architecture is described: `docker-compose.yml`.
In there you can find all services configured that you'll get.

In order to have everything work correctly you have to copy the `.env.dist`
file, located in the root folder to the `.env` in the root folder.

The `.env` file contains environment variables used for configuring 
the application and the `docker-compose.yml` architecture.

The `docker-compose.override.yml` contains the configurations for using the debugger on containers. 
Uncomment and configure a remote debug client.

For Goland IDE configure as follows:
* [Goland Remote Debug](https://www.jetbrains.com/help/go/go-remote.html)
* 
### Start environment

To start your environment, execute these commands:
```sh
make start
```

### Stop environment

To stop your environment, execute this command:
```sh
make stop
```

### Destroy environment

To destroy your environment, execute this command:
```sh
make destroy
```

### Reset environment

To restart your environment, execute this command:
```sh
make restart
```

### Monitor environment

To monitor your environment, execute this command:
```sh
make ps
```

To follow all logs, execute this command:
```sh
make logs
```

To follow single service logs, execute this command:
```sh
make service-logs
``` 
To specify the service pass the SERVICE_NAME variable:
```sh
SERVICE_NAME=<service_name> make service-logs
``` 

To monitor your environment, execute this command:
```sh
make ps
```

To run terminal in the service, execute this command:
```sh
make terminal
```
To specify the service pass the SERVICE_NAME variable:
```sh
SERVICE_NAME=<service_name> make terminal
``` 
### Code style

The project depends on `golangci-lint` to check the style of the source code.

To run the checks, execute this command:
```sh
make lint
``` 

### Mocking dependencies
The projects depends on `mockery` to mock dependencies.

To generate all mocks, execute this command:
```sh
make mock
``` 

### Testing
The projects depends both on `gomega` and `ginkgo` to test the source code.

To run all tests, execute this command:
```sh
make test
```

### Match example

To create task send this request:
```sh
curl --request POST 'http://127.0.0.1:8081/task' \
--header 'Content-Type: application/json'
--data-raw '{
    "method": "POST",
    "url": "https://webhook.site/0642f579-230f-491c-9370-819ff40766b9",
    "headers": {
        "test": [
            "test"
        ]
    },
    "body": {
        "test": "test"
    }
}'
```

To get task send this request:
```sh
curl --request GET 'http://127.0.0.1:8081/task/1' \
--header 'Content-Type: application/json'
--data-raw '{
    "id": 1,
    "status": "done",
    "httpStatusCode": 200,
    "headers": {
        "Cache-Control": [
            "no-cache, private"
        ],
        "Content-Type": [
            "text/html; charset=UTF-8"
        ],
        "Date": [
            "Sat, 06 Jan 2024 03:31:34 GMT"
        ],
        "Server": [
            "nginx"
        ],
        "Vary": [
            "Accept-Encoding"
        ],
        "X-Request-Id": [
            "50ad8cb0-d82f-4994-95e8-526c4db3c3db"
        ],
        "X-Token-Id": [
            "0642f579-230f-491c-9370-819ff40766b9"
        ]
    },
    "length": 20
}'
```

