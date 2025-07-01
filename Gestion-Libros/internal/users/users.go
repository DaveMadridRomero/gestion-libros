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

func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Email string `gorm:"unique" json:"email"`
	// ¡CAMBIO AQUÍ! Ahora el campo Password sí se decodificará del JSON entrante
	Password string `json:"password"`
}

// Struct auxiliar para la respuesta, para NO enviar la contraseña
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// Aquí NO incluimos el campo Password
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input User // input es el User que recibimos del JSON
	err := json.NewDecoder(r.Body).Decode(&input)

	log.Printf("Intento de registro: Email=%s, Name=%s, Password (recibido)=%s", input.Email, input.Name, input.Password) // Vuelve a revisar este log después del cambio

	if err != nil {
		log.Printf("Error al decodificar JSON: %v", err)
		sendJSONError(w, "Error al procesar la solicitud JSON.", http.StatusBadRequest)
		return
	}

	// Asegúrate de que todos los campos requeridos estén presentes y no vacíos.
	if input.Email == "" || input.Password == "" || input.Name == "" {
		sendJSONError(w, "Datos inválidos: nombre, email y contraseña son requeridos.", http.StatusBadRequest)
		return
	}

	// Cifrar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		sendJSONError(w, "Error al procesar contraseña.", http.StatusInternalServerError)
		return
	}

	input.Password = string(hashedPassword) // Almacenamos la contraseña cifrada en el campo Password
	result := database.DB.Create(&input)

	if result.Error != nil {
		log.Printf("Error al crear usuario en DB: %v", result.Error)
		sendJSONError(w, "No se pudo crear el usuario. El email podría ya estar en uso o hay un problema en la base de datos.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// Para la respuesta, crea un UserResponse para no enviar la contraseña
	responseUser := UserResponse{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Usuario creado con éxito",
		"user":    responseUser, // Puedes enviar los datos del usuario (sin contraseña)
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	log.Printf("Intento de login: Email=%s", creds.Email)

	var user User
	result := database.DB.Where("email = ?", creds.Email).First(&user)
	if result.Error != nil {
		log.Printf("Error al buscar usuario en DB para login: %v", result.Error)
		sendJSONError(w, "Usuario no encontrado.", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		sendJSONError(w, "Contraseña incorrecta.", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		sendJSONError(w, "Error al generar token.", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
