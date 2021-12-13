package poker

type Action struct {
	Type ActionType
	Bet int
}

type ActionType int

const (
	Fold ActionType = iota + 1
	Call
	Check
	Raise
	AllIn
)

func (a ActionType) String() string {
	switch a {
	case Fold:
		return "Fold"
	case Call:
		return "Call"
	case Check:
		return "Check"
	case Raise:
		return "Raise"
	case AllIn:
		return "AllIn"
	default:
		return "Fold"
	}
}


