package nats

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/mcrosignani/uala_challenge/tweets/internal/config"
)

type (
	MessageService interface {
		Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error)
		Drain() error
		Publish(subject string, data []byte) error
	}

	messageService struct {
		NC *nats.Conn
		JS nats.JetStreamContext
	}
)

func NewNCService(cfg config.Config) (MessageService, error) {
	var nc *nats.Conn
	var err error

	connString := fmt.Sprintf("nats://%s:%d", cfg.Nats.Host, cfg.Nats.Port)
	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(connString)
		if err == nil {
			break
		}
		log.Println("wait NATS...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("could not connect to NATS: %v", err)
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("error context JetStream: %v", err)
		return nil, err
	}

	// Crear el stream si no existe
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "TWEETS",
		Subjects: []string{"tweets"},
		Storage:  nats.FileStorage,
	})
	if err != nil {
		log.Printf("error creating stream: %v", err)
	}

	return &messageService{
		NC: nc,
		JS: js,
	}, nil
}

func (ms *messageService) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	subscription, err := ms.JS.Subscribe(subject, handler, nats.Durable("worker-1"), nats.ManualAck(), nats.AckWait(10*time.Second))
	if err != nil {
		return nil, fmt.Errorf("error subscribing to subject %s: %w", subject, err)
	}

	return subscription, nil
}

func (ms *messageService) Drain() error {
	if err := ms.NC.Drain(); err != nil {
		return fmt.Errorf("error draining NATS connection: %w", err)
	}
	return nil
}

func (ms *messageService) Publish(subject string, data []byte) error {
	_, err := ms.JS.Publish(subject, data)
	if err != nil {
		return err
	}

	return nil
}
