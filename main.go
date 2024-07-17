package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

type PaymentBody struct {
	Token             string  `json:"token"`
	IssuerId          string  `json:"issuer_id"`
	PaymentMethodId   string  `json:"payment_method_id"`
	TransactionAmount float64 `json:"transaction_amount"`
	Installments      uint64  `json:"installments"`
	Payer             struct {
		Email          string `json:"email"`
		Identification struct {
			Type   string `json:"type"`
			Number string `json:"number"`
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file!")
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	var accessToken string = os.Getenv("MERCADO_PAGO_ACCESS_TOKEN")

	router.POST("/pay", func(ctx *gin.Context) {
		var reqBody PaymentBody
		if err := ctx.ShouldBindBodyWithJSON(&reqBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

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
		}

		ctx.JSON(http.StatusOK, resource)
	})

	router.Run()
}
