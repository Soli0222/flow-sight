package api

import (
	"database/sql"

	"github.com/Soli0222/flow-sight/backend/internal/config"
	"github.com/Soli0222/flow-sight/backend/internal/handlers"
	"github.com/Soli0222/flow-sight/backend/internal/logger"
	"github.com/Soli0222/flow-sight/backend/internal/middleware"
	"github.com/Soli0222/flow-sight/backend/internal/repositories"
	"github.com/Soli0222/flow-sight/backend/internal/services"
	"github.com/Soli0222/flow-sight/backend/internal/version"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
	db     *sql.DB
	config *config.Config
	logger *logger.Logger
}

func NewServer(db *sql.DB, cfg *config.Config, appLogger *logger.Logger) *Server {
	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New() // Use gin.New() instead of gin.Default() to avoid default middleware

	// Add custom request logging middleware
	router.Use(middleware.RequestLogger(appLogger))

	// Add CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.Host}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(corsConfig))

	server := &Server{
		router: router,
		db:     db,
		config: cfg,
		logger: appLogger,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Initialize repositories
	creditCardRepo := repositories.NewCreditCardRepository(s.db)
	bankAccountRepo := repositories.NewBankAccountRepository(s.db)
	incomeSourceRepo := repositories.NewIncomeSourceRepository(s.db)
	monthlyIncomeRepo := repositories.NewMonthlyIncomeRepository(s.db)
	recurringPaymentRepo := repositories.NewRecurringPaymentRepository(s.db)
	cardMonthlyTotalRepo := repositories.NewCardMonthlyTotalRepository(s.db)
	appSettingRepo := repositories.NewAppSettingRepository(s.db)

	// Initialize services
	creditCardService := services.NewCreditCardService(creditCardRepo)
	bankAccountService := services.NewBankAccountService(bankAccountRepo)
	incomeService := services.NewIncomeService(incomeSourceRepo, monthlyIncomeRepo)
	recurringPaymentService := services.NewRecurringPaymentService(recurringPaymentRepo)
	cardMonthlyTotalService := services.NewCardMonthlyTotalService(cardMonthlyTotalRepo)
	appSettingService := services.NewAppSettingService(appSettingRepo)
	cashflowService := services.NewCashflowService(bankAccountRepo, incomeSourceRepo, monthlyIncomeRepo, recurringPaymentRepo, cardMonthlyTotalRepo, creditCardRepo, appSettingRepo)
	dashboardService := services.NewDashboardService(bankAccountRepo, creditCardRepo, incomeSourceRepo, monthlyIncomeRepo, recurringPaymentRepo, cashflowService)

	// Initialize handlers
	creditCardHandler := handlers.NewCreditCardHandler(creditCardService)
	bankAccountHandler := handlers.NewBankAccountHandler(bankAccountService)
	incomeHandler := handlers.NewIncomeHandler(incomeService)
	recurringPaymentHandler := handlers.NewRecurringPaymentHandler(recurringPaymentService)
	cardMonthlyTotalHandler := handlers.NewCardMonthlyTotalHandler(cardMonthlyTotalService)
	appSettingHandler := handlers.NewAppSettingHandler(appSettingService)
	cashflowHandler := handlers.NewCashflowHandler(cashflowService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Public routes (single-user mode, no authentication)
	api := s.router.Group("/api/v1")
	api.Use(middleware.SingleUserMiddleware())

	// Credit Card routes
	api.GET("/credit-cards", creditCardHandler.GetCreditCards)
	api.POST("/credit-cards", creditCardHandler.CreateCreditCard)
	api.GET("/credit-cards/:id", creditCardHandler.GetCreditCard)
	api.PUT("/credit-cards/:id", creditCardHandler.UpdateCreditCard)
	api.DELETE("/credit-cards/:id", creditCardHandler.DeleteCreditCard)

	// Bank Account routes
	api.GET("/bank-accounts", bankAccountHandler.GetBankAccounts)
	api.POST("/bank-accounts", bankAccountHandler.CreateBankAccount)
	api.GET("/bank-accounts/:id", bankAccountHandler.GetBankAccount)
	api.PUT("/bank-accounts/:id", bankAccountHandler.UpdateBankAccount)
	api.DELETE("/bank-accounts/:id", bankAccountHandler.DeleteBankAccount)

	// Income routes
	api.GET("/income-sources", incomeHandler.GetIncomeSources)
	api.POST("/income-sources", incomeHandler.CreateIncomeSource)
	api.GET("/income-sources/:id", incomeHandler.GetIncomeSource)
	api.PUT("/income-sources/:id", incomeHandler.UpdateIncomeSource)
	api.DELETE("/income-sources/:id", incomeHandler.DeleteIncomeSource)

	api.GET("/monthly-income-records", incomeHandler.GetMonthlyIncomeRecords)
	api.POST("/monthly-income-records", incomeHandler.CreateMonthlyIncomeRecord)
	api.GET("/monthly-income-records/:id", incomeHandler.GetMonthlyIncomeRecord)
	api.PUT("/monthly-income-records/:id", incomeHandler.UpdateMonthlyIncomeRecord)
	api.DELETE("/monthly-income-records/:id", incomeHandler.DeleteMonthlyIncomeRecord)

	// Recurring Payment routes
	api.GET("/recurring-payments", recurringPaymentHandler.GetRecurringPayments)
	api.POST("/recurring-payments", recurringPaymentHandler.CreateRecurringPayment)
	api.GET("/recurring-payments/:id", recurringPaymentHandler.GetRecurringPayment)
	api.PUT("/recurring-payments/:id", recurringPaymentHandler.UpdateRecurringPayment)
	api.DELETE("/recurring-payments/:id", recurringPaymentHandler.DeleteRecurringPayment)

	// Card Monthly Total routes
	api.GET("/card-monthly-totals", cardMonthlyTotalHandler.GetCardMonthlyTotals)
	api.POST("/card-monthly-totals", cardMonthlyTotalHandler.CreateCardMonthlyTotal)
	api.GET("/card-monthly-totals/:id", cardMonthlyTotalHandler.GetCardMonthlyTotal)
	api.PUT("/card-monthly-totals/:id", cardMonthlyTotalHandler.UpdateCardMonthlyTotal)
	api.DELETE("/card-monthly-totals/:id", cardMonthlyTotalHandler.DeleteCardMonthlyTotal)

	// App Setting routes
	api.GET("/settings", appSettingHandler.GetSettings)
	api.PUT("/settings", appSettingHandler.UpdateSettings)

	// Cashflow Projection routes
	api.GET("/cashflow-projection", cashflowHandler.GetCashflowProjection)

	// Dashboard routes
	api.GET("/dashboard/summary", dashboardHandler.GetDashboardSummary)

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Version info
	api.GET("/version", func(c *gin.Context) {
		buildInfo := version.GetBuildInfo()
		c.JSON(200, buildInfo)
	})

	// Swagger documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
