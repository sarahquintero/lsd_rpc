package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"lsd_rpc/Cliente/dto"
)

// MetadataService encapsula las llamadas REST al servidor de metadatos.
type MetadataService struct {
	BaseURL string
}

// NewMetadataService crea el servicio apuntando a la URL indicada.
func NewMetadataService(baseURL string) *MetadataService {
	return &MetadataService{BaseURL: baseURL}
}

// GetTipos obtiene la lista de tipos de audio.
func (m *MetadataService) GetTipos() ([]dto.TipoAudioDTO, error) {
	fmt.Printf("[ECO-REST] GET %s/tipos\n", m.BaseURL)
	resp, err := http.Get(m.BaseURL + "/tipos")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tipos []dto.TipoAudioDTO
	if err := json.NewDecoder(resp.Body).Decode(&tipos); err != nil {
		return nil, err
	}
	return tipos, nil
}

// GetAudiosPorTipo obtiene la lista de audios de un tipo.
func (m *MetadataService) GetAudiosPorTipo(tipoID int) ([]dto.AudioDTO, error) {
	url := fmt.Sprintf("%s/tipos/%d/audios", m.BaseURL, tipoID)
	fmt.Printf("[ECO-REST] GET %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var audios []dto.AudioDTO
	if err := json.NewDecoder(resp.Body).Decode(&audios); err != nil {
		return nil, err
	}
	return audios, nil
}

// GetAudioPorID obtiene el detalle de un audio.
func (m *MetadataService) GetAudioPorID(id string) (dto.AudioDTO, error) {
	url := fmt.Sprintf("%s/audios/%s", m.BaseURL, id)
	fmt.Printf("[ECO-REST] GET %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return dto.AudioDTO{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return dto.AudioDTO{}, fmt.Errorf("estado %d: %s", resp.StatusCode, string(body))
	}

	var audio dto.AudioDTO
	if err := json.NewDecoder(resp.Body).Decode(&audio); err != nil {
		return dto.AudioDTO{}, err
	}
	return audio, nil
}
