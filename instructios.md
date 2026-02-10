# STEPS

10/02/26

## 1. Setup

- Create the go.mod

go mod init github.com/carloscfgos1980/tech-shop-api

- Install package to handle .env

go get github.com/joho/godotenv

- We need to add and import a Postgres driver so our program knows how to talk to the database. Install it in your module:

go get github.com/lib/pq



## 2. Server

1. Set up HTTP server
2. server variable to hold the server instance
3. Start the server

## 3. Database

1. In the terminal access to postgres and create a database:
psql postgres
CREATE DATABASE techshop;

Note: Check databse exist:
 \c techshop

2. Create file (sqlc.yaml) to handle the creation of queries

3. Create schema for employees table

3. run in the terminal the command to build the table from server as root directory:

```bash
cd sql/schema
psql postgres "postgres://carlosinfante:@localhost:5432/techshop" up
```

5. Create querie to add values to employees table

6. Generate code. Run in the terminal from server as root directory:
sqlc generate
7. Create global struct to all variable of the configuration (apiConfig)
8. Load environment variables from .env file
9. Get database URL and port from environment variables
10. Connect to the database
11. Verify the database connection


## 4. Create user
1. Create functions to handle json responses (json.go)
2. Create function to hash password (internal/auth/auth.go). We use argon21 pacjage which is needed to be installed:
go get github.com/alexedwards/argon2id
3. Create the handler method to create an employee (handlerEmployeesCreate)
3.1 Define the expected parameters and response structure
3.2 Define the response structure
3.3 Decode the JSON request body into the parameters struct
3.4 Hash the password before storing it in the database



