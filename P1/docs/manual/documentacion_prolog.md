# Documentación Técnica del Sistema Prolog - MediLogic

## Índice
1. [Estructura General](#estructura-general)
2. [Módulos y Predicados](#módulos-y-predicados)
3. [Base de Conocimientos](#base-de-conocimientos)
4. [Sistema de Inferencia](#sistema-de-inferencia)
5. [Integración con el Backend](#integración-con-el-backend)

## Estructura General

El archivo `base.pl` está organizado en las siguientes secciones principales:

1. **Declaración del Módulo**
```prolog
:- module(base, [
    version/1,
    diagnosticar/5,
    diagnosticar_json/5,
    condiciones_json/1
]).
```

2. **Importación de Bibliotecas**
```prolog
:- use_module(library(lists)).
:- use_module(library(http/json)).
```

## Módulos y Predicados

### Predicados Principales

1. **diagnosticar/5**
```prolog
diagnosticar(SintomasIn, Alergias, Cronicas, TopN, Resultados)
```
- Entrada:
  - SintomasIn: Lista de síntomas con severidad
  - Alergias: Lista de alergias
  - Cronicas: Lista de enfermedades crónicas
  - TopN: Número de diagnósticos a retornar
- Salida:
  - Resultados: Lista de diagnósticos posibles

2. **diagnosticar_json/5**
```prolog
diagnosticar_json(Sintomas, Alergias, Cronicas, TopN, JsonAtom)
```
- Versión JSON del predicado diagnosticar
- Convierte los resultados a formato JSON

### Predicados de Soporte

1. **parsear_sintomas/3**
- Procesa la lista de síntomas y severidades
- Normaliza y valida los datos de entrada

2. **afinidad_w/4**
- Calcula el porcentaje de afinidad
- Considera pesos por severidad

3. **urgencia/4**
- Determina el nivel de urgencia
- Aplica reglas de priorización

4. **medicamento_seguro/5**
- Sugiere medicamentos seguros
- Evita contraindicaciones

## Base de Conocimientos

### Hechos del Dominio

1. **Síntomas**
```prolog
sintoma(fiebre).
sintoma(tos).
sintoma(dolor_cabeza).
% ...etc
```

2. **Enfermedades**
```prolog
enfermedad(gripe).
enfermedad(migrana).
enfermedad(alergia_estacional).
```

3. **Relaciones Síntoma-Enfermedad**
```prolog
sintoma_enfermedad(gripe, fiebre).
sintoma_enfermedad(gripe, tos).
% ...etc
```

4. **Medicamentos y Contraindicaciones**
```prolog
medicamento(paracetamol).
contraindicado(paracetamol, enfermedad_hepatica).
```

### Severidades y Pesos

```prolog
sev_valor(leve, 1).
sev_valor(moderado, 2).
sev_valor(severo, 3).
```

## Sistema de Inferencia

### Proceso de Diagnóstico

1. **Entrada de Datos**
- Validación de síntomas
- Normalización de severidades

2. **Cálculo de Afinidad**
- Ponderación por severidad
- Matching de síntomas

3. **Evaluación de Urgencia**
- Identificación de red flags
- Clasificación de urgencia

4. **Selección de Medicamentos**
- Verificación de contraindicaciones
- Búsqueda de alternativas

### Reglas de Inferencia

1. **Urgencia**
```prolog
urgencia(Pares, _Atoms, emergencia, r_urg_01) :-
    ( member(dolor_pecho-_, Pares)
    ; member(dificultad_respirar-_, Pares)
    ), !.
```

2. **Medicamentos Seguros**
```prolog
medicamento_seguro(Enf, Alergias, Cronicas, med(Principal, Alternativas), r_med_01) :-
    meds_seguros_para(Enf, Alergias, Cronicas, [Principal|Alternativas]),
    !.
```

## Integración con el Backend

### Formato JSON

1. **Entrada**
```json
{
    "sintomas": [{"nombre": "fiebre", "sev": "moderado"}],
    "alergias": ["alergia_ibuprofeno"],
    "cronicas": ["gastritis"],
    "topN": 3
}
```

2. **Salida**
```json
[{
    "enfermedad": "gripe",
    "score": 0.75,
    "porcentaje": 75,
    "coincidencias": ["fiebre", "tos"],
    "urgencia": {
        "nivel": "media",
        "frase": "Observacion y automanejo posibles"
    },
    "medicamento": {
        "principal": "paracetamol",
        "alternativas": []
    }
}]
```

### Predicados de Conversión

1. **resultado_a_dict/2**
- Convierte resultados Prolog a diccionarios
- Prepara la estructura para JSON

2. **condiciones_json/1**
- Genera catálogo de condiciones
- Separa alergias y crónicas