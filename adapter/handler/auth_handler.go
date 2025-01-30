package adapter

import (
	"encoding/json"
	"go-banking/domain"
	"go-banking/dto"
	"go-banking/service"
	"go-banking/utils"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlerDB struct {
	Service    service.AuthService
	AccService service.AccountService
	Validator  validator.Validate
}

func NewAuthHandlerDB(service service.AuthService, AccService service.AccountService) *AuthHandlerDB {
	return &AuthHandlerDB{Service: service, AccService: AccService, Validator: *validator.New()}
}

func (h *AuthHandlerDB) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Login")

	var req dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		errorMessage := utils.CustomValidationError(err)
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "error", errorMessage)
		return
	}

	token, err := h.Service.LoginAccount(req.Username, req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Username or password is incorrect. Failed to login")
		utils.ErrorResponse(w, http.StatusUnauthorized, "error", err.Error())
		return
	}

	resp := dto.LoginResponse{
		Token: token,
	}

	utils.ResponseJSON(w, resp, http.StatusOK, "success", "Login successful")
	log.Info().Str("username", req.Username).Msg("Login successful")
}

func (h *AuthHandlerDB) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Create a new account")

	var req dto.AccountRequest[domain.Account]
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	if err := h.Validator.Struct(&req); err != nil {
		errorMessage := utils.CustomValidationError(err)
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "error", errorMessage)
		return
	}

	account := domain.Account{
		Customer_ID: req.Customer_ID,
		Username:    req.Username,
		Password:    req.Password,
		Balance:     req.Balance,
		Currency:    req.Currency,
		Status:      req.Status,
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(password)

	newAccount, err := h.AccService.CreateAccount(account)
	if err != nil {
		if strings.Contains(err.Error(), "no customers found") {
			utils.ErrorResponse(w, http.StatusUnprocessableEntity, "error", "invalid customer ID")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "error", "Database error: "+err.Error())
		}
		log.Error().Err(err).Msg("Failed to create account")
		return
	}

	resp := dto.AccountResponse[*domain.Account]{
		ID:          newAccount.ID,
		Customer_ID: newAccount.Customer_ID,
		Username:    newAccount.Username,
		Balance:     newAccount.Balance,
		Currency:    newAccount.Currency,
		Status:      newAccount.Status,
	}

	log.Info().Msg("Account created successfully")
	utils.ResponseJSON(w, resp, http.StatusCreated, "success", "Account created successfully")
}
