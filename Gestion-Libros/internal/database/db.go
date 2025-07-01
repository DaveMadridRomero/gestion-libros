package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// DB es la conexión global a la base de datos
var DB *gorm.DB

// Connect establece la conexión con SQL Server
func Connect() {
	user := os.Getenv("DB_USER") // ej: davemadrid2
	pass := os.Getenv("DB_PASS") // ej: 123
	host := os.Getenv("DB_HOST") // ej: localhost:1433
	name := os.Getenv("DB_NAME") // ej: LibrosDB

	// Formato del DSN para SQL Server
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", user, pass, host, name)

	var err error
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a la base de datos SQL Server: %v", err)
	}

	log.Println("✅ Conexión exitosa a SQL Server con GORM")
}
