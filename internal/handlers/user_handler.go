package handlers

import (
	"database/sql"
	"encoding/json"
	"latihan_devops/internal/models"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name, email, created_at FROM users ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createdAt time.Time
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &createdAt)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		user.CreatedAt = createdAt.Format(time.RFC3339)
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// ✅ FIX: Check error from Encode
	if err := json.NewEncoder(w).Encode(users); err != nil {
		// Log error internally, but response might already be sent
		log.Printf("Error encoding users response: %v", err)
	}
}

// POST /users - Membuat user baru
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if newUser.Name == "" || newUser.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(
		"INSERT INTO users (name, email) VALUES (?, ?)",
		newUser.Name, newUser.Email,
	)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var createdUser models.User
	var createdAt time.Time
	err = h.DB.QueryRow(
		"SELECT id, name, email, created_at FROM users WHERE id = ?", id,
	).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdAt)

	if err != nil {
		http.Error(w, "Error fetching created user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdUser.CreatedAt = createdAt.Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// ✅ FIX: Check error from Encode
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		log.Printf("Error encoding created user response: %v", err)
	}
}

// Health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// ✅ FIX: Check error from Encode
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding health check response: %v", err)
	}
}
