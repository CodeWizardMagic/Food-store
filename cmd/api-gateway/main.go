package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/static", "./public")

	r.LoadHTMLGlob("public/*.html")
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", nil)
	})
	r.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order.html", nil)
	})

	inventoryAPI := r.Group("/api/products")
	{
		inventoryAPI.GET("", proxyToInventoryService)
		inventoryAPI.GET("/:id", proxyToInventoryService)
		inventoryAPI.POST("", proxyToInventoryService)
		inventoryAPI.PUT("/:id", proxyToInventoryService)
		inventoryAPI.DELETE("/:id", proxyToInventoryService)
	}

	orderAPI := r.Group("/api/orders")
	{
		orderAPI.GET("", proxyToOrderService)
		orderAPI.GET("/:id", proxyToOrderService)
		orderAPI.POST("", proxyToOrderService)
		orderAPI.PATCH("/:id", proxyToOrderService)
	}

	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway is starting on port %s...", port)
	r.Run(":" + port)
}

func proxyToInventoryService(c *gin.Context) {
	proxyRequest(c, os.Getenv("INVENTORY_SERVICE_URL"))
}

func proxyToOrderService(c *gin.Context) {
	proxyRequest(c, os.Getenv("ORDER_SERVICE_URL"))
}

func proxyRequest(c *gin.Context, serviceURL string) {
	if serviceURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service URL not configured"})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}

	targetURL := serviceURL + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}
