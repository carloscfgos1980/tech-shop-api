package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/carloscfgos1980/tech-shop-api/internal/database"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	db        *database.Queries
	port      string
	jwtSecret string
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()
	// Get database URL and port from environment variables
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	// Get the port from environment variables, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}
	// Get the JWT secret from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	// Connect to the database
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer dbConn.Close()

	// Verify the database connection
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		port:      port,
		db:        dbQueries,
		jwtSecret: jwtSecret,
	}

	log.Print(apiCfg.port)

	// Set up HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/employees", apiCfg.handlerEmployeesCreate)
	mux.HandleFunc("PUT /api/employees", apiCfg.handlerEmployeesUpdate)

	// server variable to hold the server instance
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server is running http://localhost:%s", port)

	// Start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
