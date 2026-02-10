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
goose postgres "postgres://carlosinfante:@localhost:5432/techshop" up
```

5. Create querie to add values to employees table (users.sql)

6. Generate code. Run in the terminal from server as root directory:
sqlc generate
7. Create global struct to all variable of the configuration (apiConfig)
8. Load environment variables from .env file
9. Get database URL and port from environment variables
10. Connect to the database
11. Verify the database connection


## 4. Create user

1. Create functions to handle json responses (json.go)
2. <internal/auth.go>. Create function to hash password (internal/auth/auth.go). We use argon21 pacjage which is needed to be installed:
go get github.com/alexedwards/argon2id
3. <internal/auth.go>. Create a function to compare passwords (CheckPasswordHash)
4. Create the handler method to create an employee (handlerEmployeesCreate)
4.1 Define the expected parameters and response structure
4.2 Define the response structure
4.3 Decode the JSON request body into the parameters struct
4.4 Hash the password before storing it in the database
5. <main.go>. Endpoint to create employee:
mux.HandleFunc("POST /api/employees", apiCfg.handlerEmployeesCreate)


## 5. Login

1. Create schema to save refresh token (sq/schema/002+refresh_tokens.sql)
2. run in the terminal the command to build the table from server as root directory:

```bash
cd sql/schema
goose postgres "postgres://carlosinfante:@localhost:5432/techshop" up
```

3. Create querie to input data to employees table
2. <internal/auth.go>. Create function to make JWT.
I need to install jwt package:
go get -u github.com/golang-jwt/jwt/v5
5. <internal/auth.go>. Function to validate JWT
4. <internal/auth.go>. Function to create refresh token
5.  Method to handle login (handlerLogin)
7.1 Define the expected parameters and response structure
7.2 Define the response structure
7.3 Decode the JSON request body into the parameters struct
7.4 Retrieve the employee from the database using the provided email
7.5 Check if the provided password matches the stored hash
7.6 Create a JWT token for the authenticated employee
7.7 Create a refresh token and store it in the database
7.8 Save the refresh token in the database with an expiration time of 60 days
7.9 Respond with the employee's details and the generated tokens (excluding the password)
6. <main.go>. Endpoint to login:
 mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

## 6. Create employee restrition. Only admin

1. <sql/queries.sql>. Create a querie to get employee by id.
2. Generate the code. command : sqlc generate
3. <handler_employee_create.go> get token
4. get admin Id
5. check if the Id math a admin
6. check if the user has "Admin" role

## 7. Update emplyoee data. Only by the employee self

1. <sql/queries/employees.sql>. Querie to update employee
2. <hable_employee_update.go> Method to update employee
2.1 parameter struct
2.2 struct to response
2.3 Get token
2.4 get user id
2.5 Decode json to go
2.6 hash password
2.7 update employee and get employee for the response
2.8 response with the new data of the employee

3. <main.go>. Endpoint to update employee:
 mux.HandleFunc("PUT /api/employees", apiCfg.handlerEmployeesUpdate)
