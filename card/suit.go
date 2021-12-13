package card

type CardSuit int

const (
	Spade CardSuit = iota + 1
	Heart
	Diamond
	Club
	Joker
)

func (cs CardSuit) String() string {
	switch cs {
	case Spade:
		return "Spade"
	case Heart:
		return "Heart"
	case Diamond:
		return "Diamond"
	case Club:
		return "Club"
	case Joker:
		return "Joker"
	default:
		return "NoSuit"
	}
}
