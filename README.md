FoodStore API
This is the API documentation for the FoodStore application, a microservices-based system for managing users, products, and orders. The API is exposed through an API Gateway running on http://localhost:8080. Below are all available routes, their methods, required headers, request/response formats, and descriptions of the actions they perform.
Prerequisites

Services: Ensure the following services are running:
User Service (go run cmd/user-service/main.go) on port :50052
Inventory Service (go run cmd/inventory-service/main.go) on port :50053
Order Service (go run cmd/order-service/main.go) on port :50051
API Gateway (go run cmd/api-gateway/main.go) on port :8080


Database: PostgreSQL must be running and accessible at the URL specified in the .env file (e.g., postgresql://postgres:12345678@localhost:5432/postgres?sslmode=disable).
Environment: A .env file should exist in the project root with the following variables:DB=postgresql://postgres:12345678@localhost:5432/postgres?sslmode=disable
INVENTORY_SERVICE_GRPC_URL=localhost:50053
ORDER_SERVICE_GRPC_URL=localhost:50051
USER_SERVICE_GRPC_URL=localhost:50052
API_GATEWAY_PORT=8080



Authentication
Most endpoints require an Authorization header with a token obtained via the /api/users/login endpoint. The token must be included as follows:
Authorization: <your-token>

API Routes
1. User Management
Register a User

Method: POST
URL: http://localhost:8080/api/users/register
Headers:Content-Type: application/json


Request Body:{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}


Response (Status: 201 Created):{
  "user_id": "uuid-string"
}


Description: Creates a new user account. Returns a unique user ID upon success.
Errors:
400 Bad Request: Invalid or missing fields, or user already exists.
500 Internal Server Error: Database or server issues.



Login

Method: POST
URL: http://localhost:8080/api/users/login
Headers:Content-Type: application/json


Request Body:{
  "username": "testuser",
  "password": "password123"
}


Response (Status: 200 OK):{
  "token": "auth-token-string",
  "user_id": "uuid-string"
}


Description: Authenticates a user and returns an authorization token and user ID.
Errors:
400 Bad Request: Invalid request body.
401 Unauthorized: Incorrect username or password.
500 Internal Server Error: Server issues.



2. Product Management
All product endpoints require authentication.
Create a Product

Method: POST
URL: http://localhost:8080/api/products
Headers:Content-Type: application/json
Authorization: <your-token>


Request Body:{
  "name": "Apple",
  "price": 1.99,
  "stock": 100
}


Response (Status: 201 Created):{
  "id": "product-uuid"
}


Description: Adds a new product to the inventory. Returns the product ID.
Errors:
400 Bad Request: Invalid fields (e.g., negative price or stock).
401 Unauthorized: Missing or invalid token.
500 Internal Server Error: Server issues.



List Products

Method: GET
URL: http://localhost:8080/api/products?name=apple&min_price=1&max_price=5&page=1&per_page=10
Headers:Authorization: <your-token>


Query Parameters (optional):
name: Filter by product name (partial match).
min_price: Minimum price filter.
max_price: Maximum price filter.
page: Page number (default: 1).
per_page: Items per page (default: 10).


Response (Status: 200 OK):{
  "products": [
    {
      "id": "product-uuid",
      "name": "Apple",
      "price": 1.99,
      "stock": 100
    }
  ],
  "total": 1,
  "page": 1,
  "per_page": 10
}


Description: Retrieves a paginated list of products with optional filtering.
Errors:
400 Bad Request: Invalid query parameters.
401 Unauthorized: Missing or invalid token.
500 Internal Server Error: Server issues.



Get a Product

Method: GET
URL: http://localhost:8080/api/products/<product-id>
Headers:Authorization: <your-token>


Response (Status: 200 OK):{
  "id": "product-uuid",
  "name": "Apple",
  "price": 1.99,
  "stock": 100
}


Description: Retrieves details of a specific product by its ID.
Errors:
401 Unauthorized: Missing or invalid token.
404 Not Found: Product does not exist.
500 Internal Server Error: Server issues.



Update a Product

Method: PUT
URL: http://localhost:8080/api/products/<product-id>
Headers:Content-Type: application/json
Authorization: <your-token>


Request Body:{
  "name": "Green Apple",
  "price": 2.49,
  "stock": 150
}


Response (Status: 200 OK):{
  "id": "product-uuid",
  "name": "Green Apple",
  "price": 2.49,
  "stock": 150
}


Description: Updates the details of an existing product.
Errors:
400 Bad Request: Invalid fields.
401 Unauthorized: Missing or invalid token.
404 Not Found: Product does not exist.
500 Internal Server Error: Server issues.



Delete a Product

Method: DELETE
URL: http://localhost:8080/api/products/<product-id>
Headers:Authorization: <your-token>


Response (Status: 204 No Content):
No body returned.


Description: Deletes a product from the inventory.
Errors:
401 Unauthorized: Missing or invalid token.
500 Internal Server Error: Server issues or foreign key constraint violation (e.g., product is referenced in orders).



Note: If you encounter a foreign key constraint error (order_items_product_id_fkey), ensure related order_items are handled (e.g., by setting ON DELETE CASCADE or ON DELETE SET NULL in the database schema).
3. Order Management
All order endpoints require authentication.
Create an Order

Method: POST
URL: http://localhost:8080/api/orders
Headers:Content-Type: application/json
Authorization: <your-token>


Request Body:{
  "items": [
    {
      "product_id": "product-uuid",
      "quantity": 2
    }
  ]
}


Response (Status: 201 Created):{
  "order_id": "order-uuid"
}


Description: Creates a new order for the authenticated user with specified items.
Errors:
400 Bad Request: Invalid items (e.g., non-existent product or invalid quantity).
401 Unauthorized: Missing or invalid token.
500 Internal Server Error: Server issues or insufficient stock.



List Orders

Method: GET
URL: http://localhost:8080/api/orders
Headers:Authorization: <your-token>


Response (Status: 200 OK):{
  "orders": [
    {
      "id": "order-uuid",
      "user_id": "user-uuid",
      "total_price": 3.98,
      "status": "pending",
      "created_at": 1697059200,
      "items": [
        {
          "id": "item-uuid",
          "order_id": "order-uuid",
          "product_id": "product-uuid",
          "quantity": 2,
          "price": 1.99
        }
      ]
    }
  ]
}


Description: Retrieves all orders for the authenticated user.
Errors:
401 Unauthorized: Missing or invalid token.
500 Internal Server Error: Server issues.



Get an Order

Method: GET
URL: http://localhost:8080/api/orders/<order-id>
Headers:Authorization: <your-token>


Response (Status: 200 OK):{
  "id": "order-uuid",
  "user_id": "user-uuid",
  "total_price": 3.98,
  "status": "pending",
  "created_at": 1697059200,
  "items": [
    {
      "id": "item-uuid",
      "order_id": "order-uuid",
      "product_id": "product-uuid",
      "quantity": 2,
      "price": 1.99
    }
  ]
}


Description: Retrieves details of a specific order by its ID.
Errors:
401 Unauthorized: Missing or invalid token.
404 Not Found: Order does not exist.
500 Internal Server Error: Server issues.



Update Order Status

Method: PATCH
URL: http://localhost:8080/api/orders/<order-id>
Headers:Content-Type: application/json
Authorization: <your-token>


Request Body:{
  "status": "completed"
}


Response (Status: 200 OK):{
  "status": "updated"
}


Valid Status Values:
pending
completed
cancelled


Description: Updates the status of an existing order.
Errors:
400 Bad Request: Invalid status value.
401 Unauthorized: Missing or invalid token.
404 Not Found: Order does not exist.
500 Internal Server Error: Server issues.



Testing with Postman

Set up Postman:

Create a new collection for FoodStore API.
Store the authentication token in an environment variable (e.g., {{token}}) for reuse in the Authorization header.


Sample Workflow:

Register a user (POST /api/users/register).
Log in to get a token (POST /api/users/login).
Create a product (POST /api/products).
Create an order (POST /api/orders).
List orders (GET /api/orders).
Update a product (PUT /api/products/<product-id>).
Update order status (PATCH /api/orders/<order-id>).
Delete a product (DELETE /api/products/<product-id>).


Error Handling:

Check the response body for error messages (e.g., {"error": "Invalid token"}).
If you encounter a 500 error, check the logs of the respective service (user-service, inventory-service, order-service, or api-gateway).



Notes

Database Constraints: If deleting a product fails due to a foreign key constraint (order_items_product_id_fkey), consider modifying the order_items table to include ON DELETE CASCADE or ON DELETE SET NULL for the product_id foreign key.
Logging: Enable verbose logging in services to debug issues. Logs are printed to the console for each service.
Security: The current setup uses insecure gRPC connections. For production, enable TLS for gRPC and HTTPS for the API Gateway.

For further assistance or to report issues, please check the service logs or consult the project maintainers.
