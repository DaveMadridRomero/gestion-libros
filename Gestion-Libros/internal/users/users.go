package users

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DaveMadridRomero/Gestion-Libros/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Función auxiliar para enviar mensajes de error en formato JSON
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// User representa la estructura del usuario en la base de datos
// El campo Password se recibe en JSON y se almacena cifrado
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

// UserResponse es la estructura que se envía al cliente para evitar exponer la contraseña
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// RegisterRoutes registra las rutas para registro y login de usuarios
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
}

// RegisterHandler procesa el registro de un nuevo usuario
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input User

	// Decodifico el JSON recibido al struct User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("Error al decodificar JSON: %v", err)
		sendJSONError(w, "Error al procesar la solicitud JSON.", http.StatusBadRequest)
		return
	}

	// Validar que los campos requeridos no estén vacíos
	if input.Email == "" || input.Password == "" || input.Name == "" {
		sendJSONError(w, "Datos inválidos: nombre, email y contraseña son requeridos.", http.StatusBadRequest)
		return
	}

	// Cifrar la contraseña usando bcrypt antes de guardar
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		sendJSONError(w, "Error al procesar contraseña.", http.StatusInternalServerError)
		return
	}
	input.Password = string(hashedPassword)

	// Intentar crear el usuario en la base de datos
	result := database.DB.Create(&input)
	if result.Error != nil {
		log.Printf("Error al crear usuario en DB: %v", result.Error)
		sendJSONError(w, "No se pudo crear el usuario. El email podría estar en uso o hay un problema en la base de datos.", http.StatusInternalServerError)
		return
	}

	// Responder con éxito, sin incluir la contraseña
	w.WriteHeader(http.StatusCreated)
	responseUser := UserResponse{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Usuario creado con éxito",
		"user":    responseUser,
	})
}

// LoginHandler procesa el inicio de sesión y genera un token JWT si las credenciales son correctas
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	// Busco el usuario en la base de datos según el email proporcionado
	var user User
	result := database.DB.Where("email = ?", creds.Email).First(&user)
	if result.Error != nil {
		sendJSONError(w, "Usuario no encontrado.", http.StatusUnauthorized)
		return
	}

	// Comparo la contraseña enviada con el hash almacenado
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		sendJSONError(w, "Contraseña incorrecta.", http.StatusUnauthorized)
		return
	}

	// Creo el token JWT con user_id y expiración de 72 horas
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})

	// Firmo el token con la clave secreta del entorno
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		sendJSONError(w, "Error al generar token.", http.StatusInternalServerError)
		return
	}

	// Devuelvo el token al cliente
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
