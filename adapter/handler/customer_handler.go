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
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type CustomerHandlerDB struct {
	Service   service.CustomerService
	Validator validator.Validate
}

func NewCustomerHandlerDB(service service.CustomerService) *CustomerHandlerDB {
	return &CustomerHandlerDB{Service: service, Validator: *validator.New()}
}

func (h *CustomerHandlerDB) GetCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Get all customers")

	var resp []dto.CustomerResponseDTO[domain.Customer]
	log.Info().Msg("Retrieving all customers")

	customers, err := h.Service.GetAllCustomers()
	if err != nil && !strings.Contains(err.Error(), "no customers found") {
		log.Error().Err(err).Msg("Failed to retrieve customers")
		utils.ErrorResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	for _, data := range customers {
		resp = append(resp, dto.CustomerResponseDTO[domain.Customer]{
			ID:      data.ID,
			Name:    data.Name,
			City:    data.City,
			Zipcode: data.Zipcode,
			Status:  data.Status,
		})
	}

	log.Info().Int("total", len(customers)).Msg("Customers fetched successfully")
	utils.ResponseJSON(w, resp, http.StatusOK, "success", "Customers fetched successfully")
}

func (h *CustomerHandlerDB) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Create a new customer")

	var req dto.CustomerRequestDTO[domain.Customer]
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

	customer, err := h.Service.CreateCustomer(domain.Customer(req))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create customer")
		utils.ErrorResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	resp := dto.CustomerResponseDTO[domain.Customer]{
		ID:      customer.ID,
		Name:    customer.Name,
		City:    customer.City,
		Zipcode: customer.Zipcode,
		Status:  customer.Status,
	}

	log.Info().Msg("Customer created successfully")
	utils.ResponseJSON(w, resp, http.StatusCreated, "success", "Customer created successfully")
}

func (h *CustomerHandlerDB) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Get customer by ID " + id)

	customer, err := h.Service.GetCustomerByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no customers found") {
			utils.ErrorResponse(w, http.StatusNotFound, "error", "Customer not found for the given ID")
		} else {
			utils.ErrorResponse(w, http.StatusNotFound, "error", err.Error())
		}
		return
	}

	resp := dto.CustomerResponseDTO[domain.Customer]{
		ID:      customer.ID,
		Name:    customer.Name,
		City:    customer.City,
		Zipcode: customer.Zipcode,
		Status:  customer.Status,
	}

	utils.ResponseJSON(w, resp, http.StatusOK, "success", "customer retrieved successfully")
	log.Info().Str("id", id).Interface("customer", customer).Msg("Customer retrieved successfully by ID " + id)
}

func (h *CustomerHandlerDB) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Update a new customer")

	id := mux.Vars(r)["id"]

	var req dto.CustomerRequestDTO[domain.Customer]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	if err := h.Validator.Struct(&req); err != nil {
		errorMessage := utils.CustomValidationError(err)
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "error", errorMessage)
		return
	}

	customer := domain.Customer{
		Name:        req.Name,
		City:        req.City,
		Zipcode:     req.Zipcode,
		DateOfBirth: req.DateOfBirth,
		Status:      req.Status,
	}

	updatedCustomer, err := h.Service.UpdateCustomer(id, customer)
	if err != nil {
		if strings.Contains(err.Error(), "no customers found") {
			utils.ErrorResponse(w, http.StatusNotFound, "error", "no customers found")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "error", "Something went wrong")
		}
		return
	}

	resp := dto.CustomerResponseDTO[domain.Customer]{
		ID:      updatedCustomer.ID,
		Name:    updatedCustomer.Name,
		City:    updatedCustomer.City,
		Zipcode: updatedCustomer.Zipcode,
		Status:  updatedCustomer.Status,
	}

	log.Info().Msg("Customer updated successfully")
	utils.ResponseJSON(w, resp, http.StatusCreated, "success", "Customer updated successfully")
}
