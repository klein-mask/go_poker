package card

type CardNumber int

const (
	Ace CardNumber = iota + 1
	Two
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
)

func (cn CardNumber) ToString() string {
	switch cn {
	case Ace:
		return "Ace"
	case Two:
		return "Two"
	case Three:
		return "Three"
	case Four:
		return "Four"
	case Five:
		return "Five"
	case Six:
		return "Six"
	case Seven:
		return "Seven"
	case Eight:
		return "Eight"
	case Nine:
		return "Nine"
	case Ten:
		return "Ten"
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	default:
		return "NoNumber"
	}
}

func (cn CardNumber) IsLarge(compareCardNumer CardNumber) bool {
	if cn == Ace {
		return true
	} else {
		return (cn > compareCardNumer)
	}
}