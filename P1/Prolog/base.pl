% ================================
% base.pl  (Tarea 1 - severidades + porcentaje + urgencia + alternativas)
% ================================

:- module(base, [
    version/1,
    diagnosticar/5
]).

:- use_module(library(lists)).  % member/2, append/2, predsort/3, list_to_set/2
:- use_module(library(http/json)).
version('Tarea1').

% ---------- Hechos de dominio ----------
sintoma(fiebre).
sintoma(tos).
sintoma(dolor_cabeza).
sintoma(congestion).
sintoma(lagrimeo).
sintoma(nauseas).
sintoma(fotofobia).
sintoma(dolor_pecho).          % red flag -> emergencia
sintoma(dificultad_respirar).  % red flag -> emergencia
sintoma(rigidez_nuca).         % red flag -> alta
sintoma(fiebre_alta).          % red flag -> alta (≈ fiebre-severo)

enfermedad(gripe).
enfermedad(migrana).
enfermedad(alergia_estacional).

sintoma_enfermedad(gripe, fiebre).
sintoma_enfermedad(gripe, tos).
sintoma_enfermedad(gripe, congestion).

sintoma_enfermedad(migrana, dolor_cabeza).
sintoma_enfermedad(migrana, nauseas).
sintoma_enfermedad(migrana, fotofobia).

sintoma_enfermedad(alergia_estacional, congestion).
sintoma_enfermedad(alergia_estacional, lagrimeo).

medicamento(paracetamol).
medicamento(ibuprofeno).
medicamento(loratadina).

trata(paracetamol, gripe).
trata(ibuprofeno, migrana).
trata(loratadina, alergia_estacional).

contraindicado(paracetamol, enfermedad_hepatica).
contraindicado(ibuprofeno, gastritis).
contraindicado(ibuprofeno, alergia_ibuprofeno).
contraindicado(loratadina, glaucoma).

% ---------- Severidades ----------
sev_valor(leve, 1).
sev_valor(moderado, 2).
sev_valor(severo, 3).
sev_valor(baja, 1).
sev_valor(media, 2).
sev_valor(alta, 3).
sev_valor(1, 1).
sev_valor(2, 2).
sev_valor(3, 3).

peso_maximo(3).

% ---------- Parsing de síntomas ----------
% parsear_sintomas(+In, -Pares, -Atomos)
parsear_sintomas(ListIn, ParesOut, AtomosOut) :-
    parsear_crudo(ListIn, ParesCrudos),
    consolidar_max(ParesCrudos, ParesOut),
    pares_a_sintomas(ParesOut, AtomosOut).

parsear_crudo([], []).
parsear_crudo([S|T], [S-1|R]) :-
    atom(S), !, parsear_crudo(T, R).
parsear_crudo([S-Sev|T], [S-W|R]) :- !,
    ( sev_valor(Sev, W) -> true ; W = 1 ),
    parsear_crudo(T, R).

consolidar_max([], []).
consolidar_max([S-W|T], R) :-
    consolidar_max(T, Parcial),
    insertar_par_max(S-W, Parcial, R).

insertar_par_max(S-W, [], [S-W]).
insertar_par_max(S-W, [S-W0|T], [S-W1|T]) :-  % mismo síntoma: conservar el mayor peso
    (W > W0 -> W1 = W ; W1 = W0), !.
insertar_par_max(SW, [X|T], [X|R]) :-
    insertar_par_max(SW, T, R).

pares_a_sintomas([], []).
pares_a_sintomas([S-_|T], [S|R]) :- pares_a_sintomas(T, R).

sintomas_de_enfermedad(Enf, Lista) :-
    findall(S, sintoma_enfermedad(Enf, S), Lista).

% ---------- Afinidad ponderada ----------
% suma_pesos_match(+SintomasEnf, +ParesPaciente, -Suma, -Coinc)
suma_pesos_match([], _, 0, []).
suma_pesos_match([E|T], Pares, Suma, Coinc) :-
    % Primero resuelve el resto, obteniendo S1 y C1
    suma_pesos_match(T, Pares, S1, C1),
    % Luego incorpora E si coincide
    ( member(E-W, Pares) ->
        Suma is W + S1,
        Coinc = [E|C1]
    ;   Suma = S1,
        Coinc = C1
    ).

