package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/carloscfgos1980/tech-shop-api/internal/auth"
	"github.com/carloscfgos1980/tech-shop-api/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// Define the expected parameters and response structure
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	// Define the response structure
	type response struct {
		Employee     Employee `json:"employee"`
		Token        string   `json:"token"`
		RefreshToken string   `json:"refresh_token"`
	}
	// Decode the JSON request body into the parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// Retrieve the employee from the database using the provided email
	employee, err := cfg.db.GetEmployeeByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// Check if the provided password matches the stored hash
	match, err := auth.CheckPasswordHash(params.Password, employee.Password)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// Create a JWT token for the authenticated employee
	token, err := auth.MakeJWT(
		employee.ID,
		cfg.jwtSecret,
		24*7*time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
		return
	}
	// Create a refresh token and store it in the database
	refreshToken := auth.MakeRefreshToken()

	// Save the refresh token in the database with an expiration time of 60 days
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		EmployeeID: employee.ID,
		Token:      refreshToken,
		ExpiresAt:  time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}
	// Respond with the employee's details and the generated tokens (excluding the password)
	respondWithJSON(w, http.StatusOK, response{
		Employee: Employee{
			ID:        employee.ID,
			Email:     employee.Email,
			CreatedAt: employee.CreatedAt,
			UpdatedAt: employee.UpdatedAt,
			Role:      employee.Role,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
