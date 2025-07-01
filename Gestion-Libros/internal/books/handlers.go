package books

import (
	"encoding/json"
	"net/http"
)

var repo BookRepository = &GormBookRepo{}

// GetAllBooks obtiene todos los libros desde la base de datos
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repo.GetAll()
	if err != nil {
		http.Error(w, "No se pudieron obtener los libros", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

// CreateBook permite agregar un nuevo libro a la plataforma
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Entrada inválida", http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == "" || book.Price <= 0 {
		http.Error(w, "Faltan campos requeridos", http.StatusBadRequest)
		return
	}

	// Verificamos que el stock sea positivo
	if !book.IsAvailable() {
		http.Error(w, "El libro debe tener al menos una unidad en stock", http.StatusBadRequest)
		return
	}

	err = repo.Create(&book)
	if err != nil {
		http.Error(w, ErrBookCreation.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}
