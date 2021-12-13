package poker

import (
	"errors"
	"fmt"
	"github.com/rivo/tview"
	"go_poker/card"
	"go_poker/deck"
	"go_poker/hand"
	"math/rand"
	"strconv"
	"time"
)

type Poker struct {
	Players         []*Player
	Deck            *deck.Deck
	BigBlind        int
	SmollBlind      int
	Pot             int
	Flop            []card.Card
	TurnIndex       int
	TurnBet         int
	IsHandFinished  bool
	InfomationTexts []string
	Viewer          Viewer
}

func NewPoker(bb, sb, playerInitMoney int) *Poker {
	if bb < sb {
		fmt.Println("BBはSBより大きい値を指定してください")
		return nil
	} else if playerInitMoney < bb {
		fmt.Println("プレイヤーの所持金はBBより大きい値を指定してください")
		return nil
	}
	d := deck.NewDeck().Shuffle()

	p := &Poker{
		Players: []*Player{
			NewPlayer("Player", playerInitMoney, SB),
			NewPlayer("Enemy", playerInitMoney, BB),
		},
		Deck:       d,
		BigBlind:   bb,
		SmollBlind: sb,
	}
	p.Viewer = Viewer{
		Context: p,
		App:     tview.NewApplication(),
	}

	return p
}

func (p *Poker) InitSetUp() {
	// ブラインドベット
	p.BlindBet()
	// プリフロップ
	p.PreFlop()

	// TUIに描画
	err := p.Viewer.DrawInit()
	if err != nil {
		panic(err)
	}
}

func (p *Poker) BlindBet() *Poker {
	for _, player := range p.Players {
		if player.Position == SB {
			player.Bet(p.SmollBlind)
		} else if player.Position == BB {
			player.Bet(p.BigBlind)
		}
	}
	p.TurnBet = p.BigBlind
	return p
}

func (p *Poker) PreFlop() *Poker {
	for _, player := range p.Players {
		player.Hand.Add(p.Deck.Deal(2))
	}
	return p
}

func (p *Poker) isNextTurn() bool {
	for i := 0; i < len(p.Players); i++ {
		if i >= len(p.Players)-1 {
			return true
		}
		if p.Players[i].CurrentBet != p.Players[i+1].CurrentBet {
			return false
		}
	}
	return true
}

func (p *Poker) NextTurn() error {
	notFoldPlayers := p.getNotFoldPlayers()

	if !p.isNextTurn() && len(notFoldPlayers) >= 2 {
		return errors.New("プレイヤーのベット額が一致していません")
	}

	// どちらかがフォールドした場合
	if len(notFoldPlayers) == 1 {
		notFoldPlayers[0].IsHandWin = true
		p.Finish()
		return nil
	}

	if len(p.Flop) >= 5 {
		p.ShowDown()
	} else {
		p.OpenFlop()
	}
	p.TurnIndex = 0
	p.Viewer.DrawByCurrentData()

	return nil
}

func (p *Poker) Finish() {
	losePlayerIndex := 0
	for i, player := range p.Players {
		if player.IsHandWin {
			getMoney := p.CulcPot() - player.CurrentBet
			player.Win(p.CulcPot())
			p.Viewer.WriteInfoText(fmt.Sprintf("「%s」の勝利です", player.Name))
			p.Viewer.WriteInfoText(fmt.Sprintf("獲得ドル: %s", strconv.Itoa(getMoney)))
		} else {
			losePlayerIndex = i
		}
	}
	p.Players[losePlayerIndex].Lose()
	p.Viewer.DrawByCurrentData()
}

func (p *Poker) CulcPot() (result int) {
	for _, player := range p.Players {
		result += player.CurrentBet
	}
	return result
}

func (p *Poker) OpenFlop() {
	openFlopCount := 1
	if len(p.Flop) == 0 {
		openFlopCount = 3
		p.Viewer.WriteInfoText("フロップを配布します。")
	} else {
		p.Viewer.WriteInfoText("ボードにカードを追加します。")
	}
	p.Flop = append(p.Flop, p.Deck.Deal(openFlopCount)...)
}

func (p *Poker) ShowDown() *Poker {
	winPlayer := p.Players[0]
	for _, player := range p.Players {
		player.Hand.Culc(p.Flop)
		p.Viewer.WriteInfoText(fmt.Sprintf("%s の手役は %s です", player.Name, player.Hand.Point))
		if winPlayer.Hand.Point < player.Hand.Point {
			winPlayer = player
		} else if winPlayer.Hand.Point == player.Hand.Point {
			winPlayer = p.JudgeWinPlayerFromDrawHand(winPlayer, player)
		}
	}
	winPlayer.IsHandWin = true
	p.Finish()
	if len(p.Players) >= 2 {
		p.Viewer.OpenEnemyCards(p.Players[1].GetHandStrings())
	}
	return p
}

