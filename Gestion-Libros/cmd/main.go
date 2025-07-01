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
	// Carga las variables de entorno desde el archivo .env
	config.LoadEnv()

	// Establece la conexión con la base de datos
	database.Connect()

	// Aplica la migración automática para las tablas de usuarios, libros y órdenes
	err := database.DB.AutoMigrate(&users.User{}, &books.Book{}, &orders.Order{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos: %v", err)
	}
	log.Println("Migración de la base de datos completada correctamente.")

	// Inicializa el router usando Gorilla Mux
	r := mux.NewRouter()

	// Registro de todas las rutas de la API (usuarios, libros, órdenes)
	users.RegisterRoutes(r)
	books.RegisterRoutes(r)
	orders.RegisterRoutes(r)

	// Configuro la ruta para servir archivos estáticos desde la carpeta /web/templates/static
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./web/templates/static/"))),
	)

	// Cuando el usuario accede a "/", le envío el archivo index.html de la carpeta templates
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/index.html")
	}).Methods("GET")

	// Habilito CORS para permitir que el frontend (por ejemplo desde localhost:5500) acceda a la API
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://127.0.0.1:8080", "http://127.0.0.1:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Aplico la configuración CORS al router
	handler := c.Handler(r)

	// Obtengo el puerto desde las variables de entorno, o uso 8080 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicio del servidor
	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
