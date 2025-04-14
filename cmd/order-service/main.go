package main

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/infrastructure/postgres"
	"FoodStore-AdvProg2/usecase"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	orderRepo := postgres.NewOrderPostgresRepo()
	productRepo := postgres.NewProductPostgresRepo()
	orderUC := usecase.NewOrderUseCase(orderRepo, productRepo)

	r := gin.Default()

	orders := r.Group("/api/orders")
	{
		orders.POST("", func(c *gin.Context) {
			var orderReq domain.OrderRequest
			if err := c.ShouldBindJSON(&orderReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
				return
			}

			orderID, err := orderUC.CreateOrder(orderReq)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"order_id": orderID})
		})

		orders.GET("", func(c *gin.Context) {
			userID := c.Query("user_id")
			var orders []domain.Order
			var err error

			if userID != "" {
				orders, err = orderUC.GetOrdersByUserID(userID)
			} else {
				orders, err = orderUC.GetAllOrders()
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, orders)
		})

		orders.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			order, err := orderUC.GetOrderByID(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			}
			c.JSON(http.StatusOK, order)
		})

		orders.PATCH("/:id", func(c *gin.Context) {
			id := c.Param("id")
			var statusReq domain.OrderStatusUpdateRequest
			if err := c.ShouldBindJSON(&statusReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
				return
			}

			if err := orderUC.UpdateOrderStatus(id, statusReq.Status); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "updated"})
		})
	}

	port := os.Getenv("ORDER_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Order Service is starting on port %s...", port)
	r.Run(":" + port)
}
