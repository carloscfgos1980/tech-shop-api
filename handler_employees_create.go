package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/carloscfgos1980/tech-shop-api/internal/auth"
	"github.com/carloscfgos1980/tech-shop-api/internal/database"

	"github.com/google/uuid"
)

type Employee struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}

func (cfg *apiConfig) handlerEmployeesCreate(w http.ResponseWriter, r *http.Request) {
	// Define the expected parameters and response structure
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	// Define the response structure
	type response struct {
		Employee Employee `json:"employee"`
	}

	// get token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	// get admin Id
	adminId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// check if the Id math a admin
	admin, err := cfg.db.GetAdminById(r.Context(), adminId)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "you don't have authirzation to create employee", err)
		return
	}

	log.Printf("staffer email: %s", admin.Email)
	// check if the user has "Admin" role
	if admin.Role != "admin" {
		respondWithError(w, http.StatusUnauthorized, "you don't have authirzation to create employee", err)
		return

	}

	// Decode the JSON request body into the parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Hash the password before storing it in the database
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	// Create the employee in the database and retrieve the created employee's details
	employee, err := cfg.db.CreateEmployee(r.Context(), database.CreateEmployeeParams{
		Email:    params.Email,
		Password: hashedPassword,
		Role:     params.Role,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create employee", err)
		return
	}
	// Respond with the created employee's details (excluding the password)
	respondWithJSON(w, http.StatusCreated, response{
		Employee: Employee{
			ID:        employee.ID,
			CreatedAt: employee.CreatedAt,
			UpdatedAt: employee.UpdatedAt,
			Email:     employee.Email,
			Role:      employee.Role,
		},
	})
}
