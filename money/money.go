package money

import (
	"errors"
)

type Money struct {
	Current int
}

func (m *Money) Bet(betMoney int) error {
	if m.Current < betMoney {
		return errors.New("ベット額が足りません")
	}
	m.Current -= betMoney
	return nil
}
