@echo off
setlocal
set "PROLOG_BASE_PL=%~dp0..\Prolog\base.pl"
REM Ajusta si no tienes swipl en PATH:
REM set "SWIPL_CMD=C:\Program Files\swipl\bin\swipl.exe"
set "PORT=8080"
go run .
pause