package books

type BookRepository interface {
	GetAll() ([]Book, error)
	Create(book *Book) error
}
