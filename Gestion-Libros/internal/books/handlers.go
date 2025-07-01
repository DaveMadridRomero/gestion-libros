package books

import (
	"encoding/json"
	"log" // Para logs de depuración
	"net/http"
	"strconv" // Para convertir string a uint (ID)

	"github.com/DaveMadridRomero/Gestion-Libros/internal/middlewares" // Importa tu middleware
	"github.com/gorilla/mux"                                          // Para manejar rutas con variables
)

// sendJSONError es una función de utilidad para enviar respuestas de error en formato JSON.
// La he copiado de users/handlers.go para consistencia.
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// repo es la instancia de tu BookRepository que se usará para interactuar con la DB.
var repo BookRepository = &GormBookRepo{}

// RegisterRoutes registra todas las rutas HTTP relacionadas con los libros.
// Algunas rutas están protegidas por JWTAuthMiddleware.
func RegisterRoutes(r *mux.Router) {
	// Rutas públicas (si las hubiera, aunque para libros suelen ser protegidas)
	// r.HandleFunc("/books", GetAllBooks).Methods("GET") // Si permites ver libros sin autenticar

	// Rutas protegidas por JWT
	protectedBooksRoutes := r.PathPrefix("/books").Subrouter()
	protectedBooksRoutes.Use(middlewares.JWTAuthMiddleware) // Aplica el middleware JWT a todas las rutas de este subrouter

	protectedBooksRoutes.HandleFunc("", GetAllBooks).Methods("GET")        // GET /books
	protectedBooksRoutes.HandleFunc("/create", CreateBook).Methods("POST") // POST /books/create
	protectedBooksRoutes.HandleFunc("/{id}", GetBookByID).Methods("GET")   // GET /books/{id}
	protectedBooksRoutes.HandleFunc("/{id}", DeleteBook).Methods("DELETE") // DELETE /books/{id}
}

// GetAllBooks obtiene todos los libros desde la base de datos
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repo.GetAll()
	if err != nil {
		log.Printf("Error al obtener todos los libros: %v", err)
		sendJSONError(w, "No se pudieron obtener los libros", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

// CreateBook permite agregar un nuevo libro a la plataforma
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Printf("Error al decodificar JSON para crear libro: %v", err)
		sendJSONError(w, "Entrada inválida. Asegúrate de enviar un JSON válido.", http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == "" || book.Price <= 0 {
		sendJSONError(w, "Faltan campos requeridos: título, autor y precio deben ser válidos y positivos.", http.StatusBadRequest)
		return
	}

	// Verificamos que el stock sea positivo
	if !book.IsAvailable() {
		sendJSONError(w, "El libro debe tener al menos una unidad en stock.", http.StatusBadRequest)
		return
	}

	err = repo.Create(&book)
	if err != nil {
		log.Printf("Error al crear libro en DB: %v", err)
		sendJSONError(w, ErrBookCreation.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Libro creado exitosamente"})
}

// GetBookByID obtiene un libro específico por su ID desde la base de datos.
func GetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32) // Convertir ID de string a uint
	if err != nil {
		log.Printf("Error: ID de libro inválido recibido: %s. Error: %v", idStr, err)
		sendJSONError(w, "ID de libro inválido.", http.StatusBadRequest)
		return
	}

	book, err := repo.GetByID(uint(id))
	if err != nil {
		log.Printf("Error al obtener libro con ID %d desde DB: %v", id, err)
		if err.Error() == "record not found" { // GORM devuelve "record not found" si no lo encuentra
			sendJSONError(w, ErrBookNotFound.Error(), http.StatusNotFound)
		} else {
			sendJSONError(w, "Error al obtener el libro.", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// DeleteBook elimina un libro específico por su ID de la base de datos.
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32) // Convertir ID de string a uint
	if err != nil {
		log.Printf("Error: ID de libro inválido recibido para eliminación: %s. Error: %v", idStr, err)
		sendJSONError(w, "ID de libro inválido.", http.StatusBadRequest)
		return
	}

	err = repo.Delete(uint(id))
	if err != nil {
		log.Printf("Error al eliminar libro con ID %d desde DB: %v", id, err)
		if err.Error() == "record not found" { // GORM devuelve "record not found" si no existe
			sendJSONError(w, ErrBookNotFound.Error(), http.StatusNotFound)
		} else {
			sendJSONError(w, "Error al eliminar el libro.", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(map[string]string{"message": "Libro eliminado exitosamente."})
}
