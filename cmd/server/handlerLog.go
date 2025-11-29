package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerLog() func(routing.GameLog) pubsub.AckType {
	return func(gamelog routing.GameLog) pubsub.AckType {
		defer fmt.Print("> ")
		if err := gamelogic.WriteLog(gamelog); err != nil {
			log.Printf("error writing log to disk: %v\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
