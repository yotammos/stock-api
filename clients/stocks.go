package clients

import (
	"github.com/piquette/finance-go/quote"
	"log"
)

func FetchStockPrice(symbol string, l *log.Logger) (float64, error) {
	q, err := quote.Get(symbol)
	if err != nil {
		l.Println("Failed fetching symbol")
		return 0, nil
	}

	return q.Ask, nil
}
