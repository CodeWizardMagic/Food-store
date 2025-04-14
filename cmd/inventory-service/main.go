package main

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/infrastructure/postgres"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}

	dbHost := os.Getenv("DB")
	if dbHost == "" {
		log.Fatal("DB environment variable not set")
	}
	postgres.InitDB(dbHost)
	log.Println("Connected to PostgreSQL")

	if err := postgres.InitTables(); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	productRepo := postgres.NewProductPostgresRepo()

	r := gin.Default()

	products := r.Group("/api/products")
	{
		products.GET("", func(c *gin.Context) {
			page, _ := strconv.Atoi(c.Query("page"))
			perPage, _ := strconv.Atoi(c.Query("per_page"))
			name := c.Query("name")
			minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
			maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

			filter := domain.FilterParams{
				Name:     name,
				MinPrice: minPrice,
				MaxPrice: maxPrice,
			}
			pagination := domain.PaginationParams{
				Page:    page,
				PerPage: perPage,
			}

			if pagination.Page < 1 {
				pagination.Page = 1
			}
			if pagination.PerPage < 1 {
				pagination.PerPage = 10
			}
			offset := (pagination.Page - 1) * pagination.PerPage

			products, total, err := productRepo.FindAllWithFilter(filter, pagination, offset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"products": products,
				"total":    total,
				"page":     pagination.Page,
				"per_page": pagination.PerPage,
			})
		})

		products.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			product, err := productRepo.FindByID(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			c.JSON(http.StatusOK, product)
		})

		products.POST("", func(c *gin.Context) {
			var product domain.Product
			if err := c.ShouldBindJSON(&product); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			product.ID = uuid.New().String()
			if err := productRepo.Save(product); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, product)
		})

		products.PUT("/:id", func(c *gin.Context) {
			id := c.Param("id")
			var product domain.Product
			if err := c.ShouldBindJSON(&product); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := productRepo.Update(id, product); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, product)
		})

		products.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			if err := productRepo.Delete(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})
	}

	port := os.Getenv("INVENTORY_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Inventory Service is starting on port %s...", port)
	r.Run(":" + port)
}
