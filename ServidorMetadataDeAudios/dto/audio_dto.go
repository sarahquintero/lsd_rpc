package dto

import "lsd_rpc/ServidorMetadataDeAudios/models"

// TipoAudioDTO es la representación JSON de TipoAudio.
type TipoAudioDTO struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// AudioDTO es la representación JSON de Audio.
type AudioDTO struct {
	ID       string            `json:"id"`
	TipoID   int               `json:"tipo_id"`
	Titulo   string            `json:"titulo"`
	Metadata map[string]string `json:"metadata"`
}

// ToTipoAudioDTO convierte un modelo a DTO.
func ToTipoAudioDTO(t models.TipoAudio) TipoAudioDTO {
	return TipoAudioDTO{ID: t.ID, Nombre: t.Nombre}
}

// ToAudioDTO convierte un modelo a DTO.
func ToAudioDTO(a models.Audio) AudioDTO {
	return AudioDTO{
		ID:       a.ID,
		TipoID:   a.TipoID,
		Titulo:   a.Titulo,
		Metadata: a.Metadata,
	}
}
