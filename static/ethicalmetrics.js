(function () {
  if (navigator.doNotTrack === "1") return;

  const API_URL = "/track";
  const MODULE = document.body?.dataset?.modulo || "inicio";
  const LOAD_TIME = Math.round(performance.now());

  // Evento automático de visita
  send({
    evento: "visita",
    modulo: MODULE,
    duracion_ms: LOAD_TIME
  });

  // Exponer función global
  window.ethical = {
    track: send
  };

  // Función para enviar cualquier evento personalizado
  function send(data) {
    try {
      if (navigator.doNotTrack === "1") return;

      const payload = JSON.stringify(Object.assign({
        evento: "personalizado",
        modulo: "desconocido",
        duracion_ms: 0
      }, data));

      navigator.sendBeacon(API_URL, payload);
    } catch (err) {
      console.warn("[EthicalMetrics] Error al enviar evento:", err);
    }
  }
})();
