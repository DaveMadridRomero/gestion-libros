## Autor

Madrid Romero Dave  
Universidad Internacional de Ecuador (UIDE)  
Curso: Programación Orientada a Objetos


# Gestión de Libros

## Descripción

Aplicación web para la gestión de libros que permite a los usuarios registrarse, iniciar sesión y administrar libros (crear, listar, eliminar). Utiliza autenticación con JWT y comunicación con un backend desarrollado en Go.

---

## Tecnologías utilizadas

- **Frontend:**  
  - HTML5  
  - CSS con [Tailwind CSS](https://tailwindcss.com/) (vía CDN)  
  - JavaScript (ES6) para lógica de autenticación y gestión de libros  

- **Backend:**  
  - API REST en Go (debe ejecutarse en `http://localhost:8080`)  
  - JWT para autenticación y autorización  

---

## Requisitos

- Navegador web moderno con soporte para JavaScript  
- Backend API en Go corriendo en `http://localhost:8080` con los siguientes endpoints:  
  - `POST /register` — registro de usuarios  
  - `POST /login` — inicio de sesión y obtención de token JWT  
  - `GET /books` — obtener lista de libros (requiere token JWT)  
  - `POST /books/create` — crear un nuevo libro (requiere token JWT)  
  - `DELETE /books/{id}` — eliminar libro por ID (requiere token JWT)  

---

## Instalación y ejecución


1. **Clonar el repositorio**
Utilizando git con: 

```bash
git clone https://github.com/DaveMadridRomero/gestion-libros/tree/GestionLibro
cd gestion-libros
```

O descargando el archivo zip.


Conexión a la Base de Datos
Para que la aplicación se conecte correctamente a la base de datos SQL Server (En este caso debe ser ese tipo de base de datos), se debe configurar las credenciales y datos de conexión en un archivo .env ubicado en la raíz del proyecto.

el archivo contiene las diferentes  variables de entorno:

Ejemplo:

DB_USER=davemadrid2
DB_PASS=123
DB_HOST=localhost:1433
DB_NAME=LibrosDB 

Una vez realizado ejecutaremos el archivo main.go en la terminal con el siguiente comando estando en la carpeta raíz del proyecto

```bash
go run ./cmd/main.go
```  
Si la conexión es exitosa se podrá revisar en `http://localhost:8080` 

## Funciones permitidas

Registro de nuevos usuarios

Inicio de sesión/Cierre de sesión 

Visualización de libros disponibles: Muestra una tabla con los libros existentes en la base de datos, con detalles como título, autor, género, precio y stock.

Añadir nuevos libros: Formulario para crear un nuevo libro con datos como título, autor, género, precio y cantidad en stock.

Eliminar libros: Permite eliminar libros existentes mediante un botón de "Eliminar" en cada fila de la tabla, con confirmación previa.

La aplicación muestra mensajes dinámicos para informar al usuario sobre el estado de sus acciones (éxito, error, información).

Mensajes como “Registro exitoso”, “Error al cargar libros”, “Libro eliminado”, etc.


## Funcionalidades Futuras para Mejorar el Proyecto 

Edición de libros existentes
Permitir que el usuario pueda modificar los datos de un libro ya creado (título, autor, precio, stock, etc.) mediante un formulario de edición.

Búsqueda y filtrado
Implementar un buscador para localizar libros por título, autor o género.

Añadir filtros para ordenar libros por precio, stock o género.

Paginación en la lista de libros
Para manejar grandes cantidades de libros, mostrar resultados paginados en la tabla para mejorar la experiencia y rendimiento.

Roles y permisos
Agregar roles de usuario (por ejemplo, administrador y usuario regular).

Limitar ciertas acciones (como eliminar o crear libros) solo para administradores.