package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		coins := Coins{}
		rawData, err := fetchFromAPI()
		if err != nil {
			return m, tick()
		}
		if err := json.Unmarshal(rawData, &coins); err != nil {
			fmt.Println(err.Error())
			return m, tea.Quit
		}

		m.coins = coins
		return m, tick()

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor == 0 {
				return m, tick()
			}
			m.cursor--
			return m, tick()
		case "down":
			if m.cursor >= len(m.coins)-m.height {
				return m, tick()
			}
			m.cursor++
			return m, tick()
		case "+":
			if m.height > len(m.coins)-1 {
				return m, tick()
			}
			m.height++
			return m, tick()
		case "-":
			if m.height == minHeight {
				return m, tick()
			}
			m.height--
			return m, tick()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func fetchFromAPI() ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return []byte{}, err
	}

	q := url.Values{}
	q.Add("vs_currency", "eur")
	q.Add("order", "market_cap_desc")
	q.Add("per_page", "15")
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
