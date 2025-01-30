package adapter

import (
	"encoding/json"
	"fmt"
	"go-banking/domain"
	"go-banking/dto"
	"go-banking/service"
	"go-banking/utils"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AccountHandlerDB struct {
	Service   service.AccountService
	Validator validator.Validate
}

func NewAccountHandlerDB(service service.AccountService) *AccountHandlerDB {
	return &AccountHandlerDB{Service: service, Validator: *validator.New()}
}

func (h *AccountHandlerDB) GetAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Info().Msg("Method not allowed")
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Get all accounts")

	var resp []dto.AccountWithCustomer
	log.Info().Msg("Getting all accounts")

	accounts, err := h.Service.GetAccounts()
	if err != nil && !strings.Contains(err.Error(), "no accounts found") {
		utils.ErrorResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	for _, account := range accounts {
		resp = append(resp, dto.AccountWithCustomer{
			ID:            account.ID,
			Customer_Name: account.Customer_Name,
			Username:      account.Username,
			Balance:       account.Balance,
			Currency:      account.Currency,
			Status:        account.Status,
		})
	}

	log.Info().Int("total", len(accounts)).Msg("Accounts fetched successfully")
	utils.ResponseJSON(w, resp, http.StatusOK, "success", "Accounts fetched successfully")
}

func (h *AccountHandlerDB) CreateAccount(w http.ResponseWriter, r *http.Request) {
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

	newAccount, err := h.Service.CreateAccount(account)
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

func (h *AccountHandlerDB) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	log.Info().Msg("Getting account by ID")
	account, err := h.Service.GetAccountByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no accounts found") {
			utils.ErrorResponse(w, http.StatusNotFound, "error", "Account not found for the given ID")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "error", "Database error: "+err.Error())
		}
		return
	}

	log.Info().Str("account", fmt.Sprintf("%+v", account)).Msg("Account retrieved successfully")
	utils.ResponseJSON(w, account, http.StatusOK, "success", "Account retrieved successfully")
}

func (h *AccountHandlerDB) GetAccountByCustomerID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	log.Info().Str("account_id", id).Msg("Getting Account by Customer ID")
	customer, account, err := h.Service.GetAccountByCustomerID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no customers found") {
			utils.ErrorResponse(w, http.StatusNotFound, "error", "no customers found")
		} else if strings.Contains(err.Error(), "no accounts found") {
			utils.ErrorResponse(w, http.StatusNotFound, "error", "Account not found for the given customer ID")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "error", "Database error: "+err.Error())
		}
		return
	}

	response := dto.AccountByCustomerIDResponse{
		CustomerData: customer,
		AccountData:  account,
	}

	log.Info().Str("account_id", id).Msg("Account retrieved successfully")
	utils.ResponseJSON(w, response, http.StatusOK, "Success", "Success get account by customer id")
}

func (h *AccountHandlerDB) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Update a new customer")

	id := mux.Vars(r)["id"]

	var req dto.AccountRequest[domain.Account]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	updatedAccount, err := h.Service.UpdateAccount(id, account)
	if err != nil {
		if strings.Contains(err.Error(), "no customers found") {
			utils.ErrorResponse(w, http.StatusUnprocessableEntity, "error", "invalid customer ID")
		} else if strings.Contains(err.Error(), "no accounts found") {
			utils.ErrorResponse(w, http.StatusNotFound, "error", "Account not found for the given ID")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "error", "Database error: "+err.Error())
		}
		return
	}

	resp := dto.AccountResponse[domain.Account]{
		ID:          updatedAccount.ID,
		Customer_ID: updatedAccount.Customer_ID,
		Username:    updatedAccount.Username,
		Balance:     updatedAccount.Balance,
		Currency:    updatedAccount.Currency,
		Status:      updatedAccount.Status,
	}

	log.Info().Msg("Account updated successfully")
	utils.ResponseJSON(w, resp, http.StatusCreated, "success", "Account updated successfully")
}

func (h *AccountHandlerDB) SoftDeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Delete an account")

	id := mux.Vars(r)["id"]

	account, err := h.Service.SoftDeleteAccount(id)
	if err != nil {
		if strings.Contains(err.Error(), "no accounts found") {
			utils.ErrorResponse(w, http.StatusNotFound, "error", "Account not found for the given ID")
		} else if strings.Contains(err.Error(), "account already deleted") {
			utils.ErrorResponse(w, http.StatusOK, "error", err.Error())
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "error", "Database error: "+err.Error())
		}
		return
	}

	resp := dto.AccountResponse[*domain.Account]{
		ID:          account.ID,
		Customer_ID: account.Customer_ID,
		Username:    account.Username,
		Balance:     account.Balance,
		Currency:    account.Currency,
		Status:      account.Status,
	}

	log.Info().Msg("Account deleted successfully")
	utils.ResponseJSON(w, resp, http.StatusCreated, "success", "Account deleted successfully")
}
