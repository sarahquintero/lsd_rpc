package main

import (
	"log"

	"lsd_rpc/ServidorDeEstadisticas/consumer"
)

func main() {
	const (
		rabbitURL = "amqp://guest:guest@localhost:5672/"
		colaName  = "reproducciones"
	)

	c := consumer.NewConsumidor(rabbitURL, colaName)
	if err := c.Iniciar(); err != nil {
		log.Fatalf("Error en consumidor: %v", err)
	}
}
