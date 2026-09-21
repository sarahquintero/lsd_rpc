package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// AudioStorage guarda archivos mp3 en disco.
type AudioStorage struct {
	BaseDir string
}

// NewAudioStorage construye el storage apuntando a un directorio.
func NewAudioStorage(baseDir string) *AudioStorage {
	return &AudioStorage{BaseDir: baseDir}
}

// Guardar escribe el archivo subido en BaseDir/audioID.mp3.
func (s *AudioStorage) Guardar(audioID string, archivo multipart.File) (string, error) {
	if err := os.MkdirAll(s.BaseDir, 0o755); err != nil {
		return "", err
	}
	ruta := filepath.Join(s.BaseDir, audioID+".mp3")
	out, err := os.Create(ruta)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, archivo); err != nil {
		return "", err
	}
	return ruta, nil
}

// Existe indica si ya hay un archivo con ese id.
func (s *AudioStorage) Existe(audioID string) bool {
	ruta := filepath.Join(s.BaseDir, audioID+".mp3")
	if _, err := os.Stat(ruta); err == nil {
		return true
	}
	return false
}

// FormatoID construye el id final en base al tipo.
func FormatoID(prefijo string, contador int) string {
	return fmt.Sprintf("%s-%03d", prefijo, contador)
}
