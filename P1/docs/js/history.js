// Historial temporal de diagnósticos
const diagnosticHistory = {
    history: [],
    
    // Agregar un nuevo diagnóstico al historial
    add(diagnosis) {
        this.history.push({
            ...diagnosis,
            timestamp: new Date()
        });
        this.save();
    },
    
    // Obtener todo el historial
    getAll() {
        return this.history;
    },
    
    // Guardar en localStorage
    save() {
        try {
            localStorage.setItem('diagnoseHistory', JSON.stringify(this.history));
        } catch (e) {
            console.error('Error guardando historial:', e);
        }
    },
    
    // Cargar desde localStorage
    load() {
        try {
            const saved = localStorage.getItem('diagnoseHistory');
            if (saved) {
                this.history = JSON.parse(saved);
            }
        } catch (e) {
            console.error('Error cargando historial:', e);
        }
    },
    
    // Limpiar historial
    clear() {
        this.history = [];
        this.save();
    }
};

// Cargar historial al inicio
diagnosticHistory.load();

// Modificar la función postDiag para guardar el diagnóstico
async function postDiag() {
    const status = document.getElementById('status');
    status.textContent = "Consultando...";
    try {
        const res = await fetch(API_DIAG, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(buildPayload()),
        });
        const txt = await res.text();
        if (!res.ok) throw new Error(txt);
        const diagnosis = JSON.parse(txt);
        render(diagnosis);
        status.textContent = "OK";
        
        // Guardar en historial
        diagnosticHistory.add(diagnosis);
        
        // Actualizar vista del historial si existe
        updateHistoryView();
    } catch (e) {
        status.textContent = "Error: " + e.message;
    }
}

// Función para actualizar la vista del historial
function updateHistoryView() {
    const historyContainer = document.getElementById('history');
    if (!historyContainer) return;

    const history = diagnosticHistory.getAll();
    historyContainer.innerHTML = "";

    history.forEach((diagnosis, index) => {
        const div = document.createElement('div');
        div.className = 'history-item';
        const date = new Date(diagnosis.timestamp);
        
        div.innerHTML = `
            <div class="history-header">
                <span>Diagnóstico #${index + 1} - ${date.toLocaleString()}</span>
                <button onclick="showHistoryDiagnosis(${index})">Ver</button>
            </div>
        `;
        
        historyContainer.appendChild(div);
    });
}

// Función para mostrar un diagnóstico del historial
function showHistoryDiagnosis(index) {
    const diagnosis = diagnosticHistory.getAll()[index];
    if (diagnosis) {
        render(diagnosis);
    }
}

// Agregar botones y contenedor del historial al HTML
document.getElementById('go').insertAdjacentHTML('afterend', `
    <button id="showHistory">Ver Historial</button>
    <button onclick="diagnosticHistory.clear(); updateHistoryView()">Limpiar Historial</button>
    <div id="history" class="history-container" style="display: none;"></div>
`);

// Manejar visualización del historial
document.getElementById('showHistory').addEventListener('click', function() {
    const historyContainer = document.getElementById('history');
    if (historyContainer.style.display === 'none') {
        historyContainer.style.display = 'block';
        updateHistoryView();
        this.textContent = 'Ocultar Historial';
    } else {
        historyContainer.style.display = 'none';
        this.textContent = 'Ver Historial';
    }
});