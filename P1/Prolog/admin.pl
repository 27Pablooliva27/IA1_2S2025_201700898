% =============================
% Operaciones Administrativas
% =============================

:- module(admin, [
    % Sintomas
    listar_sintomas_json/1,
    agregar_sintoma/1,
    eliminar_sintoma/1,
    % Enfermedades
    listar_enfermedades_json/1,
    agregar_enfermedad/2,
    eliminar_enfermedad/1,
    % Medicamentos
    listar_medicamentos_json/1,
    agregar_medicamento/3,
    eliminar_medicamento/1
]).

:- use_module(library(http/json)).

:- dynamic sintoma/1.
:- dynamic enfermedad/1.
:- dynamic medicamento/1.
:- dynamic sintoma_enfermedad/2.
:- dynamic trata/2.
:- dynamic contraindicado/2.

% Gestión de Síntomas
agregar_sintoma(Nombre) :-
    \+ sintoma(Nombre),
    assertz(sintoma(Nombre)).

eliminar_sintoma(Nombre) :-
    retractall(sintoma(Nombre)),
    retractall(sintoma_enfermedad(_, Nombre)).

% Gestión de Enfermedades
agregar_enfermedad(Nombre, Sintomas) :-
    \+ enfermedad(Nombre),
    assertz(enfermedad(Nombre)),
    agregar_sintomas_enfermedad(Nombre, Sintomas).

agregar_sintomas_enfermedad(_, []).
agregar_sintomas_enfermedad(Enfermedad, [Sintoma|Resto]) :-
    assertz(sintoma_enfermedad(Enfermedad, Sintoma)),
    agregar_sintomas_enfermedad(Enfermedad, Resto).

eliminar_enfermedad(Nombre) :-
    retractall(enfermedad(Nombre)),
    retractall(sintoma_enfermedad(Nombre, _)),
    retractall(trata(_, Nombre)).

% Gestión de Medicamentos
agregar_medicamento(Nombre, Enfermedades, Contraindicaciones) :-
    \+ medicamento(Nombre),
    assertz(medicamento(Nombre)),
    agregar_tratamientos(Nombre, Enfermedades),
    agregar_contraindicaciones(Nombre, Contraindicaciones).

agregar_tratamientos(_, []).
agregar_tratamientos(Medicamento, [Enfermedad|Resto]) :-
    assertz(trata(Medicamento, Enfermedad)),
    agregar_tratamientos(Medicamento, Resto).

agregar_contraindicaciones(_, []).
agregar_contraindicaciones(Medicamento, [Contraindicacion|Resto]) :-
    assertz(contraindicado(Medicamento, Contraindicacion)),
    agregar_contraindicaciones(Medicamento, Resto).

eliminar_medicamento(Nombre) :-
    retractall(medicamento(Nombre)),
    retractall(trata(Nombre, _)),
    retractall(contraindicado(Nombre, _)).

% Consultas para la interfaz administrativa
listar_sintomas_json(JsonAtom) :-
    findall(_{nombre: S}, sintoma(S), List),
    atom_json_dict(JsonAtom, List, [as(atom)]).

listar_enfermedades_json(JsonAtom) :-
    findall(
        _{
            nombre: E,
            sintomas: Sintomas
        },
        (
            enfermedad(E),
            findall(S, sintoma_enfermedad(E, S), Sintomas)
        ),
        List
    ),
    atom_json_dict(JsonAtom, List, [as(atom)]).

listar_medicamentos_json(JsonAtom) :-
    findall(
        _{
            nombre: M,
            trata: Enfermedades,
            contraindicaciones: Contraindicaciones
        },
        (
            medicamento(M),
            findall(E, trata(M, E), Enfermedades),
            findall(C, contraindicado(M, C), Contraindicaciones)
        ),
        List
    ),
    atom_json_dict(JsonAtom, List, [as(atom)]).