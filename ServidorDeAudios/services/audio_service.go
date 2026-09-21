package services

import (
	"fmt"
	"mime/multipart"

	"lsd_rpc/ServidorDeAudios/storage"
	"lsd_rpc/comun/catalogo"
)

type AudioService struct {
	storage      *storage.AudioStorage
	rutaCatalogo string
}

func NewAudioService(s *storage.AudioStorage, rutaCatalogo string) *AudioService {
	return &AudioService{storage: s, rutaCatalogo: rutaCatalogo}
}

func (s *AudioService) Subir(tipoID int, titulo string, archivo multipart.File) (string, error) {
	prefijo := prefijoPorTipo(tipoID)
	if prefijo == "" {
		return "", fmt.Errorf("tipo_id inválido: %d", tipoID)
	}

	for i := 1; i < 1000; i++ {
		id := storage.FormatoID(prefijo, i)
		if !s.storage.Existe(id) {
			ruta, err := s.storage.Guardar(id, archivo)
			if err != nil {
				return "", err
			}
			fmt.Printf("[ECO-REST][Audios] Archivo guardado en %s (titulo=%q)\n", ruta, titulo)

			// Registrar en el catálogo compartido
			if err := s.registrarEnCatalogo(id, tipoID, titulo); err != nil {
				return id, fmt.Errorf("archivo guardado pero no registrado en catálogo: %w", err)
			}
			return id, nil
		}
	}
	return "", fmt.Errorf("no hay espacio para nuevos audios del tipo %d", tipoID)
}

func (s *AudioService) registrarEnCatalogo(id string, tipoID int, titulo string) error {
	c, err := catalogo.Leer(s.rutaCatalogo)
	if err != nil {
		return err
	}
	c.Audios = append(c.Audios, catalogo.Audio{
		ID:       id,
		TipoID:   tipoID,
		Titulo:   titulo,
		Metadata: map[string]string{"titulo": titulo},
	})
	fmt.Printf("[ECO-REST][Audios] Registrado en catálogo: id=%s tipo=%d\n", id, tipoID)
	return catalogo.Escribir(s.rutaCatalogo, c)
}

func prefijoPorTipo(tipoID int) string {
	switch tipoID {
	case 1:
		return "mus"
	case 2:
		return "pod"
	case 3:
		return "aud"
	case 4:
		return "rbl"
	}
	return ""
}
