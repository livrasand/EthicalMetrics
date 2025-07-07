package api

// misspell:ignore instruccion calcular

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/livrasand/ethicalmetrics/internal/db"
	"github.com/oschwald/geoip2-golang"
	"golang.org/x/time/rate"
)

var geoDB *geoip2.Reader

// --- Inicio: Middleware de Rate Limiting Mejorado ---

// Visitor representa a un usuario con su propio limitador y hora de última visita.
type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Global store para los limitadores de cada visitante.
// La clave es la dirección IP del visitante.
var visitors = make(map[string]*Visitor)
var mtx sync.Mutex

// Inicializa y ejecuta la limpieza periódica de visitantes antiguos.
func init() {
	go cleanupVisitors()
}

// RateLimitMiddleware envuelve un http.Handler para aplicar el límite de peticiones.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener la dirección IP real del visitante.
		ip := getIP(r)

		mtx.Lock()
		v, exists := visitors[ip]
		if !exists {
			// Permitir 30 eventos por minuto, con ráfagas de hasta 10.
			limiter := rate.NewLimiter(rate.Every(2*time.Second), 10)
			v = &Visitor{limiter: limiter}
			visitors[ip] = v
		}
		v.lastSeen = time.Now()
		mtx.Unlock()

		// Verificar si el visitante ha excedido el límite.
		if !v.limiter.Allow() {
			http.Error(w, "Demasiadas peticiones, espera un minuto.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getIP obtiene la dirección IP real del cliente, considerando proxies.
func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

// cleanupVisitors elimina periódicamente los visitantes inactivos para no agotar la memoria.
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mtx.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mtx.Unlock()
	}
}

// --- Fin: Middleware de Rate Limiting ---

// Event representa la estructura de datos para un evento de seguimiento.
type Event struct {
	EventType string `json:"evento"`
	Module    string `json:"modulo"`
	SiteID    string `json:"site_id"`
	Duration  int    `json:"duracion_ms"`
}

// processedStats es una estructura interna para contener los datos agregados
// después de procesar todos los eventos sin procesar de Redis.
type processedStats struct {
	modCount         map[string]int
	dayCount         map[string]int
	browserCount     map[string]int
	refererCount     map[string]int
	pageCount        map[string]int
	browserLangCount map[string]int
	osCount          map[string]int
	cityCount        map[string]int
	countryCount     map[string]int
	deviceCount      map[string]int
	totalDuration    int
	sessionCount     int
	activeUsers      int
	weekData         map[string][2]int
	monthData        map[string][2]int
	pageVisits       map[string]map[string]bool
	newVisitorsCount int
	newSessionsCount int
	uniqueViewsCount int
	regionCount      map[string]int
}

func InitGeoIP(path string) error {
	var err error
	geoDB, err = geoip2.Open(path)
	return err
}


func NuevoHandler(w http.ResponseWriter, r *http.Request) {
	siteID := uuid.NewString()
	adminToken := generarToken()

	err := db.RDB.HSet(db.Ctx, "site:"+siteID, map[string]interface{}{
		"name":        "Sitio sin nombre",
		"admin_token": adminToken,
		"created_at":  time.Now().Format(time.RFC3339),
	}).Err()

	if err != nil {
		http.Error(w, "Error creando sitio", http.StatusInternalServerError)
		return
	}

	siteName, _ := db.RDB.HGet(db.Ctx, "site:"+siteID, "name").Result()

	resp := map[string]interface{}{
		"site_id":     siteID,
		"admin_token": adminToken,
		"instruccion": `<script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js?id=` + siteID + `"></script>`,
		"site_name":   siteName,
	}
	json.NewEncoder(w).Encode(resp)
}

func generarToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 24)

	// Usando crypto/rand para mayor seguridad
	if _, err := rand.Read(b); err != nil {
		log.Printf("Error generando token seguro: %v", err)
		panic("No se pudo generar token seguro") // O podrías usar log.Fatal dependiendo de tu caso
	}

	// Mapear los bytes aleatorios a nuestros caracteres permitidos
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo body", http.StatusBadRequest)
		return
	}

	var e Event
	err = json.Unmarshal(body, &e)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validación estricta de campos esperados
	if e.SiteID == "" || len(e.SiteID) > 64 {
		http.Error(w, "site_id inválido", http.StatusBadRequest)
		return
	}
	if e.EventType == "" || len(e.EventType) > 32 {
		http.Error(w, "evento inválido", http.StatusBadRequest)
		return
	}
	if e.Module == "" || len(e.Module) > 64 {
		http.Error(w, "modulo inválido", http.StatusBadRequest)
		return
	}
	if e.Duration < 0 || e.Duration > 86400000 {
		http.Error(w, "duracion_ms inválido", http.StatusBadRequest)
		return
	}

	// Validar campos extra si existen en el body
	var extra map[string]interface{}
	json.Unmarshal(body, &extra)
	if v, ok := extra["browser"]; ok {
		if browser, ok := v.(string); !ok || len(browser) > 64 {
			http.Error(w, "browser inválido", http.StatusBadRequest)
			return
		}
	}
	if v, ok := extra["browser_lang"]; ok {
		if lang, ok := v.(string); !ok || len(lang) > 16 {
			http.Error(w, "browser_lang inválido", http.StatusBadRequest)
			return
		}
	}
	if v, ok := extra["os"]; ok {
		if os, ok := v.(string); !ok || len(os) > 64 {
			http.Error(w, "os inválido", http.StatusBadRequest)
			return
		}
	}
	if v, ok := extra["referer"]; ok {
		if ref, ok := v.(string); !ok || len(ref) > 256 {
			http.Error(w, "referer inválido", http.StatusBadRequest)
			return
		}
	}
	if v, ok := extra["page"]; ok {
		if page, ok := v.(string); !ok || len(page) > 256 {
			http.Error(w, "page inválido", http.StatusBadRequest)
			return
		}
	}
	if v, ok := extra["device"]; ok {
		if device, ok := v.(string); !ok || len(device) > 32 {
			http.Error(w, "device inválido", http.StatusBadRequest)
			return
		}
	}

	if v, ok := extra["is_new_session"]; ok {
		if _, ok := v.(bool); !ok {
			http.Error(w, "is_new_session debe ser booleano", http.StatusBadRequest)
			return
		}
	}
	if v, ok := extra["is_new_visit"]; ok {
		if _, ok := v.(bool); !ok {
			http.Error(w, "is_new_visit debe ser booleano", http.StatusBadRequest)
			return
		}
	}
	if v, ok := extra["is_unique"]; ok {
		if _, ok := v.(bool); !ok {
			http.Error(w, "is_unique debe ser booleano", http.StatusBadRequest)
			return
		}
	}

	siteKey := "site:" + e.SiteID

	// Detectar origen real
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = r.Header.Get("Referer")
	}
	if origin == "" {
		origin = r.Host
	}

	// Leer dominio registrado
	registeredDomain, _ := db.RDB.HGet(db.Ctx, siteKey, "domain").Result()
	if registeredDomain == "" {
		db.RDB.HSet(db.Ctx, siteKey, "domain", origin)
		registeredDomain = origin
	}

	// CORS estricto: solo permite el dominio registrado
	w.Header().Set("Access-Control-Allow-Origin", registeredDomain)
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Validar dominio
	if !strings.Contains(origin, registeredDomain) {
		http.Error(w, "Este script solo puede usarse en el dominio registrado: "+registeredDomain, http.StatusForbidden)
		return
	}

	// Verifica si el sitio existe
	exists, err := db.RDB.Exists(db.Ctx, "site:"+e.SiteID).Result()
	if err != nil || exists == 0 {
		http.Error(w, "Site ID inválido", http.StatusForbidden)
		return
	}

	// Guarda el evento como JSON en una lista
	eventMap := map[string]interface{}{
		"type":        e.EventType,
		"module":      e.Module,
		"duration_ms": e.Duration,
		"timestamp":   time.Now().Format(time.RFC3339),
	}
	// Agregar campos extra si existen en el body
	json.Unmarshal(body, &extra)
	if v, ok := extra["browser"]; ok {
		eventMap["browser"] = v
	}
	if v, ok := extra["browser_lang"]; ok {
		eventMap["browser_lang"] = v
	}
	if v, ok := extra["os"]; ok {
		eventMap["os"] = v
	}
	if v, ok := extra["referer"]; ok {
		eventMap["referer"] = v
	}
	if v, ok := extra["page"]; ok {
		eventMap["page"] = v
	}
	if v, ok := extra["device"]; ok {
		eventMap["device"] = v
	}
	if v, ok := extra["is_new_session"]; ok {
		eventMap["is_new_session"] = v
	}
	if v, ok := extra["is_new_visit"]; ok {
		eventMap["is_new_visit"] = v
	}
	if v, ok := extra["is_unique"]; ok {
		eventMap["is_unique"] = v
	}
	// Detectar ciudad y país por IP (GeoLite2)
	clientIP := func(r *http.Request) string {
		// 1) X-Forwarded-For puede contener "IP1, IP2, ...", elige la primera
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			ip := strings.TrimSpace(parts[0])
			return ip
		}
		// 2) Fallback a RemoteAddr
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			return host
		}
		return r.RemoteAddr
	}

	userIP := clientIP(r)

	eventMap["city"] = cityFromIP(userIP)
	eventMap["country"] = countryFromIP(userIP)
	eventMap["region"] = regionFromIP(userIP)
	// La IP del usuario no se guarda, solo los datos obtenidos.
	eventJSON, _ := json.Marshal(eventMap)

	err = db.RDB.RPush(db.Ctx, "events:"+e.SiteID, eventJSON).Err()
	if err != nil {
		http.Error(w, "Error guardando evento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// StatsHandler maneja la solicitud de estadísticas. Su complejidad se reduce
// al delegar el procesamiento a funciones auxiliares.
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	siteID := r.URL.Query().Get("site")
	token := r.URL.Query().Get("token")

	if siteID == "" || token == "" {
		http.Error(w, "Faltan parámetros", http.StatusBadRequest)
		return
	}

	storedToken, err := db.RDB.HGet(db.Ctx, "site:"+siteID, "admin_token").Result()
	if err != nil || storedToken != token {
		http.Error(w, "Token inválido", http.StatusForbidden)
		return
	}

	eventsRaw, err := db.RDB.LRange(db.Ctx, "events:"+siteID, 0, -1).Result()
	if err != nil {
		http.Error(w, "Error obteniendo eventos", http.StatusInternalServerError)
		return
	}

	// Procesar todos los eventos en una sola pasada.
	stats := processEvents(eventsRaw)

	// Construir la respuesta final.
	resp := buildResponse(stats)

	json.NewEncoder(w).Encode(resp)
}

// processEvents itera sobre los eventos sin procesar una vez y los agrega en la estructura processedStats.
func processEvents(eventsRaw []string) *processedStats {
	stats := &processedStats{
		modCount:         make(map[string]int),
		dayCount:         make(map[string]int),
		browserCount:     make(map[string]int),
		refererCount:     make(map[string]int),
		pageCount:        make(map[string]int),
		browserLangCount: make(map[string]int),
		osCount:          make(map[string]int),
		cityCount:        make(map[string]int),
		countryCount:     make(map[string]int),
		deviceCount:      make(map[string]int),
		weekData:         make(map[string][2]int),
		monthData:        make(map[string][2]int),
		pageVisits:       make(map[string]map[string]bool),
		newVisitorsCount: 0,
		newSessionsCount: 0,
		uniqueViewsCount: 0,
		regionCount: 	make(map[string]int),
	}

	now := time.Now()
	usuariosActivosWindow := now.Add(-5 * time.Minute)
	currentYear, currentWeek := now.ISOWeek()
	prevYear, prevWeek := now.AddDate(0, 0, -7).ISOWeek()
	currentMonth := now.Month()
	prevMonth := now.AddDate(0, -1, 0).Month()

	for _, raw := range eventsRaw {
		var evt map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &evt); err != nil {
			continue // Omitir evento malformado
		}

		// Extraer datos comunes del evento
		tsStr, _ := evt["timestamp"].(string)
		t, err := time.Parse(time.RFC3339, tsStr)
		if err != nil {
			continue
		}

		// Métricas básicas
		if mod, ok := evt["module"].(string); ok {
			stats.modCount[mod]++
		}
		if browser, ok := evt["browser"].(string); ok {
			stats.browserCount[browser]++
		}
		if referer, ok := evt["referer"].(string); ok {
			stats.refererCount[referer]++
		}
		if page, ok := evt["page"].(string); ok {
			stats.pageCount[page]++
		}
		if bl, ok := evt["browser_lang"].(string); ok {
			stats.browserLangCount[bl]++
		}
		if os, ok := evt["os"].(string); ok {
			stats.osCount[os]++
		}
		if city, ok := evt["city"].(string); ok {
			stats.cityCount[city]++
		}
		if country, ok := evt["country"].(string); ok {
			stats.countryCount[country]++
		}
		if region, ok := evt["region"].(string); ok {
			stats.regionCount[region]++
		}
		if dev, ok := evt["device"].(string); ok {
			stats.deviceCount[dev]++
		}
		if dur, ok := evt["duration_ms"].(float64); ok && dur > 0 {
			stats.totalDuration += int(dur)
			stats.sessionCount++
		}

		if isNew, ok := evt["is_new_visit"].(bool); ok && isNew {
			stats.newVisitorsCount++
		}
		if isNew, ok := evt["is_new_session"].(bool); ok && isNew {
			stats.newSessionsCount++
		}
		if isUnique, ok := evt["is_unique"].(bool); ok && isUnique {
			stats.uniqueViewsCount++
		}

		// Métricas por tiempo
		day := t.Format("2006-01-02")
		stats.dayCount[day]++

		if t.After(usuariosActivosWindow) {
			stats.activeUsers++
		}

		// Comparativas semanales y mensuales
		year, week := t.ISOWeek()
		weekLabel := t.Format("Mon")
		if year == currentYear && week == currentWeek {
			val := stats.weekData[weekLabel]
			val[0]++
			stats.weekData[weekLabel] = val
		} else if year == prevYear && week == prevWeek {
			val := stats.weekData[weekLabel]
			val[1]++
			stats.weekData[weekLabel] = val
		}

		monthLabel := t.Format("02")
		if t.Year() == now.Year() && t.Month() == currentMonth {
			val := stats.monthData[monthLabel]
			val[0]++
			stats.monthData[monthLabel] = val
		} else if t.Year() == now.Year() && t.Month() == prevMonth {
			val := stats.monthData[monthLabel]
			val[1]++
			stats.monthData[monthLabel] = val
		}

		// Datos para retención
		if page, ok := evt["page"].(string); ok && page != "" {
			if stats.pageVisits[page] == nil {
				stats.pageVisits[page] = make(map[string]bool)
			}
			stats.pageVisits[page][day] = true
		}
	}

	return stats
}

