package main

import (
	"encoding/json"
	"net/http"

	"github.com/carloscfgos1980/tech-shop-api/internal/auth"
	"github.com/carloscfgos1980/tech-shop-api/internal/database"
)

func (cfg *apiConfig) handlerEmployeesUpdate(w http.ResponseWriter, r *http.Request) {
	// parameter struct
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	// struct to response
	type response struct {
		Employee
	}

	// Get token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	// get user id
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// Decode json to go
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// update employee and get employee for the response
	employee, err := cfg.db.UpdateEmployee(r.Context(), database.UpdateEmployeeParams{
		ID:       userID,
		Email:    params.Email,
		Password: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	// response with the new data of the employee
	respondWithJSON(w, http.StatusOK, response{
		Employee: Employee{
			ID:        employee.ID,
			CreatedAt: employee.CreatedAt,
			UpdatedAt: employee.UpdatedAt,
			Email:     employee.Email,
			Role:      employee.Role,
		},
	})
}
