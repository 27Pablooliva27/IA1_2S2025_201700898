package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Sintoma struct {
	Nombre string `json:"nombre"`
	Sev    string `json:"sev"` // leve|moderado|severo (opcional)
}
type Req struct {
	Sintomas []Sintoma `json:"sintomas"`
	Alergias []string  `json:"alergias"`
	Cronicas []string  `json:"cronicas"`
	TopN     int       `json:"topN"`
}

var atom = regexp.MustCompile(`^[a-z_]+$`)
var sevOK = map[string]bool{"leve": true, "moderado": true, "severo": true}

func toPrologListAtoms(list []string) (string, error) {
	items := make([]string, 0, len(list))
	for _, s := range list {
		if !atom.MatchString(s) {
			return "", fmt.Errorf("valor invalido: %q (use minusculas y _)", s)
		}
		items = append(items, s)
	}
	return "[" + strings.Join(items, ",") + "]", nil
}
func toPrologSintomas(ss []Sintoma) (string, error) {
	items := make([]string, 0, len(ss))
	for _, s := range ss {
		if !atom.MatchString(s.Nombre) {
			return "", fmt.Errorf("sintoma invalido: %q", s.Nombre)
		}
		if s.Sev == "" {
			items = append(items, s.Nombre)
		} else {
			if !sevOK[s.Sev] {
				return "", fmt.Errorf("severidad invalida: %q", s.Sev)
			}
			items = append(items, fmt.Sprintf("%s-%s", s.Nombre, s.Sev))
		}
	}
	return "[" + strings.Join(items, ",") + "]", nil
}

func allowCORS(w http.ResponseWriter, r *http.Request) bool {
	origin := os.Getenv("ALLOW_ORIGIN")
	if origin == "" {
		origin = "*" // en dev; en prod pon tu URL de Pages
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Vary", "Origin")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func diagnosticarHandler(w http.ResponseWriter, r *http.Request) {
	if allowCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "solo POST", http.StatusMethodNotAllowed)
		return
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON invalido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.TopN <= 0 {
		req.TopN = 3
	}

	sintStr, err := toPrologSintomas(req.Sintomas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	alrgStr, err := toPrologListAtoms(req.Alergias)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cronStr, err := toPrologListAtoms(req.Cronicas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Config: ruta a swipl y a base.pl
	swipl := os.Getenv("SWIPL_CMD")
	if swipl == "" {
		swipl = "swipl"
	}
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		// por defecto: ../prolog/base.pl relativo al backend
		p, _ := filepath.Abs("../prolog/base.pl")
		basePl = p
	}
	plPath := strings.ReplaceAll(basePl, `\`, `/`)

	goal := fmt.Sprintf("['%s'], base:diagnosticar_json(%s,%s,%s,%d,J), writeln(J)",
		plPath, sintStr, alrgStr, cronStr, req.TopN,
	)

	// Timeout para el proceso Prolog (8s)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, swipl, "-q", "-g", goal, "-t", "halt")
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		msg := "Prolog error: " + err.Error()
		if e := strings.TrimSpace(errb.String()); e != "" {
			msg += "\n" + e
		}
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Validar salida JSON
	var parsed any
	if err := json.Unmarshal(out.Bytes(), &parsed); err != nil {
		http.Error(w, "Salida no-JSON desde Prolog:\n"+out.String(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Write(out.Bytes())
}

func catalogosHandler(w http.ResponseWriter, r *http.Request) {
	if allowCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "solo GET", http.StatusMethodNotAllowed)
		return
	}

	swipl := os.Getenv("SWIPL_CMD")
	if swipl == "" {
		swipl = "swipl"
	}

	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	plPath := strings.ReplaceAll(basePl, `\`, `/`)

	goal := fmt.Sprintf("['%s'], base:condiciones_json(J), writeln(J)", plPath)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, swipl, "-q", "-g", goal, "-t", "halt")
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		msg := "Prolog error: " + err.Error()
		if e := strings.TrimSpace(errb.String()); e != "" {
			msg += "\n" + e
		}
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	// validar JSON
	var parsed any
	if err := json.Unmarshal(out.Bytes(), &parsed); err != nil {
		http.Error(w, "Salida no-JSON desde Prolog:\n"+out.String(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out.Bytes())
}

func main() {
	http.HandleFunc("/diagnosticar", diagnosticarHandler)
	http.HandleFunc("/catalogos", catalogosHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API escuchando en http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
