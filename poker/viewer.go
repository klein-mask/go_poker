package poker

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"go_poker/card"
	"strings"
)

type Viewer struct {
	Context         *Poker
	App             *tview.Application
	turnText        *tview.TextView
	potText         *tview.TextView
	flopText        *tview.TextView
	flopCardTable   *tview.Table
	infoText        *tview.TextView
	playerCardTable *tview.Table
	enemyCardTable  *tview.Table
	playerMoneyText *tview.TextView
	playerBetText   *tview.TextView
	enemyMoneyText  *tview.TextView
	enemyBetText    *tview.TextView
	playerActions   *tview.List
}

func (v *Viewer) DrawInit() error {
	if len(v.Context.Players) < 2 {
		return errors.New("プレイヤー数が不足しています。")
	}

	cp := v.Context.getCurrentPlayer()
	// ターン経過用テキスト
	v.turnText = tview.NewTextView().
		SetText(cp.Name).
		SetTextColor(tcell.ColorWhite).
		SetChangedFunc(func() {
			v.App.Draw()
		}).SetTextAlign(tview.AlignCenter)
	v.turnText.SetTitle("Turn").SetBorder(true).SetTitleColor(tcell.ColorNavy)

	// ポット確認用テキスト
	v.potText = tview.NewTextView().
		SetText(v.Context.GetPotString()).
		SetTextColor(tcell.ColorWhite).
		SetChangedFunc(func() {
			v.App.Draw()
		}).SetTextAlign(tview.AlignCenter)
	v.potText.SetTitle("Pot").SetBorder(true).SetTitleColor(tcell.ColorYellow)

	// フロップカード
	v.flopCardTable = v.createCardTable(v.Context.GetFlopStrings())

	// Information用テキスト
	v.infoText = tview.NewTextView().
		SetText(fmt.Sprintf("%sのターンです。アクションを選択してください。", cp.Name)).
		SetTextColor(tcell.ColorOrange).
		SetChangedFunc(func() {
			v.App.Draw()
		}).SetTextAlign(tview.AlignCenter)
	v.infoText.SetTitle("Infomation").SetTitleColor(tcell.ColorRed).SetBorder(true)

	// プレイヤー
	player := v.Context.Players[0]
	playerTableText := tview.NewTextView().SetText("Player's Table").SetTextAlign(tview.AlignCenter)
	v.playerCardTable = v.createCardTable(player.GetHandStrings())

	v.playerMoneyText = tview.NewTextView().
		SetText(player.GetMoneyString()).
		SetTextColor(tcell.ColorWhite).
		SetChangedFunc(func() {
			v.App.Draw()
		})
	v.playerBetText = tview.NewTextView().
		SetText(player.GetBetString()).
		SetTextColor(tcell.ColorWhite).
		SetChangedFunc(func() {
			v.App.Draw()
		})

	v.playerActions = tview.NewList().
		AddItem(Fold.String(), "You hand fold", '1', func() {
			v.infoText.SetText("")
			v.Context.Action(Action{
				Type: Fold,
			})
			v.Context.NextTurn()
			v.DrawByCurrentData()
		}).
		AddItem(Call.String(), "You hand call", '2', func() {
			v.infoText.SetText("")

			err := v.Context.Action(Action{
				Type: Call,
			})
			if err != nil {
				v.infoText.Write([]byte(err.Error()))
			} else {
				v.infoText.Write([]byte(fmt.Sprintf("%sがCallを選択しました。\n", v.Context.getCurrentPlayer().Name)))
				v.Context.NextPlayer()
				v.DrawByCurrentData()
			}
		}).
		AddItem(Raise.String(), "You hand raise", '3', func() {
			v.infoText.SetText("")

			err := v.Context.Action(Action{
				Type: Raise,
				Bet: v.Context.TurnBet * 2,
			})
			if err != nil {
				v.infoText.Write([]byte(err.Error()))
			} else {
				v.infoText.Write([]byte(fmt.Sprintf("%sがRaiseを選択しました。\n", v.Context.getCurrentPlayer().Name)))
				v.Context.NextPlayer()
				v.DrawByCurrentData()
			}
		}).
		AddItem(Check.String(), "You hand check", '4', func() {
			v.infoText.SetText("")

			err := v.Context.Action(Action{
				Type: Check,
			})
			if err != nil {
				v.infoText.Write([]byte(err.Error()))
			} else {
				v.infoText.Write([]byte(fmt.Sprintf("%sがCheckを選択しました。\n", v.Context.getCurrentPlayer().Name)))
				v.Context.NextPlayer()
				v.DrawByCurrentData()
			}
		}).
		AddItem("Quit", "Press to exit this game", 'q', func() {
			v.App.Stop()
		})

	// エネミー
	enemy := v.Context.Players[1]
	enemyTableText := tview.NewTextView().SetText("Enemy's Table").SetTextAlign(tview.AlignCenter)

	v.enemyCardTable = v.createCardTable([]string{"？", "？"})

	v.enemyMoneyText = tview.NewTextView().
		SetText(enemy.GetMoneyString()).
		SetTextColor(tcell.ColorWhite).
		SetChangedFunc(func() {
			v.App.Draw()
		})
	v.enemyBetText = tview.NewTextView().
		SetText(enemy.GetBetString()).
		SetTextColor(tcell.ColorWhite).
		SetChangedFunc(func() {
			v.App.Draw()
		})

	playerMainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	enemyMainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	playerMainFlex.SetBorder(true)
	enemyMainFlex.SetBorder(true)

	rootFlex := tview.NewFlex().
		AddItem(playerMainFlex.
			AddItem(playerTableText, 5, 1, false).
			AddItem(tview.NewTextView().SetText("Hand").SetTextColor(tcell.ColorGreen), 2, 1, false).
			AddItem(v.playerCardTable, 5, 1, false).
			AddItem(tview.NewTextView().SetText("Total Money").SetTextColor(tcell.ColorYellow), 2, 1, false).
			AddItem(v.playerMoneyText, 3, 1, false).
			AddItem(tview.NewTextView().SetText("Bet").SetTextColor(tcell.ColorYellow), 2, 1, false).
			AddItem(v.playerBetText, 3, 1, false).
			AddItem(tview.NewTextView().SetText("Action").SetTextColor(tcell.ColorRed), 2, 1, false).
			AddItem(v.playerActions, 20, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(v.turnText, 3, 1, false).
			AddItem(v.potText, 3, 1, false).
			AddItem(tview.NewTextView().SetText("Flop").SetTextColor(tcell.ColorGreen).SetTextAlign(tview.AlignCenter), 10, 1, false).
			AddItem(v.flopCardTable, 5, 1, false).
			AddItem(v.infoText, 0, 1, false), 0, 2, false).
		AddItem(enemyMainFlex.
			AddItem(enemyTableText, 5, 1, false).
			AddItem(tview.NewTextView().SetText("Hand").SetTextColor(tcell.ColorGreen), 2, 1, false).
			AddItem(v.enemyCardTable, 5, 1, false).
			AddItem(tview.NewTextView().SetText("Total Money").SetTextColor(tcell.ColorYellow), 2, 1, false).
			AddItem(v.enemyMoneyText, 3, 1, false).
			AddItem(tview.NewTextView().SetText("Bet").SetTextColor(tcell.ColorYellow), 2, 1, false).
			AddItem(v.enemyBetText, 3, 1, false).
			AddItem(tview.NewBox(), 0, 1, false), 0, 1, false)
	v.DrawByCurrentData()
	if err := v.App.SetRoot(rootFlex, true).Run(); err != nil {
		return err
	}
	return nil
}

func (v *Viewer) DrawByCurrentData() {
	v.potText.SetText(v.Context.GetPotString())
	v.playerMoneyText.SetText(v.Context.Players[0].GetMoneyString())
	v.playerBetText.SetText(v.Context.Players[0].GetBetString())
	v.enemyMoneyText.SetText(v.Context.Players[1].GetMoneyString())
	v.enemyBetText.SetText(v.Context.Players[1].GetBetString())
	for i, flopStr := range v.Context.GetFlopStrings() {
		s := strings.Split(flopStr, " ")[0]
		v.flopCardTable.
			SetCell(0, i, tview.NewTableCell(flopStr).
				SetTextColor(v.getCardTableCellColor(s)).
				SetAlign(tview.AlignCenter))
	}
}

func (v *Viewer) createCardTable(cardStrings []string) *tview.Table {
	cardTable := tview.NewTable().SetBorders(true)
	for i, cardStr := range cardStrings {
		s := strings.Split(cardStr, " ")[0]
		cardTable.
			SetCell(0, i, tview.NewTableCell(cardStr).
				SetTextColor(v.getCardTableCellColor(s)).
				SetAlign(tview.AlignCenter))
	}
	return cardTable
}

func (v *Viewer) WriteInfoText(text string) {
	v.infoText.Write([]byte(text + "\n"))
}

func (v *Viewer) OpenEnemyCards(cardStrings []string) {
	for i, cardStr := range cardStrings {
		s := strings.Split(cardStr, " ")[0]
		v.enemyCardTable.
			SetCell(0, i, tview.NewTableCell(cardStr).
				SetTextColor(v.getCardTableCellColor(s)).
				SetAlign(tview.AlignCenter))
	}
}

func (v *Viewer) getCardTableCellColor(s string) tcell.Color {
	switch s {
	case card.Spade.String():
		return tcell.ColorWhite
	case card.Heart.String():
		return tcell.ColorRed
	case card.Diamond.String():
		return tcell.ColorYellow
	case card.Club.String():
		return tcell.ColorGreen
	default:
		return tcell.ColorWhite
	}
}
