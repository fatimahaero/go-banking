package routes

import (
	"fmt"
	hand "go-banking/adapter/handler"
	repo "go-banking/adapter/repository"
	conf "go-banking/config"
	"go-banking/domain"
	"go-banking/middleware"
	serv "go-banking/service"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewRouter(router *mux.Router, db *sqlx.DB) {
	// apply middleware to all routes
	router.Use(middleware.ApiKeyMiddleware)

	customerRepo := repo.NewCustomerRepositoryDB(db)
	CustomerService := serv.NewCustomerService(customerRepo)
	customerHandler := hand.NewCustomerHandlerDB(CustomerService)

	router.Handle("/customers", middleware.AuthMiddleware(http.HandlerFunc(customerHandler.GetCustomers))).Methods("GET")
	router.Handle("/customers/add", http.HandlerFunc(customerHandler.CreateCustomer)).Methods("POST")
	router.Handle("/customers/{id}", middleware.AuthMiddleware(http.HandlerFunc(customerHandler.GetCustomerByID))).Methods("GET")
	router.Handle("/customers/{id}/edit", middleware.AuthMiddleware(http.HandlerFunc(customerHandler.UpdateCustomer))).Methods("PUT")

	accountRepo := repo.NewAccountRepositoryDB(db)
	accountService := serv.NewAccountService(accountRepo, customerRepo)
	accountHandler := hand.NewAccountHandlerDB(accountService)

	router.Handle("/accounts", middleware.AuthMiddleware(http.HandlerFunc(accountHandler.GetAccounts))).Methods("GET")
	router.Handle("/accounts/add", middleware.AuthMiddleware(http.HandlerFunc(accountHandler.CreateAccount))).Methods("POST")
	router.Handle("/accounts/{id}", middleware.AuthMiddleware(http.HandlerFunc(accountHandler.GetAccountByID))).Methods("GET")
	router.Handle("/accounts/customer/{id}", middleware.AuthMiddleware(http.HandlerFunc(accountHandler.GetAccountByCustomerID))).Methods("GET")
	router.Handle("/accounts/{id}/edit", middleware.AuthMiddleware(http.HandlerFunc(accountHandler.UpdateAccount))).Methods("PUT")
	router.Handle("/accounts/{id}/delete", middleware.AuthMiddleware(http.HandlerFunc(accountHandler.SoftDeleteAccount))).Methods("PUT")

	transactionRepo := repo.NewTransactionRepositoryDB(db)
	transactionService := serv.NewTransactionService(transactionRepo, accountRepo)
	transactionHandler := hand.NewTransactionHandlerDB(transactionService)

	router.Handle("/transactions/add", middleware.AuthMiddleware(http.HandlerFunc(transactionHandler.CreateTransaction))).Methods("POST")
	router.Handle("/transactions", middleware.AuthMiddleware(http.HandlerFunc(transactionHandler.GetAllTransaction))).Methods("GET")
	router.Handle("/transactions/account/{id}", middleware.AuthMiddleware(http.HandlerFunc(transactionHandler.GetTransactionByAccountID))).Methods("GET")

	authService := serv.NewAuthService(accountRepo)
	authHandler := hand.NewAuthHandlerDB(authService, accountService)

	router.Handle("/login", http.HandlerFunc(authHandler.Login)).Methods("POST")
	router.Handle("/register", http.HandlerFunc(authHandler.Register)).Methods("POST")
	/*
	 * Datanya diambil dari mock data
	 */
	repoCustMock := repo.NewCustomerRepositoryMock()
	svcCustMock := serv.NewCustomerService(repoCustMock)
	handCustMock := hand.NewCustomerHandler(svcCustMock)

	router.Handle("/mock/customers", middleware.AuthMiddleware(http.HandlerFunc(handCustMock.GetCustomers))).Methods("GET")
	router.Handle("/mock/customers/add", middleware.AuthMiddleware(http.HandlerFunc(handCustMock.AddCustomer))).Methods("POST")
}

/*
 * implementasi routing dari Kang Ari
 */
func StartServer() {

	// Start of log setup
	conf.InitiateLog()
	defer conf.CloseLog() // Close log when application is stopped
	// End of log setup

	config, _ := domain.GetConfig()
	port := config.Server.Port

	db, _ := conf.NewDBConnectionENV()

	defer db.Close()

	router := mux.NewRouter()

	NewRouter(router, db)

	fmt.Println("starting server on port " + port)

	http.ListenAndServe(":"+port, router)
}
