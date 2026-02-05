package router

import (
	"go-login-api-task/controller"
	"go-login-api-task/middleware"
	"go-login-api-task/repository"
	"go-login-api-task/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	
	currencyRepo := repository.NewCurrencyRepository()
	exchangeRateRepo := repository.NewExchangeRateRepository()

	
	currencyService := service.NewCurrencyService(currencyRepo)
	exchangeRateService := service.NewExchangeRateService(exchangeRateRepo, currencyRepo)
	conversionService := service.NewConversionService(currencyRepo, exchangeRateRepo)

	
	authController := controller.NewAuthController()
	currencyController := controller.NewCurrencyController(currencyService)
	exchangeRateController := controller.NewExchangeRateController(exchangeRateService)
	conversionController := controller.NewConversionController(conversionService)

	
	r.POST("/login", authController.UserLogin)
	r.POST("/register", authController.RegisterUser)

	
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/currencies", currencyController.CreateCurrency)
	protected.GET("/currencies", currencyController.GetAllCurrencies)
	protected.GET("/currencies/:id", currencyController.GetCurrencyByID)
	protected.PUT("/currencies/:id", currencyController.UpdateCurrency)
	protected.DELETE("/currencies/:id", currencyController.DeleteCurrency)

	protected.POST("/exchange-rates", exchangeRateController.CreateExchangeRate)
	protected.GET("/exchange-rates", exchangeRateController.GetAllExchangeRate)
	protected.GET("/exchange-rates/:id", exchangeRateController.GetExchangeRateByID)
	protected.PUT("/exchange-rates/:id", exchangeRateController.UpdateExchangeRate)
	protected.DELETE("/exchange-rates/:id", exchangeRateController.DeleteExchangeRate)

	protected.GET("/convert", conversionController.ConvertCurrency)
}
