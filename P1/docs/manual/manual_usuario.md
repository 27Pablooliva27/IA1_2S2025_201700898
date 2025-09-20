# Manual de Usuario - MediLogic

## Índice
1. [Introducción](#introducción)
2. [Requisitos del Sistema](#requisitos-del-sistema)
3. [Acceso al Sistema](#acceso-al-sistema)
4. [Uso del Sistema](#uso-del-sistema)
5. [Interpretación de Resultados](#interpretación-de-resultados)

## Introducción
MediLogic es un sistema experto de diagnóstico médico preliminar que permite a los usuarios ingresar sus síntomas, alergias y condiciones crónicas para recibir posibles diagnósticos y recomendaciones de medicamentos.

## Requisitos del Sistema
- Navegador web moderno (Chrome, Firefox, Edge, etc.)
- Conexión a internet
- Resolución de pantalla mínima recomendada: 1024x768

## Acceso al Sistema
1. Asegúrese que el servidor backend esté en ejecución
2. Abra el archivo `index.html` en su navegador web
3. El sistema debería mostrar la interfaz principal de diagnóstico

## Uso del Sistema

### Ingreso de Síntomas
1. Haga clic en el botón "+ Agregar síntoma"
2. Seleccione el síntoma de la lista desplegable
3. Especifique la severidad del síntoma:
   - Sin severidad
   - Leve
   - Moderado
   - Severo

### Especificación de Condiciones Médicas
1. En la sección "Alergias":
   - Escriba la alergia o selecciónela del menú desplegable
   - Presione Enter para agregarla
   - Use el botón "×" para eliminar una alergia

2. En la sección "Condiciones crónicas":
   - Escriba la condición o selecciónela del menú desplegable
   - Presione Enter para agregarla
   - Use el botón "×" para eliminar una condición

### Obtención de Diagnóstico
1. Seleccione el número de diagnósticos a mostrar en "Top N"
2. Presione el botón "Diagnosticar"
3. También puede usar el botón "Carga demo" para ver un ejemplo

### Limpiar Formulario
- Use el botón "Limpiar" para reiniciar todos los campos

## Interpretación de Resultados

### Diagnósticos
Por cada posible enfermedad, se muestra:
- Nombre de la enfermedad
- Porcentaje de afinidad
- Barra visual de afinidad
- Síntomas coincidentes

### Niveles de Urgencia
El sistema clasifica la urgencia en:
- **Emergencia** (rojo): Requiere atención médica inmediata
- **Alta** (amarillo): Atención prioritaria en las próximas horas
- **Media**: Observación y automanejo posibles
- **Baja**: Autocuidado y monitoreo

### Medicamentos
Para cada diagnóstico se muestra:
- Medicamento principal recomendado
- Medicamentos alternativos (si están disponibles)
- No se recomiendan medicamentos contraindicados según sus alergias y condiciones crónicas

### Reglas Activadas
El sistema muestra las reglas que se utilizaron para llegar al diagnóstico, proporcionando transparencia en el proceso de inferencia.