package books

// BookRepository define las operaciones que cualquier repositorio de libros debe implementar.
// Esto permite separar la lógica del controlador de la lógica de acceso a datos,
// lo que facilita pruebas, mantenimiento y posibilidad de cambiar de implementación (por ejemplo, usar una base en memoria en vez de GORM).
type BookRepository interface {
	// Obtiene todos los libros disponibles en la base de datos
	GetAll() ([]Book, error)

	// Crea un nuevo libro en la base de datos
	Create(book *Book) error

	// Busca un libro por su ID
	GetByID(id uint) (*Book, error)

	// Elimina un libro según su ID
	Delete(id uint) error
}
