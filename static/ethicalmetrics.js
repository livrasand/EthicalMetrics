(function () {
  const API_URL = "https://ethicalmetrics.onrender.com/track";
  const LOAD_TIME = Math.round(performance.now());

  // --- DNT y consentimiento ---
  function hasConsent() {
    return sessionStorage.getItem("ethicalmetrics_consent") === "true";
  }
  function setConsent(val) {
    sessionStorage.setItem("ethicalmetrics_consent", val ? "true" : "false");
  }
  function showConsentBanner(onAccept) {
    if (document.getElementById("ethicalmetrics-banner")) return;
    const SCRIPT = document.querySelector('script[src*="ethicalmetrics.js"]');
    // Lee opciones del script
    const bg = SCRIPT?.getAttribute("data-banner-bg") || "#fff";
    const color = SCRIPT?.getAttribute("data-banner-color") || "#222";
    const btnBg = SCRIPT?.getAttribute("data-banner-btn-bg") || "#4a90e2";
    const btnColor = SCRIPT?.getAttribute("data-banner-btn-color") || "#fff";
    const text = SCRIPT?.getAttribute("data-banner-text") ||
      "We respect your privacy. 'Do Not Track' is enabled in your browser, so we do not collect analytics by default.<br>Only if you accept, we will collect anonymous data to improve this site.";
    const btnText = SCRIPT?.getAttribute("data-banner-btn-text") || "Accept";

    const banner = document.createElement("div");
    banner.id = "ethicalmetrics-banner";
    banner.style.cssText = `
      position:fixed;bottom:0;left:0;right:0;z-index:9999;
      background:${bg};color:${color};padding:1em;text-align:center;
      box-shadow:0 -2px 8px #0002;font-family:sans-serif;
    `;
    banner.innerHTML = `
      <span style="display:inline-block;margin-bottom:0.5em;">
        ${text}
      </span>
      <button id="ethicalmetrics-accept" style="margin-left:1em;padding:0.5em 1.5em;background:${btnBg};color:${btnColor};border:none;border-radius:4px;cursor:pointer;">
        ${btnText}
      </button>
    `;
    document.body.appendChild(banner);
    document.getElementById("ethicalmetrics-accept").onclick = function () {
      setConsent(true);
      banner.remove();
      main(true); // relanza main con consentimiento
    };
  }

  // Esperar a que el DOM esté listo para asegurar que <body> existe
  if (!document.body) {
    document.addEventListener("DOMContentLoaded", checkConsentAndRun);
  } else {
    checkConsentAndRun();
  }

  function checkConsentAndRun() {
    const DNT = (
      navigator.doNotTrack == "1" ||
      window.doNotTrack == "1" ||
      navigator.msDoNotTrack == "1"
    );
    if (DNT && !hasConsent()) {
      showConsentBanner(() => main(true));
      return;
    }
    main();
  }

  function main(force) {
    // Si DNT está activo y no hay consentimiento, no hacer nada
    const DNT = (
      navigator.doNotTrack == "1" ||
      window.doNotTrack == "1" ||
      navigator.msDoNotTrack == "1"
    );
    if (DNT && !hasConsent() && !force) return;

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
    // Navegador detallado
    let BROWSER = "Otro";
    if (/edg\/\d+/i.test(NAVIGATOR)) BROWSER = "Edge (Chromium)";
    else if (/edge|edgios/i.test(NAVIGATOR)) BROWSER = "Edge (iOS)";
    else if (/crios/i.test(NAVIGATOR)) BROWSER = "Chrome (iOS)";
    else if (/fxios/i.test(NAVIGATOR)) BROWSER = "Firefox (iOS)";
    else if (/samsung/i.test(NAVIGATOR)) BROWSER = "Samsung";
    else if (/miui/i.test(NAVIGATOR)) BROWSER = "MIUI";
    else if (/yabrowser|yandex/i.test(NAVIGATOR)) BROWSER = "Yandex";
    else if (/opr|opera/i.test(NAVIGATOR)) BROWSER = "Opera";
    else if (/chrome\/[.0-9]* mobile/i.test(NAVIGATOR)) BROWSER = "Chrome (webview)";
    else if (/chromium/i.test(NAVIGATOR) && /webview/i.test(NAVIGATOR)) BROWSER = "Chromium Webview";
    else if (/chromium/i.test(NAVIGATOR)) BROWSER = "Chromium";
    else if (/chrome/i.test(NAVIGATOR)) BROWSER = "Chrome";
    else if (/firefox/i.test(NAVIGATOR)) BROWSER = "Firefox";
    else if (/safari/i.test(NAVIGATOR) && !/chrome|crios|chromium|edg|opr|opera|fxios|yabrowser|miui/i.test(NAVIGATOR)) BROWSER = "Safari";
    else if (/ios/i.test(NAVIGATOR)) BROWSER = "iOS";
    else if (/webview/i.test(NAVIGATOR) && /iphone|ipad|ipod/i.test(NAVIGATOR)) BROWSER = "iOS (webview)";

    // Sistemas operativos detallado
    let OS = "Otro";
    if (/windows nt 10|windows nt 11/i.test(NAVIGATOR)) OS = "Windows 10/11";
    else if (/windows nt 6\.1/i.test(NAVIGATOR)) OS = "Windows 7";
    else if (/windows/i.test(NAVIGATOR)) OS = "Windows";
    else if (/macintosh|mac os x/i.test(NAVIGATOR)) OS = "macOS";
    else if (/cros/i.test(NAVIGATOR)) OS = "ChromeOS";
    else if (/android/i.test(NAVIGATOR)) OS = "Android";
    else if (/iphone|ipad|ipod/i.test(NAVIGATOR)) OS = "iOS";
    else if (/linux/i.test(NAVIGATOR)) OS = "Linux";

    // Detección de idioma del navegador
    const BROWSER_LANG = navigator.language || (navigator.languages && navigator.languages[0]) || "desconocido";

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
      browser_lang: BROWSER_LANG,
      referer: REFERRER,
      page: PAGE,
      device: DEVICE,
      os: OS
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
            if (r.status === 403) {
              alert("⚠️ Este script es propiedad de otro sitio web y no funcionará aquí.");
            }
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
