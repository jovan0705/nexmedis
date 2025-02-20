package controllers

import (
	"log"
    "context"
    "encoding/json"
    "net/http"
	"io/ioutil"
    "github.com/jovan0705/nexmedis/models"
    "github.com/jovan0705/nexmedis/config"
	"github.com/gorilla/mux"
	"fmt"
)

// CartController defines the controller for managing cart operations
type CartController struct{}

func (cc *CartController) GetCart(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["userId"]
    conn := config.GetDB()
    rows, err := conn.Query(context.Background(), "SELECT user_id, product_id, quantity, price FROM cart_items WHERE user_id = $1 AND isPaid = false", userID)
    if err != nil {
        http.Error(w, "Error fetching cart items", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var cartItems []models.CartItem
    for rows.Next() {
        var item models.CartItem
		
        if err := rows.Scan(&item.UserID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			log.Println(err)
            http.Error(w, "Error reading cart item", http.StatusInternalServerError)
            return
        }
		log.Println(item)
        cartItems = append(cartItems, item)
    }

    if len(cartItems) == 0 {
        http.Error(w, "Cart is empty", http.StatusNotFound)
        return
    }

    if err := rows.Err(); err != nil {
        http.Error(w, "Error with database rows", http.StatusInternalServerError)
        return
    }

	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += float64(item.Quantity) * item.Price
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
		"totalPrice": totalPrice,
		"cartItems": cartItems,
	})
}

func (cc *CartController) AddToCart(w http.ResponseWriter, r *http.Request) {
    var item models.CartItem
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    conn := config.GetDB()
    defer conn.Close(context.Background())

    _, err := conn.Exec(context.Background(), 
        "INSERT INTO cart_items (user_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
        item.UserID, item.ProductID, item.Quantity, item.Price)
    if err != nil {
        http.Error(w, "Error adding item to cart", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(item)
}

func (cc *CartController) PayCartItems(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.Body)
	// log.Println(json.NewDecoder(r.Body))
	// log.Println("hehe")
	vars := mux.Vars(r)
    userID := vars["userId"]
    conn := config.GetDB()
	
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error amount not found", http.StatusInternalServerError)
        return
	}
	var jsonData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

    rows, err := conn.Query(context.Background(), "SELECT user_id, product_id, quantity, price FROM cart_items WHERE user_id = $1 AND isPaid = false", userID)
    if err != nil {
        http.Error(w, "Error fetching cart items", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

	if !rows.Next() {
		http.Error(w, "Cart is empty", http.StatusNotFound)
		return
	}

    var cartItems []models.CartItem
    for rows.Next() {
        var item models.CartItem
		log.Println(item)
        if err := rows.Scan(&item.UserID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			log.Println(err)
            http.Error(w, "Error reading cart item", http.StatusInternalServerError)
            return
        }
        cartItems = append(cartItems, item)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, "Error with database rows", http.StatusInternalServerError)
        return
    }

	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += float64(item.Quantity) * item.Price
	}

	if totalPrice > jsonData["amount"].(float64) {
		http.Error(w, "Amount paid not enough", http.StatusInternalServerError)
        return
	}

	_, err = conn.Exec(context.Background(), 
	"UPDATE cart_items SET isPaid=true WHERE user_id=$1", userID)
    if err != nil {
        http.Error(w, "Error adding item to cart", http.StatusInternalServerError)
        return
    }

	message := map[string]string{
		"message": fmt.Sprintf("Cart Paid successfully with change %.0f", jsonData["amount"].(float64)-totalPrice),
	}
	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
}
