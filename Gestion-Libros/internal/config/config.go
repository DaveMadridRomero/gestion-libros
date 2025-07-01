package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv intenta cargar las variables de entorno desde un archivo .env.
// Si no encuentra el archivo, solo muestra un mensaje en el log, pero no detiene la ejecución.
// Esto es útil en desarrollo, pero en producción puede que no sea necesario el .env.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env (puede ser normal en producción)")
	}
}
