package views

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lsd_rpc/Cliente/dto"
)

// LeerOpcion muestra un prompt y devuelve la opción elegida.
func LeerOpcion(reader *bufio.Reader, prompt string) int {
	fmt.Print(prompt)
	linea, _ := reader.ReadString('\n')
	linea = strings.TrimSpace(linea)
	n, err := strconv.Atoi(linea)
	if err != nil {
		return -1
	}
	return n
}

// MostrarTipos imprime la lista de tipos disponibles.
func MostrarTipos(tipos []dto.TipoAudioDTO) {
	fmt.Println("Spotify")
	for i, t := range tipos {
		fmt.Printf("%d. %s\n", i+1, t.Nombre)
	}
	fmt.Printf("%d. Atrás\n", len(tipos)+1)
}

// MostrarAudios imprime la lista de audios de un tipo.
func MostrarAudios(audios []dto.AudioDTO) {
	for i, a := range audios {
		fmt.Printf("%d. %s\n", i+1, a.Titulo)
	}
	fmt.Printf("%d. Atrás\n", len(audios)+1)
}

// MostrarDetalleAudio imprime los metadatos de un audio.
func MostrarDetalleAudio(a dto.AudioDTO) {
	fmt.Printf("Audio libro: %s\n\n", a.Titulo)
	for k, v := range a.Metadata {
		fmt.Printf("- %s: %s\n", k, v)
	}
	fmt.Println()
	fmt.Println("1. Reproducir")
	fmt.Println("2. Atrás")
}

// LimpiarPantalla intenta limpiar la consola (funciona en Linux).
func LimpiarPantalla() {
	fmt.Print("\033[H\033[2J")
	_ = os.Stdout.Sync()
}
