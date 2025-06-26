const urlParams = new URLSearchParams(window.location.search);
const site = urlParams.get("site");
const token = urlParams.get("token");

if (!site || !token) {
  document.body.innerHTML = "<h2>🔒 Acceso denegado</h2>";
} else {
  cargarDatos(site, token);
}

async function cargarDatos(site, token) {
  try {
    const res = await fetch(`/stats?site=${site}&token=${token}`);
    if (!res.ok) {
      document.body.innerHTML = "<h2>🔒 Acceso denegado o inválido</h2>";
      return;
    }

    const data = await res.json();
    console.log("Datos recibidos:", data);

    const porModulo = Array.isArray(data.por_modulo) ? data.por_modulo : [];
    const porDia = Array.isArray(data.por_dia) ? data.por_dia : [];
    const navegadores = Array.isArray(data.navegadores) ? data.navegadores : [];
    const referencias = Array.isArray(data.referencias) ? data.referencias : [];
    const paginas = Array.isArray(data.paginas) ? data.paginas : [];
    const duracionMedia = data.duracion_media || 0;

    new Chart(document.getElementById("modulosChart"), {
      type: 'bar',
      data: {
        labels: porModulo.map(d => d.modulo),
        datasets: [{
          label: 'Usos por módulo',
          data: porModulo.map(d => d.total),
          backgroundColor: '#3cbf8e'
        }]
      }
    });

    new Chart(document.getElementById("visitasChart"), {
      type: 'line',
      data: {
        labels: porDia.map(d => d.dia),
        datasets: [{
          label: 'Visitas por día',
          data: porDia.map(d => d.total),
          borderColor: '#1a2330',
          fill: false
        }]
      }
    });

    // Navegadores
    new Chart(document.getElementById("navegadoresChart"), {
      type: 'doughnut',
      data: {
        labels: navegadores.map(d => d.navegador),
        datasets: [{
          label: 'Navegadores',
          data: navegadores.map(d => d.total),
          backgroundColor: ['#3cbf8e', '#1a2330', '#f5a623', '#e94e77', '#4a90e2']
        }]
      }
    });

    // Referencias
    new Chart(document.getElementById("referenciasChart"), {
      type: 'pie',
      data: {
        labels: referencias.map(d => d.referencia),
        datasets: [{
          label: 'Referencias',
          data: referencias.map(d => d.total),
          backgroundColor: ['#3cbf8e', '#1a2330', '#f5a623', '#e94e77', '#4a90e2']
        }]
      }
    });

    // Páginas más vistas
    new Chart(document.getElementById("paginasChart"), {
      type: 'bar',
      data: {
        labels: paginas.map(d => d.pagina),
        datasets: [{
          label: 'Páginas más vistas',
          data: paginas.map(d => d.total),
          backgroundColor: '#4a90e2'
        }]
      }
    });

    // Duración media de sesión
    const duracionDiv = document.getElementById("duracionMedia");
    if (duracionDiv) {
      const segundos = Math.round(duracionMedia / 1000);
      duracionDiv.textContent = segundos + 's';
    }

    // Dispositivos
    const dispositivos = Array.isArray(data.dispositivos) ? data.dispositivos : [];
    new Chart(document.getElementById("dispositivosChart"), {
      type: 'doughnut',
      data: {
        labels: dispositivos.map(d => d.dispositivo),
        datasets: [{
          label: 'Dispositivos',
          data: dispositivos.map(d => d.total),
          backgroundColor: ['#3cbf8e', '#1a2330', '#f5a623']
        }]
      }
    });

    // Países
    const paises = Array.isArray(data.paises) ? data.paises : [];
    new Chart(document.getElementById("paisesChart"), {
      type: 'pie',
      data: {
        labels: paises.map(d => d.pais),
        datasets: [{
          label: 'Países',
          data: paises.map(d => d.total),
          backgroundColor: ['#4a90e2', '#3cbf8e', '#f5a623', '#e94e77', '#1a2330']
        }]
      }
    });

    // Usuarios activos
    const usuariosActivos = typeof data.usuarios_activos === "number" ? data.usuarios_activos : 0;
    const usuariosActivosDiv = document.getElementById("usuariosActivos");
    if (usuariosActivosDiv) {
      usuariosActivosDiv.textContent = usuariosActivos;
    }

  } catch (e) {
    console.error("❌ Error en fetch:", e);
    document.body.innerHTML = "<h2>❌ Error de conexión</h2>";
  }
}