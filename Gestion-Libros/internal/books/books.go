package books

import (
	"errors" // Necesario para definir los errores

	"github.com/DaveMadridRomero/Gestion-Libros/internal/database"
)

type Book struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Genre  string  `json:"genre"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

// Errores específicos del paquete books
var (
	ErrInvalidDiscount = errors.New("descuento inválido: el porcentaje debe estar entre 0 y 100")
	ErrBookNotFound    = errors.New("libro no encontrado") // Asegúrate de que este error esté aquí
	ErrBookCreation    = errors.New("error al crear el libro en la base de datos")
)

func AutoMigrate() {
	database.DB.AutoMigrate(&Book{})
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
