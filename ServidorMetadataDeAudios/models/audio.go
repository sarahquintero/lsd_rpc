package models

// TipoAudio representa una categoría de audio registrada en el sistema.
type TipoAudio struct {
	ID     int
	Nombre string
}

// Audio representa un audio con sus metadatos.
type Audio struct {
	ID       string
	TipoID   int
	Titulo   string
	Metadata map[string]string
}
