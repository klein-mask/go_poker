package poker

type Position int

const (
	Dealer Position = iota + 1
	SB
	BB
)

func (p Position) String() string {
	switch p {
	case Dealer:
		return "Dealer"
	case SB:
		return "SB"
	case BB:
		return "BB"
	default:
		return "Normal"
	}
}

func (p Position) Next() Position {
	switch p {
	case SB:
		return BB
	case BB:
		return SB
	default:
		return SB
	}
}
