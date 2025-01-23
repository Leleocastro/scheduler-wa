package validatorhdl

import (
	"bytes"
	"complete-api/internal/core/ports"
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

type ValidatorHandler struct {
	paymentSrv ports.PaymentService
}

func NewValidatorHandler(paymentSrv ports.PaymentService) *ValidatorHandler {
	return &ValidatorHandler{
		paymentSrv: paymentSrv,
	}
}

func (h *ValidatorHandler) ValidateSignature(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Erro ao ler o corpo da requisição: %v", err)
		body, err := json.MarshalIndent(c.Request.Body, "", "  ")
		if err != nil {
			log.Printf("Erro ao converter price para JSON: %v", err)
		} else {
			log.Println("Price: ", string(body))
		}
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(payload))

	sigHeader := c.Request.Header.Get("stripe-signature")

	if err := h.paymentSrv.ValidateSignature(payload, sigHeader); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}
}
