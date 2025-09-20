package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Consulta síntomas desde Prolog
func consultarSintomas() ([]SintomaAdmin, error) {
	// Obtener ruta al archivo base.pl
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	adminPl := strings.ReplaceAll(filepath.Join(filepath.Dir(basePl), "admin.pl"), `\`, `/`)
	goal := fmt.Sprintf("['%s'], ['%s'], listar_sintomas_json(J), writeln(J)", basePl, adminPl)
	output, err := ejecutarProlog(goal)
	if err != nil {
		return nil, err
	}

	var sintomas []SintomaAdmin
	if err := json.Unmarshal([]byte(output), &sintomas); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v - output: %s", err, output)
	}
	return sintomas, nil
}

// Agrega un síntoma a Prolog
func agregarSintoma(s *SintomaAdmin) error {
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	adminPl := strings.ReplaceAll(filepath.Join(filepath.Dir(basePl), "admin.pl"), `\`, `/`)
	goal := fmt.Sprintf("['%s'], ['%s'], agregar_sintoma('%s')", basePl, adminPl, s.Nombre)
	if s.SeveridadDefault != "" {
		goal = fmt.Sprintf("['%s'], ['%s'], agregar_sintoma('%s', '%s')", 
			basePl, adminPl, s.Nombre, s.SeveridadDefault)
	}
	_, err := ejecutarProlog(goal)
	return err
}

// Elimina un síntoma de Prolog
func eliminarSintoma(nombre string) error {
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	adminPl := strings.ReplaceAll(filepath.Join(filepath.Dir(basePl), "admin.pl"), `\`, `/`)
	goal := fmt.Sprintf("['%s'], ['%s'], eliminar_sintoma('%s')", basePl, adminPl, nombre)
	_, err := ejecutarProlog(goal)
	return err
}

// Consulta enfermedades desde Prolog
func consultarEnfermedades() ([]EnfermedadAdmin, error) {
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	adminPl := strings.ReplaceAll(filepath.Join(filepath.Dir(basePl), "admin.pl"), `\`, `/`)
	goal := fmt.Sprintf("['%s'], ['%s'], listar_enfermedades_json(J), writeln(J)", basePl, adminPl)
	output, err := ejecutarProlog(goal)
	if err != nil {
		return nil, err
	}

	var enfermedades []EnfermedadAdmin
	if err := json.Unmarshal([]byte(output), &enfermedades); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v - output: %s", err, output)
	}
	return enfermedades, nil
}

// Agrega una enfermedad a Prolog
func agregarEnfermedad(e *EnfermedadAdmin) error {
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	adminPl := strings.ReplaceAll(filepath.Join(filepath.Dir(basePl), "admin.pl"), `\`, `/`)
	sintomas := "[" + strings.Join(quoted(e.Sintomas), ",") + "]"
	goal := fmt.Sprintf("['%s'], ['%s'], agregar_enfermedad('%s', %s)", basePl, adminPl, e.Nombre, sintomas)
	_, err := ejecutarProlog(goal)
	return err
}

// Elimina una enfermedad de Prolog
func eliminarEnfermedad(nombre string) error {
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	adminPl := strings.ReplaceAll(filepath.Join(filepath.Dir(basePl), "admin.pl"), `\`, `/`)
	goal := fmt.Sprintf("['%s'], ['%s'], eliminar_enfermedad('%s')", basePl, adminPl, nombre)
	_, err := ejecutarProlog(goal)
	return err
}

// Consulta medicamentos desde Prolog
func consultarMedicamentos() ([]MedicamentoAdmin, error) {
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	adminPl := strings.ReplaceAll(filepath.Join(filepath.Dir(basePl), "admin.pl"), `\`, `/`)
	goal := fmt.Sprintf("['%s'], ['%s'], listar_medicamentos_json(J), writeln(J)", basePl, adminPl)
	output, err := ejecutarProlog(goal)
	if err != nil {
		return nil, err
	}

	var medicamentos []MedicamentoAdmin
	if err := json.Unmarshal([]byte(output), &medicamentos); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v - output: %s", err, output)
	}
	return medicamentos, nil
}

// Agrega un medicamento a Prolog
func agregarMedicamento(m *MedicamentoAdmin) error {
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	adminPl := strings.ReplaceAll(filepath.Join(filepath.Dir(basePl), "admin.pl"), `\`, `/`)
	enfermedades := "[" + strings.Join(quoted(m.Trata), ",") + "]"
	contraindicaciones := "[" + strings.Join(quoted(m.Contraindicaciones), ",") + "]"
	goal := fmt.Sprintf("['%s'], ['%s'], agregar_medicamento('%s', %s, %s)",
		basePl, adminPl, m.Nombre, enfermedades, contraindicaciones)
	_, err := ejecutarProlog(goal)
	return err
}

// Elimina un medicamento de Prolog
func eliminarMedicamento(nombre string) error {
	basePl := os.Getenv("PROLOG_BASE_PL")
	if basePl == "" {
		p, _ := filepath.Abs("../Prolog/base.pl")
		basePl = p
	}
	adminPl := strings.ReplaceAll(filepath.Join(filepath.Dir(basePl), "admin.pl"), `\`, `/`)
	goal := fmt.Sprintf("['%s'], ['%s'], eliminar_medicamento('%s')", basePl, adminPl, nombre)
	_, err := ejecutarProlog(goal)
	return err
}

// Utilidades

// Ejecuta un goal en Prolog
func ejecutarProlog(goal string) (string, error) {
	swipl := os.Getenv("SWIPL_CMD")
	if swipl == "" {
		swipl = "swipl"
	}

	fmt.Println("\n=== Ejecutando consulta Prolog ===")
	fmt.Printf("Comando: %s -q -g '%s' -t halt\n", swipl, goal)
	fmt.Print("¿Desea continuar con la iteración? [S/n]: ")
	
	// Leer respuesta del usuario
	var respuesta string
	fmt.Scanln(&respuesta)
	respuesta = strings.ToLower(strings.TrimSpace(respuesta))
	if respuesta == "n" || respuesta == "no" {
		return "", fmt.Errorf("operación cancelada por el usuario")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, swipl, "-q", "-g", goal, "-t", "halt")
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))
	
	fmt.Println("\n=== Resultado de la consulta ===")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Salida: %s\n", outputStr)
		fmt.Println("===============================")
		return "", fmt.Errorf("error executing Prolog: %v - output: %s", err, outputStr)
	}

	fmt.Printf("Salida: %s\n", outputStr)
	fmt.Println("===============================")
	return outputStr, nil
}

// Añade comillas a cada string en un slice
func quoted(strs []string) []string {
	quoted := make([]string, len(strs))
	for i, s := range strs {
		quoted[i] = "'" + s + "'"
	}
	return quoted
}