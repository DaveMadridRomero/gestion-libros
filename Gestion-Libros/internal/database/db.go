package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Variable global que mantiene la conexión activa a la base de datos
var DB *gorm.DB

// Connect se encarga de establecer la conexión a SQL Server usando GORM
func Connect() {
	// Leo las variables de entorno necesarias para la conexión
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")

	// Verifico que ninguna de las variables críticas esté vacía
	if user == "" || pass == "" || host == "" || name == "" {
		log.Fatal("❌ Error: Asegúrate de que DB_USER, DB_PASS, DB_HOST y DB_NAME estén configuradas en el archivo .env")
	}

	// Construyo el DSN (cadena de conexión) para SQL Server
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=true&trustServerCertificate=true", user, pass, host, name)

	// Configuro el logger de GORM para que muestre consultas en consola (útil en desarrollo)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // Tiempo que define una consulta como lenta
			LogLevel:      logger.Info, // Nivel de detalle del log (usar Silent en producción)
			Colorful:      true,        // Habilito colores para facilitar lectura en consola
		},
	)

	// Intento abrir la conexión usando el driver de SQL Server y GORM
	var err error
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a la base de datos: %v", err)
	}

	log.Println("✅ Conexión exitosa a SQL Server")

	// Ajustes del pool de conexiones (uso *sql.DB para configurarlo)
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Error al obtener instancia para configuración del pool: %v", err)
	}

	// Limito la cantidad de conexiones inactivas para evitar consumo innecesario
	sqlDB.SetMaxIdleConns(10)

	// Defino el número máximo de conexiones abiertas al mismo tiempo
	sqlDB.SetMaxOpenConns(100)

	// Establezco cuánto tiempo puede mantenerse viva una conexión antes de cerrarse
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Pool de conexiones configurado correctamente")
}

// CloseDB permite cerrar la conexión a la base de datos manualmente.
// Útil si la aplicación termina o se necesita reiniciar la conexión.
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("⚠️ No se pudo obtener la instancia para cerrar: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Printf("❌ Error al cerrar la base de datos: %v", err)
		} else {
			log.Println("✅ Conexión a la base de datos cerrada correctamente")
		}
	}
}