% Score = (sumatoria pesos de match) / (3 * #sintomas_enf)
afinidad_w(Enf, ParesPac, Score, Coinc) :-
    enfermedad(Enf),
    sintomas_de_enfermedad(Enf, SE),
    SE \= [],
    suma_pesos_match(SE, ParesPac, SumW, Coinc),
    length(SE, Tot),
    peso_maximo(Max),
    Den is Max * Tot,
    ( Den =:= 0 -> Score = 0.0 ; Score is SumW / Den ).

% ---------- Urgencia + frase ----------
urgencia(Pares, _Atoms, emergencia, r_urg_01) :-
    ( member(dolor_pecho-_, Pares)
    ; member(dificultad_respirar-_, Pares)
    ), !.

urgencia(Pares, _Atoms, alta, r_urg_02) :-
    ( member(rigidez_nuca-_, Pares)
    ; member(fiebre_alta-_, Pares)
    ; member(fiebre-W, Pares), W >= 3
    ), !.

urgencia(Pares, _Atoms, media, r_urg_03) :- Pares \= [], !.
urgencia(_Pares, _Atoms, baja,  r_urg_03).

urgencia_texto(emergencia, 'Consulta medica inmediata sugerida').
urgencia_texto(alta,       'Atencion prioritaria en las proximas horas').
urgencia_texto(media,      'Observacion y automanejo posibles').
urgencia_texto(baja,       'Autocuidado y monitoreo').

% ---------- Medicamento principal + alternativas ----------
meds_seguros_para(Enf, Alergias, Cronicas, SegurosOrdenados) :-
    findall(M,
        ( trata(M, Enf),
          \+ contraindicado_para_paciente(M, Alergias, Cronicas)
        ),
        Seguros),
    sort(Seguros, SegurosOrdenados).

medicamento_seguro(Enf, Alergias, Cronicas,
                   med(Principal, Alternativas), r_med_01) :-
    meds_seguros_para(Enf, Alergias, Cronicas, [Principal|Alternativas]),
    !.
medicamento_seguro(Enf, Alergias, Cronicas,
                   med(ninguno, []), r_med_02) :-
    meds_seguros_para(Enf, Alergias, Cronicas, []).

contraindicado_para_paciente(Med, Alergias, Cronicas) :-
    contraindicado(Med, Cond),
    ( member(Cond, Alergias)
    ; member(Cond, Cronicas)
    ).

% ---------- Utiles ----------
pct(Score, Pct) :- Pct is round(Score * 100).

tomar_k(_, 0, []) :- !.
tomar_k([], _, []) :- !.
tomar_k([H|T], K, [H|R]) :-
    K > 0, K1 is K - 1,
    tomar_k(T, K1, R).

% Comparador para ordenar por Score DESC (para predsort/3)
% A y B son de la forma Score-Enf-Coinc
cmp_score_desc(Orden, A, B) :-
    A = SA-_-_,
    B = SB-_-_,
    compare(Orden0, SB, SA),  % invertimos para descendente
    ( Orden0 = '=' -> Orden = '=' ; Orden = Orden0 ).

% Construye item de salida
armar_resultado(Alergias, Cronicas, Urg, ReglaUrg, Score-Enf-Coinc,
                resultado(Enf, Score, Pct, Coinc,
                          urg(Urg, Frase, ReglaUrg),
                          med(MedPrincipal, Alternativas, ReglaMed),
                          Reglas)) :-
    ReglasBase = [r_aff_01],
    pct(Score, Pct),
    urgencia_texto(Urg, Frase),
    medicamento_seguro(Enf, Alergias, Cronicas,
                       med(MedPrincipal, Alternativas), ReglaMed),
    append([ReglasBase, [ReglaMed, ReglaUrg]], Flat),
    list_to_set(Flat, Reglas).

% ---------- Consulta principal ----------
diagnosticar(SintomasIn, Alergias, Cronicas, TopN, Resultados) :-
    parsear_sintomas(SintomasIn, Pares, Atomos),
    urgencia(Pares, Atomos, Urg, ReglaUrg),
    findall(Score-Enf-Coinc, afinidad_w(Enf, Pares, Score, Coinc), ParesEnf),
    predsort(cmp_score_desc, ParesEnf, Ordenados),   % Score DESC
    tomar_k(Ordenados, TopN, Top),
    maplist(armar_resultado(Alergias, Cronicas, Urg, ReglaUrg), Top, Resultados).




resultado_a_dict(
  resultado(Enf,Score,Pct,Coinc,urg(Nivel,Frase,RegUrg),med(MP,Alt,RegMed),Reglas),
  _{
    enfermedad: Enf,
    score: Score,
    porcentaje: Pct,
    coincidencias: Coinc,
    urgencia: _{nivel: Nivel, frase: Frase, regla: RegUrg},
    medicamento: _{principal: MP, alternativas: Alt, regla: RegMed},
    reglas: Reglas
  }
).

diagnosticar_json(Sintomas,Alergias,Cronicas,TopN,JsonAtom) :-
  diagnosticar(Sintomas,Alergias,Cronicas,TopN,Rs),
  maplist(resultado_a_dict, Rs, Dicts),
  atom_json_dict(JsonAtom, Dicts, [as(atom)]).