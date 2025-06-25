(function () {
  const API_URL = "https://ethicalmetrics.onrender.com/track";
  const LOAD_TIME = Math.round(performance.now());

  // Esperar a que el DOM esté listo para asegurar que <body> existe
  if (!document.body) {
    document.addEventListener("DOMContentLoaded", main);
  } else {
    main();
  }

  function main() {
    const SCRIPT = document.querySelector('script[src*="ethicalmetrics.js"]');
    const URL_PARAMS = new URL(SCRIPT?.src || "").searchParams;
    let SITE_ID = SCRIPT?.getAttribute("data-site-id") || new URL(SCRIPT?.src || "").searchParams.get("id");

    const MODULE = document.body?.dataset?.modulo || "inicio";

    if (!SITE_ID) {
      console.warn("[EthicalMetrics] No se proporcionó site_id.");
      return;
    }

    console.log("EthicalMetrics loaded, enviando evento de visita...");

    // Manejo de cola tipo GA
    window.ethicalData = window.ethicalData || [];
    const queue = window.ethicalData;

    window.ethicalData = {
      push: handleCommand
    };

    queue.forEach(handleCommand);

    // Evento automático de visita
    send({
      evento: "visita",
      modulo: MODULE,
      duracion_ms: LOAD_TIME,
      site_id: SITE_ID
    });

    // Función global alternativa
    window.ethical = {
      track: send
    };

    function handleCommand(args) {
      if (!Array.isArray(args)) return;

      const [command, payload] = args;

      if (command === "init") {
        // Se ignora por ahora, se puede usar en el futuro
      }

      else if (command === "config" && typeof payload === "string") {
        SITE_ID = payload;
      }

      else if (command === "event" && typeof payload === "object") {
        send(Object.assign({
          evento: "personalizado",
          modulo: MODULE,
          duracion_ms: 0,
          site_id: SITE_ID
        }, payload));
      }
    }

    function send(data) {
      try {
        if (navigator.doNotTrack === "1") {
          console.warn("[EthicalMetrics] Do Not Track activado, no se enviará evento.");
          return;
        }

        const payload = JSON.stringify(Object.assign({
          evento: "personalizado",
          modulo: "desconocido",
          duracion_ms: 0,
          site_id: SITE_ID
        }, data));

        const ok = navigator.sendBeacon(API_URL, payload);
        console.log("Intentando sendBeacon a", API_URL);
        console.log("[EthicalMetrics] Evento enviado:", payload, "Resultado:", ok);
        if (!ok) {
          // Fallback a fetch si sendBeacon falla
          fetch(API_URL, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: payload,
            keepalive: true
          }).then(r => {
            console.log("[EthicalMetrics] Evento enviado por fetch, status:", r.status);
          }).catch(e => {
            console.warn("[EthicalMetrics] Error al enviar evento por fetch:", e);
          });
        }
      } catch (err) {
        console.warn("[EthicalMetrics] Error al enviar evento:", err);
      }
    }
  }
})();
