package hand

type HandPoint int

const (
	OnePair HandPoint = iota + 1
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	AFullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (hp HandPoint) String() string {
	switch hp {
	case OnePair:
		return "OnePair"
	case TwoPair:
		return "TwoPair"
	case ThreeOfAKind:
		return "ThreeOfAKind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case AFullHouse:
		return "AFullHouse"
	case FourOfAKind:
		return "FourOfAKind"
	case StraightFlush:
		return "StraightFlush"
	case RoyalFlush:
		return "RoyalFlush"
	default:
		return "HighCard"
	}
}