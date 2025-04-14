package main

import (
	"FoodStore-AdvProg2/infrastructure/grpc"
	"FoodStore-AdvProg2/proto"
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type APIGateway struct {
	clients *grpc.Clients
}

func NewAPIGateway(clients *grpc.Clients) *APIGateway {
	return &APIGateway{clients: clients}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}

	inventoryAddr := os.Getenv("INVENTORY_SERVICE_GRPC_URL")
	orderAddr := os.Getenv("ORDER_SERVICE_GRPC_URL")
	userAddr := os.Getenv("USER_SERVICE_GRPC_URL")
	if inventoryAddr == "" || orderAddr == "" || userAddr == "" {
		log.Fatal("Service gRPC URLs must be set in .env")
	}

	clients, err := grpc.NewClients(inventoryAddr, orderAddr, userAddr)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer clients.Close()

	gateway := NewAPIGateway(clients)

	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gateway.MetricsMiddleware())
	r.Use(gateway.AuthMiddleware())

	// Static files and HTML
	r.Static("/static", "./public")
	r.LoadHTMLGlob("public/*.html")
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", nil)
	})
	r.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order.html", nil)
	})

	// Inventory API
	inventoryAPI := r.Group("/api/products")
	{
		inventoryAPI.POST("", gateway.CreateProduct)
		inventoryAPI.GET("", gateway.ListProducts)
		inventoryAPI.GET("/:id", gateway.GetProduct)
		inventoryAPI.PUT("/:id", gateway.UpdateProduct)
		inventoryAPI.DELETE("/:id", gateway.DeleteProduct)
	}

	// Order API
	orderAPI := r.Group("/api/orders")
	{
		orderAPI.POST("", gateway.CreateOrder)
		orderAPI.GET("", gateway.GetUserOrders)
		orderAPI.GET("/:id", gateway.GetOrder)
		orderAPI.PATCH("/:id", gateway.UpdateOrderStatus)
	}

	// User API
	userAPI := r.Group("/api/users")
	{
		userAPI.POST("/register", gateway.RegisterUser)
		userAPI.POST("/login", gateway.AuthenticateUser)
	}

	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway is starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}


func (g *APIGateway) MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
	}
}

