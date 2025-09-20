package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Estructuras para la administración
type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"` // En producción usar hash
}

type Session struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type SintomaAdmin struct {
	Nombre           string `json:"nombre"`
	SeveridadDefault string `json:"severidad_default,omitempty"`
}

type EnfermedadAdmin struct {
	Nombre   string   `json:"nombre"`
	Sintomas []string `json:"sintomas"`
}

type MedicamentoAdmin struct {
	Nombre             string   `json:"nombre"`
	Trata              []string `json:"trata"`
	Contraindicaciones []string `json:"contraindicaciones"`
}

type RegistroHistorico struct {
	Fecha       time.Time `json:"fecha"`
	Sintomas    []Sintoma `json:"sintomas"`
	Diagnostico struct {
		Enfermedad string `json:"enfermedad"`
		Porcentaje int    `json:"porcentaje"`
	} `json:"diagnostico"`
	Medicamento string `json:"medicamento"`
}

// Variables globales para la administración
var (
	admins    = map[string]Admin{"admin": {Username: "admin", Password: "admin"}} // En producción usar base de datos
	sessions  = map[string]Session{}
	historico = []RegistroHistorico{}
	mutex     sync.RWMutex
)

// Middleware de autenticación
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		mutex.RLock()
		session, valid := sessions[token]
		mutex.RUnlock()

		if !valid || time.Now().After(session.ExpiresAt) {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// Handler de login
func loginHandler(w http.ResponseWriter, r *http.Request) {

	var creds Admin
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	admin, valid := admins[creds.Username]
	if !valid || admin.Password != creds.Password {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	token := generateToken()
	session := Session{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	mutex.Lock()
	sessions[token] = session
	mutex.Unlock()

	json.NewEncoder(w).Encode(session)
}

// Handlers CRUD para síntomas
func adminSintomasHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		sintomas, err := consultarSintomas()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(sintomas)

	case http.MethodPost:
		var sintoma SintomaAdmin
		if err := json.NewDecoder(r.Body).Decode(&sintoma); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := agregarSintoma(&sintoma); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		nombre := r.URL.Query().Get("nombre")
		if nombre == "" {
			http.Error(w, "Nombre requerido", http.StatusBadRequest)
			return
		}

		if err := eliminarSintoma(nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Handlers CRUD para enfermedades
func adminEnfermedadesHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		enfermedades, err := consultarEnfermedades()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(enfermedades)

	case http.MethodPost:
		var enfermedad EnfermedadAdmin
		if err := json.NewDecoder(r.Body).Decode(&enfermedad); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := agregarEnfermedad(&enfermedad); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		nombre := r.URL.Query().Get("nombre")
		if nombre == "" {
			http.Error(w, "Nombre requerido", http.StatusBadRequest)
			return
		}

		if err := eliminarEnfermedad(nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Handlers CRUD para medicamentos
func adminMedicamentosHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		medicamentos, err := consultarMedicamentos()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(medicamentos)

	case http.MethodPost:
		var medicamento MedicamentoAdmin
		if err := json.NewDecoder(r.Body).Decode(&medicamento); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := agregarMedicamento(&medicamento); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		nombre := r.URL.Query().Get("nombre")
		if nombre == "" {
			http.Error(w, "Nombre requerido", http.StatusBadRequest)
			return
		}

		if err := eliminarMedicamento(nombre); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// Handler para histórico
func historicoHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	desde := r.URL.Query().Get("desde")
	hasta := r.URL.Query().Get("hasta")

	mutex.RLock()
	defer mutex.RUnlock()

	// Si no hay filtros, devolver todo el histórico
	if desde == "" && hasta == "" {
		json.NewEncoder(w).Encode(historico)
		return
	}

	// Convertir fechas de string a time.Time
	var desdeTime, hastaTime time.Time
	var err error

	if desde != "" {
		desdeTime, err = time.Parse("2006-01-02", desde)
		if err != nil {
			http.Error(w, "Formato de fecha 'desde' inválido", http.StatusBadRequest)
			return
		}
	}

	if hasta != "" {
		hastaTime, err = time.Parse("2006-01-02", hasta)
		if err != nil {
			http.Error(w, "Formato de fecha 'hasta' inválido", http.StatusBadRequest)
			return
		}
	}

	// Filtrar registros por fecha
	filtrado := make([]RegistroHistorico, 0)
	for _, registro := range historico {
		if desde != "" && registro.Fecha.Before(desdeTime) {
			continue
		}
		if hasta != "" && registro.Fecha.After(hastaTime) {
			continue
		}
		filtrado = append(filtrado, registro)
	}

	json.NewEncoder(w).Encode(filtrado)
}

// Utilidades
func generateToken() string {
	// En producción usar un método más seguro
	return "admin-token"
}
