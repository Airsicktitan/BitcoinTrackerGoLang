package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type BitCoin struct {
	Time       Time   `json:"time"`
	Disclaimer string `json:"disclaimer"`
	ChartName  string `json:"chartName"`
	BPI        BPI    `json:"bpi"`
}

type Time struct {
	Updated    string `json:"updated"`
	UpdatedISO string `json:"updatedISO"`
	Updateduk  string `json:"updateduk"`
}

type BPI struct {
	USD  USD `json:"USD"`
	GBP  GBP `json:"GBP"`
	EURO EUR `json:"EUR"`
}

type USD struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	Rate_float  float64 `json:"rate_float"`
}

type GBP struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	Rate_float  float64 `json:"rate_float"`
}

type EUR struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	Rate_float  float64 `json:"rate_float"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Missing command-line argument")
		return
	}
	
	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Println("Command-line argument not a number")
		return
	}

	title := "Bitcoin Price Tracker"
	fmt.Println("\n", title)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			fmt.Println("\nProgram terminated.")
			fmt.Print("\n")
		}
		os.Exit(0)
	}()

	for i := 0; i < len(title); i++ {
		print("-")
	}
	fmt.Print("\n\n")

	for {
		resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("API unavailable")
			return
		}

		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			errMsg := errors.New("something went wrong")
			fmt.Println(errMsg)
			return
		}

		var bitcoin BitCoin
		err = json.Unmarshal(respData, &bitcoin)
		if err != nil {
			fmt.Println(err)
			return
		}

		price := bitcoin.BPI.USD.Rate_float * amount
		fmt.Printf("$%.02f\n", price)

		time.Sleep(18 * time.Second)
	}
}
