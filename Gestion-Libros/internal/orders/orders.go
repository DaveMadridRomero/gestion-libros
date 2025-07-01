package orders

import (
	"encoding/json"
	"net/http"

	"github.com/DaveMadridRomero/Gestion-Libros/internal/database"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/middlewares"

	"github.com/gorilla/mux"
)

// Order representa una compra hecha por un usuario
type Order struct {
	ID     uint   `gorm:"primaryKey"` // ID de la orden
	UserID uint   // ID del usuario que realiza la compra
	BookID uint   // ID del libro comprado
	Status string // Estado de la orden (ej. "completed")
}

// RegisterRoutes define las rutas HTTP relacionadas con órdenes de compra
func RegisterRoutes(r *mux.Router) {
	// Ruta protegida para crear una orden (compra)
	r.Handle("/orders", middlewares.JWTAuthMiddleware(http.HandlerFunc(CreateOrder))).Methods("POST")

	// Ruta protegida para obtener las órdenes de un usuario (su biblioteca)
	r.Handle("/library", middlewares.JWTAuthMiddleware(http.HandlerFunc(GetUserLibrary))).Methods("GET")
}

// CreateOrder registra una nueva orden de compra para el usuario autenticado
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Estructura para recibir el ID del libro desde el cuerpo de la petición
	var input struct {
		BookID uint `json:"book_id"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	// Obtengo el user_id desde el contexto (ya validado por el middleware JWT)
	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	// Creo la orden con estado "completed" por defecto
	order := Order{
		UserID: userID,
		BookID: input.BookID,
		Status: "completed",
	}

	// Intento guardar la orden en la base de datos
	result := database.DB.Create(&order)
	if result.Error != nil {
		http.Error(w, "No se pudo registrar la orden", http.StatusInternalServerError)
		return
	}

	// Respondo con mensaje de éxito
	json.NewEncoder(w).Encode(map[string]string{"message": "Compra registrada"})
}

// GetUserLibrary devuelve todas las órdenes asociadas al usuario autenticado
func GetUserLibrary(w http.ResponseWriter, r *http.Request) {
	// Obtengo el user_id desde el contexto
	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	// Busco todas las órdenes donde user_id coincida con el usuario autenticado
	var orders []Order
	err := database.DB.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		http.Error(w, "Error al recuperar biblioteca", http.StatusInternalServerError)
		return
	}

	// Envío la lista de órdenes como respuesta
	json.NewEncoder(w).Encode(orders)
}
