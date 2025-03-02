package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type ProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
}

var products []Product

func main() {
	router := gin.Default()
	router.GET("/products", getProducts)
	router.GET("/product/:id", getProductByID)
	router.POST("/product", postNewProduct)

	router.PUT("/product/:id", putUpdateProduct)
	router.DELETE("/product/:id", deleteProduct)

	router.POST("/product/:id/image", uploadProductImage)
	router.GET("/product/:id/image", getProductImage)

	router.Run("localhost:8080")
}

func getProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, products)
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")

	for _, product := range products {
		if product.ID == id {
			c.IndentedJSON(http.StatusOK, product)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "product not found"})
}

func postNewProduct(c *gin.Context) {
	var request ProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if request.Name == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "name of product can't be empty"})
		return
	}

	for _, p := range products {
		if p.Name == request.Name && p.Description == request.Description {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "product already exists"})
			return
		}
	}

	newProduct := Product{
		ID:          uuid.New().String(),
		Name:        request.Name,
		Description: request.Description,
	}

	products = append(products, newProduct)
	c.IndentedJSON(http.StatusCreated, newProduct)
}

func putUpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var request ProductUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	for i, product := range products {
		if product.ID == id {
			if request.Name != nil {
				products[i].Name = *request.Name
			}
			if request.Description != nil {
				products[i].Description = *request.Description
			}
			if request.Icon != nil {
				products[i].Icon = *request.Icon
			}
			c.IndentedJSON(http.StatusOK, products[i])
			return
		}
	}
	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "product with this id doesn't exist"})
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "product not found"})
}

func uploadProductImage(c *gin.Context) {
	id := c.Param("id")

	var product *Product
	for i, p := range products {
		if p.ID == id {
			product = &products[i]
			break
		}
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	file, err := c.FormFile("icon")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	filePath := "uploads/" + id + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	product.Icon = filePath
	c.JSON(http.StatusOK, gin.H{"message": "image uploaded", "icon": filePath})
}

func getProductImage(c *gin.Context) {
	id := c.Param("id")

	var product *Product
	for _, p := range products {
		if p.ID == id {
			product = &p
			break
		}
	}
	if product == nil || product.Icon == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	c.File(product.Icon)
}
