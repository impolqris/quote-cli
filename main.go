package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	// errors.
	ErrGetQuote = iota + 1
	ErrBodyQuote
	ErrBodyJson

	// error strings.
	ErrCallApi   = "There was an error calling the api."
	ErrQuoteJson = "There was an error getting the JSON of the quote. Please try again later."

	// api url.
	ApiURL = "https://favqs.com/api/qotd"
)

// QuoteInfo is required since the favqs API returns a nested JSON object, so
// our struct must be nested in order to match the JSON object.
type QuoteInfo struct {
	QuoteID  int64  `json:"id"`
	QuoteURL string `json:"url"`
	Author   string `json:"author"`
	Body     string `json:"body"`
}

type FullQuote struct {
	QotdDate string    `json:"qotd_date"`
	Quote    QuoteInfo `json:"quote"`
}

// CallApi requires an API url with the proper credentials. It will then return
// the full body of the API call, only returning an integer when there is an error
// with either the http request or there is an error reading from the response.
func CallApi(url string) ([]byte, error) {
	resp, rerr := http.Get(url)

	if rerr != nil {
		return nil, rerr
	}

	body, berr := io.ReadAll(resp.Body)

	if berr != nil {
		return nil, berr
	}

	return body, nil
}

func main() {
	body, berr := CallApi(ApiURL)

	if berr != nil {
		fmt.Println(ErrCallApi)

		// My reason for not using panic() is because I'm not good with concurrency
		// yet, but this will be updated in the future.

		os.Exit(ErrGetQuote)
	}

	quote := FullQuote{}

	jerr := json.Unmarshal(body, &quote)

	if jerr != nil {
		fmt.Println(ErrQuoteJson)

		os.Exit(ErrBodyJson)
	}

	fmt.Println(quote.Quote.Body + "\n\nWritten By: " + quote.Quote.Author)
}
