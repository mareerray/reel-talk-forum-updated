package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/model"
)

// GetCategoriesHandler handles requests for fetching all categories
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Query the database for all categories
	rows, err := DB.Query("SELECT id, name, emoji FROM categories")
	if err != nil {
		log.Printf("Error querying categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the categories
	var categories []model.Category

	// Iterate through the rows and populate the categories slice
	for rows.Next() {
		var cat model.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Emoji); err != nil {
			log.Printf("Error scanning category row: %v", err)
			continue
		}
		categories = append(categories, cat)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating category rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode the categories as JSON and send the response
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		log.Printf("Error encoding categories to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
