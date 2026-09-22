package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Reproduccion es el mensaje que viaja por la cola.
type Reproduccion struct {
	AudioID   string `json:"audio_id"`
	Titulo    string `json:"titulo"`
	ClienteIP string `json:"cliente_ip"`
	Timestamp string `json:"timestamp"`
}

// Consumidor escucha la cola de reproducciones y las imprime.
type Consumidor struct {
	url      string
	colaName string
}

// NewConsumidor construye el consumidor.
func NewConsumidor(url, colaName string) *Consumidor {
	return &Consumidor{url: url, colaName: colaName}
}

// Iniciar bloquea y consume la cola.
func (c *Consumidor) Iniciar() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("no se pudo conectar a RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declarar la cola (idempotente)
	if _, err := ch.QueueDeclare(
		c.colaName, // nombre
		true,       // durable
		false,      // auto-delete
		false,      // exclusive
		false,      // no-wait
		nil,        // args
	); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.colaName, // cola
		"",         // consumidor
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Printf("[Cola] Esperando mensajes en %q...\n", c.colaName)
	for msg := range msgs {
		var r Reproduccion
		if err := json.Unmarshal(msg.Body, &r); err != nil {
			log.Printf("Mensaje inválido: %v", err)
			continue
		}
		// ECO obligatorio de consumo desde cola
		fmt.Printf("[ECO-COLA][Estadisticas] Reproducción recibida: audio_id=%s titulo=%q cliente=%s ts=%s\n",
			r.AudioID, r.Titulo, r.ClienteIP, r.Timestamp)
	}
	return nil
}
