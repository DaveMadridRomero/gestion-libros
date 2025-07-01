package database

import (
	"fmt"
	"log"
	"os"
	"time" // Necesario para time.Hour en SetConnMaxLifetime

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Importar el paquete logger de GORM
)

// DB es la conexión global a la base de datos
var DB *gorm.DB

// Connect establece la conexión con SQL Server
func Connect() {
	// Obtener los valores de las variables de entorno
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST") // ej: localhost:1433 o IP:Puerto
	name := os.Getenv("DB_NAME")

	// Validar que las variables de entorno críticas no estén vacías
	if user == "" || pass == "" || host == "" || name == "" {
		log.Fatal("❌ Error: Asegúrate de que DB_USER, DB_PASS, DB_HOST y DB_NAME estén configuradas en tu archivo .env")
	}

	// Formato del DSN para SQL Server. Incluye parámetros de encriptación recomendados.
	// encrypt=true: Fuerza el cifrado de la conexión.
	// trustServerCertificate=true: Permite conexiones con certificados autofirmados (común en desarrollo).
	// ¡CUIDADO con trustServerCertificate=true en producción sin un certificado válido!
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=true&trustServerCertificate=true", user, pass, host, name)

	// Configurar el logger de GORM para mostrar queries en desarrollo
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io.Writer: os.Stdout para imprimir en consola
		logger.Config{
			SlowThreshold: time.Second, // Umbral de tiempo para considerar una consulta "lenta"
			LogLevel:      logger.Info, // Cambia a logger.Silent en producción para no mostrar queries
			Colorful:      true,        // Activa los colores en la salida del log
		},
	)

	var err error
	// Abrir la conexión a la base de datos con GORM y el driver de SQL Server
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: newLogger, // Asigna el logger configurado a GORM
	})
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a la base de datos SQL Server: %v", err)
	}

	log.Println("✅ Conexión exitosa a SQL Server con GORM")

	// Configuración adicional del pool de conexiones
	sqlDB, err := DB.DB() // Obtener la instancia *sql.DB subyacente
	if err != nil {
		log.Fatalf("❌ Error al obtener la instancia SQL DB para configurar el pool: %v", err)
	}

	// Establecer el número máximo de conexiones inactivas en el pool
	sqlDB.SetMaxIdleConns(10)

	// Establecer el número máximo de conexiones abiertas al DB
	// Es importante no tener demasiadas conexiones abiertas simultáneamente.
	sqlDB.SetMaxOpenConns(100)

	// Establecer el tiempo máximo que una conexión puede ser reutilizada.
	// Después de este tiempo, la conexión será cerrada y recreada si es necesario.
	// Ayuda a manejar el reciclaje de conexiones y problemas de red transitorios.
	sqlDB.SetConnMaxLifetime(time.Hour) // Ejemplo: las conexiones durarán como máximo 1 hora

	log.Println("✅ Pool de conexiones configurado.")
}

// CloseDB cierra la conexión a la base de datos.
// Es buena práctica llamar a esta función cuando la aplicación se detiene
// (ej. usando un defer en main o en un manejador de señales de cierre).
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("⚠️ Error al obtener la instancia SQL DB para cerrar: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Printf("❌ Error al cerrar la conexión a la base de datos: %v", err)
		} else {
			log.Println("✅ Conexión a la base de datos cerrada.")
		}
	}
}
