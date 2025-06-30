# Manejo de la dirección IP en EthicalMetrics

EthicalMetrics utiliza la dirección IP del usuario exclusivamente para obtener información de geolocalización (ciudad y país).  Este proceso se realiza de la siguiente manera:

1. **Obtención de la IP:** La dirección IP del usuario se obtiene de la cabecera `X-Forwarded-For` (si está presente, indicando un proxy) o, en su defecto, de `r.RemoteAddr` (la dirección remota de la conexión).

2. **Geolocalización:** Se utiliza una base de datos GeoIP2 para determinar la ciudad y el país asociados a la dirección IP.  Esta búsqueda se realiza mediante las funciones `cityFromIP` y `countryFromIP`.

3. **Almacenamiento:**  **La dirección IP en sí NO se almacena ni se incluye en los datos del evento que se guardan en Redis.**  Solo se guardan la ciudad y el país, que se añaden al mapa `eventMap` junto con otros datos del evento (tipo de evento, módulo, duración, etc.).  Este mapa se serializa a JSON y se almacena en Redis.

4. **Descarte:**  Después de obtener la ciudad y el país, la dirección IP se descarta y no se utiliza para ningún otro propósito.  No se guarda en la base de datos ni se asocia de ninguna manera con los datos del evento.

Este enfoque asegura que EthicalMetrics recopile información de ubicación a nivel de ciudad y país para análisis, sin comprometer la privacidad del usuario al almacenar o procesar su dirección IP más allá de este uso puntual.