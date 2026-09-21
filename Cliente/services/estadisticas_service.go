package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EstadisticasService publica eventos de reproducción en la cola.
type EstadisticasService struct {
	url      string
	colaName string
}

// NewEstadisticasService construye el servicio.
func NewEstadisticasService(url, colaName string) *EstadisticasService {
	return &EstadisticasService{url: url, colaName: colaName}
}

// Mensaje es el payload publicado.
type Mensaje struct {
	AudioID   string `json:"audio_id"`
	Titulo    string `json:"titulo"`
	ClienteIP string `json:"cliente_ip"`
	Timestamp string `json:"timestamp"`
}

// Publicar envía el evento a la cola sin bloquear al cliente.
func (s *EstadisticasService) Publicar(audioID, titulo string) error {
	conn, err := amqp.Dial(s.url)
	if err != nil {
		return fmt.Errorf("no se pudo conectar a RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(s.colaName, true, false, false, false, nil); err != nil {
		return err
	}

	msg := Mensaje{
		AudioID:   audioID,
		Titulo:    titulo,
		ClienteIP: "127.0.0.1",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(msg)

	fmt.Printf("[ECO-COLA] Publicando en %q -> audio_id=%s\n", s.colaName, audioID)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(
		ctx,
		"",         // exchange
		s.colaName, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
