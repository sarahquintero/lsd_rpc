package controllers

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/sys/unix"

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
	fmt.Println("Presione cualquier tecla para detener")
	fmt.Println()

	// Publicar en la cola (asíncrono, no bloquea)
	go func() {
		if err := c.estadisticas.Publicar(audio.ID, audio.Titulo); err != nil {
			fmt.Println("[Cola] Error publicando:", err)
		}
	}()

	fd := int(os.Stdin.Fd())

	// Modo cbreak: sin ICANON ni ECHO, pero conservando OPOST para que \n funcione
	oldState, err := hacerCbreak(fd)
	if err != nil {
		fmt.Println("[Audio] (No se pudo activar modo cbreak; no se podrá detener)")
		if _, err := c.streaming.Reproducir(audio.ID, nil); err != nil {
			fmt.Println("Error:", err)
		}
		return
	}
	defer restaurarTerminal(fd, oldState)

	stop := make(chan struct{})
	done := make(chan struct{})

	go func() {
		defer close(done)
		if _, err := c.streaming.Reproducir(audio.ID, stop); err != nil {
			fmt.Println("Error:", err)
		}
	}()

	if esperarTeclaOFinal(done) {
		close(stop)
		<-done
	}

	drenarStdin(fd)
}

// hacerCbreak pone el terminal en modo cbreak (sin ICANON ni ECHO),
// pero mantiene OPOST para que \n se siga traduciendo a \r\n.
func hacerCbreak(fd int) (*unix.Termios, error) {
	old, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}
	estado := *old
	estado.Lflag &^= unix.ICANON | unix.ECHO
	estado.Cc[unix.VMIN] = 1
	estado.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &estado); err != nil {
		return nil, err
	}
	return old, nil
}

// restaurarTerminal devuelve el terminal a su estado anterior.
func restaurarTerminal(fd int, old *unix.Termios) {
	unix.IoctlSetTermios(fd, unix.TCSETS, old)
}

// esperarTeclaOFinal espera una tecla o el fin de la reproducción.
// Devuelve true si el usuario presionó una tecla.
func esperarTeclaOFinal(done <-chan struct{}) bool {
	fd := int(os.Stdin.Fd())
	for {
		select {
		case <-done:
			return false
		default:
		}

		var fds unix.FdSet
		fds.Bits[fd/64] |= 1 << (uint(fd) % 64)
		tv := unix.Timeval{Usec: 100000} // 100 ms

		n, err := unix.Select(fd+1, &fds, nil, nil, &tv)
		if err == nil && n > 0 {
			buf := make([]byte, 1)
			unix.Read(fd, buf)
			return true
		}
	}
}

// drenarStdin consume bytes pendientes para no romper el siguiente bufio.Reader.
func drenarStdin(fd int) {
	for {
		var fds unix.FdSet
		fds.Bits[fd/64] |= 1 << (uint(fd) % 64)
		tv := unix.Timeval{Usec: 10000} // 10 ms

		n, err := unix.Select(fd+1, &fds, nil, nil, &tv)
		if err != nil || n == 0 {
			return
		}
		buf := make([]byte, 64)
		unix.Read(fd, buf)
	}
}
