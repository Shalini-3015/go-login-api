package controller

import (
	"net/http"
	"strconv"

	"go-login-api-task/service"

	"github.com/gin-gonic/gin"
)

type ConversionController struct {
	service *service.ConversionService
}

func NewConversionController(service *service.ConversionService) *ConversionController {
	return &ConversionController{service: service}
}

func (c *ConversionController) ConvertCurrency(ctx *gin.Context) {
	from := ctx.Query("from")
	to := ctx.Query("to")
	amountStr := ctx.Query("amount")

	if from == "" || to == "" || amountStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "from, to and amount are required"})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	result, err := c.service.ConvertCurrencyAmt(from, to, amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
