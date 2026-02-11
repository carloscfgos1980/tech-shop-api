package main

import (
	"log"
	"net/http"

	"github.com/carloscfgos1980/tech-shop-api/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerEmployeesDelete(w http.ResponseWriter, r *http.Request) {
	// get the employee id and parse it
	employeeIdString := r.PathValue("employeeId")
	employeeId, err := uuid.Parse(employeeIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID", err)
		return
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
	// Delete employee from database
	err = cfg.db.DeleteEmployee(r.Context(), employeeId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete employee", err)
		return
	}
	// response with 204 and not body
	w.WriteHeader(http.StatusNoContent)
}
