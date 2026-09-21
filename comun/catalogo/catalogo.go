package catalogo

import (
	"encoding/json"
	"os"
	"sync"
)

// TipoAudio tal como se guarda en el JSON.
type TipoAudio struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// Audio tal como se guarda en el JSON.
type Audio struct {
	ID       string            `json:"id"`
	TipoID   int               `json:"tipo_id"`
	Titulo   string            `json:"titulo"`
	Metadata map[string]string `json:"metadata"`
}

// Catalogo es el documento completo.
type Catalogo struct {
	Tipos  []TipoAudio `json:"tipos"`
	Audios []Audio     `json:"audios"`
}

var mu sync.Mutex

// Ruta por defecto del catálogo.
const RutaPorDefecto = "catalogo.json"

// Leer carga el catálogo desde disco.
func Leer(ruta string) (*Catalogo, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(ruta)
	if err != nil {
		return nil, err
	}
	var c Catalogo
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Escribir guarda el catálogo en disco.
func Escribir(ruta string, c *Catalogo) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ruta, data, 0o644)
}
