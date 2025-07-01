package orders

import (
	"encoding/json"
	"net/http"

	"github.com/DaveMadridRomero/Gestion-Libros/internal/database"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/middlewares"

	"github.com/gorilla/mux"
)

type Order struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	BookID uint
	Status string // "completed", etc.
}

func RegisterRoutes(r *mux.Router) {
	r.Handle("/orders", middlewares.JWTAuthMiddleware(http.HandlerFunc(CreateOrder))).Methods("POST")
	r.Handle("/library", middlewares.JWTAuthMiddleware(http.HandlerFunc(GetUserLibrary))).Methods("GET")
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BookID uint `json:"book_id"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	order := Order{
		UserID: userID,
		BookID: input.BookID,
		Status: "completed",
	}

	result := database.DB.Create(&order)
	if result.Error != nil {
		http.Error(w, "No se pudo registrar la orden", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Compra registrada"})
}

func GetUserLibrary(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	var orders []Order
	err := database.DB.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		http.Error(w, "Error al recuperar biblioteca", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
