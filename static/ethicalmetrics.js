(function () {
  if (navigator.doNotTrack === "1") return;

  const evento = {
    evento: "visita",
    modulo: "inicio",
    duracion_ms: performance.now()
  };

  navigator.sendBeacon("/track", JSON.stringify(evento));
})();
