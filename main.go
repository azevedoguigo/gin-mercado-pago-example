package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

const ACCESS_TOKEN string = "TEST-4452616593758483-070808-16b13697005f228aa2b359b6e84727a2-1534160175"

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

	router.POST("/pay", func(ctx *gin.Context) {
		config, err := config.New(ACCESS_TOKEN)
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
