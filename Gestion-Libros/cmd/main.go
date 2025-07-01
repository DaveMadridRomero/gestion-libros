package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DaveMadridRomero/Gestion-Libros/internal/books"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/config"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/database"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/orders"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/users" // Asegúrate de que users esté importado

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Cargar configuración
	config.LoadEnv()

	// Conectar base de datos
	database.Connect()

	// --- ¡MUY IMPORTANTE! REALIZAR LA MIGRACIÓN DE LA BASE DE DATOS ---
	// Esto creará las tablas 'users', 'books', 'orders' (y otras si las agregas)
	// si no existen, basándose en tus structs.
	err := database.DB.AutoMigrate(&users.User{}, &books.Book{}, &orders.Order{}) //
	if err != nil {
		log.Fatalf("Error al migrar la base de datos: %v", err) //
	}
	log.Println("Migración de la base de datos completada.") //
	// --- FIN DE MIGRACIÓN ---

	// Crear router
	r := mux.NewRouter()

	// Registrar rutas de usuarios
	users.RegisterRoutes(r)
	books.RegisterRoutes(r)
	orders.RegisterRoutes(r)

	// --- Configuración CORS ---
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                       // Permite cualquier origen (para desarrollo). En producción, especifica tus dominios.
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Métodos permitidos
		AllowedHeaders:   []string{"Authorization", "Content-Type"},           // Encabezados permitidos
		ExposedHeaders:   []string{"Link"},                                    // Otros encabezados que quieras exponer
		AllowCredentials: true,                                                // Permite el envío de credenciales (cookies, encabezados de autorización)
		MaxAge:           300,                                                 // Tiempo de caché de la pre-verificación OPTIONS
	})

	// Envuelve el router con el middleware CORS
	handler := c.Handler(r)
	// --- Fin de Configuración CORS ---

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