// buildResponse construye el mapa de respuesta final a partir de los datos procesados.
func buildResponse(stats *processedStats) map[string]interface{} {
	var duracionMedia int
	if stats.sessionCount > 0 {
		duracionMedia = stats.totalDuration / stats.sessionCount
	}

	return map[string]interface{}{
		"por_modulo":           processMapToSlice(stats.modCount, "modulo", "total"),
		"por_dia":              processMapToSlice(stats.dayCount, "dia", "total"),
		"navegadores":          processMapToSlice(stats.browserCount, "navegador", "total"),
		"referencias":          processMapToSlice(stats.refererCount, "referencia", "total"),
		"paginas":              processMapToSlice(stats.pageCount, "pagina", "total"),
		"duracion_media":       duracionMedia,
		"dispositivos":         processMapToSlice(stats.deviceCount, "dispositivo", "total"),
		"paises":               processMapToSlice(stats.countryCount, "pais", "total"),
		"usuarios_activos":     stats.activeUsers,
		"browser_langs":        processMapToSlice(stats.browserLangCount, "lang", "total"),
		"os":                   processMapToSlice(stats.osCount, "os", "total"),
		"cities":               processMapToSlice(stats.cityCount, "city", "total"),
	    "regiones":             processMapToSlice(stats.regionCount, "region", "total"),
		"week_compare":         processComparison(stats.weekData),
		"month_compare":        processComparison(stats.monthData),
		"retention":            processRetention(stats.pageVisits),
		"funnel":               processFunnel(stats.modCount),
		"nuevos_visitantes":    stats.newVisitorsCount,
		"sesiones_totales":     stats.newSessionsCount,
		"vistas_unicas_pagina": stats.uniqueViewsCount,
	}
}

