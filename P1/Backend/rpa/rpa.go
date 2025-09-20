package rpa

import (
	"bufio"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"strings"

	"backend/admin"
)

// Config contiene la configuración para el envío de correos
type Config struct {
	SmtpServer   string   `json:"smtp_server"`
	SmtpPort     int      `json:"smtp_port"`
	SmtpUser     string   `json:"smtp_user"`
	SmtpPassword string   `json:"smtp_password"`
	AdminEmails  []string `json:"admin_emails"`
}

// LoadDisease carga una enfermedad desde un archivo de texto
func LoadDisease(reader *bufio.Reader) (*admin.Enfermedad, error) {
	enfermedad := &admin.Enfermedad{}

	// Leer nombre
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error leyendo nombre: %v", err)
	}
	enfermedad.Nombre = strings.TrimSpace(line)

	// Leer descripción
	line, err = reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error leyendo descripción: %v", err)
	}
	enfermedad.Descripcion = strings.TrimSpace(line)

	// Leer síntomas asociados
	line, err = reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error leyendo síntomas: %v", err)
	}
	sintomas := strings.Split(strings.TrimSpace(line), ",")
	for i := range sintomas {
		sintomas[i] = strings.TrimSpace(sintomas[i])
	}
	enfermedad.SintomasAsociados = sintomas

	// Leer medicamentos
	line, err = reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error leyendo medicamentos: %v", err)
	}
	medicamentos := strings.Split(strings.TrimSpace(line), ",")
	for i := range medicamentos {
		medicamentos[i] = strings.TrimSpace(medicamentos[i])
	}
	enfermedad.Medicamentos = medicamentos

	// Leer contraindicaciones
	line, err = reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error leyendo contraindicaciones: %v", err)
	}
	contraindicaciones := strings.Split(strings.TrimSpace(line), ",")
	for i := range contraindicaciones {
		contraindicaciones[i] = strings.TrimSpace(contraindicaciones[i])
	}
	enfermedad.Contraindicaciones = contraindicaciones

	// Leer sistema
	line, err = reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			if len(line) > 0 {
				enfermedad.Sistema = strings.TrimSpace(line)
			}
			return enfermedad, nil
		}
		return nil, fmt.Errorf("error leyendo sistema: %v", err)
	}
	enfermedad.Sistema = strings.TrimSpace(line)

	// Leer tipo
	line, err = reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			if len(line) > 0 {
				enfermedad.Tipo = strings.TrimSpace(line)
			}
			return enfermedad, nil
		}
		return nil, fmt.Errorf("error leyendo tipo: %v", err)
	}
	enfermedad.Tipo = strings.TrimSpace(line)

	return enfermedad, nil
}

// ProcessFile procesa un archivo de texto con información de enfermedades
func ProcessFile(filePath string) ([]admin.Enfermedad, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var enfermedades []admin.Enfermedad

	for {
		enfermedad, err := LoadDisease(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		enfermedades = append(enfermedades, *enfermedad)
	}

	return enfermedades, nil
}

// SendReport envía un informe por correo electrónico a los administradores
func SendReport(config Config, enfermedades []admin.Enfermedad) error {
	// Crear el contenido del informe
	var report strings.Builder
	report.WriteString("Informe de carga de enfermedades\n\n")
	report.WriteString("Enfermedades cargadas:\n")

	for _, e := range enfermedades {
		report.WriteString(fmt.Sprintf("\nNombre: %s\n", e.Nombre))
		report.WriteString(fmt.Sprintf("Sistema: %s\n", e.Sistema))
		report.WriteString(fmt.Sprintf("Tipo: %s\n", e.Tipo))
		report.WriteString("------------------\n")
	}

	// Configurar correo
	auth := smtp.PlainAuth("", config.SmtpUser, config.SmtpPassword, config.SmtpServer)
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Informe de carga de enfermedades\r\n"+
		"\r\n"+
		"%s\r\n", strings.Join(config.AdminEmails, ","), report.String()))

	// Enviar correo
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", config.SmtpServer, config.SmtpPort),
		auth,
		config.SmtpUser,
		config.AdminEmails,
		msg,
	)

	if err != nil {
		return fmt.Errorf("error enviando correo: %v", err)
	}

	return nil
}
