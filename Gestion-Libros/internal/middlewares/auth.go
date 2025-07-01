package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Tipo personalizado para las claves del contexto (evita colisiones con otras claves)
type contextKey string

// Clave con la que guardaré el user_id extraído del token en el contexto de la petición
const UserIDKey contextKey = "user_id"

// JWTAuthMiddleware valida el token JWT incluido en la cabecera Authorization.
// Si el token es válido, extrae el user_id y lo guarda en el contexto.
// Si no es válido, devuelve un error 401 (no autorizado).
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Leo la cabecera Authorization
		authHeader := r.Header.Get("Authorization")

		// Verifico que el token esté presente y comience con "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Token no proporcionado", http.StatusUnauthorized)
			return
		}

		// Extraigo el token quitando el prefijo "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Clave secreta para verificar el token
		secret := []byte(os.Getenv("JWT_SECRET"))

		// Parsea y verifica el token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Aseguro que se use un método de firma HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})

		// Si el token no es válido o hubo error, rechazo la petición
		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Obtengo los claims del token (datos asociados al usuario)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Claims inválidos", http.StatusUnauthorized)
			return
		}

		// Extraigo el user_id desde los claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "user_id no encontrado en token", http.StatusUnauthorized)
			return
		}

		// Guardo el user_id en el contexto para que esté disponible en los handlers
		ctx := context.WithValue(r.Context(), UserIDKey, uint(userID))

		// Paso la petición al siguiente handler con el contexto actualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
