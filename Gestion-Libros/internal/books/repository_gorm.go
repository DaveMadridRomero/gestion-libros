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

func (g *GormBookRepo) GetByID(id uint) (*Book, error) {
	var book Book
	err := database.DB.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Delete implementa el método para eliminar un libro por su ID de la base de datos.
func (g *GormBookRepo) Delete(id uint) error {
	// GORM permite eliminar un registro simplemente pasando el ID
	// o un modelo con el ID. Usaremos el ID directamente.
	// El método Delete devuelve un *gorm.DB que contiene el error si hubo uno.
	result := database.DB.Delete(&Book{}, id)
	return result.Error
}
