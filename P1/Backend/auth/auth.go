package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Admin representa un usuario administrador
type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"` // Almacenar hash en producción
}

// Session representa una sesión de usuario
type Session struct {
	Token     string
	Username  string
	ExpiresAt time.Time
}

var (
	// En producción, esto debería estar en una base de datos
	admins = map[string]string{
		"admin": "admin123", // En producción usar hashes
	}

	// Almacén de sesiones activas
	sessions     = make(map[string]Session)
	sessionMutex sync.RWMutex
)

// hashPassword crea un hash SHA-256 de la contraseña
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// Login maneja la autenticación de administradores
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var admin Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Verificar credenciales
	storedPassword, exists := admins[admin.Username]
	if !exists || storedPassword != admin.Password { // En producción comparar hashes
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	// Crear sesión
	token := base64.StdEncoding.EncodeToString([]byte(time.Now().String() + admin.Username))
	session := Session{
		Token:     token,
		Username:  admin.Username,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Guardar sesión
	sessionMutex.Lock()
	sessions[token] = session
	sessionMutex.Unlock()

	// Devolver token
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// VerifyAuth es un middleware para verificar la autenticación
func VerifyAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Token no proporcionado", http.StatusUnauthorized)
			return
		}

		sessionMutex.RLock()
		session, exists := sessions[token]
		sessionMutex.RUnlock()

		if !exists || time.Now().After(session.ExpiresAt) {
			http.Error(w, "Sesión inválida o expirada", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// Logout cierra la sesión del usuario
func Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Token no proporcionado", http.StatusUnauthorized)
		return
	}

	sessionMutex.Lock()
	delete(sessions, token)
	sessionMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}
