package books

import (
	"github.com/DaveMadridRomero/Gestion-Libros/internal/database"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Genre  string  `json:"genre"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

func AutoMigrate() {
	database.DB.AutoMigrate(&Book{})
}
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/books", GetAllBooks).Methods("GET")
	r.HandleFunc("/books/create", CreateBook).Methods("POST")
}

// IsAvailable verifica si el libro tiene stock disponible
func (b *Book) IsAvailable() bool {
	return b.Stock > 0
}

// ApplyDiscount aplica un descuento si es válido
func (b *Book) ApplyDiscount(percent float64) error {
	if percent < 0 || percent > 100 {
		return ErrInvalidDiscount
	}
	b.Price -= b.Price * (percent / 100)
	return nil
}
