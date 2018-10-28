package api

import (
	"context"
	"encoding/json"
	"github.com/aweris/stp/internal/models"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
)

func (ah *ApiHandler) registerInventoryRoutes() {
	inv := ah.router.PathPrefix("/inv").Subrouter()

	cat := inv.PathPrefix("/categories").Subrouter()

	cat.HandleFunc("", ah.createCategoryHandler).Methods("PUT")
	cat.HandleFunc("", ah.updateCategoryHandler).Methods("POST")
	cat.HandleFunc("", ah.fetchCategoryHandler).Methods("GET")
	cat.HandleFunc("/{id}", ah.deleteCategoryHandler).Methods("DELETE")
	cat.HandleFunc("/{id}", ah.getCategoryByIdHandler).Methods("GET")
	cat.HandleFunc("/name/{name}", ah.getCategoryByNameHandler).Methods("GET")

	it := inv.PathPrefix("/items").Subrouter()

	it.HandleFunc("", ah.createItemHandler).Methods("PUT")
	it.HandleFunc("", ah.updateItemHandler).Methods("POST")
	it.HandleFunc("", ah.fetchItemHandler).Methods("GET")
	it.HandleFunc("/{id}", ah.deleteItemHandler).Methods("DELETE")
	it.HandleFunc("/{id}", ah.getItemByIdHandler).Methods("GET")
	it.HandleFunc("/category/{category}", ah.getItemByCategoryIdHandler).Methods("GET")
}

func (ah *ApiHandler) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	nc, err := ah.server.InventoryService.CreateCategory(r.Context(), &c)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nc)
}

func (ah *ApiHandler) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	uc, err := ah.server.InventoryService.UpdateCategory(r.Context(), &c)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uc)
}

func (ah *ApiHandler) fetchCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	categories, err := ah.server.InventoryService.FetchAllCategories(r.Context())

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (ah *ApiHandler) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	itemId := vars[`id`]

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	id, err := uuid.FromString(itemId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	t, err := ah.server.InventoryService.DeleteCategory(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if t == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (ah *ApiHandler) getCategoryByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	categoryId := vars[`id`]

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	id, err := uuid.FromString(categoryId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	t, err := ah.server.InventoryService.GetCategoryByID(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if t == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (ah *ApiHandler) getCategoryByNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	categoryName := vars[`name`]

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	c, err := ah.server.InventoryService.GetCategoryByName(r.Context(), categoryName)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if c == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (ah *ApiHandler) createItemHandler(w http.ResponseWriter, r *http.Request) {
	var i models.InventoryItem
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	ni, err := ah.server.InventoryService.CreateItem(r.Context(), &i)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ni)
}

func (ah *ApiHandler) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	var i models.InventoryItem
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	ui, err := ah.server.InventoryService.UpdateItem(r.Context(), &i)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ui)
}

func (ah *ApiHandler) fetchItemHandler(w http.ResponseWriter, r *http.Request) {
	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	categories, err := ah.server.InventoryService.FetchAllItems(r.Context())

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (ah *ApiHandler) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	itemId := vars[`id`]

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	id, err := uuid.FromString(itemId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	t, err := ah.server.InventoryService.DeleteItem(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if t == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (ah *ApiHandler) getItemByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	itemId := vars[`id`]

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	id, err := uuid.FromString(itemId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	t, err := ah.server.InventoryService.GetItemByID(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if t == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (ah *ApiHandler) getItemByCategoryIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	categoryId := vars[`category`]

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	id, err := uuid.FromString(categoryId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	t, err := ah.server.InventoryService.GetItemsByCategoryID(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if t == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
