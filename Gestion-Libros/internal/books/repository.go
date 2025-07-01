package books

type BookRepository interface {
	GetAll() ([]Book, error)
	Create(book *Book) error
	GetByID(id uint) (*Book, error)
	Delete(id uint) error
}
