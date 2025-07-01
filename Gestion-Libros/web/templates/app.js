// URL base de la API
const API_BASE_URL = 'http://localhost:8080';

// Intento obtener el token JWT guardado en el almacenamiento local (localStorage)
let authToken = localStorage.getItem('jwtToken');

// Referencias a elementos del DOM para manipular la UI
const authSection = document.getElementById('auth-section'); // Sección de login/registro
const appSection = document.getElementById('app-section');   // Sección principal de la app (libros)
const messageContainer = document.getElementById('message-container'); // Contenedor para mostrar mensajes

const loginForm = document.getElementById('login-form');         // Formulario de login
const registerForm = document.getElementById('register-form');   // Formulario de registro
const showLoginBtn = document.getElementById('show-login-btn');  // Botón para mostrar formulario login
const showRegisterBtn = document.getElementById('show-register-btn'); // Botón para mostrar formulario registro
const logoutBtn = document.getElementById('logout-btn');         // Botón para cerrar sesión

const createBookForm = document.getElementById('create-book-form'); // Formulario para crear libro
const booksTableBody = document.getElementById('books-table-body'); // Cuerpo de la tabla de libros
const refreshBooksBtn = document.getElementById('refresh-books-btn'); // Botón para recargar libros

// Función para mostrar mensajes en la UI con estilos según tipo (info, error, success)
function showMessage(message, type = 'info') {
    messageContainer.textContent = message;
    // Remuevo todas las clases de estilo para asegurar que solo quede el estilo deseado
    messageContainer.classList.remove('hidden', 'bg-red-100', 'text-red-700', 'bg-green-100', 'text-green-700', 'bg-blue-100', 'text-blue-700');
    
    // Agrego clase según el tipo de mensaje
    if (type === 'error') {
        messageContainer.classList.add('bg-red-100', 'text-red-700');
    } else if (type === 'success') {
        messageContainer.classList.add('bg-green-100', 'text-green-700');
    } else {
        messageContainer.classList.add('bg-blue-100', 'text-blue-700');
    }

    // Muestro el mensaje
    messageContainer.classList.remove('hidden');

    // Oculto el mensaje automáticamente después de 5 segundos
    setTimeout(() => {
        messageContainer.classList.add('hidden');
    }, 5000);
}

// Alterna la visibilidad entre la sección de autenticación y la app principal según el estado del token
function toggleAppVisibility() {
    if (authToken) {
        // Si hay token, oculto la autenticación y muestro la app principal
        authSection.classList.add('hidden');
        appSection.classList.remove('hidden');
        loadBooks(); // Cargo los libros cuando la sesión está activa
    } else {
        // Si no hay token, muestro la autenticación y oculto la app principal
        authSection.classList.remove('hidden');
        appSection.classList.add('hidden');
        loginForm.classList.remove('hidden');   // Por defecto muestro el login
        registerForm.classList.add('hidden');   // Oculto el registro
    }
}

// --- Manejo de vistas de login y registro ---

// Mostrar formulario de login cuando el usuario presiona el botón correspondiente
showLoginBtn.addEventListener('click', () => {
    loginForm.classList.remove('hidden');
    registerForm.classList.add('hidden');
});

// Mostrar formulario de registro cuando el usuario presiona el botón correspondiente
showRegisterBtn.addEventListener('click', () => {
    registerForm.classList.remove('hidden');
    loginForm.classList.add('hidden');
});

// --- Registro de usuario ---

registerForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    // Obtengo los valores del formulario de registro
    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    try {
        // Envío petición POST para registrar usuario
        const response = await fetch(`${API_BASE_URL}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password })
        });
        const data = await response.json();

        if (response.ok) {
            // Registro exitoso, muestro mensaje y cambio a formulario login
            showMessage(data.message || 'Registro exitoso. ¡Ahora puedes iniciar sesión!', 'success');
            registerForm.reset();
            loginForm.classList.remove('hidden');
            registerForm.classList.add('hidden');
        } else {
            // Error en el registro, muestro mensaje de error
            showMessage(data.message || 'Error en el registro', 'error');
        }
    } catch (error) {
        // Error de red u otro error inesperado
        showMessage('Error de red al registrar: ' + error.message, 'error');
    }
});

// --- Inicio de sesión ---

loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    // Obtengo los valores del formulario de login
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        // Envío petición POST para login
        const response = await fetch(`${API_BASE_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        const data = await response.json();

        if (response.ok && data.token) {
            // Login exitoso: guardo token y muestro la app principal
            authToken = data.token;
            localStorage.setItem('jwtToken', authToken);
            showMessage('Inicio de sesión exitoso', 'success');
            loginForm.reset();
            toggleAppVisibility();
        } else {
            // Credenciales inválidas u otro error
            showMessage(data.message || 'Error al iniciar sesión. Credenciales inválidas.', 'error');
        }
    } catch (error) {
        // Error de red u otro error inesperado
        showMessage('Error de red al iniciar sesión: ' + error.message, 'error');
    }
});

// --- Cerrar sesión ---

logoutBtn.addEventListener('click', () => {
    // Limpio el token y vuelvo a la vista de autenticación
    authToken = null;
    localStorage.removeItem('jwtToken');
    showMessage('Sesión cerrada correctamente', 'info');
    toggleAppVisibility();
    booksTableBody.innerHTML = ''; // Limpio la tabla de libros
});

// --- Manejo y carga de libros ---

