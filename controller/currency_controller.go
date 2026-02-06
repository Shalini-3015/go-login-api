package controller

import (
	"net/http"
	"strconv"

	"go-login-api-task/models"
	"go-login-api-task/service"

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
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var updateCurrReq struct {
		Name     *string `json:"name"`
		Symbol   *string `json:"symbol"`
		IsActive bool   `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&updateCurrReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateCurrency(uint(id), updateCurrReq.Name, updateCurrReq.Symbol, &updateCurrReq.IsActive); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if err := c.service.DeactivateCurrency(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "currency deactivated successfully"})
}