func (g *APIGateway) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		path := strings.TrimSuffix(c.Request.URL.Path, "/")
		log.Printf("Processing request: %s %s", c.Request.Method, path)

		if path == "/api/users/register" || path == "/api/users/login" {
			log.Printf("Skipping auth for open endpoint: %s", path)
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		log.Printf("Authorization header: %s", token)
		if token == "" {
			log.Printf("No Authorization token provided for %s", path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		resp, err := g.clients.UserClient.ValidateToken(context.Background(), &proto.ValidateTokenRequest{Token: token})
		if err != nil {
			log.Printf("Invalid token for %s: %v", path, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		log.Printf("Valid token for user_id: %s", resp.UserId)
		c.Set("user_id", resp.UserId)
		c.Next()
	}
}

func (g *APIGateway) CreateProduct(c *gin.Context) {
	var req struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required,gt=0"`
		Stock int32   `json:"stock" binding:"required,gte=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	log.Printf("Creating product by user_id: %s", userID)

	resp, err := g.clients.InventoryClient.CreateProduct(context.Background(), &proto.CreateProductRequest{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	})
	if err != nil {
		log.Printf("Failed to create product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": resp.Id})
}

func (g *APIGateway) GetProduct(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Fetching product with id: %s", id)

	resp, err := g.clients.InventoryClient.GetProduct(context.Background(), &proto.GetProductRequest{Id: id})
	if err != nil {
		log.Printf("Failed to get product: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    resp.Id,
		"name":  resp.Name,
		"price": resp.Price,
		"stock": resp.Stock,
	})
}

func (g *APIGateway) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required,gt=0"`
		Stock int32   `json:"stock" binding:"required,gte=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Updating product with id: %s", id)
	resp, err := g.clients.InventoryClient.UpdateProduct(context.Background(), &proto.UpdateProductRequest{
		Id:    id,
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	})
	if err != nil {
		log.Printf("Failed to update product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    resp.Id,
		"name":  resp.Name,
		"price": resp.Price,
		"stock": resp.Stock,
	})
}

func (g *APIGateway) DeleteProduct(c *gin.Context) {
    id := c.Param("id")
    log.Printf("Deleting product with id: %s", id)

    ctx := context.Background()

    _, err := g.clients.OrderClient.DeleteOrderItemsByProduct(ctx, &proto.DeleteOrderItemsByProductRequest{ProductId: id})
    if err != nil {
        log.Printf("Failed to delete related order items: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related order items"})
        return
    }

    _, err = g.clients.InventoryClient.DeleteProduct(ctx, &proto.DeleteProductRequest{Id: id})
    if err != nil {
        log.Printf("Failed to delete product: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

func (g *APIGateway) ListProducts(c *gin.Context) {
	var req struct {
		Name     string  `form:"name"`
		MinPrice float64 `form:"min_price"`
		MaxPrice float64 `form:"max_price"`
		Page     int32   `form:"page"`
		PerPage  int32   `form:"per_page"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("Invalid query params: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PerPage < 1 {
		req.PerPage = 10
	}

	log.Printf("Listing products with filter: name=%s, min_price=%f, max_price=%f, page=%d, per_page=%d",
		req.Name, req.MinPrice, req.MaxPrice, req.Page, req.PerPage)

	resp, err := g.clients.InventoryClient.ListProducts(context.Background(), &proto.ListProductsRequest{
		Filter: &proto.FilterParams{
			Name:     req.Name,
			MinPrice: req.MinPrice,
			MaxPrice: req.MaxPrice,
		},
		Pagination: &proto.PaginationParams{
			Page:    req.Page,
			PerPage: req.PerPage,
		},
	})
	if err != nil {
		log.Printf("Failed to list products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	products := make([]gin.H, len(resp.Products))
	for i, p := range resp.Products {
		products[i] = gin.H{
			"id":    p.Id,
			"name":  p.Name,
			"price": p.Price,
			"stock": p.Stock,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"products":  products,
		"total":     resp.Total,
		"page":      resp.Page,
		"per_page":  resp.PerPage,
	})
}

// Order Handlers
func (g *APIGateway) CreateOrder(c *gin.Context) {
	var req struct {
		Items []struct {
			ProductID string `json:"product_id" binding:"required"`
			Quantity  int32  `json:"quantity" binding:"required,gt=0"`
		} `json:"items" binding:"required,dive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	log.Printf("Creating order for user_id: %s", userID)

	items := make([]*proto.OrderItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = &proto.OrderItemRequest{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	resp, err := g.clients.OrderClient.CreateOrder(context.Background(), &proto.CreateOrderRequest{
		UserId: userID.(string),
		Items:  items,
	})
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_id": resp.OrderId})
}

func (g *APIGateway) GetOrder(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Fetching order with id: %s", id)

	resp, err := g.clients.OrderClient.GetOrder(context.Background(), &proto.GetOrderRequest{OrderId: id})
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	items := make([]gin.H, len(resp.Items))
	for i, item := range resp.Items {
		items[i] = gin.H{
			"id":         item.Id,
			"order_id":   item.OrderId,
			"product_id": item.ProductId,
			"quantity":   item.Quantity,
			"price":      item.Price,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          resp.Id,
		"user_id":     resp.UserId,
		"total_price": resp.TotalPrice,
		"status":      resp.Status,
		"created_at":  resp.CreatedAt,
		"items":       items,
	})
}

func (g *APIGateway) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=pending completed cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Updating order status for id: %s to %s", id, req.Status)
	_, err := g.clients.OrderClient.UpdateOrderStatus(context.Background(), &proto.UpdateOrderStatusRequest{
		OrderId: id,
		Status:  req.Status,
	})
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (g *APIGateway) GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	log.Printf("Fetching orders for user_id: %s", userID)

	resp, err := g.clients.OrderClient.GetUserOrders(context.Background(), &proto.GetUserOrdersRequest{UserId: userID.(string)})
	if err != nil {
		log.Printf("Failed to get user orders: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orders := make([]gin.H, len(resp.Orders))
	for i, order := range resp.Orders {
		items := make([]gin.H, len(order.Items))
		for j, item := range order.Items {
			items[j] = gin.H{
				"id":         item.Id,
				"order_id":   item.OrderId,
				"product_id": item.ProductId,
				"quantity":   item.Quantity,
				"price":      item.Price,
			}
		}
		orders[i] = gin.H{
			"id":          order.Id,
			"user_id":     order.UserId,
			"total_price": order.TotalPrice,
			"status":      order.Status,
			"created_at":  order.CreatedAt,
			"items":       items,
		}
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// User Handlers
func (g *APIGateway) RegisterUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Registering user: %s", req.Username)
	resp, err := g.clients.UserClient.Register(context.Background(), &proto.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": resp.UserId})
}

func (g *APIGateway) AuthenticateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Authenticating user: %s", req.Username)
	resp, err := g.clients.UserClient.Authenticate(context.Background(), &proto.AuthenticateRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Failed to authenticate user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   resp.Token,
		"user_id": resp.UserId,
	})
}