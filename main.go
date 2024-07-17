package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

type PaymentBody struct {
	Token             string
	IssuerId          string
	PaymentMethodId   string
	TransactionAmount float64
	Installments      uint64
	Payer             struct {
		Email          string
		Identification struct {
			Type   string
			Number string
		}
	}
}

func main() {
	router := gin.Default()

	var reqBody PaymentBody

	var accessToken string = os.Getenv("MERCADO_PAGO_ACCESS_TOKEN")

	router.POST("/pay", func(ctx *gin.Context) {
		config, err := config.New(accessToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		client := payment.NewClient(config)

		request := payment.Request{
			TransactionAmount: reqBody.TransactionAmount,
			PaymentMethodID:   reqBody.PaymentMethodId,
			Payer: &payment.PayerRequest{
				Email: reqBody.Payer.Email,
			},
			Token:        reqBody.Token,
			Installments: int(reqBody.Installments),
		}

		resource, err := client.Create(context.Background(), request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, resource)
	})

	router.Run()
}
