# Guía de Instalación - MediLogic

## Índice
1. [Requisitos Previos](#requisitos-previos)
2. [Instalación de Componentes](#instalación-de-componentes)
3. [Configuración del Sistema](#configuración-del-sistema)
4. [Verificación de la Instalación](#verificación-de-la-instalación)
5. [Solución de Problemas](#solución-de-problemas)

## Requisitos Previos

### Software Necesario
1. **Go (Golang)**
   - Versión recomendada: 1.16 o superior
   - [Descargar Go](https://golang.org/dl/)

2. **SWI-Prolog**
   - Versión recomendada: 8.0 o superior
   - [Descargar SWI-Prolog](https://www.swi-prolog.org/download/stable)

### Requisitos de Sistema
- Sistema Operativo: Windows, Linux o MacOS
- RAM: 4GB mínimo recomendado
- Espacio en disco: 1GB mínimo

## Instalación de Componentes

### 1. Instalación de Go
1. Descargue el instalador de Go para su sistema operativo
2. Ejecute el instalador y siga las instrucciones
3. Verifique la instalación abriendo una terminal y ejecutando:
   ```
   go version
   ```

### 2. Instalación de SWI-Prolog
1. Descargue el instalador de SWI-Prolog
2. Ejecute el instalador
3. Asegúrese de marcar la opción "Add to PATH" durante la instalación
4. Verifique la instalación ejecutando:
   ```
   swipl --version
   ```

## Configuración del Sistema

### 1. Estructura de Directorios
Asegúrese de mantener la siguiente estructura de directorios:
```
P1/
├── Backend/
│   ├── go.mod
│   ├── main.go
│   └── run.bat
├── docs/
│   └── index.html
└── Prolog/
    └── base.pl
```

### 2. Configuración del Backend
El archivo `run.bat` incluye las siguientes variables de entorno:
```bat
set "PROLOG_BASE_PL=%~dp0..\Prolog\base.pl"
set "PORT=8080"
```

Puede modificar:
- `PORT`: Puerto donde se ejecutará el servidor
- `SWIPL_CMD`: Ruta al ejecutable de SWI-Prolog si no está en PATH

## Verificación de la Instalación

### 1. Iniciar el Servidor
1. Abra una terminal en la carpeta `Backend`
2. Ejecute:
   ```
   .\run.bat
   ```
3. Debería ver el mensaje:
   ```
   API escuchando en http://localhost:8080
   ```

### 2. Probar la Interfaz Web
1. Abra `docs/index.html` en su navegador
2. La interfaz debería cargarse correctamente
3. Pruebe agregar síntomas y realizar un diagnóstico

## Solución de Problemas

### Error: Puerto en uso
Si ve el error "puerto 8080 en uso":
1. Cambie el puerto en `run.bat`
2. O cierre la aplicación que esté usando ese puerto

### Error: SWI-Prolog no encontrado
Si el sistema no encuentra SWI-Prolog:
1. Asegúrese de que está instalado correctamente
2. Agregue la ruta completa en `run.bat`:
   ```bat
   set "SWIPL_CMD=C:\Program Files\swipl\bin\swipl.exe"
   ```

### Error: No se pueden cargar los catálogos
Si los catálogos no cargan:
1. Verifique que el servidor backend esté corriendo
2. Compruebe que la ruta a `base.pl` sea correcta en `run.bat`