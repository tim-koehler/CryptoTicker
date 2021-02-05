package main

import (
	"fmt"
	"strings"

	aurora "github.com/logrusorgru/aurora/v3"
)

func (m model) View() string {
	if m.coins == nil {
		return "Loading..."
	}
	if m.width < 75 {
		return "Terminal to narrow. Pleas resize..."
	}
	output := aurora.Underline("    Coin                         1H        24H         7D            Price\n").String()
	for index := m.cursor; index < m.cursor+m.height; index++ {
		if len(m.coins) > index {
			coin := m.coins[index]
			output += buildLine(&coin, index)
			output += fmt.Sprintf("    %-67s\n", coin.Name)
		}
	}
	output += fmt.Sprintf("%75s\n", fmt.Sprintf("Updated: %s\n", m.coins[0].LastUpdated))
	output += fmt.Sprintf("%-75s", aurora.Gray(15, fmt.Sprintf("\n[q: Exit]  [▲/▼: Scroll]  [+/-: Height]")).String())
	return output
}

func buildLine(coin *Coin, index int) string {
	return fmt.Sprintf("%s%-23s%-10s%-10s%-10s%10.3f%s\n",
		aurora.Bold(fmt.Sprintf("%2d. ", index+1)),
		aurora.Bold(fmt.Sprintf("%s", strings.ToUpper(coin.Symbol))).White(),
		getColorOfPercentChange(coin.PriceChangePercentage1HInCurrency),
		getColorOfPercentChange(coin.PriceChangePercentage24HInCurrency),
		getColorOfPercentChange(coin.PriceChangePercentage7DInCurrency),
		coin.CurrentPrice,
		" EUR")
}

func getColorOfPercentChange(change float64) aurora.Value {
	if change > 0 {
		return aurora.Green(aurora.Bold(fmt.Sprintf("▲%6.2f%%   ", change)))
	}
	return aurora.Red(aurora.Bold(fmt.Sprintf("▼%s%%   ", strings.TrimPrefix(fmt.Sprintf("%6.2f", change), "-"))))
}
