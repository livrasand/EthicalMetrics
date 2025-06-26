package api

import (
	"encoding/json"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/livrasand/ethicalmetrics/internal/db"
	"github.com/oschwald/geoip2-golang"
)

var geoDB *geoip2.Reader

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

	resp := map[string]string{
		"site_id":     siteID,
		"admin_token": adminToken,
		"instruccion": `<script src="https://ethicalmetrics.onrender.com/ethicalmetrics.js?id=` + siteID + `"></script>`,
	}
	json.NewEncoder(w).Encode(resp)
}

func generarToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 24)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type Event struct {
	EventType string `json:"evento"`
	Module    string `json:"modulo"`
	SiteID    string `json:"site_id"`
	Duration  int    `json:"duracion_ms"`
}

type moduloStat struct {
	Modulo string `json:"modulo"`
	Total  int    `json:"total"`
}

type diaStat struct {
	Dia   string `json:"dia"`
	Total int    `json:"total"`
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
	var extra map[string]interface{}
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
	// Detectar ciudad por IP (GeoLite2)
	userIP := r.Header.Get("X-Forwarded-For")
	if userIP == "" {
		userIP, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	eventMap["city"] = cityFromIP(userIP)
	eventJSON, _ := json.Marshal(eventMap)

	err = db.RDB.RPush(db.Ctx, "events:"+e.SiteID, eventJSON).Err()
	if err != nil {
		http.Error(w, "Error guardando evento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func countryFromIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil || geoDB == nil {
		return "Desconocido"
	}
	record, err := geoDB.Country(ip)
	if err != nil || record == nil || record.Country.Names["en"] == "" {
		return "Desconocido"
	}
	return record.Country.Names["en"]
}

func cityFromIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil || geoDB == nil {
		return "Desconocido"
	}
	record, err := geoDB.City(ip)
	if err != nil || record == nil || record.City.Names["es"] == "" {
		return "Desconocido"
	}
	return record.City.Names["es"]
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	siteID := r.URL.Query().Get("site")
	token := r.URL.Query().Get("token")

	if siteID == "" || token == "" {
		http.Error(w, "Faltan parámetros", http.StatusBadRequest)
		return
	}

	// Verifica si el token es válido
	storedToken, err := db.RDB.HGet(db.Ctx, "site:"+siteID, "admin_token").Result()
	if err != nil || storedToken != token {
		http.Error(w, "Token inválido", http.StatusForbidden)
		return
	}

	// Obtener todos los eventos del sitio
	eventsRaw, err := db.RDB.LRange(db.Ctx, "events:"+siteID, 0, -1).Result()
	if err != nil {
		http.Error(w, "Error obteniendo eventos", http.StatusInternalServerError)
		return
	}

	modCount := map[string]int{}
	dayCount := map[string]int{}
	browserCount := map[string]int{}
	refererCount := map[string]int{}
	pageCount := map[string]int{}
	totalDuration := 0
	sessionCount := 0
	browserLangCount := map[string]int{}
	osCount := map[string]int{}
	cityCount := map[string]int{}

	for _, raw := range eventsRaw {
		var evt map[string]interface{}
		json.Unmarshal([]byte(raw), &evt)

		mod, _ := evt["module"].(string)
		ts, _ := evt["timestamp"].(string)

		// agrupar por módulo
		if mod != "" {
			modCount[mod]++
		}

		// agrupar por día
		t, _ := time.Parse(time.RFC3339, ts)
		day := t.Format("2006-01-02")
		dayCount[day]++

		// agrupar por navegador
		if browser, ok := evt["browser"].(string); ok && browser != "" {
			browserCount[browser]++
		}
		// agrupar por referencia
		if referer, ok := evt["referer"].(string); ok && referer != "" {
			refererCount[referer]++
		}
		// agrupar por página
		if page, ok := evt["page"].(string); ok && page != "" {
			pageCount[page]++
		}
		// calcular duración media de sesión
		if dur, ok := evt["duration_ms"].(float64); ok && dur > 0 {
			totalDuration += int(dur)
			sessionCount++
		}
		if bl, ok := evt["browser_lang"].(string); ok {
			browserLangCount[bl]++
		}
		if os, ok := evt["os"].(string); ok {
			osCount[os]++
		}
		if city, ok := evt["city"].(string); ok {
			cityCount[city]++
		}
	}

	var porModulo []moduloStat
	for m, total := range modCount {
		porModulo = append(porModulo, moduloStat{Modulo: m, Total: total})
	}

	var porDia []diaStat
	for d, total := range dayCount {
		porDia = append(porDia, diaStat{Dia: d, Total: total})
	}

	// Navegadores
	var navegadores []map[string]interface{}
	for b, total := range browserCount {
		navegadores = append(navegadores, map[string]interface{}{"navegador": b, "total": total})
	}
	if navegadores == nil {
		navegadores = []map[string]interface{}{}
	}
	// Referencias
	var referencias []map[string]interface{}
	for r, total := range refererCount {
		referencias = append(referencias, map[string]interface{}{"referencia": r, "total": total})
	}
	if referencias == nil {
		referencias = []map[string]interface{}{}
	}
	// Páginas
	var paginas []map[string]interface{}
	for p, total := range pageCount {
		paginas = append(paginas, map[string]interface{}{"pagina": p, "total": total})
	}
	if paginas == nil {
		paginas = []map[string]interface{}{}
	}
	// Duración media de sesión
	var duracionMedia int
	if sessionCount > 0 {
		duracionMedia = totalDuration / sessionCount
	}

	// Nuevas métricas
	dispositivos := map[string]int{}
	paises := map[string]int{}
	usuariosActivos := 0
	usuariosActivosWindow := time.Now().Add(-5 * time.Minute) // últimos 5 minutos

	// Detecta país de la IP de la petición actual (en memoria, no guardar)
	userIP := r.Header.Get("X-Forwarded-For")
	if userIP == "" {
		userIP, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	pais := countryFromIP(userIP)
	paises[pais]++

	for _, ev := range eventsRaw {
		var evt map[string]interface{}
		json.Unmarshal([]byte(ev), &evt)

		// Dispositivo
		if dev, ok := evt["device"].(string); ok {
			dispositivos[dev]++
		}
		// Usuarios activos (por timestamp)
		if tsStr, ok := evt["timestamp"].(string); ok {
			if ts, err := time.Parse(time.RFC3339, tsStr); err == nil && ts.After(usuariosActivosWindow) {
				usuariosActivos++
			}
		}
	}

	// Convertir a slice para el frontend
	var dispositivosArr []map[string]interface{}
	for k, v := range dispositivos {
		dispositivosArr = append(dispositivosArr, map[string]interface{}{"dispositivo": k, "total": v})
	}
	var paisesArr []map[string]interface{}
	for k, v := range paises {
		paisesArr = append(paisesArr, map[string]interface{}{"pais": k, "total": v})
	}
	var browserLangs []map[string]interface{}
	for k, v := range browserLangCount {
		browserLangs = append(browserLangs, map[string]interface{}{"lang": k, "total": v})
	}
	var osArr []map[string]interface{}
	for k, v := range osCount {
		osArr = append(osArr, map[string]interface{}{"os": k, "total": v})
	}
	var cities []map[string]interface{}
	for k, v := range cityCount {
		cities = append(cities, map[string]interface{}{"city": k, "total": v})
	}

	resp := map[string]interface{}{
		"por_modulo":     porModulo,
		"por_dia":        porDia,
		"navegadores":    navegadores,
		"referencias":    referencias,
		"paginas":        paginas,
		"duracion_media": duracionMedia,
		"dispositivos":     dispositivosArr,
		"paises":           paisesArr,
		"usuarios_activos": usuariosActivos,
		"browser_langs": browserLangs,
		"os":            osArr,
		"cities":        cities,
	}
	json.NewEncoder(w).Encode(resp)
}
