package handler

import (
	"go/src/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type productHandler struct {
	productService product.Service
}

func NewProductHandler(productService product.Service) *productHandler {
	return &productHandler{productService}
}

func (ph *productHandler) Create(c *gin.Context) {
	var input product.InputProduct
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := ProductResponse{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newProduct, err := ph.productService.Create(input)
	if err != nil {
		response := ProductResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := ProductResponse{
		Success: true,
		Message: "New product created",
		Data:    newProduct,
	}
	c.JSON(http.StatusCreated, response)
}

func (ph *productHandler) GetAll(c *gin.Context) {
	products, err := ph.productService.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, ProductResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Data:    products,
	})
}

func (ph *productHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ProductResponse{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	product, err := ph.productService.GetById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ProductResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Data:    product,
	})
}

func (ph *productHandler) Update(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, ProductResponse{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	var input product.InputProduct
	err = c.ShouldBindJSON(&input)
	if err != nil {
		response := ProductResponse{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := ph.productService.Update(id, input)
	if err != nil {
		response := ProductResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := ProductResponse{
		Success: true,
		Message: "Product updated",
		Data:    product,
	}
	c.JSON(http.StatusCreated, response)
}

func (ph *productHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, ProductResponse{
			Success: false,
			Message: "Product not found",
			Data:    err.Error(),
		})
		return
	}

	err = ph.productService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ProductResponse{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ProductResponse{
		Success: true,
		Message: "Product deleted",
	})
}
