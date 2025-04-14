package handler

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/usecase"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	UC *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{UC: uc}
}

func (h *OrderHandler) respondJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *OrderHandler) parseID(r *http.Request) string {
	return mux.Vars(r)["id"]
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq domain.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	orderID, err := h.UC.CreateOrder(orderReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{"order_id": orderID}
	h.respondJSON(w, response, http.StatusCreated)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := h.parseID(r)
	order, err := h.UC.GetOrderByID(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	h.respondJSON(w, order, http.StatusOK)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := h.parseID(r)

	var statusReq domain.OrderStatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&statusReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := h.UC.UpdateOrderStatus(id, statusReq.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.respondJSON(w, map[string]string{"status": "updated"}, http.StatusOK)
}

func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		orders, err := h.UC.GetAllOrders()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.respondJSON(w, orders, http.StatusOK)
		return
	}

	orders, err := h.UC.GetOrdersByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, orders, http.StatusOK)
}
