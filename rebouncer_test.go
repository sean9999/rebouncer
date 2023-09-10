package rebouncer // import rebouncer

import (
	"fmt"
	"math/rand"
	"testing"

	frenchDeck "github.com/sean9999/GoCards/deck/french"
	"github.com/sean9999/GoCards/game/easypoker"

	"time"
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
		//done := make(chan bool)
		cardsChan, done := frenchDeck.StreamCards(randy)
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

	t.Run("create a rebouncer with three structs and no user-defined functions", func(t *testing.T) {

		rebecca := NewRebouncer[pokerInfo](
			ingestCards,
			removeLowHands,
			pushItRealGood,
			1024,
		)

		for hand := range rebecca.Subscribe() {
			fmt.Println(hand)
		}

	})

}

type PokerInfo struct {
	RowId int
	Cards easypoker.Cards
	Hand  easypoker.PokerHand
}

func (pi PokerInfo) String() string {
	return fmt.Sprintf("row:\t%d\nhand:\t%s\t(%s)\n", pi.RowId, pi.Cards.Strand(), pi.Hand.Grade)
}

func ExampleRebouncer() {

	//	This example ingests a source of randomly shuffled cards,
	//	excludes the jokers,
	//	batches them into hands of 5,
	//	and emits those hands which beat 🂷🃇🂧🃞🂣 (three sevens).

	//	Consume a stream of cards. Reject jokers. Make piles of 5. Send them to incomingEvents
	ingestFunc := func(incoming chan<- PokerInfo) {

		//done := make(chan bool)
		randy := rand.NewSource(time.Now().UnixNano())
		cardsChan, done := frenchDeck.StreamCards(randy)

		piles := make(chan easypoker.Card, 5)
		i := 0
		for card := range cardsChan {
			i++
			goodCard, err := easypoker.CardFromFrench(card)
			if err == nil {
				piles <- goodCard
				if len(piles) == 5 {
					fiveCards := []easypoker.Card{
						<-piles, <-piles, <-piles, <-piles, <-piles,
					}
					pi := PokerInfo{
						RowId: i,
						Cards: fiveCards,
					}
					incoming <- pi
					if i > 10_000 {
						done <- true // signal to StreamCards
					}
				}
			}
		}
	}

	//	reducer. Omit any hand that doesn't beat 3 sevens
	reduceFunc := func(oldcards []PokerInfo) []PokerInfo {
		goodHands := make([]PokerInfo, 0, len(oldcards))
		lowHand, _ := easypoker.HandFromString("🂷🃇🂧🃞🂣")
		for _, thisHand := range oldcards {
			if thisHand.Cards.Beats(lowHand) {
				thisHand.Hand = easypoker.HighestPokerHand(thisHand.Cards)
				goodHands = append(goodHands, thisHand)
			}
		}
		return goodHands
	}

	//	quantize. Wait before flushing
	quantizeFunc := func(stuff []PokerInfo) bool {
		time.Sleep(time.Millisecond * 100)
		return (len(stuff) > 0)
	}

	//	invoke rebouncer
	streamOfPokerHands := NewRebouncer[PokerInfo](
		ingestFunc,
		reduceFunc,
		quantizeFunc,
		1024,
	)

	//	subscribe to rebouncer's OutgoingEvents channel
	for pokerHand := range streamOfPokerHands.Subscribe() {
		fmt.Println(pokerHand)
		fmt.Println("------------------")
	}

}
