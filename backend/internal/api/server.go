package api

import (
	"database/sql"
	"flow-sight-backend/internal/config"
	"flow-sight-backend/internal/handlers"
	"flow-sight-backend/internal/logger"
	"flow-sight-backend/internal/middleware"
	"flow-sight-backend/internal/repositories"
	"flow-sight-backend/internal/services"
	"flow-sight-backend/internal/version"

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
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
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
	userRepo := repositories.NewUserRepository(s.db)
	creditCardRepo := repositories.NewCreditCardRepository(s.db)
	bankAccountRepo := repositories.NewBankAccountRepository(s.db)
	incomeSourceRepo := repositories.NewIncomeSourceRepository(s.db)
	monthlyIncomeRepo := repositories.NewMonthlyIncomeRepository(s.db)
	recurringPaymentRepo := repositories.NewRecurringPaymentRepository(s.db)
	cardMonthlyTotalRepo := repositories.NewCardMonthlyTotalRepository(s.db)
	appSettingRepo := repositories.NewAppSettingRepository(s.db)

	// Initialize services
	authService := services.NewAuthService(userRepo, s.config)
	creditCardService := services.NewCreditCardService(creditCardRepo)
	bankAccountService := services.NewBankAccountService(bankAccountRepo)
	incomeService := services.NewIncomeService(incomeSourceRepo, monthlyIncomeRepo)
	recurringPaymentService := services.NewRecurringPaymentService(recurringPaymentRepo)
	cardMonthlyTotalService := services.NewCardMonthlyTotalService(cardMonthlyTotalRepo)
	appSettingService := services.NewAppSettingService(appSettingRepo)
	cashflowService := services.NewCashflowService(bankAccountRepo, incomeSourceRepo, monthlyIncomeRepo, recurringPaymentRepo, cardMonthlyTotalRepo, creditCardRepo, appSettingRepo)
	dashboardService := services.NewDashboardService(bankAccountRepo, creditCardRepo, incomeSourceRepo, monthlyIncomeRepo, recurringPaymentRepo, cashflowService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, s.config)
	creditCardHandler := handlers.NewCreditCardHandler(creditCardService)
	bankAccountHandler := handlers.NewBankAccountHandler(bankAccountService)
	incomeHandler := handlers.NewIncomeHandler(incomeService)
	recurringPaymentHandler := handlers.NewRecurringPaymentHandler(recurringPaymentService)
	cardMonthlyTotalHandler := handlers.NewCardMonthlyTotalHandler(cardMonthlyTotalService)
	appSettingHandler := handlers.NewAppSettingHandler(appSettingService)
	cashflowHandler := handlers.NewCashflowHandler(cashflowService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Public routes (no authentication required)
	api := s.router.Group("/api/v1")

	// Auth routes
	api.GET("/auth/google", authHandler.GoogleLogin)
	api.GET("/auth/google/callback", authHandler.GoogleCallback)

	// Protected routes (authentication required)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService))

	// User info
	protected.GET("/auth/me", authHandler.GetMe)

	// Credit Card routes
	protected.GET("/credit-cards", creditCardHandler.GetCreditCards)
	protected.POST("/credit-cards", creditCardHandler.CreateCreditCard)
	protected.GET("/credit-cards/:id", creditCardHandler.GetCreditCard)
	protected.PUT("/credit-cards/:id", creditCardHandler.UpdateCreditCard)
	protected.DELETE("/credit-cards/:id", creditCardHandler.DeleteCreditCard)

	// Bank Account routes
	protected.GET("/bank-accounts", bankAccountHandler.GetBankAccounts)
	protected.POST("/bank-accounts", bankAccountHandler.CreateBankAccount)
	protected.GET("/bank-accounts/:id", bankAccountHandler.GetBankAccount)
	protected.PUT("/bank-accounts/:id", bankAccountHandler.UpdateBankAccount)
	protected.DELETE("/bank-accounts/:id", bankAccountHandler.DeleteBankAccount)

	// Income routes
	protected.GET("/income-sources", incomeHandler.GetIncomeSources)
	protected.POST("/income-sources", incomeHandler.CreateIncomeSource)
	protected.GET("/income-sources/:id", incomeHandler.GetIncomeSource)
	protected.PUT("/income-sources/:id", incomeHandler.UpdateIncomeSource)
	protected.DELETE("/income-sources/:id", incomeHandler.DeleteIncomeSource)

	protected.GET("/monthly-income-records", incomeHandler.GetMonthlyIncomeRecords)
	protected.POST("/monthly-income-records", incomeHandler.CreateMonthlyIncomeRecord)
	protected.GET("/monthly-income-records/:id", incomeHandler.GetMonthlyIncomeRecord)
	protected.PUT("/monthly-income-records/:id", incomeHandler.UpdateMonthlyIncomeRecord)
	protected.DELETE("/monthly-income-records/:id", incomeHandler.DeleteMonthlyIncomeRecord)

	// Recurring Payment routes
	protected.GET("/recurring-payments", recurringPaymentHandler.GetRecurringPayments)
	protected.POST("/recurring-payments", recurringPaymentHandler.CreateRecurringPayment)
	protected.GET("/recurring-payments/:id", recurringPaymentHandler.GetRecurringPayment)
	protected.PUT("/recurring-payments/:id", recurringPaymentHandler.UpdateRecurringPayment)
	protected.DELETE("/recurring-payments/:id", recurringPaymentHandler.DeleteRecurringPayment)

	// Card Monthly Total routes
	protected.GET("/card-monthly-totals", cardMonthlyTotalHandler.GetCardMonthlyTotals)
	protected.POST("/card-monthly-totals", cardMonthlyTotalHandler.CreateCardMonthlyTotal)
	protected.GET("/card-monthly-totals/:id", cardMonthlyTotalHandler.GetCardMonthlyTotal)
	protected.PUT("/card-monthly-totals/:id", cardMonthlyTotalHandler.UpdateCardMonthlyTotal)
	protected.DELETE("/card-monthly-totals/:id", cardMonthlyTotalHandler.DeleteCardMonthlyTotal)

	// App Setting routes
	protected.GET("/settings", appSettingHandler.GetSettings)
	protected.PUT("/settings", appSettingHandler.UpdateSettings)

	// Cashflow Projection routes
	protected.GET("/cashflow-projection", cashflowHandler.GetCashflowProjection)

	// Dashboard routes
	protected.GET("/dashboard/summary", dashboardHandler.GetDashboardSummary)

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
