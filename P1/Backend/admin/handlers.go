package admin

import (
	"encoding/json"
	"net/http"
)

// Enfermedad representa una enfermedad en el sistema
type Enfermedad struct {
	Nombre             string   `json:"nombre"`
	Descripcion        string   `json:"descripcion"`
	SintomasAsociados  []string `json:"sintomas_asociados"`
	Medicamentos       []string `json:"medicamentos"`
	Contraindicaciones []string `json:"contraindicaciones"`
	Sistema            string   `json:"sistema"` // respiratorio, digestivo, etc.
	Tipo               string   `json:"tipo"`    // viral, crónico, etc.
}

// Sintoma representa un síntoma en el sistema
type Sintoma struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

// Medicamento representa un medicamento en el sistema
type Medicamento struct {
	Nombre             string   `json:"nombre"`
	Descripcion        string   `json:"descripcion"`
	TrataPara          []string `json:"trata_para"`
	Contraindicaciones []string `json:"contraindicaciones"`
}

// EnfermedadesHandler maneja las operaciones CRUD para enfermedades
func EnfermedadesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: Implementar lectura de enfermedades desde Prolog
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodPost:
		var enfermedad Enfermedad
		if err := json.NewDecoder(r.Body).Decode(&enfermedad); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		// TODO: Implementar guardado en Prolog
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodPut:
		var enfermedad Enfermedad
		if err := json.NewDecoder(r.Body).Decode(&enfermedad); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		// TODO: Implementar actualización en Prolog
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodDelete:
		// TODO: Implementar eliminación en Prolog
		w.WriteHeader(http.StatusNotImplemented)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// SintomasHandler maneja las operaciones CRUD para síntomas
func SintomasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: Implementar lectura de síntomas desde Prolog
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodPost:
		var sintoma Sintoma
		if err := json.NewDecoder(r.Body).Decode(&sintoma); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		// TODO: Implementar guardado en Prolog
		w.WriteHeader(http.StatusNotImplemented)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// MedicamentosHandler maneja las operaciones CRUD para medicamentos
func MedicamentosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: Implementar lectura de medicamentos desde Prolog
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodPost:
		var medicamento Medicamento
		if err := json.NewDecoder(r.Body).Decode(&medicamento); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		// TODO: Implementar guardado en Prolog
		w.WriteHeader(http.StatusNotImplemented)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
