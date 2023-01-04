package handler

import (
	"go/src/broadcaster"
	"go/src/payment"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentResponse struct {
	Success bool
	Message string
	Data    interface{}
}

type paymentHandler struct {
	paymentService payment.Service
	broadcaster    broadcaster.Broadcaster
}

func NewPaymentHandler(paymentService payment.Service, broadcaster broadcaster.Broadcaster) *paymentHandler {
	return &paymentHandler{
		paymentService,
		broadcaster,
	}
}

func (ph *paymentHandler) Stream(c *gin.Context) {
	listener := make(chan interface{})
	ph.broadcaster.Register(listener)
	defer ph.broadcaster.Unregister(listener)

	c.Stream(func(w io.Writer) bool {
		for {
			select {
			case pmt := <-listener:
				c.SSEvent("Payment created/updated", pmt)
				return true
			}
		}
	})
}

func (ph *paymentHandler) Create(c *gin.Context) {
	var input payment.InputPayment
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := &PaymentResponse{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newPayment, err := ph.paymentService.Create(input)
	if err != nil {
		response := &PaymentResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &PaymentResponse{
		Success: true,
		Message: "New payment created",
		Data:    newPayment,
	}
	c.JSON(http.StatusCreated, response)

	//On envoie le payment au broadcaster
	ph.broadcaster.Submit(newPayment)
}

func (ph *paymentHandler) GetAll(c *gin.Context) {
	payments, err := ph.paymentService.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, PaymentResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaymentResponse{
		Success: true,
		Data:    payments,
	})
}

func (ph *paymentHandler) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, PaymentResponse{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	payment, err := ph.paymentService.GetById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, PaymentResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, PaymentResponse{
		Success: true,
		Data:    payment,
	})
}

func (ph *paymentHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, PaymentResponse{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	var input payment.InputPayment
	err = c.ShouldBindJSON(&input)
	if err != nil {
		response := &PaymentResponse{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	payment, err := ph.paymentService.Update(id, input)
	if err != nil {
		response := &PaymentResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &PaymentResponse{
		Success: true,
		Message: "Payment updated",
		Data:    payment,
	}
	c.JSON(http.StatusCreated, response)

	//On envoie le payment au broadcaster
	ph.broadcaster.Submit(payment)
}

func (ph *paymentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, PaymentResponse{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	err = ph.paymentService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, PaymentResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaymentResponse{
		Success: true,
		Message: "Payment successfully deleted",
	})
}
