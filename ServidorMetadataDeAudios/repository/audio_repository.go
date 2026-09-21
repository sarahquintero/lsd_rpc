package repository

import (
	"sync"

	"lsd_rpc/ServidorMetadataDeAudios/models"
	"lsd_rpc/comun/catalogo"
)

type AudioRepository interface {
	GetTipos() []models.TipoAudio
	GetAudiosPorTipo(tipoID int) []models.Audio
	GetAudioPorID(id string) (models.Audio, bool)
	Refrescar() error
}

type jsonRepository struct {
	mu     sync.RWMutex
	ruta   string
	tipos  []models.TipoAudio
	audios []models.Audio
}

// NewAudioRepository carga el catálogo desde disco.
func NewAudioRepository(ruta string) (AudioRepository, error) {
	r := &jsonRepository{ruta: ruta}
	if err := r.Refrescar(); err != nil {
		return nil, err
	}
	return r, nil
}

// Refrescar recarga el catálogo desde el archivo.
func (r *jsonRepository) Refrescar() error {
	c, err := catalogo.Leer(r.ruta)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tipos = make([]models.TipoAudio, 0, len(c.Tipos))
	for _, t := range c.Tipos {
		r.tipos = append(r.tipos, models.TipoAudio{ID: t.ID, Nombre: t.Nombre})
	}

	r.audios = make([]models.Audio, 0, len(c.Audios))
	for _, a := range c.Audios {
		r.audios = append(r.audios, models.Audio{
			ID:       a.ID,
			TipoID:   a.TipoID,
			Titulo:   a.Titulo,
			Metadata: a.Metadata,
		})
	}
	return nil
}

func (r *jsonRepository) GetTipos() []models.TipoAudio {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]models.TipoAudio(nil), r.tipos...)
}

func (r *jsonRepository) GetAudiosPorTipo(tipoID int) []models.Audio {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtrados []models.Audio
	for _, a := range r.audios {
		if a.TipoID == tipoID {
			filtrados = append(filtrados, a)
		}
	}
	return filtrados
}

func (r *jsonRepository) GetAudioPorID(id string) (models.Audio, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.audios {
		if a.ID == id {
			return a, true
		}
	}
	return models.Audio{}, false
}
