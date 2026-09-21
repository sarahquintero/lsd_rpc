package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"lsd_rpc/ServidorDeStreaming/controllers"
	"lsd_rpc/ServidorDeStreaming/repository"
	pb "lsd_rpc/proto/streaming"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("No se pudo escuchar en :50051: %v", err)
	}

	repo := repository.NewAudioRepository("audios")
	ctrl := controllers.NewStreamingController(repo)

	grpcServer := grpc.NewServer()
	pb.RegisterStreamingServiceServer(grpcServer, ctrl)

	fmt.Println("Servidor de streaming gRPC escuchando en :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error en el servidor gRPC: %v", err)
	}
}
