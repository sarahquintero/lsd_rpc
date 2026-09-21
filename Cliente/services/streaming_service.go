package services

import (
	"context"
	"fmt"
	"io"
	"os"

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
func (s *StreamingService) Reproducir(audioID string) (string, error) {
	fmt.Printf("[ECO-GRPC] Conectando a %s para audio %s\n", s.Direccion, audioID)

	conn, err := grpc.NewClient(s.Direccion, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	defer out.Close()

	var total int64
	var chunks int
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error recibiendo chunk: %w", err)
		}
		if _, err := out.Write(chunk.Data); err != nil {
			return "", err
		}
		total += int64(len(chunk.Data))
		chunks++
	}

	fmt.Printf("[ECO-GRPC] Recibidos %d chunks, %d bytes\n", chunks, total)
	return rutaSalida, nil
}
