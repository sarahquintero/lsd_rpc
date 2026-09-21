package main

import (
	"lsd_rpc/Cliente/controllers"
	"lsd_rpc/Cliente/services"
)

func main() {
	metadata := services.NewMetadataService("http://localhost:8081")
	streaming := services.NewStreamingService("localhost:50051")
	estadisticas := services.NewEstadisticasService(
		"amqp://guest:guest@localhost:5672/",
		"reproducciones",
	)

	ctrl := controllers.NewClienteController(metadata, streaming, estadisticas)
	ctrl.Ejecutar()
}
