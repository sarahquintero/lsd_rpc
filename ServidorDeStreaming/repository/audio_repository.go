package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

// AudioRepository resuelve ids de audio a rutas de archivos mp3.
type AudioRepository struct {
	BaseDir string
}

// NewAudioRepository crea el repositorio apuntando al directorio base.
func NewAudioRepository(baseDir string) *AudioRepository {
	return &AudioRepository{BaseDir: baseDir}
}

// Abrir devuelve el archivo correspondiente al audioID.
func (r *AudioRepository) Abrir(audioID string) (*os.File, int64, error) {
	ruta := filepath.Join(r.BaseDir, audioID+".mp3")
	f, err := os.Open(ruta)
	if err != nil {
		return nil, 0, fmt.Errorf("no se pudo abrir %s: %w", ruta, err)
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, 0, err
	}
	return f, info.Size(), nil
}
