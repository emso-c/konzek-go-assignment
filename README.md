# Task definiton

Assignment: Building a Concurrent Web Service in Go

Scenario: You are working on a project that requires building a high-performance web service in Go. The service will accept incoming HTTP requests, process data concurrently, and store it in a database. The application should be scalable, maintainable, and well-tested.

# General Information

This project is built for Konzek's Junior Backend Developer Assignment. It is aimed to provide a simple and easy to use data storage API for tasks. It is written in Go and uses PostgreSQL as the database. The API is designed to be RESTful and ready to be consumed by any client.

The Project implements the following features:
- Setup and Basic HTTP Server
    - [x] Using the standard Go net/http package
    - [x] Implement basic request/response handling
    - [x] Create endpoints for handling POST and GET requests
      > Also implemented DELETE, PUT and OPTIONS methods
- Data Model and Persistence
    - [x] Define a data model for the application
    - [x] Set up a database connection using a Go database library
      > Used system database/sql package and github.com/lib/pq driver for PostgreSQL database
    - [x] Implement CRUD operations for the data model
- Concurrent Processing
    - [x] Modify HTTP handlers to process incoming requests concurrently using Goroutines
    - [x] Implement a worker pool or a concurrent mechanism to handle concurrent requests efficiently
    - [x] Ensure proper synchronization and error handling
- Validation and Error Handling
    - [x] Implement request validation
    - [x] Handle errors gracefully and provide meaningful error messages to clients
- API Documentation
    - [x] Create clear and concise documentation for the API
    - [x] Use GoDoc to generate API documentation
- Testing
    - [x] Write unit tests and integration tests for the application
        > Need to enable SSL for local postgres container or point to a database that has SSL enabled to run integration tests
    - [x] Use testing libraries like testing and httptest to test HTTP handlers and database operations
        > Used go-sqlmock for mocking database operations
    - [x] Ensure good test coverage
- Logging and Monitoring
    - [x] Implement logging to record important events and errors in the application
      > The logs are stored in `./app/logs/app.log` and `./app/logs/error.log`
    - [ ] Set up monitoring and metrics collection using Prometheus and Grafana
- Security
    - [x] Implement basic security measures such as input validation
    - [ ] Implement authentication and authorization
    - [x] Ensure the application is protected against common web security vulnerabilities
      > Implemented protection against SQL Injection, XSS, CSRF, and other common web security vulnerabilities and enabled CORS
- Deployment
    - [x] Prepare the application for deployment in a production environment
    - [x] Document the deployment process
- Bonus (Optional)
    - [x] Implement pagination for listing endpoints
    - [ ] Add authentication using OAuth2 or JWT
    - [x] Secure sensitive data
      >  Used environment variables in `.env` to store sensitive data. Though the file is not ignored for demonstration purposes.
    - [x] Implement rate limiting to protect against abuse

> Note: Most of the configurations can be managed by modifying the `./app/config.toml` file.


# Running the Tests

All the tests are passing as of the last commit. The tests can be run in various ways:

> ⚠️ Need to enable SSL for local postgres container or point to a database that has SSL enabled using `./app/.env` file to run integration tests

1. Using the `go test` command:
    ```bash
    cd app && go test ./...  # Discover and run all tests in the app directory
    ```
2. Using the sh file:
    ```bash
    sh ./bin/run_tests.sh  # This produces a coverage report in `./app/coverage.html` file
    ```
3. Using Docker: 
    ```bash
    cd app && docker build -t task-api-test -f Dockerfile.test . && docker run task-api-test
    ```
4. Using Makefile:
    ```bash
    make run-tests
    ```

# Viewing the Documentation

Follow the steps mentioned in `./app/README.md`

# Production Deployment

The application is containerized using Docker and can be deployed in a production environment. The application can be deployed either with or without a local database. The database can be deployed separately and the application can be configured to connect to it.

## 1. Deploying the Application without database

### Using Docker
   1. cd into the `app` directory
       ```bash
       cd app
       ```
   2. Edit the `.env` file to a point to a valid PostgreSQL database
       ```bash
       vi .env
       ```
   3. Build the Docker image
       ```bash
       docker build -t task-api .
       ```
   4. Run the Docker container
       ```bash
       docker run -p 8080:8080 task-api
       ```
   5. The application should be running on `http://localhost:8080`

### Using go run
   1. cd into the `app` directory
       ```bash
       cd app
       ```
   2. Run the application
       ```bash
       go run main.go
       ```
   3. The application should be running on `http://localhost:8080`

## 2. Deploying the Application with Database

1. Use docker-compose.yml file to deploy the application with the database with a single command
    ```bash
    docker-compose up -d --build
    ```
2. The application should be running on `http://localhost:8080`
