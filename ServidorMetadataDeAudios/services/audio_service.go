package services

import (
	"lsd_rpc/ServidorMetadataDeAudios/dto"
	"lsd_rpc/ServidorMetadataDeAudios/repository"
)

// AudioService contiene la lógica de negocio del servidor de metadatos.
type AudioService struct {
	repo repository.AudioRepository
}

// NewAudioService crea el servicio a partir de un repositorio.
func NewAudioService(repo repository.AudioRepository) *AudioService {
	return &AudioService{repo: repo}
}

// ObtenerTipos devuelve todos los tipos de audio como DTOs.
func (s *AudioService) ObtenerTipos() []dto.TipoAudioDTO {
	tipos := s.repo.GetTipos()
	salida := make([]dto.TipoAudioDTO, 0, len(tipos))
	for _, t := range tipos {
		salida = append(salida, dto.ToTipoAudioDTO(t))
	}
	return salida
}

// ObtenerAudiosPorTipo retorna los audios de un tipo dado.
func (s *AudioService) ObtenerAudiosPorTipo(tipoID int) []dto.AudioDTO {
	audios := s.repo.GetAudiosPorTipo(tipoID)
	salida := make([]dto.AudioDTO, 0, len(audios))
	for _, a := range audios {
		salida = append(salida, dto.ToAudioDTO(a))
	}
	return salida
}

// ObtenerAudioPorID devuelve un audio o (nil, false) si no existe.
func (s *AudioService) ObtenerAudioPorID(id string) (dto.AudioDTO, bool) {
	a, ok := s.repo.GetAudioPorID(id)
	if !ok {
		return dto.AudioDTO{}, false
	}
	return dto.ToAudioDTO(a), true
}

// Refrescar recarga el catálogo si el repositorio lo soporta.
func (s *AudioService) Refrescar() error {
	type refrescable interface{ Refrescar() error }
	if r, ok := s.repo.(refrescable); ok {
		return r.Refrescar()
	}
	return nil
}
