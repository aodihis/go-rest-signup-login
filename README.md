# Sign Up and Login using Golang

This project implements a simple sign-up and login service using Golang. It supports user registration, authentication, and basic database interactions.

## Getting Started

To run this project, follow the steps below:

1. **Copy the environment file:**
   - Copy `.env.example` and rename it to `.env`.
   - Update the credentials in the `.env` file, such as your database connection information.

2. **Run migrations:**
   - To migrate the tables, run the following command:
     ```bash
     make migrate
     ```

3. **Run the application:**
   - To run the application, execute:
     ```bash
     make run
     ```

4. **Build the project:**
   - To build the project, use:
     ```bash
     make build
     ```

5. **Run tests:**
   - To execute the tests, run:
     ```bash
     make test
     ```

6. **Run linting:**
   - To run linting, use:
     ```bash
     make lint
     ```
   - Ensure that `golangci-lint` is installed on your system before running this command.

## Project Structure

Here's an overview of the folder structure and key components of the project:

```
cmd/
    server/
        main.go          : Main program that starts the application
    migration/
        main.go          : Program for running database migrations
config/
    config.go           : Loads environment variables
database/
    database.go         : Database configuration to connect and close the connection
internal/
    handlers/
        auth.go          : REST API handler for authentication
        auth_test.go     : Test cases for authentication handler
    models/
        user.go          : User model definition
    repository/
        user.go          : Logic for saving user data into the database
    services/
        auth.go          : Authentication logic (e.g., password hashing, token generation)
    utils/
        utils.go         : Helper functions (e.g., validation, formatting)
```