package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	frenchDeck "github.com/sean9999/GoCards/deck/french"
	"github.com/sean9999/GoCards/game/easypoker"
	"github.com/sean9999/rebouncer"
)

type PokerInfo struct {
	RowId int
	Cards easypoker.Cards
	Hand  easypoker.PokerHand
}

func (pi PokerInfo) String() string {
	return fmt.Sprintf("row:\t%d\nhand:\t%s\t(%s)\n", pi.RowId, pi.Cards.Strand(), pi.Hand.Grade)
}

func main() {

	//	our dirty source of events
	//	define this outside of ingestFunc so we can interrupt it with Ctrl+C
	randy := rand.NewSource(time.Now().UnixNano())
	cardsChan, done := frenchDeck.StreamCards(randy)

	//	Consume a stream of cards. Reject jokers. Make piles of 5. Send them to incomingEvents
	ingestFunc := func(incoming chan<- PokerInfo) {
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
	streamOfPokerHands := rebouncer.NewRebouncer[PokerInfo](
		ingestFunc,
		reduceFunc,
		quantizeFunc,
		1024,
	)

	//	Listen for Ctrl+C and perform shutdown
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		for range sigs {
			done <- true
			//streamOfPokerHands.Interrupt()
		}
	}()

	//	subscribe to rebouncer's OutgoingEvents channel
	for pokerHand := range streamOfPokerHands.Subscribe() {
		fmt.Println(pokerHand)
		fmt.Println("------------------")
	}

}
