package data

import (
	"encoding/json"
	"fmt"
	"io"
)

type StockMetadata struct {
	Name	string `json:"name"`
	Symbol  string `json:"symbol"`
}

type Stock struct {
	Name	string `json:"name"`
	Symbol  string `json:"symbol"`
	Value 	float64 `json:"value"`
}

func (s *StockMetadata) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(s)
}

type Stocks []*Stock

func (s *Stocks) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(s)
}

func GetStocks() Stocks {
	return stockList
}

func AddStock(s *Stock) {
	stockList = append(stockList, s)
}

func UpdateStockValue(symbol string, newValue float64) error {
	index, err := findStockBySymbol(symbol)
	if err != nil {
		return err
	}

	stockList[index].Value = newValue
	return nil
}

var ErrStockNotFound = fmt.Errorf("stock not found")

func findStockBySymbol(symbol string) (int, error) {
	for i, s := range stockList {
		if s.Symbol == symbol {
			return i, nil
		}
	}
	return -1, ErrStockNotFound
}

var stockList = Stocks{}
