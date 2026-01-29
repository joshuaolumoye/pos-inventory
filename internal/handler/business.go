package handler

import (
	"encoding/json"
	"net/http"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/internal/usecase"
	"github.com/joshuaolumoye/pos-backend/pkg/dto"
)

// These should be set by main.go DI
var BusinessUC *usecase.BusinessUsecase
var AuthRepo domain.AuthRepository

func RegisterBusinessHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	b := &domain.Business{
		Name:             req.BusinessName,
		OwnerFullName:    req.OwnerFullName,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		StoreAddress:     req.StoreAddress,
		BusinessCategory: req.BusinessCategory,
		Currency:         req.Currency,
		StoreIcon:        req.StoreIcon,
	}
	err := BusinessUC.RegisterBusiness(b, req.Password)
	if err != nil {
		msg := err.Error()
		if msg == "email already registered" ||
			msg == "Error 1062 (23000): Duplicate entry '"+b.Email+"' for key 'email'" ||
			(len(msg) > 0 && (msg == "duplicate key value violates unique constraint \"email\"" || msg == "UNIQUE constraint failed: businesses.email")) {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}
		http.Error(w, "Registration failed: "+msg, http.StatusBadRequest)
		return
	}
	resp := dto.RegisterBusinessResponse{
		BusinessID: b.ID,
		Message:    "Business registered successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	access, refresh, businessID, role, err := BusinessUC.Login(req.Email, req.Password, AuthRepo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	resp := dto.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		BusinessID:   businessID,
		Role:         role,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
