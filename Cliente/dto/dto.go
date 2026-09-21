package dto

// TipoAudioDTO representa un tipo de audio recibido por REST.
type TipoAudioDTO struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// AudioDTO representa un audio recibido por REST.
type AudioDTO struct {
	ID       string            `json:"id"`
	TipoID   int               `json:"tipo_id"`
	Titulo   string            `json:"titulo"`
	Metadata map[string]string `json:"metadata"`
}
