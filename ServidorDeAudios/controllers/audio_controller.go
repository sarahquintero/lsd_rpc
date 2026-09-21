package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"lsd_rpc/ServidorDeAudios/services"
)

// AudioController expone el endpoint para que el admin suba audios.
type AudioController struct {
	servicio *services.AudioService
}

// NewAudioController construye el controlador.
func NewAudioController(s *services.AudioService) *AudioController {
	return &AudioController{servicio: s}
}

// RegistrarRutas asocia las rutas del controlador al router.
func (c *AudioController) RegistrarRutas(r *mux.Router) {
	r.HandleFunc("/admin/audios", c.SubirAudio).Methods("POST")
}

func echo(r *http.Request, msg string) {
	fmt.Printf("[ECO-REST][Audios] %s -> %s %s\n", r.RemoteAddr, r.Method, msg)
}

// SubirAudio espera multipart/form-data con:
//   - tipo_id (int, 1-4)
//   - titulo  (string)
//   - archivo (mp3)
func (c *AudioController) SubirAudio(w http.ResponseWriter, r *http.Request) {
	echo(r, "/admin/audios")

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "formulario inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	tipoIDStr := r.FormValue("tipo_id")
	titulo := r.FormValue("titulo")
	if tipoIDStr == "" || titulo == "" {
		http.Error(w, "faltan campos tipo_id o titulo", http.StatusBadRequest)
		return
	}

	var tipoID int
	if _, err := fmt.Sscanf(tipoIDStr, "%d", &tipoID); err != nil {
		http.Error(w, "tipo_id debe ser entero", http.StatusBadRequest)
		return
	}

	archivo, _, err := r.FormFile("archivo")
	if err != nil {
		http.Error(w, "falta archivo mp3: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer archivo.Close()

	id, err := c.servicio.Subir(tipoID, titulo, archivo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":     id,
		"titulo": titulo,
	})
}
