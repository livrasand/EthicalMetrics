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

    const MODULE = document.body?.dataset?.modulo || "visita";

    // --- NUEVO: detectar info extra ---
    const NAVIGATOR = navigator.userAgent;
    const REFERRER = document.referrer || "directo";
    const PAGE = location.pathname;
    // Detección simple de dispositivo
    let DEVICE = "desktop";
    if (/Mobi|Android/i.test(NAVIGATOR)) DEVICE = "mobile";
    else if (/Tablet|iPad/i.test(NAVIGATOR)) DEVICE = "tablet";
    // Navegador simplificado
    let BROWSER = "Otro";
    if (/edg/i.test(NAVIGATOR)) BROWSER = "Edge";
    else if (/opr|opera/i.test(NAVIGATOR)) BROWSER = "Opera";
    else if (/brave/i.test(NAVIGATOR)) BROWSER = "Brave";
    else if (/chrome|crios/i.test(NAVIGATOR)) BROWSER = "Chrome";
    else if (/firefox|fxios/i.test(NAVIGATOR)) BROWSER = "Firefox";
    else if (/safari/i.test(NAVIGATOR) && !/chrome|crios|edg|opr|opera|brave/i.test(NAVIGATOR)) BROWSER = "Safari";
    const BROWSER_LANG = navigator.language || navigator.userLanguage || "desconocido";

    let OS = "Otro";
    if (/windows/i.test(NAVIGATOR)) OS = "Windows";
    else if (/macintosh|mac os x/i.test(NAVIGATOR)) OS = "MacOS";
    else if (/linux/i.test(NAVIGATOR)) OS = "Linux";
    else if (/android/i.test(NAVIGATOR)) OS = "Android";
    else if (/iphone|ipad|ipod/i.test(NAVIGATOR)) OS = "iOS";

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
      site_id: SITE_ID,
      browser: BROWSER,
      browser_lang: BROWSER_LANG, // <--- nuevo
      referer: REFERRER,
      page: PAGE,
      device: DEVICE,
      os: OS // <--- nuevo, ver abajo
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
          site_id: SITE_ID,
          browser: BROWSER,
          referer: REFERRER,
          page: PAGE,
          device: DEVICE
        }, payload));
      }
    }

    function send(data) {
      try {
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
