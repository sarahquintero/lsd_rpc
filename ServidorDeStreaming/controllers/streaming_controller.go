package controllers

import (
	"fmt"
	"io"

	"lsd_rpc/ServidorDeStreaming/repository"
	pb "lsd_rpc/proto/streaming"
)

const tamanoChunk = 64 * 1024 // 64 KB por chunk

// StreamingController implementa el servicio gRPC de streaming.
type StreamingController struct {
	pb.UnimplementedStreamingServiceServer
	repo *repository.AudioRepository
}

// NewStreamingController construye el controlador.
func NewStreamingController(repo *repository.AudioRepository) *StreamingController {
	return &StreamingController{repo: repo}
}

// Reproducir transmite el audio en chunks al cliente.
func (c *StreamingController) Reproducir(req *pb.StreamRequest, stream pb.StreamingService_ReproducirServer) error {
	fmt.Printf("[ECO-GRPC][Streaming] Reproducir -> audio_id=%s\n", req.GetAudioId())

	f, total, err := c.repo.Abrir(req.GetAudioId())
	if err != nil {
		fmt.Printf("[ECO-GRPC][Streaming] Error abriendo audio: %v\n", err)
		return err
	}
	defer f.Close()

	buffer := make([]byte, tamanoChunk)
	var chunkID int32
	for {
		n, err := f.Read(buffer)
		if n > 0 {
			chunk := &pb.AudioChunk{
				Data:      buffer[:n],
				ChunkId:   chunkID,
				TotalSize: total,
			}
			if err := stream.Send(chunk); err != nil {
				fmt.Printf("[ECO-GRPC][Streaming] Error enviando chunk %d: %v\n", chunkID, err)
				return err
			}
			chunkID++
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("[ECO-GRPC][Streaming] Error leyendo archivo: %v\n", err)
			return err
		}
	}

	fmt.Printf("[ECO-GRPC][Streaming] Fin de reproducción -> %d chunks enviados\n", chunkID)
	return nil
}
