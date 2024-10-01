// this should be handlers
package repository

import (
	"bobot/database"
	model "bobot/models"
	"bobot/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetAllEntry(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	orderBy := r.URL.Query().Get("order_by")

	entry, err := model.FindAll(database.DB, order, orderBy)
	if err != nil {
		http.Error(w, "Failed to create entry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

func FindEntryByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	entry, err := model.FindByID(database.DB, id)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entry)
}

func GetEntriesByPattern(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")
	order := r.URL.Query().Get("order")

	entries, err := model.SearchEntries(database.DB, pattern, order)
	if err != nil {
		http.Error(w, "Failed to retrieve entries", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entries)
}

func FindEntryByNumber(w http.ResponseWriter, r *http.Request) {
	num := chi.URLParam(r, "num")
	order := r.URL.Query().Get("order")

	entry, err := model.FindByEntryNumber(database.DB, num, order)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entry)
}

func CreateEntry(w http.ResponseWriter, r *http.Request) {
	var entry model.Entry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := entry.Create(database.DB); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

func UpdateEntry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	entry, err := model.FindByID(database.DB, id)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	entry.ID = id

	nentry, err := entry.Update(database.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mentry := model.Entry{
		ID:           nentry.ID,
		Entry_number: nentry.Entry_number,
		Content:      nentry.Content,
		SearchVector: nentry.SearchVector,
		UpdatedAt:    nentry.UpdatedAt,
		CreatedAt:    nentry.CreatedAt,
	}

	entry_map, err := utils.StructToMap(mentry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	std_ret := utils.StandardResponse(true, entry_map)
	log.Print(std_ret)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(std_ret)
}

func DeleteEntry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	entry, err := model.FindByID(database.DB, id)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	if err := entry.Delete(database.DB, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
