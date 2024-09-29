package repository

import (
	"bobot/database"
	model "bobot/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword) // Set the hashed password

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := database.DB.Where("username = ?", loginDetails.Username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare the password
	if err := user.ComparePassword(loginDetails.Password); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	var tokenAuth *jwtauth.JWTAuth
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID, "username": user.Username})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "token": tokenString})
}
