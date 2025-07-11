(function() {
    const API_URL = "https://ethicalmetrics.onrender.com/track";
    const LOAD_TIME = Math.round(performance.now());

    function sanitizeForCSS(value) {
        return value.replace(/[^a-zA-Z0-9\s,#-]/g, '');
    }

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
        const bg = sanitizeForCSS(SCRIPT?.getAttribute("data-banner-bg") || "#ffffffee");
        const color = sanitizeForCSS(SCRIPT?.getAttribute("data-banner-color") || "#111");
        const btnBg = sanitizeForCSS(SCRIPT?.getAttribute("data-banner-btn-bg") || "#111");
        const btnColor = sanitizeForCSS(SCRIPT?.getAttribute("data-banner-btn-color") || "#fff");
        const text = SCRIPT?.getAttribute("data-banner-text") ||
            "Anonymous analytics help improve this site. Enable?";
        const btnText = SCRIPT?.getAttribute("data-banner-btn-text") || "Allow";

        // Banner element
        const banner = document.createElement("div");
        banner.id = "ethicalmetrics-banner";

        const topRow = document.createElement("div");
        topRow.className = "top-row";

        const closeButton = document.createElement("button");
        closeButton.id = "ethicalmetrics-close";
        closeButton.setAttribute("aria-label", "Close banner");
        closeButton.textContent = "×";

        const textSpan = document.createElement("span");
        textSpan.textContent = text;

        const acceptButton = document.createElement("button");
        acceptButton.id = "ethicalmetrics-accept";
        acceptButton.textContent = btnText;

        topRow.appendChild(closeButton);
        topRow.appendChild(textSpan);
        topRow.appendChild(acceptButton);

        const promoDiv = document.createElement("div");
        promoDiv.id = "ethicalmetrics-promo";
        const promoLink = document.createElement("a");
        promoLink.href = "https://ethicalmetrics.onrender.com/";
        promoLink.target = "_blank";
        promoLink.rel = "noopener";
        promoLink.textContent = "Powered by EthicalMetrics";
        promoDiv.appendChild(promoLink);

        banner.appendChild(topRow);
        banner.appendChild(promoDiv);

        document.body.appendChild(banner);

        banner.style.position = 'fixed';
        banner.style.bottom = '1.5em';
        banner.style.right = '1.5em';
        banner.style.zIndex = '9999';
        banner.style.background = bg;
        banner.style.color = color;
        banner.style.fontFamily = 'system-ui, sans-serif';
        banner.style.fontSize = '0.85rem';
        banner.style.padding = '0.75em 1em';
        banner.style.borderRadius = '6px';
        banner.style.display = 'flex';
        banner.style.flexDirection = 'column';
        banner.style.alignItems = 'flex-start';
        banner.style.gap = '0.75em';
        banner.style.boxShadow = '0 4px 12px rgba(0,0,0,0.06)';
        banner.style.backdropFilter = 'blur(6px)';
        banner.style.transition = 'opacity 0.3s ease';

        topRow.style.width = '100%';
        topRow.style.display = 'flex';
        topRow.style.justifyContent = 'space-between';
        topRow.style.alignItems = 'flex-start';
        topRow.style.gap = '1em';

        closeButton.style.background = 'none';
        closeButton.style.color = btnBg;
        closeButton.style.padding = '0.4em 0';
        closeButton.style.border = 'none';
        closeButton.style.cursor = 'pointer';

        acceptButton.style.background = btnBg;
        acceptButton.style.color = btnColor;
        acceptButton.style.border = 'none';
        acceptButton.style.padding = '0.4em 1em';
        acceptButton.style.borderRadius = '4px';
        acceptButton.style.fontSize = '0.6rem';
        acceptButton.style.cursor = 'pointer';

        promoDiv.style.fontSize = '0.65rem';
        promoDiv.style.color = '#666';
        promoDiv.style.alignSelf = 'flex-end';
        promoDiv.style.marginTop = '0.25em';

        promoLink.style.textDecoration = 'none';
        promoLink.style.color = 'inherit';
        promoLink.style.opacity = '0.6';

        acceptButton.onclick = function() {
            setConsent(true);
            banner.style.opacity = "0";
            setTimeout(() => banner.remove(), 300);
            main(true);
        };

        closeButton.onclick = function() {
            setConsent(false);
            banner.style.opacity = "0";
            setTimeout(() => banner.remove(), 300);
            main(false);
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

        // UTM tracking
        function getUTMParams() {
            const params = {};
            const urlParams = new URLSearchParams(window.location.search);
            ["utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content"].forEach(key => {
                if (urlParams.has(key)) params[key] = urlParams.get(key);
            });
            return params;
        }
        const UTM_PARAMS = getUTMParams();
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

        // Sistemas operations detallado
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

        // Lógica para is_new_session
        const is_new_session = !sessionStorage.getItem('ethicalmetrics_session');
        if (is_new_session) {
            sessionStorage.setItem('ethicalmetrics_session', 'true');
        }

        // Lógica para is_new_visit
        const is_new_visit = !localStorage.getItem('ethicalmetrics_visited');
        if (is_new_visit) {
            localStorage.setItem('ethicalmetrics_visited', 'true');
        }

        // Lógica para is_unique (vista de página única por sesión)
        const pagesViewedStr = sessionStorage.getItem('ethicalmetrics_pages_viewed') || '[]';
        const pagesViewed = JSON.parse(pagesViewedStr);
        const is_unique = !pagesViewed.includes(PAGE);
        if (is_unique) {
            pagesViewed.push(PAGE);
            sessionStorage.setItem('ethicalmetrics_pages_viewed', JSON.stringify(pagesViewed));
        }

        if (!SITE_ID) {
            console.warn("[EthicalMetrics] No se proportionó site_id.");
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
            os: OS,
            is_new_session: is_new_session,
            is_new_visit: is_new_visit,
            is_unique: is_unique,
            ...UTM_PARAMS // Añade los UTM si existen
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
            } else if (command === "config" && typeof payload === "string") {
                SITE_ID = payload;
            } else if (command === "event" && typeof payload === "object") {
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
                        headers: {
                            "Content-Type": "application/json"
                        },
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