func (p *Poker) Action(a Action) error {
	turnPlayer := p.getCurrentPlayer()
	turnPlayer.CurrentAction = a

	switch a.Type {
	case Call:
		diff := p.TurnBet - turnPlayer.CurrentBet
		if diff == 0 {
			return errors.New("ベット額が既に足りています。RaiseかCheckを選択してください。")
		} else {
			err := turnPlayer.Bet(diff)
			if err != nil {
				return err
			}
			//p.Pot += diff
		}
	case Check:
		diff := p.TurnBet - turnPlayer.CurrentBet
		if diff > 0 {
			return errors.New("ベット額が足りていません。CallかRaiseを選択してください。")
		}
	case Raise:
		diff := a.Bet - turnPlayer.CurrentBet
		err := turnPlayer.Bet(diff)
		if err != nil {
			return err
		}
		//p.Pot += diff
	case AllIn:
		turnPlayer.AllIn()
	}
	p.TurnBet = turnPlayer.CurrentBet
	return nil
}

func (p *Poker) NextPlayer() *Poker {
	if p.TurnIndex < len(p.Players)-1 {
		p.TurnIndex += 1
	} else {
		p.TurnIndex = 0
	}

	cp := p.getCurrentPlayer()
	p.Viewer.WriteInfoText(fmt.Sprintf("次は、%sのアクションです。", cp.Name))

	if cp.Name == "Enemy" {
		a := p.RandomAction()
		cp.CurrentAction = a
		p.Viewer.WriteInfoText(fmt.Sprintf("%sは%sを選択しました。", cp.Name, a.Type))
		if a.Type == Fold {
			p.NextTurn()
			return p
		}
		if p.isNextTurn() {
			p.Viewer.WriteInfoText(fmt.Sprintf("次のターンへ進みます。"))
			p.NextTurn()
			return p
		}

		switch a.Type {
		case Raise:
			diff := a.Bet - cp.CurrentBet
			cp.Bet(diff)
			p.TurnBet = cp.CurrentBet
			p.NextPlayer()
		case AllIn:
			cp.AllIn()
			p.TurnBet = cp.CurrentBet
			p.NextPlayer()
		case Call:
			diff := p.TurnBet - cp.CurrentBet
			cp.Bet(diff)
			p.NextTurn()
		}
	}
	return p
}

func (p *Poker) JudgeWinPlayerFromDrawHand(playerA, playerB *Player) *Player {
	switch playerA.Hand.Point {
	case hand.AFullHouse:
		playerAThreeCardNumber := playerA.Hand.GetHandInnerThreeCardGroup()[0].Number
		playerBThreeCardNumber := playerB.Hand.GetHandInnerThreeCardGroup()[0].Number
		if playerAThreeCardNumber.IsLarge(playerBThreeCardNumber) {
			return playerA
		} else {
			return playerB
		}
	// 役とマークの判定を追加する
	default:
		return playerA
	}
}

func (p *Poker) GetPotString() string {
	return "＄" + strconv.Itoa(p.CulcPot())
}

func (p *Poker) getCurrentPlayer() *Player {
	return p.Players[p.TurnIndex]
}

func (p *Poker) getNextPlayer() *Player {
	idx := p.TurnIndex
	if p.TurnIndex < len(p.Players)-1 {
		idx += 1
	} else {
		idx = 0
	}
	return p.Players[idx]
}

func (p *Poker) getNotFoldPlayers() []*Player {
	results := make([]*Player, 0, len(p.Players))
	for _, player := range p.Players {
		if player.CurrentAction.Type != Fold {
			results = append(results, player)
		}
	}
	return results
}

func (p *Poker) GetFlopStrings() (results []string) {
	for _, card := range p.Flop {
		results = append(results, fmt.Sprintf("%s : %s", card.Suit, strconv.Itoa(int(card.Number))))
	}
	return results
}

func (p *Poker) RandomAction() Action {
	actions := []Action{
		//Action{
		//	Type: Fold,
		//},
		//Action{
		//	Type: AllIn,
		//},
	}
	cp := p.getCurrentPlayer()
	if p.TurnBet-cp.CurrentBet > 0 && cp.Money >= p.TurnBet-cp.CurrentBet {
		actions = append(actions, Action{
			Type: Call,
		})
	}
	if p.TurnBet-cp.CurrentBet <= 0 {
		actions = append(actions, Action{
			Type: Check,
		})
	}
	if p.TurnBet-cp.CurrentBet > 0 && cp.Money >= (p.TurnBet*2-cp.CurrentBet) {
		actions = append(actions, Action{
			Type: Raise,
			Bet:  (p.TurnBet*2 - cp.CurrentBet),
		})
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(actions), func(i, j int) {
		actions[i], actions[j] = actions[j], actions[i]
	})

	return actions[0]
}
