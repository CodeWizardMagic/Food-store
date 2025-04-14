package handler

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	UC *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{UC: uc}
}

func (h *ProductHandler) respondJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *ProductHandler) parseID(r *http.Request) string {
	return mux.Vars(r)["id"]
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	_ = json.NewDecoder(r.Body).Decode(&p)
	if err := h.UC.Create(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, p, http.StatusCreated)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	p, err := h.UC.GetByID(h.parseID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	h.respondJSON(w, p, http.StatusOK)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	_ = json.NewDecoder(r.Body).Decode(&p)
	if err := h.UC.Update(h.parseID(r), p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.respondJSON(w, p, http.StatusOK)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.UC.Delete(h.parseID(r)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
    name := r.URL.Query().Get("name")
    minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64)
    maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64)

    filter := domain.FilterParams{
        Name:     name,
        MinPrice: minPrice,
        MaxPrice: maxPrice,
    }
    pagination := domain.PaginationParams{
        Page:    page,
        PerPage: perPage,
    }

    products, total, err := h.UC.List(filter, pagination)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := struct {
        Products []domain.Product `json:"products"`
        Total    int              `json:"total"`
        Page     int              `json:"page"`
        PerPage  int              `json:"per_page"`
    }{
        Products: products,
        Total:    total,
        Page:     pagination.Page,
        PerPage:  pagination.PerPage,
    }

    h.respondJSON(w, response, http.StatusOK)
}


