package users

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/DaveMadridRomero/Gestion-Libros/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Email == "" || input.Password == "" {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Cifrar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al procesar contraseña", http.StatusInternalServerError)
		return
	}

	input.Password = string(hashedPassword)
	result := database.DB.Create(&input)

	if result.Error != nil {
		http.Error(w, "No se pudo crear el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario creado con éxito"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	var user User
	result := database.DB.Where("email = ?", creds.Email).First(&user)
	if result.Error != nil {
		http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
		return
	}

	// Crear token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Error al generar token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
