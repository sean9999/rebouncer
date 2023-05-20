package rebouncer_test

type Suit uint8

const (
	IllegalSuit Suit = iota // zero-value should be illegal to protect against accidental values
	Diamonds
	Clubs
	Hearts
	Spades
)

type Face uint8

const (
	ZeroIsIllegal Face = iota // zero-value is illegal
	Joker
	Two // face-values correspond to integer values, allowing easy calculation
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace // we allow ourselves to assume that ace is high
)

type Card struct {
	Suit
	Face
}

type Deck [52]Card

type Hand []Card

type Pile []Card
