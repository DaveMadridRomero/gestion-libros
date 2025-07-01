package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DaveMadridRomero/Gestion-Libros/internal/books"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/config"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/database"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/orders"
	"github.com/DaveMadridRomero/Gestion-Libros/internal/users"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	config.LoadEnv()
	database.Connect()

	err := database.DB.AutoMigrate(&users.User{}, &books.Book{}, &orders.Order{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos: %v", err)
	}
	log.Println("Migración de la base de datos completada.")

	r := mux.NewRouter()

	// ====================================================================
	// --- REGISTRO DE RUTAS API (¡DEBE IR ANTES DE LAS RUTAS ESTÁTICAS!) ---
	// ====================================================================
	users.RegisterRoutes(r)
	books.RegisterRoutes(r)
	orders.RegisterRoutes(r)
	// Agrega aquí cualquier otra ruta API (ej. admin, reviews)

	// ====================================================================
	// --- CONFIGURACIÓN PARA SERVIR ARCHIVOS ESTÁTICOS ---
	// ====================================================================

	// 1. Sirve la carpeta 'static' dentro de 'web/templates'
	// Es decir, 'your_project_root/web/templates/static/'
	// La URL en el navegador '/static/...' se mapea a 'web/templates/static/...'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/templates/static/"))))

	// 2. Sirve el index.html y cualquier otro archivo directamente en 'web/templates'
	// Es decir, 'your_project_root/web/templates/index.html'
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/index.html")
	}).Methods("GET")

	// ====================================================================
	// --- Configuración CORS ---
	// ====================================================================
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://127.0.0.1:8080", "http://127.0.0.1:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
