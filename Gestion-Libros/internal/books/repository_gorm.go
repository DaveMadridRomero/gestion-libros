package books

import "github.com/DaveMadridRomero/Gestion-Libros/internal/database"

// Estructura vacía que implementa la interfaz BookRepository usando GORM
type GormBookRepo struct{}

// Obtiene todos los libros almacenados en la base de datos
func (g *GormBookRepo) GetAll() ([]Book, error) {
	var books []Book
	err := database.DB.Find(&books).Error
	return books, err
}

// Crea un nuevo libro en la base de datos
func (g *GormBookRepo) Create(book *Book) error {
	return database.DB.Create(book).Error
}

// Busca un libro por su ID y lo retorna si existe
func (g *GormBookRepo) GetByID(id uint) (*Book, error) {
	var book Book
	err := database.DB.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Elimina un libro por su ID desde la base de datos
func (g *GormBookRepo) Delete(id uint) error {
	result := database.DB.Delete(&Book{}, id)
	return result.Error
}
