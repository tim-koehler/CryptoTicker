package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		updateModelData(m)
		return m, tick()

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.coins)-m.height {
				m.cursor++
			}
		case "left":
			if m.fiatIndex > 0 {
				m.fiatIndex--
				updateModelData(m)
			}
		case "right":
			if m.fiatIndex < len(fiatCurrencies)-1 {
				m.fiatIndex++
				updateModelData(m)
			}
		case "+":
			if m.height < len(m.coins)-1 {
				m.height++
			}
		case "-":
			if m.height > minHeight {
				m.height--
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func updateModelData(m *model) {
	coins := Coins{}
	rawData, err := callAPI(fiatCurrencies[m.fiatIndex])
	if err != nil {
		fmt.Println(err.Error())
	}
	if err := json.Unmarshal(rawData, &coins); err != nil {
		fmt.Println(err.Error())
	}
	m.coins = coins
}

func callAPI(fiatCurrency string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return []byte{}, err
	}

	q := url.Values{}
	q.Add("vs_currency", fiatCurrency)
	q.Add("order", "market_cap_desc")
	q.Add("per_page", "50")
	q.Add("page", "1")
	q.Add("sparkline", "false")
	q.Add("price_change_percentage", "1h,24h,7d,30d")

	req.Header.Set("Accepts", "application/json")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody, nil
}
