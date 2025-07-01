package books

import (
	"errors"

	"github.com/DaveMadridRomero/Gestion-Libros/internal/database"
)

// Estructura que representa un libro dentro del sistema
type Book struct {
	ID     uint    `gorm:"primaryKey" json:"id"` // ID único, clave primaria
	Title  string  `json:"title"`                // Título del libro
	Author string  `json:"author"`               // Autor del libro
	Genre  string  `json:"genre"`                // Género o categoría
	Price  float64 `json:"price"`                // Precio actual del libro
	Stock  int     `json:"stock"`                // Cantidad disponible en inventario
}

// Definición de errores específicos relacionados con libros
var (
	ErrInvalidDiscount = errors.New("descuento inválido: el porcentaje debe estar entre 0 y 100")
	ErrBookNotFound    = errors.New("libro no encontrado")
	ErrBookCreation    = errors.New("error al crear el libro en la base de datos")
)

// Ejecuta la migración de la tabla de libros a la base de datos
func AutoMigrate() {
	database.DB.AutoMigrate(&Book{})
}

// Revisa si el libro está disponible en stock
func (b *Book) IsAvailable() bool {
	return b.Stock > 0
}

// Aplica un descuento al libro, si el porcentaje es válido
func (b *Book) ApplyDiscount(percent float64) error {
	if percent < 0 || percent > 100 {
		return ErrInvalidDiscount
	}
	b.Price -= b.Price * (percent / 100)
	return nil
}
