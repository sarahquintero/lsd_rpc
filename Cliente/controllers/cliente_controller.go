package controllers

import (
	"bufio"
	"fmt"
	"os"

	"lsd_rpc/Cliente/dto"
	"lsd_rpc/Cliente/services"
	"lsd_rpc/Cliente/views"
)

// ClienteController orquesta el menú y las llamadas a servicios.
type ClienteController struct {
	metadata     *services.MetadataService
	streaming    *services.StreamingService
	estadisticas *services.EstadisticasService
	reader       *bufio.Reader
}

// NewClienteController construye el controlador del cliente.
func NewClienteController(m *services.MetadataService, s *services.StreamingService, e *services.EstadisticasService) *ClienteController {
	return &ClienteController{
		metadata:     m,
		streaming:    s,
		estadisticas: e,
		reader:       bufio.NewReader(os.Stdin),
	}
}

// Ejecutar arranca el bucle principal del cliente.
func (c *ClienteController) Ejecutar() {
	for {
		views.LimpiarPantalla()
		fmt.Println("Spotify")
		fmt.Println("1. Ver tipos de audio")
		fmt.Println("2. Salir")

		switch views.LeerOpcion(c.reader, "> ") {
		case 1:
			c.menuTipos()
		case 2:
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func (c *ClienteController) menuTipos() {
	tipos, err := c.metadata.GetTipos()
	if err != nil {
		fmt.Println("Error consultando tipos:", err)
		return
	}

	for {
		views.LimpiarPantalla()
		views.MostrarTipos(tipos)
		op := views.LeerOpcion(c.reader, "> ")

		if op == len(tipos)+1 || op <= 0 {
			return
		}
		if op > len(tipos) {
			fmt.Println("Opción inválida")
			continue
		}
		c.menuAudios(tipos[op-1])
	}
}

func (c *ClienteController) menuAudios(tipo dto.TipoAudioDTO) {
	audios, err := c.metadata.GetAudiosPorTipo(tipo.ID)
	if err != nil {
		fmt.Println("Error consultando audios:", err)
		return
	}

	for {
		views.LimpiarPantalla()
		fmt.Println("Spotify")
		views.MostrarAudios(audios)
		op := views.LeerOpcion(c.reader, "> ")

		if op == len(audios)+1 || op <= 0 {
			return
		}
		if op > len(audios) {
			fmt.Println("Opción inválida")
			continue
		}
		c.menuDetalle(audios[op-1])
	}
}

func (c *ClienteController) menuDetalle(audio dto.AudioDTO) {
	for {
		views.LimpiarPantalla()
		views.MostrarDetalleAudio(audio)
		op := views.LeerOpcion(c.reader, "> ")

		switch op {
		case 1:
			c.menuReproducir(audio)
		case 2:
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func (c *ClienteController) menuReproducir(audio dto.AudioDTO) {
	views.LimpiarPantalla()
	fmt.Printf("Spotify\n%s\n\n", audio.Titulo)
	fmt.Println("Reproduciendo audio")
	fmt.Println()
	fmt.Println("1. Salir")

	go func() {
		if err := c.estadisticas.Publicar(audio.ID, audio.Titulo); err != nil {
			fmt.Println("[Cola] Error publicando:", err)
		}
	}()

	ruta, err := c.streaming.Reproducir(audio.ID)
	if err != nil {
		fmt.Println("Error de streaming:", err)
	} else {
		fmt.Println("Audio guardado en:", ruta)
	}
	views.LeerOpcion(c.reader, "> ")
}
