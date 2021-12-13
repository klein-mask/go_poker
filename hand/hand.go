package hand

import (
	"go_poker/card"
	"go_poker/util"
	"reflect"
	"sort"
)

type Hand struct {
	Cards []card.Card
	Point HandPoint
	AddedFlopCards []card.Card
}

func NewHand(cards []card.Card) *Hand {
	h := &Hand{ Cards: cards }
	return h
}

func (h *Hand) Add(addCards []card.Card) *Hand {
	h.Cards = append(h.Cards, addCards...)
	return h
}

func (h *Hand) Numbers() []card.CardNumber {
	results := make([]card.CardNumber, 0, len(h.AddedFlopCards))

	for _, card := range h.AddedFlopCards {
		results = append(results, card.Number)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})
	return results
}

func (h *Hand) Suits() []card.CardSuit {
	results := make([]card.CardSuit, 0, len(h.AddedFlopCards))

	for _, card := range h.AddedFlopCards {
		results = append(results, card.Suit)
	}
	return results
}

func (h *Hand) Culc(flopCards []card.Card) *Hand {
	h.AddedFlopCards = append(h.Cards, flopCards...)
	gloupByNumberCards := h.GloupByNumber()
	gloupByNumberCardLengths := make([]int, 0, len(gloupByNumberCards))
	for _, v := range gloupByNumberCards {
		gloupByNumberCardLengths = append(gloupByNumberCardLengths, len(v))
	}

	var maxpairOfNumberCards int
	if len(gloupByNumberCards) > 0 {
		sort.Sort(sort.IntSlice(gloupByNumberCardLengths))
		maxpairOfNumberCards = gloupByNumberCardLengths[len(gloupByNumberCardLengths) - 1]
	}

	if h.IsRoyalFlush() {
		h.Point = RoyalFlush
	} else if h.IsFlush() && h.IsStraight() {
		h.Point = StraightFlush
	} else if maxpairOfNumberCards == 4 {
		h.Point = FourOfAKind
	} else if util.Contains(gloupByNumberCardLengths, 3) && util.Contains(gloupByNumberCardLengths, 2) {
		h.Point = AFullHouse
	} else if h.IsFlush() {
		h.Point = Flush
	} else if h.IsStraight() {
		h.Point = Straight
	} else if maxpairOfNumberCards == 3 {
		h.Point = ThreeOfAKind
	} else if len(gloupByNumberCards) == 2 {
		h.Point = TwoPair
	} else if len(gloupByNumberCards) == 1 {
		h.Point = OnePair
	}

	return h
}

func (h *Hand) GloupByNumber() map[string][]card.Card {
	multipleCheckNumbers := make(map[string][]card.CardNumber, len(h.AddedFlopCards))
	results := make(map[string][]card.Card, len(h.AddedFlopCards))
	for _, card := range h.AddedFlopCards {
		cardNumMapKey := card.Number.ToString()
		if len(results[cardNumMapKey]) <= 0 || util.Contains(multipleCheckNumbers[cardNumMapKey], card.Number) {
			multipleCheckNumbers[cardNumMapKey] = append(multipleCheckNumbers[cardNumMapKey], card.Number)
			results[cardNumMapKey] = append(results[cardNumMapKey], card)
		}
	}

	for k, v := range results {
		if len(v) <= 1 {
			delete(results, k)
		}
	}

	return results
}

func (h *Hand) GloupBySUit() map[string][]card.Card {
	results := make(map[string][]card.Card, len(h.AddedFlopCards))
	for _, card := range h.AddedFlopCards {
		results[card.Suit.String()] = append(results[card.Suit.String()], card)
	}
	return results
}

func (h *Hand) IsStraight() bool {
	numbers := h.Numbers()
	for i := 0; i < len(numbers); i++ {
		if i >= len(numbers) - 1 {
			return true
		}
		if numbers[i] != (numbers[i + 1] - 1) && !(numbers[i] == card.Ace && numbers[i + 1] == card.Ten) {
			return false
		}
	}
	return false
}

func (h *Hand) IsFlush() bool {
	suits := h.Suits()
	for i := 0; i < len(suits); i++ {
		if i >= len(suits) - 1 {
			return true
		}
		if suits[i] != suits[i + 1] {
			return false
		}
	}
	return false
}

func (h *Hand) IsRoyalFlush() bool {
	if !h.IsStraight() || !h.IsFlush() {
		return false
	}
	numbers := h.Numbers()
	royalFlushNumbers := []card.CardNumber{
		card.Ace,
		card.Ten,
		card.Jack,
		card.Queen,
		card.King,
	}
	return reflect.DeepEqual(numbers, royalFlushNumbers)
}

func (h *Hand) GetHandInnerThreeCardGroup() []card.Card {
	// result := make([]card.Card, 0, len(h.AddedFlopCards))
	var results []card.Card
	for _, v := range h.GloupByNumber() {
		if len(v) >= 3 {
			results = v
		}
	}
	return results
}

//func (h *Hand) CulcForDrawHandPoint() {
//
//}