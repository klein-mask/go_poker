package poker

import (
	"errors"
	"fmt"
	"go_poker/hand"
	"strconv"
)

type Player struct {
	Name          string
	Hand          hand.Hand
	Position      Position
	Money         int
	CurrentBet    int
	CurrentAction Action
	IsHandWin     bool
}

func NewPlayer(name string, initMoney int, position Position) *Player {
	return &Player{
		Name:     name,
		Money:    initMoney,
		Position: position,
	}
}

func (p *Player) NextHand() *Player {
	p.Hand = hand.Hand{}
	p.Position = p.Position.Next()
	p.CurrentBet = 0
	p.CurrentAction = Action{Type: Check}
	p.IsHandWin = false
	return p
}

func (p *Player) Win(addMoney int) {
	p.Money += addMoney
	p.CurrentBet = 0
	p.IsHandWin = true
}

func (p *Player) Bet(betMoney int) error {
	if p.Money < betMoney {
		return errors.New("ベット額が足りません")
	}
	p.Money -= betMoney
	p.CurrentBet += betMoney
	return nil
}

func (p *Player) Lose() error {
	p.CurrentBet = 0
	return nil
}

func (p *Player) AllIn() *Player {
	p.Bet(p.Money)
	return p
}

func (p Player) GetMoneyString() string {
	return "＄" + strconv.Itoa(p.Money)
}

func (p Player) GetBetString() string {
	return "＄" + strconv.Itoa(p.CurrentBet)
}

func (p Player) GetHandStrings() (results []string) {
	for _, card := range p.Hand.Cards {
		results = append(results, fmt.Sprintf("%s : %s", card.Suit, strconv.Itoa(int(card.Number))))
	}
	return results
}
