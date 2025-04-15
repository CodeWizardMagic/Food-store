# 🍔 FoodStore API Documentation

Welcome to the **FoodStore API**, a microservices-based system for managing users, products, and orders. This documentation provides all the information you need to set up, run, and interact with the API.

---

## 🚀 Prerequisites

### Services
Ensure the following services are up and running:

- **User Service** — `go run cmd/user-service/main.go` → Port `:50052`
- **Inventory Service** — `go run cmd/inventory-service/main.go` → Port `:50053`
- **Order Service** — `go run cmd/order-service/main.go` → Port `:50051`
- **API Gateway** — `go run cmd/api-gateway/main.go` → Port `:8080`

### Database
- PostgreSQL should be accessible at the URL defined in the `.env` file, e.g.:
  ```env
  postgresql://postgres:12345678@localhost:5432/postgres?sslmode=disable
  ```

### Environment Variables
Ensure a `.env` file exists in the project root with the following:
```env
DB=postgresql://postgres:12345678@localhost:5432/postgres?sslmode=disable
INVENTORY_SERVICE_GRPC_URL=localhost:50053
ORDER_SERVICE_GRPC_URL=localhost:50051
USER_SERVICE_GRPC_URL=localhost:50052
API_GATEWAY_PORT=8080
```

---

## 🔐 Authentication
Most endpoints require an `Authorization` header with a token obtained from the `/api/users/login` endpoint:
```
Authorization: <your-token>
```

---

## 👤 1. User Management

### ✉️ Register a User
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/users/register`
- **Headers:** `Content-Type: application/json`
- **Request Body:**
```json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```
- **Response (201):**
```json
{
  "user_id": "uuid-string"
}
```
- **Errors:** `400`, `500`

### 🔑 Login
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/users/login`
- **Headers:** `Content-Type: application/json`
- **Request Body:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```
- **Response (200):**
```json
{
  "token": "auth-token-string",
  "user_id": "uuid-string"
}
```
- **Errors:** `400`, `401`, `500`

---

## 🍎 2. Product Management *(Requires Authentication)*

### ➕ Create a Product
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/products`
- **Headers:**
  - `Content-Type: application/json`
  - `Authorization: <your-token>`
- **Request Body:**
```json
{
  "name": "Apple",
  "price": 1.99,
  "stock": 100
}
```
- **Response (201):** `{ "id": "product-uuid" }`
- **Errors:** `400`, `401`, `500`

### 📝 List Products
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/products`
- **Headers:** `Authorization: <your-token>`
- **Query Parameters (optional):**
  - `name`, `min_price`, `max_price`, `page`, `per_page`
- **Response (200):**
```json
{
  "products": [ { "id": "...", "name": "Apple", "price": 1.99, "stock": 100 } ],
  "total": 1,
  "page": 1,
  "per_page": 10
}
```
- **Errors:** `400`, `401`, `500`

### 🔎 Get a Product
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/products/<product-id>`
- **Headers:** `Authorization: <your-token>`
- **Response (200):**
```json
{ "id": "...", "name": "Apple", "price": 1.99, "stock": 100 }
```
- **Errors:** `401`, `404`, `500`

### ✏️ Update a Product
- **Method:** `PUT`
- **URL:** `http://localhost:8080/api/products/<product-id>`
- **Headers:** `Content-Type: application/json`, `Authorization`
- **Request Body:**
```json
{ "name": "Green Apple", "price": 2.49, "stock": 150 }
```
- **Response (200):** Updated product object
- **Errors:** `400`, `401`, `404`, `500`

### ❌ Delete a Product
- **Method:** `DELETE`
- **URL:** `http://localhost:8080/api/products/<product-id>`
- **Headers:** `Authorization`
- **Response (204):** No content
- **Errors:** `401`, `500`

---

## 🛒 3. Order Management *(Requires Authentication)*

### 📅 Create an Order
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/orders`
- **Headers:** `Content-Type: application/json`, `Authorization`
- **Request Body:**
```json
{
  "items": [ { "product_id": "product-uuid", "quantity": 2 } ]
}
```
- **Response (201):** `{ "order_id": "order-uuid" }`
- **Errors:** `400`, `401`, `500`

### 📃 List Orders
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/orders`
- **Headers:** `Authorization`
- **Response (200):** Array of order objects
- **Errors:** `401`, `500`

### 🔍 Get an Order
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/orders/<order-id>`
- **Headers:** `Authorization`
- **Response (200):** Order details with items
- **Errors:** `401`, `404`, `500`

### ✅ Update Order Status
- **Method:** `PATCH`
- **URL:** `http://localhost:8080/api/orders/<order-id>`
- **Headers:** `Content-Type: application/json`, `Authorization`
- **Request Body:**
```json
{ "status": "completed" }
```
- **Response (200):** `{ "status": "updated" }`
- **Errors:** `400`, `401`, `404`, `500`

---

## 💡 Testing with Postman

1. Create a new Postman collection
2. Save the token from login in an environment variable (e.g. `{{token}}`)
3. Use `Authorization: {{token}}` in headers

**Sample Workflow:**
- Register → Login → Create Product → Create Order → List Orders

---

## ⚠️ Notes
- **Foreign Key Errors:** If deleting a product fails, consider setting `ON DELETE CASCADE` or `SET NULL` on foreign keys
- **Logs:** All services print logs to the console. Use these for debugging
- **Security:** For production, use TLS for gRPC and HTTPS for the API Gateway

---

For issues or questions, check the logs or contact the maintainers.

---

📅 Last Updated: April 2025

