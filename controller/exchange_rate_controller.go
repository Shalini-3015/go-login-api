package controller

import (
	
	"go-login-api-task/models"
	"go-login-api-task/service"
	"net/http"
	"strconv"
	"strings"
	"go-login-api-task/dto/exc_rate"
	"github.com/gin-gonic/gin"
)

type ExchangeRateController struct {
	service *service.ExchangeRateService
}

func NewExchangeRateController(service *service.ExchangeRateService) *ExchangeRateController {
	return &ExchangeRateController{service: service}
}

func (c *ExchangeRateController) CreateExchangeRate(ctx *gin.Context) {
	var rate models.ExchangeRate

	if err := ctx.ShouldBindJSON(&rate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateExchangeRate(&rate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, rate)
}

func (c *ExchangeRateController) GetAllExchangeRate(ctx *gin.Context) {
	rates, err := c.service.GetActiveExchangeRates()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rates)
}

func (c *ExchangeRateController) GetExchangeRateByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rate, err := c.service.GetExchangeRateByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rate)
}

func (c *ExchangeRateController) UpdateExchangeRate(ctx *gin.Context) {

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var dto exchange_rate.UpdateExchangeRateDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateExchangeRate(uint(id), dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "exchange rate updated successfully",
	})
}


func (c *ExchangeRateController) DeleteExchangeRate(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.service.DeactivateExchangeRate(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "exchange rate deactivated successfully"})
}

func (c *ExchangeRateController) SyncExchangeRates(ctx *gin.Context) {

	var req struct {
		Base string `json:"base"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	rates, err := c.service.FetchAndSyncRates(req.Base)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		
		"base":  strings.ToUpper(req.Base),
		"rates": rates,
	})
}
