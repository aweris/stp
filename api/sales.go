package api

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
)

func (ah *ApiHandler) registerSalesRoutes() {
	sale := ah.router.PathPrefix("/sales").Subrouter()

	br := sale.PathPrefix("/basket").Subrouter()

	br.HandleFunc("", ah.createBasketHandler).Methods("POST")
	br.HandleFunc("/{id}", ah.getBasketHandler).Methods("GET")
	br.HandleFunc("/{id}/item", ah.addItemToBasketHandler).Methods("POST")
	br.HandleFunc("/{id}/item", ah.deleteItemFromBasketHandler).Methods("DELETE")
	br.HandleFunc("/{id}/cancel", ah.cancelBasketHandler).Methods("POST")
	br.HandleFunc("/{id}/close", ah.closeBasketHandler).Methods("POST")

	rr := sale.PathPrefix("/receipt").Subrouter()

	rr.HandleFunc("", ah.fetchAllReceiptsHandler).Methods("GET")
	rr.HandleFunc("/{id}", ah.getReceiptHandler).Methods("GET")
}

type BasketDTO struct {
	Id uuid.UUID `json:"id"`
}

type BasketItemDTO struct {
	ItemId uuid.UUID `json:"item_id"`
	Count  int       `json:"count"`
}

func (ah *ApiHandler) createBasketHandler(w http.ResponseWriter, r *http.Request) {
	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	bid, err := ah.server.SaleService.CreateBasket(r.Context())

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&BasketDTO{Id: bid})
}

func (ah *ApiHandler) getBasketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	basketId := vars[`id`]

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	id, err := uuid.FromString(basketId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	t, err := ah.server.SaleService.GetBasketByID(r.Context(), id)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if t == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (ah *ApiHandler) addItemToBasketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	basketId := vars[`id`]

	id, err := uuid.FromString(basketId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}
	var b BasketItemDTO
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	err = ah.server.SaleService.AddItem(r.Context(), id, b.ItemId, b.Count)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (ah *ApiHandler) deleteItemFromBasketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	basketId := vars[`id`]

	id, err := uuid.FromString(basketId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}
	var b BasketItemDTO
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	err = ah.server.SaleService.RemoveItem(r.Context(), id, b.ItemId, b.Count)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (ah *ApiHandler) cancelBasketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	basketId := vars[`id`]

	id, err := uuid.FromString(basketId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	err = ah.server.SaleService.CancelBasket(r.Context(), id)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (ah *ApiHandler) closeBasketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	basketId := vars[`id`]

	id, err := uuid.FromString(basketId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	receipt, err := ah.server.SaleService.CloseBasket(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if receipt == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

func (ah *ApiHandler) fetchAllReceiptsHandler(w http.ResponseWriter, r *http.Request) {
	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	receipts, err := ah.server.SaleService.FetchAllReceipts(r.Context())

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if receipts == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipts)
}

func (ah *ApiHandler) getReceiptHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	basketId := vars[`id`]

	id, err := uuid.FromString(basketId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	receipt, err := ah.server.SaleService.GetReceiptByID(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if receipt == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}
