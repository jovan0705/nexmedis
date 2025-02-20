package controllers

import (
	"context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/jovan0705/nexmedis/models"
    "github.com/jovan0705/nexmedis/config"
    "github.com/jovan0705/nexmedis/helpers"
    "github.com/jackc/pgx/v4"
)

var users []models.User
type UserController struct{}

// RegisterUser handles user registration
func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    db := config.GetDB()

    // Insert user into the database
    _, err := db.Exec(context.Background(), "INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", user.Username, user.Password, user.Email)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// LoginUser handles user login
func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
    var loginData models.User
    if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    db := config.GetDB()
    var user models.User

    // Query user from the database
    row := db.QueryRow(context.Background(), "SELECT id, username, email FROM users WHERE username = $1 AND password = $2", loginData.Username, loginData.Password)
    err := row.Scan(&user.ID, &user.Username, &user.Email)

    if err != nil {
        if err == pgx.ErrNoRows {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        } else {
            http.Error(w, fmt.Sprintf("Error logging in: %v", err), http.StatusInternalServerError)
        }
        return
    }

    token, err := helpers.GenerateJWT(user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
