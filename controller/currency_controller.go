package controller

import (
	"go-login-api-task/models"
	"go-login-api-task/service"
	"log"
	"net/http"
	"strconv"
"go-login-api-task/dto/currency"
	"github.com/gin-gonic/gin"
)

type CurrencyController struct {
	service *service.CurrencyService
}

func NewCurrencyController(service *service.CurrencyService) *CurrencyController {
	return &CurrencyController{service: service}
}

func (c *CurrencyController) CreateCurrency(ctx *gin.Context) {
	var currency models.Currency

	if err := ctx.ShouldBindJSON(&currency); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateCurrency(&currency); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, currency)
}

func (c *CurrencyController) GetAllCurrencies(ctx *gin.Context) {
	currencies, err := c.service.GetAllActiveCurrencies()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, currencies)
}

func (c *CurrencyController) GetCurrencyByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	currency, err := c.service.GetCurrencyByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, currency)
}

func (c *CurrencyController) UpdateCurrency(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.CurrencyUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Name == nil && req.Symbol == nil && req.IsActive == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided"})
		return
	}

	err = c.service.UpdateCurrency(ctx.Request.Context(), uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "currency updated successfully"})
}


func (c *CurrencyController) DeleteCurrency(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.service.DeactivateCurrency(ctx.Request.Context(), uint(id)); err != nil {
		log.Println(">>> DELETE /currencies hit")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "currency deactivated successfully"})
}
