package main

import (
	"fmt"
	"go/src/broadcaster"
	"go/src/handler"
	"go/src/payment"
	"go/src/product"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dbURL := "user:password@tcp(127.0.0.1:3309)/goapi?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&payment.Payment{}, &product.Product{})

	broadcaster := broadcaster.NewBroadcaster(10)

	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	paymentRepository := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepository)
	paymentHandler := handler.NewPaymentHandler(paymentService, broadcaster)

	r := gin.Default()
	api := r.Group("/api")
	{
		products := api.Group("/products")
		{
			products.POST("/", productHandler.Create)
			products.GET("/", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
		}
		payments := api.Group("/payments")
		{
			payments.POST("/", paymentHandler.Create)
			payments.GET("/", paymentHandler.GetAll)
			payments.GET("/stream", paymentHandler.Stream)
			payments.GET("/:id", paymentHandler.GetById)
			payments.PUT("/:id", paymentHandler.Update)
			payments.DELETE("/:id", paymentHandler.Delete)
		}
	}

	r.Run(":3333")

	fmt.Println("Ca tourne !")
}
