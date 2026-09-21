package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"lsd_rpc/ServidorMetadataDeAudios/services"
)

// AudioController expone los endpoints REST del servidor de metadatos.
type AudioController struct {
	servicio *services.AudioService
}

// NewAudioController construye el controlador.
func NewAudioController(s *services.AudioService) *AudioController {
	return &AudioController{servicio: s}
}

// RegistrarRutas asocia las rutas del controlador al router.
func (c *AudioController) RegistrarRutas(r *mux.Router) {
	r.HandleFunc("/tipos", c.GetTipos).Methods("GET")
	r.HandleFunc("/tipos/{id}/audios", c.GetAudiosPorTipo).Methods("GET")
	r.HandleFunc("/audios/{id}", c.GetAudio).Methods("GET")
	r.HandleFunc("/refrescar", c.Refrescar).Methods("POST")
}

// echo imprime el eco obligatorio de cada llamada REST recibida.
func echo(r *http.Request, msg string) {
	fmt.Printf("[ECO-REST][Metadata] %s -> %s %s\n", r.RemoteAddr, r.Method, msg)
}

// GetTipos devuelve la lista de tipos de audio.
func (c *AudioController) GetTipos(w http.ResponseWriter, r *http.Request) {
	echo(r, "/tipos")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c.servicio.ObtenerTipos())
}

// GetAudiosPorTipo devuelve los audios de un tipo.
func (c *AudioController) GetAudiosPorTipo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}
	echo(r, fmt.Sprintf("/tipos/%d/audios", id))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c.servicio.ObtenerAudiosPorTipo(id))
}

// GetAudio devuelve el detalle de un audio.
func (c *AudioController) GetAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	echo(r, fmt.Sprintf("/audios/%s", id))
	audio, ok := c.servicio.ObtenerAudioPorID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(audio)
}

// Refrescar recarga el catálogo en memoria desde el archivo JSON.
func (c *AudioController) Refrescar(w http.ResponseWriter, r *http.Request) {
	echo(r, "/refrescar")
	if err := c.servicio.Refrescar(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}
