package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

// Reproducir solicita el audio por streaming y lo escribe en un archivo temporal.
// Devuelve la ruta del archivo reproducido.
// Reproducir solicita el audio por streaming, lo guarda y lo reproduce.
func (s *StreamingService) Reproducir(audioID string) (string, error) {
	fmt.Printf("[ECO-GRPC] Conectando a %s para audio %s\n", s.Direccion, audioID)

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

	// Cerrar el archivo antes de abrirlo para reproducir
	out.Close()

	fmt.Printf("[ECO-GRPC] Recibidos %d chunks, %d bytes\n", chunks, total)

	// ---- Reproducción con beep ----
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

	fmt.Println("[Audio] Reproduciendo...")

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done

	fmt.Println("[Audio] Reproducción terminada.")
	return rutaSalida, nil
}
