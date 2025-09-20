// Función para generar y descargar PDF
function downloadPDF(diagnosis) {
    // Importar jsPDF (asegúrate de incluir el script en el HTML)
    if (typeof jsPDF === 'undefined') {
        const script = document.createElement('script');
        script.src = 'https://cdnjs.cloudflare.com/ajax/libs/jspdf/2.5.1/jspdf.umd.min.js';
        script.onload = () => generatePDF(diagnosis);
        document.head.appendChild(script);
        return;
    }
    
    generatePDF(diagnosis);
}

function generatePDF(diagnosis) {
    const { jsPDF } = window.jspdf;
    const doc = new jsPDF();
    
    // Configurar fuente y tamaños
    doc.setFont('helvetica');
    doc.setFontSize(20);
    
    // Título
    doc.text('MediLogic - Informe de Diagnóstico', 105, 20, { align: 'center' });
    doc.setFontSize(12);
    doc.text(`Fecha: ${new Date().toLocaleString()}`, 20, 30);
    
    // Línea separadora
    doc.line(20, 35, 190, 35);
    
    let y = 45; // Posición Y inicial
    
    // Enfermedades diagnosticadas
    doc.setFontSize(16);
    doc.text('Diagnósticos Sugeridos:', 20, y);
    y += 10;
    
    doc.setFontSize(12);
    diagnosis.forEach((d, index) => {
        // Verificar si necesitamos una nueva página
        if (y > 250) {
            doc.addPage();
            y = 20;
        }
        
        doc.text(`${index + 1}. ${d.enfermedad} (${d.porcentaje}% de afinidad)`, 25, y);
        y += 7;
        
        if (d.coincidencias && d.coincidencias.length) {
            doc.setFontSize(10);
            doc.text(`Síntomas coincidentes: ${d.coincidencias.join(", ")}`, 30, y);
            doc.setFontSize(12);
            y += 7;
        }
        
        if (d.urgencia) {
            doc.text(`Nivel de urgencia: ${d.urgencia.nivel}`, 30, y);
            y += 7;
            if (d.urgencia.frase) {
                doc.text(`Recomendación: ${d.urgencia.frase}`, 30, y);
                y += 7;
            }
        }
        
        if (d.medicamento) {
            doc.text(`Medicamento sugerido: ${d.medicamento.principal}`, 30, y);
            y += 7;
            if (d.medicamento.alternativas && d.medicamento.alternativas.length) {
                doc.setFontSize(10);
                doc.text(`Alternativas: ${d.medicamento.alternativas.join(", ")}`, 30, y);
                doc.setFontSize(12);
                y += 7;
            }
        }
        
        if (d.reglas && d.reglas.length) {
            doc.setFontSize(10);
            doc.text(`Reglas aplicadas: ${d.reglas.join(", ")}`, 30, y);
            doc.setFontSize(12);
            y += 10;
        }
    });
    
    // Agregar pie de página
    const pageCount = doc.internal.getNumberOfPages();
    for (let i = 1; i <= pageCount; i++) {
        doc.setPage(i);
        doc.setFontSize(10);
        doc.text(`Página ${i} de ${pageCount}`, 105, 290, { align: 'center' });
        doc.text('Este es un diagnóstico preliminar. Consulte a un profesional de la salud.', 105, 285, { align: 'center' });
    }
    
    // Descargar el PDF
    doc.save('MediLogic-Diagnostico.pdf');
}

// Agregar botón de descarga al renderizar diagnóstico
const originalRender = render;
render = function(data) {
    originalRender(data);
    
    // Agregar botón de descarga si hay resultados
    if (Array.isArray(data) && data.length) {
        const out = document.getElementById('out');
        const downloadButton = document.createElement('button');
        downloadButton.textContent = 'Descargar PDF';
        downloadButton.onclick = () => downloadPDF(data);
        downloadButton.style.marginTop = '12px';
        out.appendChild(downloadButton);
    }
}