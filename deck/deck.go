package deck

import (
	"go_poker/card"
	"math/rand"
	"time"
)

type Deck struct {
	Cards []card.Card
}

func NewDeck() *Deck {
	suits := []card.CardSuit{
		card.Spade,
		card.Heart,
		card.Diamond,
		card.Club,
	}

	numbers := []card.CardNumber{
		card.Ace,
		card.Two,
		card.Three,
		card.Four,
		card.Five,
		card.Six,
		card.Seven,
		card.Eight,
		card.Nine,
		card.Ten,
		card.Jack,
		card.Queen,
		card.King,
	}

	maxCardsCount := len(suits)*len(numbers) + 2
	cards := make([]card.Card, 0, maxCardsCount)

	for _, suit := range suits {
		for _, number := range numbers {
			cards = append(cards, card.Card{
				Suit:   suit,
				Number: number,
			})
		}
	}

	// cards = append(cards, card.Card{Suit: card.Joker}, card.Card{Suit: card.Joker})
	return &Deck{
		Cards: cards,
	}
}

func (d *Deck) Shuffle() *Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
	return d
}

func (d *Deck) Deal(n int) []card.Card {
	dealCards := d.Cards[0:n]
	d.Cards = d.Cards[n:]
	return dealCards
}

func (d *Deck) Count() int {
	return len(d.Cards)
}
