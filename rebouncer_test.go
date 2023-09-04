package rebouncer_test

import (
	"fmt"
	"math/rand"
	"testing"

	frenchDeck "github.com/sean9999/GoCards/deck/french"
	"github.com/sean9999/GoCards/game/easypoker"
	"github.com/sean9999/rebouncer"
)

func TestNewRebouncer(t *testing.T) {

	//	final form for outputting
	type pokerInfo struct {
		Cards     easypoker.Cards
		PokerHand easypoker.PokerHand
		N         int64
	}

	randy := rand.NewSource(0)

	// eat up random cards from a random-card source
	// group them into hands of five and push those hands to the queue
	ingestCards := func(q chan<- pokerInfo) {
		done := make(chan bool)
		cardsChan := frenchDeck.StreamCards(randy, done)
		cardsBuffer := make([]easypoker.Card, 0, 5)
		n := int64(0)
		for c := range cardsChan {
			n++
			easyPokerCard, err := easypoker.CardFromFrench(c)
			//	drop the Jokers
			if err == nil {
				cardsBuffer = append(cardsBuffer, easyPokerCard)
			}
			// group them into hands of five
			if len(cardsBuffer) == 5 {
				e := pokerInfo{
					cardsBuffer, easypoker.HighestPokerHand(cardsBuffer), n,
				}
				q <- e
				cardsBuffer = cardsBuffer[:0]
			}
			// give up after 100_000 iterations
			if n > int64(1000) {
				done <- true
			}
		}
	}

	//	we're not interested in low hands
	removeLowHands := func(queue []pokerInfo) []pokerInfo {
		newQueue := make([]pokerInfo, 0, len(queue))
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
	pushItRealGood := func(queue []pokerInfo) bool {
		okToEmit := (len(queue) > 0)
		return okToEmit
	}

	//	emit doesn't need to do anything special
	passThrough := func(e pokerInfo) pokerInfo {
		return e
	}

	t.Run("create a rebouncer with three structs and no user-defined functions", func(t *testing.T) {

		rebecca := rebouncer.NewRebouncer[pokerInfo](
			ingestCards,
			removeLowHands,
			pushItRealGood,
			passThrough,
			1024,
		)

		for hand := range rebecca.Subscribe() {
			fmt.Println(hand)
		}

	})

}
