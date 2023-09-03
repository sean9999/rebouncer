package rebouncer_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/sean9999/GoCards/deck/french"
	"github.com/sean9999/GoCards/game/easypoker"
	"github.com/sean9999/rebouncer"
)

func TestNewRebouncer(t *testing.T) {

	//	noisy source of cards
	type naughtyEvent french.Card

	//	semi-structured for for analyizing
	type niceEvent struct {
		Cards easypoker.Cards
		N     int64
		Poker easypoker.PokerHand
	}

	//	final form for outputting
	type beautifulEvent struct {
		Hand      easypoker.Cards
		PokerHand easypoker.PokerHand
		N         int64
	}

	randy := rand.NewSource(0)

	// eat up random cards from a random-card source
	// group them into hands of five and push those hands to the queue
	ingestCards := func(q chan<- niceEvent, doneChan chan bool) {
		done := make(chan bool)
		cardsChan := french.StreamCards(randy, done)
		cardsBuffer := make([]easypoker.Card, 0, 5)
		n := int64(0)
		for c := range cardsChan {
			n++
			easyPokerCard, err := easypoker.CardFromFrench(c)
			//	drop the Jokers
			if err != nil {
				cardsBuffer = append(cardsBuffer, easyPokerCard)
			}
			// group them into hands of five
			if len(cardsBuffer) == 5 {
				e := niceEvent{
					cardsBuffer, n, easypoker.HighestPokerHand(cardsBuffer),
				}
				q <- e
				cardsBuffer = cardsBuffer[:0]
			}
			// give up after 100_000 iterations
			if n > 100_000 {
				done <- true
			}
		}
		doneChan <- true
	}

	//	we're not interested in low hands
	removeLowHands := func(queue []niceEvent) []niceEvent {
		newQueue := make([]niceEvent, 0, len(queue))
		for _, hand := range queue {
			if easypoker.HighestPokerHand(hand.Cards).Grade >= easypoker.ThreeOfAKind {
				newQueue = append(newQueue, hand)
			}
		}
		return newQueue
	}

	//	quantizer is dead simple here:
	//	if there's anything in the queue,
	//	push it out
	// type QuantizeFunction[NICE any] func(chan<- bool, Queue[NICE])
	pushItRealGood := func(okchan chan<- bool, queue []niceEvent) {
		if len(queue) > 0 {
			okchan <- true
		}
	}

	//	emit doesn't need to do anything special
	passThrough := func(e niceEvent) niceEvent {
		return e
	}

	t.Run("create a rebouncer with three structs and no user-defined functions", func(t *testing.T) {

		rebecca := rebouncer.NewRebouncer[niceEvent](
			ingestCards,
			removeLowHands,
			pushItRealGood,
			passThrough,
			2048,
		)

		for hand := range rebecca.Subscribe() {
			fmt.Println(hand)
		}

	})

}
