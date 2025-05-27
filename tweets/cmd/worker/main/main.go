package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	natsgo "github.com/nats-io/nats.go"

	"github.com/mcrosignani/uala_challenge/tweets/internal/container"
)

func main() {
	ctx := context.Background()
	deps, err := container.Build()
	if err != nil {
		panic(err)
	}
	defer deps.MessageService.Drain()

	_, err = deps.MessageService.Subscribe("tweets", func(msg *natsgo.Msg) {
		errHandle := deps.TweetHandler.HandleNewTweet(ctx, msg.Data)
		if errHandle != nil {
			log.Printf("error processing msg: %v", errHandle)
			return
		}
		msg.Ack()
	})
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
