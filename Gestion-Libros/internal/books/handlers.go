package books

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/DaveMadridRomero/Gestion-Libros/internal/middlewares"
	"github.com/gorilla/mux"
)

// Función auxiliar para enviar mensajes de error en formato JSON
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// Defino el repositorio que usará las funciones de acceso a datos de libros
var repo BookRepository = &GormBookRepo{}

// Registro de todas las rutas relacionadas con libros
func RegisterRoutes(r *mux.Router) {
	// Agrupo todas las rutas bajo el prefijo /books y les aplico el middleware JWT
	protectedBooksRoutes := r.PathPrefix("/books").Subrouter()
	protectedBooksRoutes.Use(middlewares.JWTAuthMiddleware)

	protectedBooksRoutes.HandleFunc("", GetAllBooks).Methods("GET")
	protectedBooksRoutes.HandleFunc("/create", CreateBook).Methods("POST")
	protectedBooksRoutes.HandleFunc("/{id}", GetBookByID).Methods("GET")
	protectedBooksRoutes.HandleFunc("/{id}", DeleteBook).Methods("DELETE")
}

// Retorna todos los libros disponibles en la base de datos
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repo.GetAll()
	if err != nil {
		log.Printf("Error al obtener libros: %v", err)
		sendJSONError(w, "No se pudieron obtener los libros", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

// Crea un nuevo libro a partir de los datos recibidos en el cuerpo de la solicitud
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Printf("Error al parsear el JSON: %v", err)
		sendJSONError(w, "Entrada inválida. Asegúrate de enviar un JSON válido.", http.StatusBadRequest)
		return
	}

	// Valido campos obligatorios antes de guardar
	if book.Title == "" || book.Author == "" || book.Price <= 0 {
		sendJSONError(w, "Faltan campos requeridos o los valores son inválidos.", http.StatusBadRequest)
		return
	}

	// Verifico que haya al menos un libro en stock
	if !book.IsAvailable() {
		sendJSONError(w, "El libro debe tener al menos una unidad en stock.", http.StatusBadRequest)
		return
	}

	err = repo.Create(&book)
	if err != nil {
		log.Printf("Error al crear el libro: %v", err)
		sendJSONError(w, ErrBookCreation.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Libro creado exitosamente"})
}

// Busca un libro por su ID y lo devuelve como respuesta
func GetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("ID inválido recibido: %s - %v", idStr, err)
		sendJSONError(w, "ID de libro inválido.", http.StatusBadRequest)
		return
	}

	book, err := repo.GetByID(uint(id))
	if err != nil {
		log.Printf("No se encontró el libro con ID %d: %v", id, err)
		if err.Error() == "record not found" {
			sendJSONError(w, ErrBookNotFound.Error(), http.StatusNotFound)
		} else {
			sendJSONError(w, "Error al obtener el libro.", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// Elimina un libro según su ID
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("ID inválido recibido para eliminación: %s - %v", idStr, err)
		sendJSONError(w, "ID de libro inválido.", http.StatusBadRequest)
		return
	}

	err = repo.Delete(uint(id))
	if err != nil {
		log.Printf("Error al eliminar libro con ID %d: %v", id, err)
		if err.Error() == "record not found" {
			sendJSONError(w, ErrBookNotFound.Error(), http.StatusNotFound)
		} else {
			sendJSONError(w, "Error al eliminar el libro.", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Libro eliminado exitosamente."})
}
