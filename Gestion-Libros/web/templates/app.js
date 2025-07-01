// web/templates/static/js/app.js

const API_BASE_URL = 'http://localhost:8080'; // 
let authToken = localStorage.getItem('jwtToken'); // Intentar obtener el token del almacenamiento local

const authSection = document.getElementById('auth-section');
const appSection = document.getElementById('app-section');
const messageContainer = document.getElementById('message-container');

const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const showLoginBtn = document.getElementById('show-login-btn');
const showRegisterBtn = document.getElementById('show-register-btn');
const logoutBtn = document.getElementById('logout-btn');

const createBookForm = document.getElementById('create-book-form');
const booksTableBody = document.getElementById('books-table-body');
const refreshBooksBtn = document.getElementById('refresh-books-btn');

// Función para mostrar mensajes al usuario
function showMessage(message, type = 'info') {
    messageContainer.textContent = message;
    messageContainer.classList.remove('hidden', 'bg-red-100', 'text-red-700', 'bg-green-100', 'text-green-700', 'bg-blue-100', 'text-blue-700');
    if (type === 'error') {
        messageContainer.classList.add('bg-red-100', 'text-red-700');
    } else if (type === 'success') {
        messageContainer.classList.add('bg-green-100', 'text-green-700');
    } else {
        messageContainer.classList.add('bg-blue-100', 'text-blue-700');
    }
    messageContainer.classList.remove('hidden');
    setTimeout(() => {
        messageContainer.classList.add('hidden');
    }, 5000); // Ocultar mensaje después de 5 segundos
}

// Función para alternar entre la vista de autenticación y la aplicación
function toggleAppVisibility() {
    if (authToken) {
        authSection.classList.add('hidden');
        appSection.classList.remove('hidden');
        loadBooks(); // Cargar libros al iniciar sesión
    } else {
        authSection.classList.remove('hidden');
        appSection.classList.add('hidden');
        loginForm.classList.remove('hidden'); // Por defecto, mostrar login al no estar autenticado
        registerForm.classList.add('hidden');
    }
}

// --- Manejo de Autenticación ---

showLoginBtn.addEventListener('click', () => {
    loginForm.classList.remove('hidden');
    registerForm.classList.add('hidden');
});

showRegisterBtn.addEventListener('click', () => {
    registerForm.classList.remove('hidden');
    loginForm.classList.add('hidden');
});

// Registrar Usuario
registerForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    try {
        const response = await fetch(`${API_BASE_URL}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password })
        });
        const data = await response.json();
        if (response.ok) {
            showMessage(data.message || 'Registro exitoso. ¡Ahora puedes iniciar sesión!', 'success');
            registerForm.reset();
            loginForm.classList.remove('hidden'); // Mostrar formulario de login
            registerForm.classList.add('hidden');
        } else {
            showMessage(data.message || 'Error en el registro', 'error');
        }
    } catch (error) {
        showMessage('Error de red al registrar: ' + error.message, 'error');
    }
});

// Iniciar Sesión
loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch(`${API_BASE_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        const data = await response.json();
        if (response.ok && data.token) {
            authToken = data.token;
            localStorage.setItem('jwtToken', authToken); // Guardar token
            showMessage('Inicio de sesión exitoso', 'success');
            loginForm.reset();
            toggleAppVisibility(); // Mostrar la app principal
        } else {
            showMessage(data.message || 'Error al iniciar sesión. Credenciales inválidas.', 'error');
        }
    } catch (error) {
        showMessage('Error de red al iniciar sesión: ' + error.message, 'error');
    }
});

// Cerrar Sesión
logoutBtn.addEventListener('click', () => {
    authToken = null;
    localStorage.removeItem('jwtToken'); // Eliminar token
    showMessage('Sesión cerrada correctamente', 'info');
    toggleAppVisibility(); // Volver a la vista de autenticación
    booksTableBody.innerHTML = ''; // Limpiar tabla de libros
});


// --- Manejo de Libros ---

// Cargar Libros
async function loadBooks() {
    if (!authToken) {
        showMessage('No autenticado. Por favor, inicie sesión para ver los libros.', 'error');
        return;
    }
    try {
        const response = await fetch(`${API_BASE_URL}/books`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${authToken}` // Enviar token JWT
            }
        });
        const books = await response.json();
        
        if (response.ok) {
            booksTableBody.innerHTML = ''; // Limpiar tabla
            if (books.length === 0) {
                booksTableBody.innerHTML = '<tr><td colspan="6" class="px-6 py-4 whitespace-nowrap text-center text-gray-500">No hay libros disponibles.</td></tr>';
                return;
            }
            books.forEach(book => {
                const row = `
                    <tr>
                        <td class="px-6 py-4 whitespace-nowrap">${book.id}</td>
                        <td class="px-6 py-4 whitespace-nowrap">${book.title}</td>
                        <td class="px-6 py-4 whitespace-nowrap">${book.author}</td>
                        <td class="px-6 py-4 whitespace-nowrap">${book.genre}</td>
                        <td class="px-6 py-4 whitespace-nowrap">$${book.price.toFixed(2)}</td>
                        <td class="px-6 py-4 whitespace-nowrap">${book.stock}</td>
                    </tr>
                `;
                booksTableBody.innerHTML += row;
            });
            showMessage('Libros cargados exitosamente', 'success');
        } else {
            showMessage(books.message || 'Error al cargar los libros', 'error');
            // Si el token es inválido/expirado, forzar logout
            if (response.status === 401) {
                 showMessage('Sesión expirada o inválida. Por favor, inicie sesión de nuevo.', 'error');
                 logoutBtn.click(); // Simular clic en logout
            }
        }
    } catch (error) {
        showMessage('Error de red al cargar libros: ' + error.message, 'error');
    }
}

// Crear Libro
createBookForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    if (!authToken) {
        showMessage('No autenticado. Por favor, inicie sesión para crear libros.', 'error');
        return;
    }

    const title = document.getElementById('book-title').value;
    const author = document.getElementById('book-author').value;
    const genre = document.getElementById('book-genre').value;
    const price = parseFloat(document.getElementById('book-price').value);
    const stock = parseInt(document.getElementById('book-stock').value);

    // Validación básica
    if (!title || !author || !genre || isNaN(price) || price <= 0 || isNaN(stock) || stock <= 0) {
        showMessage('Por favor, complete todos los campos y asegúrese de que Precio y Stock sean valores válidos y positivos.', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/books/create`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${authToken}` // Enviar token JWT
            },
            body: JSON.stringify({ title, author, genre, price, stock })
        });
        const data = await response.json();
        if (response.ok) {
            showMessage(data.message || 'Libro creado exitosamente', 'success');
            createBookForm.reset();
            loadBooks(); // Recargar la lista de libros
        } else {
            showMessage(data.message || 'Error al crear el libro', 'error');
            if (response.status === 401) {
                 showMessage('Sesión expirada o inválida. Por favor, inicie sesión de nuevo.', 'error');
                 logoutBtn.click();
            }
        }
    } catch (error) {
        showMessage('Error de red al crear libro: ' + error.message, 'error');
    }
});

// Actualizar Libros
refreshBooksBtn.addEventListener('click', loadBooks);

// --- Inicialización ---
document.addEventListener('DOMContentLoaded', toggleAppVisibility); // Verificar estado de autenticación al cargar la página