// Función para cargar todos los libros desde la API
async function loadBooks() {
    if (!authToken) {
        showMessage('No autenticado. Por favor, inicie sesión para ver los libros.', 'error');
        return;
    }

    try {
        // Solicito la lista de libros enviando el token JWT en la cabecera Authorization
        const response = await fetch(`${API_BASE_URL}/books`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        const books = await response.json();

        if (response.ok) {
            booksTableBody.innerHTML = ''; // Limpio la tabla antes de llenarla

            if (books.length === 0) {
                // Si no hay libros, muestro mensaje en la tabla
                booksTableBody.innerHTML = '<tr><td colspan="7" class="px-6 py-4 whitespace-nowrap text-center text-gray-500">No hay libros disponibles.</td></tr>';
                return;
            }

            // Recorro la lista de libros y agrego filas con datos y botones
            books.forEach(book => {
                const row = booksTableBody.insertRow();

                // Cada celda representa un campo del libro
                const idCell = row.insertCell(0);
                idCell.classList.add('px-6', 'py-4', 'whitespace-nowrap');
                idCell.textContent = book.id;

                const titleCell = row.insertCell(1);
                titleCell.classList.add('px-6', 'py-4', 'whitespace-nowrap');
                titleCell.textContent = book.title;

                const authorCell = row.insertCell(2);
                authorCell.classList.add('px-6', 'py-4', 'whitespace-nowrap');
                authorCell.textContent = book.author;

                const genreCell = row.insertCell(3);
                genreCell.classList.add('px-6', 'py-4', 'whitespace-nowrap');
                genreCell.textContent = book.genre;

                const priceCell = row.insertCell(4);
                priceCell.classList.add('px-6', 'py-4', 'whitespace-nowrap');
                priceCell.textContent = `$${book.price.toFixed(2)}`;

                const stockCell = row.insertCell(5);
                stockCell.classList.add('px-6', 'py-4', 'whitespace-nowrap');
                stockCell.textContent = book.stock;

                // Celda para las acciones (botón eliminar)
                const actionsCell = row.insertCell(6);
                actionsCell.classList.add('px-6', 'py-4', 'whitespace-nowrap', 'text-right', 'text-sm', 'font-medium');

                // Crear botón eliminar y asociar ID del libro
                const deleteButton = document.createElement('button');
                deleteButton.textContent = 'Eliminar';
                deleteButton.classList.add('btn-danger', 'btn-sm', 'ml-2'); // Clases para estilos
                deleteButton.dataset.id = book.id;

                // Asignar función al hacer clic para eliminar libro
                deleteButton.addEventListener('click', handleDeleteBook);

                // Agrego el botón a la celda de acciones
                actionsCell.appendChild(deleteButton);
            });

            showMessage('Libros cargados exitosamente', 'success');
        } else {
            // Si hubo error al cargar libros, muestro mensaje
            showMessage(books.message || 'Error al cargar los libros', 'error');

            // Si el error es por token inválido, forzar cierre de sesión
            if (response.status === 401) {
                 showMessage('Sesión expirada o inválida. Por favor, inicie sesión de nuevo.', 'error');
                 logoutBtn.click();
            }
        }
    } catch (error) {
        // Error de red o inesperado al hacer la petición
        showMessage('Error de red al cargar libros: ' + error.message, 'error');
    }
}

// --- Función para eliminar un libro ---

async function handleDeleteBook(event) {
    const bookId = event.target.dataset.id; // Obtengo el ID del libro desde el atributo data-id

    // Confirmo que el usuario quiere eliminar el libro
    if (!confirm(`¿Estás seguro de que quieres eliminar el libro con ID ${bookId}?`)) {
        return; // Si cancela, no hago nada
    }

    if (!authToken) {
        showMessage('No autenticado. Por favor, inicie sesión para eliminar libros.', 'error');
        return;
    }

    try {
        // Envío petición DELETE a la API con el token JWT
        const response = await fetch(`${API_BASE_URL}/books/${bookId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            }
        });

        const data = await response.json();

        if (response.ok) {
            // Eliminación exitosa: muestro mensaje y recargo la lista
            showMessage(data.message || 'Libro eliminado exitosamente.', 'success');
            loadBooks();
        } else {
            // Error al eliminar libro
            showMessage(data.message || 'Error al eliminar el libro.', 'error');

            // Si el token expiró o es inválido, cierro sesión
            if (response.status === 401) {
                showMessage('Sesión expirada o inválida. Por favor, inicie sesión de nuevo.', 'error');
                logoutBtn.click();
            }
        }
    } catch (error) {
        // Error en la conexión o inesperado
        showMessage('Error de conexión al intentar eliminar el libro.', 'error');
    }
}

// --- Crear nuevo libro ---

createBookForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    if (!authToken) {
        showMessage('No autenticado. Por favor, inicie sesión para crear libros.', 'error');
        return;
    }

    // Obtengo valores del formulario para el nuevo libro
    const title = document.getElementById('book-title').value;
    const author = document.getElementById('book-author').value;
    const genre = document.getElementById('book-genre').value;
    const price = parseFloat(document.getElementById('book-price').value);
    const stock = parseInt(document.getElementById('book-stock').value);

    // Validación básica para asegurar que los campos estén completos y correctos
    if (!title || !author || !genre || isNaN(price) || price <= 0 || isNaN(stock) || stock <= 0) {
        showMessage('Por favor, complete todos los campos y asegúrese de que Precio y Stock sean valores válidos y positivos.', 'error');
        return;
    }

    try {
        // Envío petición POST para crear libro con datos y token JWT
        const response = await fetch(`${API_BASE_URL}/books/create`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}`
            },
            body: JSON.stringify({ title, author, genre, price, stock })
        });

        const data = await response.json();

        if (response.ok) {
            // Libro creado correctamente, muestro mensaje, limpio formulario y recargo libros
            showMessage(data.message || 'Libro creado exitosamente', 'success');
            createBookForm.reset();
            loadBooks();
        } else {
            // Error al crear libro
            showMessage(data.message || 'Error al crear el libro', 'error');

            // Si token inválido o expirado, cierro sesión
            if (response.status === 401) {
                 showMessage('Sesión expirada o inválida. Por favor, inicie sesión de nuevo.', 'error');
                 logoutBtn.click();
            }
        }
    } catch (error) {
        // Error de red o inesperado
        showMessage('Error de red al crear libro: ' + error.message, 'error');
    }
});

// Permite recargar la lista de libros manualmente con el botón
refreshBooksBtn.addEventListener('click', loadBooks);

// --- Inicialización ---
// Al cargar la página, verifico el estado de autenticación para mostrar la vista correcta
document.addEventListener('DOMContentLoaded', toggleAppVisibility);
