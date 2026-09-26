package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"

	pb "lsd_rpc/proto/streaming"
)

// StreamingService encapsula la llamada gRPC al servidor de streaming.
type StreamingService struct {
	Direccion string
}

// NewStreamingService crea el servicio apuntando al servidor de streaming.
func NewStreamingService(direccion string) *StreamingService {
	return &StreamingService{Direccion: direccion}
}

// Reproducir descarga el audio por streaming y lo reproduce.
// Si se cierra el canal `stop`, la reproducción se detiene inmediatamente.
func (s *StreamingService) Reproducir(audioID string, stop <-chan struct{}) (string, error) {
	fmt.Printf("[ECO-GRPC] Conectando a %s para audio %s\n", s.Direccion, audioID)

	// ---- 1. Descarga por gRPC ----
	conn, err := grpc.NewClient(
		s.Direccion,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return "", fmt.Errorf("no se pudo conectar: %w", err)
	}
	defer conn.Close()

	cliente := pb.NewStreamingServiceClient(conn)
	stream, err := cliente.Reproducir(
		context.Background(),
		&pb.StreamRequest{AudioId: audioID},
	)
	if err != nil {
		return "", fmt.Errorf("error al iniciar stream: %w", err)
	}

	rutaSalida := fmt.Sprintf("/tmp/%s_reproducido.mp3", audioID)
	out, err := os.Create(rutaSalida)
	if err != nil {
		return "", err
	}

	var total int64
	var chunks int
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			out.Close()
			return "", fmt.Errorf("error recibiendo chunk: %w", err)
		}
		if _, err := out.Write(chunk.Data); err != nil {
			out.Close()
			return "", err
		}
		total += int64(len(chunk.Data))
		chunks++
	}
	out.Close()

	fmt.Printf("[ECO-GRPC] Recibidos %d chunks, %d bytes\n", chunks, total)

	// ---- 2. Reproducción ----
	f, err := os.Open(rutaSalida)
	if err != nil {
		return rutaSalida, fmt.Errorf("no se pudo abrir el archivo descargado: %w", err)
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return rutaSalida, fmt.Errorf("no se pudo decodificar el mp3: %w", err)
	}
	defer streamer.Close()

	fmt.Println("[Audio] Reproduciendo... (presione una tecla para detener)")

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	// ---- 3. Esperar fin de reproducción o señal de stop ----
	if stop == nil {
		<-done
		fmt.Println("\n[Audio] Reproducción terminada.")
		return rutaSalida, nil
	}

	select {
	case <-done:
		fmt.Println("\n[Audio] Reproducción terminada.")
	case <-stop:
		speaker.Clear()
		fmt.Println("\n[Audio] Reproducción detenida por el usuario.")
	}
	return rutaSalida, nil
}