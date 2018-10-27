package api

import (
	"context"
	"encoding/json"
	"github.com/aweris/stp/internal/models"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"net/http"
)

func (ah *ApiHandler) registerTaxRoutes() {
	sub := ah.router.PathPrefix("/taxes").Subrouter()

	sub.HandleFunc("", ah.createTaxHandler).Methods("PUT")
	sub.HandleFunc("", ah.updateTaxHandler).Methods("POST")
	sub.HandleFunc("", ah.fetchTaxHandler).Methods("GET")
	sub.HandleFunc("/{id}", ah.deleteTaxHandler).Methods("DELETE")
	sub.HandleFunc("/{id}", ah.getTaxByIdHandler).Methods("GET")
}

type TaxDTO struct {
	Id         uuid.UUID           `json:"id"`
	Name       string              `json:"name"`
	Rate       decimal.Decimal     `json:"rate"`
	Origin     models.TaxOrigin    `json:"origin"`
	Condition  models.TaxCondition `json:"condition"`
	Categories []uuid.UUID         `json:"categories"`
}

func fromTaxToDTO(tax *models.Tax) *TaxDTO {

	categories := make([]uuid.UUID, 0, len(tax.Categories))
	for k := range tax.Categories {
		categories = append(categories, k)
	}

	return &TaxDTO{
		Id:         tax.Id,
		Name:       tax.Name,
		Rate:       tax.Rate,
		Origin:     tax.Origin,
		Condition:  tax.Condition,
		Categories: categories,
	}

}

func (t *TaxDTO) toTax() *models.Tax {
	categories := make(map[uuid.UUID]bool, len(t.Categories))
	for _, v := range t.Categories {
		categories[v] = true
	}

	return &models.Tax{
		Id:         t.Id,
		Name:       t.Name,
		Rate:       t.Rate,
		Origin:     t.Origin,
		Condition:  t.Condition,
		Categories: categories,
	}
}

func (ah *ApiHandler) createTaxHandler(w http.ResponseWriter, r *http.Request) {
	var dto TaxDTO
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	nt, err := ah.server.TaxService.CreateTax(r.Context(), dto.toTax())

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	data, err := json.Marshal(fromTaxToDTO(nt))
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Write(data)
}

func (ah *ApiHandler) updateTaxHandler(w http.ResponseWriter, r *http.Request) {
	var dto TaxDTO
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	nt, err := ah.server.TaxService.UpdateTax(r.Context(), dto.toTax())

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	data, err := json.Marshal(fromTaxToDTO(nt))
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Write(data)
}

func (ah *ApiHandler) fetchTaxHandler(w http.ResponseWriter, r *http.Request) {
	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	taxes, err := ah.server.TaxService.FetchAllTaxes(r.Context())

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	result := make([]*TaxDTO, 0, len(taxes))

	for _, v := range taxes {
		result = append(result, fromTaxToDTO(v))
	}

	// put tenant to bucket
	data, err := json.Marshal(taxes)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Write(data)
}

func (ah *ApiHandler) deleteTaxHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taxId := vars[`id`]

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	id, err := uuid.FromString(taxId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	t, err := ah.server.TaxService.DeleteTax(r.Context(), id)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if t == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data, err := json.Marshal(fromTaxToDTO(t))
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Write(data)
}

func (ah *ApiHandler) getTaxByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taxId := vars[`id`]

	// Timeout in context
	context.WithTimeout(
		r.Context(),
		ah.timeout,
	)

	id, err := uuid.FromString(taxId)
	if err != nil {
		http.Error(w, "Invalid id format", 500)
		return
	}

	t, err := ah.server.TaxService.GetTaxByID(r.Context(), id)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if t == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data, err := json.Marshal(fromTaxToDTO(t))
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Write(data)
}
