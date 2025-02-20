package controllers

import (
	"context"
    "encoding/json"
    "net/http"
   	"github.com/jovan0705/nexmedis/models"
    "github.com/jovan0705/nexmedis/config"
	"github.com/jackc/pgx/v4"
	"github.com/gorilla/mux"
)

type ProductController struct{}

// GetAllProducts handles fetching all products
func (pc *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) {
    conn := config.GetDB()
    search := r.URL.Query().Get("search")
    query := "SELECT id, name, price, description FROM products WHERE name ILIKE $1"
    rows, err := conn.Query(context.Background(), query, "%"+search+"%")
    if err != nil {
        http.Error(w, "Error fetching products", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var p models.Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description); err != nil {
            http.Error(w, "Error reading product data", http.StatusInternalServerError)
            return
        }
        products = append(products, p)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, "Error with database rows", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

// GetProduct handles fetching a single product by its ID
func (pc *ProductController) GetProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    conn := config.GetDB()
    row := conn.QueryRow(context.Background(), "SELECT id, name, price, description FROM products WHERE id = $1", id)

    var p models.Product
    if err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Description); err != nil {
        if err == pgx.ErrNoRows {
            http.Error(w, "Product not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error fetching product", http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(p)
}
