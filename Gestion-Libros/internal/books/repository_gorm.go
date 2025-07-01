package books

import "github.com/DaveMadridRomero/Gestion-Libros/internal/database"

type GormBookRepo struct{}

func (g *GormBookRepo) GetAll() ([]Book, error) {
	var books []Book
	err := database.DB.Find(&books).Error
	return books, err
}

func (g *GormBookRepo) Create(book *Book) error {
	return database.DB.Create(book).Error
}
