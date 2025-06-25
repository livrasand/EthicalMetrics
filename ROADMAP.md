## Roadmap: Funcionalidades y Objetivos

### 📊 Métricas básicas

* [x] Visitas y vistas por página
* [x] Referencias (referrers)
* [x] Tiempo de carga
* [x] Páginas más visitadas
* [x] Duración media de sesión
* [x] Eventos personalizados
* [ ] Bounce rate (por página y global)
* [ ] Dispositivos, navegadores y sistemas operativos
* [ ] Idiomas del navegador
* [ ] Ubicación (ciudad, región, país) con MaxMind GeoLite2
* [ ] UTM tracking (`utm_source`, `utm_medium`, `utm_campaign`, etc.)
* [ ] Comparaciones temporales (semana vs semana, mes vs mes)
* [ ] Retención (usuarios que regresan, frecuencia, duración media de visitas repetidas)
* [ ] Funnels básicos (progresión de eventos o páginas)

### 🔬 Métricas avanzadas y segmentación

* [ ] Filtros por país, navegador, URL, dispositivo, fuente, etc.
* [ ] Segmentación avanzada (criterios múltiples, condiciones AND/OR)
* [ ] Custom dimensions y variables personalizadas
* [ ] Visualización de journeys o navegación secuencial
* [ ] Attribution de conversiones/eventos multi-canal
* [ ] Objetivos (goals) definidos por el usuario, con condiciones flexibles
* [ ] Cohortes de usuarios y comparación histórica
* [ ] Evolución temporal por fila (row evolution)
* [ ] Eventos enriquecidos: clics en enlaces externos, descargas, banners, errores 404, scroll, etc.
* [ ] Duración exacta por página
* [ ] Páginas de entrada y salida
* [ ] Búsqueda interna del sitio
* [ ] CTR por elemento (textos, imágenes, banners)
* [ ] Seguimiento de formularios (envíos, abandonos, campos completados)

---

## 🛍️ 2. Ecommerce y conversiones

* [ ] Seguimiento de ecommerce: productos vistos, añadidos, comprados, ingresos
* [ ] Conversiones multi-canal (ads, orgánico, directo, referers)
* [ ] Valor monetario asignado a eventos
* [ ] Visualización de funnels de conversión
* [ ] Campañas con tracking automático (Google Ads, Facebook Ads, etc.)

---

## 🎨 3. Experiencia de usuario (UX/UI)

* [ ] SPA (Single Page App) con Vue.js, React o Svelte
* [ ] Interfaz responsive con modo claro/oscuro
* [ ] Vista de overview por defecto con métricas clave
* [ ] Dashboard en tiempo real
* [ ] Paneles personalizados (drag & drop)
* [ ] Segmentos y filtros personalizados con UI intuitiva
* [ ] Dashboard multi-sitio
* [ ] Visualización de user flow (flujo de navegación)
* [ ] Page transitions (qué hizo el usuario antes/después de una página)
* [ ] Page overlay (estadísticas visuales sobre el sitio real)
* [ ] Anotaciones en gráficos
* [ ] Alertas automáticas personalizadas

---

## 🔥 4. Funciones Premium (UX analytics visuales)

* [ ] Heatmaps y scrollmaps
* [ ] Grabación de sesiones (reproducción de interacciones reales)
* [ ] A/B Testing nativo
* [ ] Análisis de formularios detallado
* [ ] Media analytics (videos/audio, pausas, porcentaje reproducido)

---

## ⚙️ 5. Backend e infraestructura

* [ ] Base de datos optimizada (índices, consultas eficientes)
* [ ] Carga de datos asincrónica (websockets o polling)
* [ ] Procesamiento eficiente de eventos (batch o streaming)
* [ ] Importación/exportación de datos
* [ ] API RESTful o GraphQL pública
* [ ] Soporte para múltiples sitios con roll-up reporting
* [ ] Control de usuarios por sitio (roles y permisos)
* [ ] Auditoría y logs de actividades
* [ ] Soporte para intranets (tracking por logs Apache/Nginx)
* [ ] Multiidioma en la interfaz
* [ ] Escalabilidad para alto tráfico y balanceo de carga
* [ ] White labeling (marca blanca)

---

## 🔧 6. Integraciones y extensibilidad

* [ ] Sistema de plugins/extensiones
* [ ] API pública de administración y tracking
* [ ] SDKs oficiales para:

  * JavaScript
  * PHP
  * Python
  * Android/iOS
* [ ] Plugins para CMS (WordPress, Joomla, Drupal)
* [ ] Integraciones con frameworks modernos (Vue, React, Nuxt, Next, Astro, etc.)
* [ ] Tracking sin conexión (offline tracking que se sincroniza después)
* [ ] Tracking por logs (Apache, Nginx, IIS)

---

## 🛡️ 7. Privacidad y cumplimiento legal

* [x] Sin cookies
* [x] Sin IPs
* [x] Autohospedado
* [x] Compatible con SQLCipher
* [ ] Cumplimiento automático con DNT (Do Not Track)
* [ ] Certificación explícita GDPR, CCPA, PECR
* [ ] Herramientas legales prehechas (DPA, términos, políticas)
* [ ] Formulario de opt-out embebible (iFrame)
* [ ] Gestión de consentimiento granular por categoría
* [ ] Anonimización de referrer y geodatos
* [ ] Funciones ARCO (acceso, rectificación, cancelación, oposición)
* [ ] Eliminación de datos de visitantes específicos ("derecho al olvido")
* [ ] Control de cookies (1ra y 3ra parte)
* [ ] Reemplazo de User ID por pseudónimos

---

## 🧪 8. Herramientas de análisis y administración

* [ ] Exportación de datos: JSON, XML, CSV, Excel
* [ ] Reportes automáticos programables (PDF, HTML, PNG)
* [ ] Soporte de zonas horarias por sitio
* [ ] Exclusión de IPs o rangos de IPs
* [ ] Exclusión de parámetros de URL
* [ ] Dashboard embebible en apps o sitios
* [ ] Aplicación móvil oficial
* [ ] Soporte para múltiples monedas
* [ ] Gestión de múltiples usuarios y sitios
* [ ] Panel de configuración multisitio