// processMapToSlice es una función genérica para convertir un mapa de contadores a un slice de mapas.
func processMapToSlice(m map[string]int, keyName, valueName string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(m))
	for k, v := range m {
		result = append(result, map[string]interface{}{keyName: k, valueName: v})
	}
	if result == nil {
		return []map[string]interface{}{}
	}
	return result
}

// processComparison convierte los datos de comparación (semanal/mensual) al formato de respuesta.
func processComparison(data map[string][2]int) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(data))
	for label, vals := range data {
		result = append(result, map[string]interface{}{
			"label":    label,
			"current":  vals[0],
			"previous": vals[1],
		})
	}
	return result
}

// processRetention calcula la retención a partir de las visitas por página.
func processRetention(pageVisits map[string]map[string]bool) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(pageVisits))
	for page, days := range pageVisits {
		result = append(result, map[string]interface{}{
			"label": page,
			"value": len(days), // días distintos con visitas
		})
	}
	return result
}

// processFunnel calcula el funnel básico de conversión.
func processFunnel(modCount map[string]int) []map[string]interface{} {
	funnel := []map[string]interface{}{}
	modOrder := []string{"landing", "signup", "checkout", "thanks"}
	for _, mod := range modOrder {
		funnel = append(funnel, map[string]interface{}{
			"step":  mod,
			"value": modCount[mod],
		})
	}
	return funnel
}

func countryFromIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil || geoDB == nil {
		log.Printf("countryFromIP: IP inválida %q", ipStr)
		return "Desconocido"
	}
	record, err := geoDB.Country(ip)
	if err != nil || record == nil || record.Country.Names["en"] == "" {
		log.Printf("countryFromIP: sin datos para %q (err=%v)", ipStr, err)
		return "Desconocido"
	}
	return record.Country.Names["en"]
}

func cityFromIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil || geoDB == nil {
		log.Printf("cityFromIP: IP inválida %q", ipStr)
		return "Desconocido"
	}
	record, err := geoDB.City(ip)
	if err != nil || record == nil || record.City.Names["es"] == "" {
		log.Printf("cityFromIP: sin datos para %q (err=%v)", ipStr, err)
		return "Desconocido"
	}
	return record.City.Names["es"]
}

// regionFromIP devuelve la primera subdivisión (región/estado) en español
func regionFromIP(ipStr string) string {
    ip := net.ParseIP(ipStr)
    if ip == nil || geoDB == nil {
		log.Printf("regionFromIP: IP inválida %q", ipStr)
        return "Desconocido"
    }
    record, err := geoDB.City(ip)
    if err != nil || record == nil {
		log.Printf("regionFromIP: sin datos para %q (err=%v)", ipStr, err)
        return "Desconocido"
    }
    if len(record.Subdivisions) > 0 {
        // Usa el nombre en español si está disponible
        if name, ok := record.Subdivisions[0].Names["es"]; ok && name != "" {
            return name
        }
        // Sino en inglés
        return record.Subdivisions[0].Names["en"]
    }
    return "Desconocido"
}
