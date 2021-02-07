package main

import (
	"fmt"
	"strings"

	"github.com/muesli/termenv"
)

func (m model) View() string {
	if m.coins == nil {
		return "Loading..."
	}
	if m.width < 72 {
		return "Terminal to narrow. Pleas resize..."
	}
	output := termenv.String("    Coin                      1H        24H         7D            Price\n").Underline().String()
	for index := m.cursor; index < m.cursor+m.height; index++ {
		if len(m.coins) > index {
			coin := m.coins[index]
			output += buildLine(&coin, index, m.fiatIndex)
			output += fmt.Sprintf("    %-67s\n", coin.Name)
		}
	}
	output += fmt.Sprintf("%72s\n", fmt.Sprintf("Updated: %s\n", m.coins[0].LastUpdated.Local()))
	output += fmt.Sprintf("%-72s", termenv.String(fmt.Sprintf("\n[q: Exit]  [▲/▼: Scroll]  [+/-: Height] [←/→: USD/EUR]")).
		Foreground(term.Color("#9e9e9e")).String())
	return output
}

func buildLine(coin *Coin, index int, fiatIndex int) string {
	return fmt.Sprintf("%s%-23s%-10s%-10s%-10s%10.3f %s\n",
		termenv.String(fmt.Sprintf("%2d. ", index+1)).Foreground(term.Color("#ffffff")),
		termenv.String(fmt.Sprintf("%-20s", strings.ToUpper(coin.Symbol))).Foreground(term.Color("#ffffff")).Bold(),
		getColorOfPercentChange(coin.PriceChangePercentage1HInCurrency),
		getColorOfPercentChange(coin.PriceChangePercentage24HInCurrency),
		getColorOfPercentChange(coin.PriceChangePercentage7DInCurrency),
		coin.CurrentPrice,
		strings.ToUpper(fiatCurrencies[fiatIndex]))
}

func getColorOfPercentChange(change float64) string {
	if change > 0 {
		return termenv.String(fmt.Sprintf("▲%6.2f%%   ", change)).
			Foreground(term.Color("#0ff00")).Bold().String()
	}
	return termenv.String(fmt.Sprintf("▼%s%%   ", strings.TrimPrefix(fmt.Sprintf("%6.2f", change), "-"))).
		Foreground(term.Color("#ff0000")).Bold().String()
}
