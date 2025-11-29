package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)
		var message string
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			message = fmt.Sprintf("%s won a war against %s", winner, loser)
			return gamelogic.PublishGameLog(ch, routing.GameLog{
				CurrentTime: time.Now(),
				Message:     message,
				Username:    rw.Attacker.Username,
			})
		case gamelogic.WarOutcomeYouWon:
			message = fmt.Sprintf("%s won a war against %s", winner, loser)
			return gamelogic.PublishGameLog(ch, routing.GameLog{
				CurrentTime: time.Now(),
				Message:     message,
				Username:    rw.Attacker.Username,
			})
		case gamelogic.WarOutcomeDraw:
			message = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			return gamelogic.PublishGameLog(ch, routing.GameLog{
				CurrentTime: time.Now(),
				Message:     message,
				Username:    rw.Attacker.Username,
			})
		default:
			fmt.Println("unknown war outcome")
			return pubsub.NackDiscard
		}
	}
}
