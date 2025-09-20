# Documentación de Pruebas - MediLogic

## Índice
1. [Casos de Prueba](#casos-de-prueba)
2. [Escenarios de Prueba](#escenarios-de-prueba)
3. [Resultados Esperados](#resultados-esperados)
4. [Validación de Reglas](#validación-de-reglas)

## Casos de Prueba

### 1. Diagnóstico de Gripe

#### Entrada
```json
{
    "sintomas": [
        {"nombre": "fiebre", "sev": "moderado"},
        {"nombre": "tos", "sev": "leve"},
        {"nombre": "congestion", "sev": "leve"}
    ],
    "alergias": ["alergia_ibuprofeno"],
    "cronicas": ["gastritis"],
    "topN": 3
}
```

#### Resultado Esperado
- Enfermedad principal: gripe
- Porcentaje alto de afinidad (>70%)
- Medicamento sugerido: paracetamol
- Nivel de urgencia: media

#### Reglas Activadas
- r_aff_01: Cálculo de afinidad base
- r_med_01: Selección de medicamento principal
- r_urg_03: Determinación de nivel de urgencia media

### 2. Caso de Emergencia

#### Entrada
```json
{
    "sintomas": [
        {"nombre": "dolor_pecho", "sev": "severo"},
        {"nombre": "dificultad_respirar", "sev": "moderado"}
    ],
    "alergias": [],
    "cronicas": [],
    "topN": 3
}
```

#### Resultado Esperado
- Nivel de urgencia: emergencia
- Mensaje de consulta médica inmediata
- Cualquier diagnóstico es secundario a la urgencia

#### Reglas Activadas
- r_urg_01: Detección de síntomas de emergencia
- r_aff_01: Cálculo de afinidad base

### 3. Prueba de Contraindicaciones

#### Entrada
```json
{
    "sintomas": [
        {"nombre": "dolor_cabeza", "sev": "severo"},
        {"nombre": "nauseas", "sev": "moderado"}
    ],
    "alergias": ["alergia_ibuprofeno"],
    "cronicas": ["gastritis"],
    "topN": 3
}
```

#### Resultado Esperado
- Diagnóstico: migraña
- Medicamentos: evitar ibuprofeno
- Sugerir alternativas seguras
- Nivel de urgencia: media/alta

#### Reglas Activadas
- r_med_02: Manejo de contraindicaciones
- r_urg_02: Evaluación de severidad alta

## Escenarios de Prueba

### 1. Validación de Entrada

- Síntomas inválidos
- Severidades fuera de rango
- Formatos de entrada incorrectos
- Caracteres especiales en nombres

### 2. Pruebas de Carga

- Múltiples síntomas simultáneos
- Listas largas de alergias/crónicas
- Solicitudes concurrentes al servidor

### 3. Casos Límite

- Sin síntomas ingresados
- Todas las severidades máximas
- Todas las severidades mínimas
- TopN mayor que diagnósticos posibles

## Resultados Esperados

### 1. Formato de Respuesta

- JSON válido
- Todos los campos requeridos presentes
- Valores dentro de rangos esperados
- Mensajes de error apropiados

### 2. Cálculos

- Porcentajes correctamente calculados
- Priorización adecuada de síntomas
- Ponderación correcta por severidad
- Ordenamiento apropiado de resultados

### 3. Seguridad

- Validación de entrada
- Manejo de errores
- Timeout en consultas largas
- Sanitización de datos

## Validación de Reglas

### 1. Reglas de Urgencia

- r_urg_01: Síntomas de emergencia
- r_urg_02: Condiciones de alta prioridad
- r_urg_03: Casos regulares

### 2. Reglas de Medicación

- r_med_01: Selección de medicamento principal
- r_med_02: Manejo de contraindicaciones

### 3. Reglas de Afinidad

- r_aff_01: Cálculo base de coincidencias
- Ponderación por severidad
- Normalización de porcentajes

### 4. Casos de Prueba Específicos

#### Test de Severidad
```prolog
?- parsear_sintomas([fiebre-severo, tos-leve], Pares, Atoms).
```
Resultado esperado:
- Pares = [fiebre-3, tos-1]
- Atoms = [fiebre, tos]

#### Test de Contraindicaciones
```prolog
?- medicamento_seguro(gripe, [alergia_paracetamol], [gastritis], Med, Regla).
```
Resultado esperado:
- Med = med(ibuprofeno, [])
- Regla = r_med_01