package api

import (
	"database/sql"
	"flow-sight-backend/internal/config"
	"flow-sight-backend/internal/handlers"
	"flow-sight-backend/internal/repositories"
	"flow-sight-backend/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
	db     *sql.DB
	config *config.Config
}

func NewServer(db *sql.DB, cfg *config.Config) *Server {
	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	server := &Server{
		router: router,
		db:     db,
		config: cfg,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Initialize repositories
	assetRepo := repositories.NewAssetRepository(s.db)
	bankAccountRepo := repositories.NewBankAccountRepository(s.db)
	incomeSourceRepo := repositories.NewIncomeSourceRepository(s.db)
	monthlyIncomeRepo := repositories.NewMonthlyIncomeRepository(s.db)
	recurringPaymentRepo := repositories.NewRecurringPaymentRepository(s.db)
	cardMonthlyTotalRepo := repositories.NewCardMonthlyTotalRepository(s.db)
	appSettingRepo := repositories.NewAppSettingRepository(s.db)

	// Initialize services
	assetService := services.NewAssetService(assetRepo)
	bankAccountService := services.NewBankAccountService(bankAccountRepo)
	incomeService := services.NewIncomeService(incomeSourceRepo, monthlyIncomeRepo)
	recurringPaymentService := services.NewRecurringPaymentService(recurringPaymentRepo)
	cardMonthlyTotalService := services.NewCardMonthlyTotalService(cardMonthlyTotalRepo)
	appSettingService := services.NewAppSettingService(appSettingRepo)
	cashflowService := services.NewCashflowService(bankAccountRepo, incomeSourceRepo, monthlyIncomeRepo, recurringPaymentRepo, cardMonthlyTotalRepo, assetRepo)

	// Initialize handlers
	assetHandler := handlers.NewAssetHandler(assetService)
	bankAccountHandler := handlers.NewBankAccountHandler(bankAccountService)
	incomeHandler := handlers.NewIncomeHandler(incomeService)
	recurringPaymentHandler := handlers.NewRecurringPaymentHandler(recurringPaymentService)
	cardMonthlyTotalHandler := handlers.NewCardMonthlyTotalHandler(cardMonthlyTotalService)
	appSettingHandler := handlers.NewAppSettingHandler(appSettingService)
	cashflowHandler := handlers.NewCashflowHandler(cashflowService)

	// API routes
	api := s.router.Group("/api/v1")

	// Asset routes
	api.GET("/assets", assetHandler.GetAssets)
	api.POST("/assets", assetHandler.CreateAsset)
	api.GET("/assets/:id", assetHandler.GetAsset)
	api.PUT("/assets/:id", assetHandler.UpdateAsset)
	api.DELETE("/assets/:id", assetHandler.DeleteAsset)

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

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
