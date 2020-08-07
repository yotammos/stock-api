package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"stock-api/clients"
	"stock-api/data"
)

type KeyStock struct {}

type Stocks struct {
	l *log.Logger
}

func NewStocks(l *log.Logger) *Stocks {
	return &Stocks{l}
}

func (s *Stocks) GetStocks(rw http.ResponseWriter, _ *http.Request) {
	enableCors(&rw)
	sl := data.GetStocks()
	err := sl.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (s *Stocks) AddStocks(rw http.ResponseWriter, r *http.Request) {
	stockMetadata := r.Context().Value(KeyStock{}).(data.StockMetadata)

	newStockValue, err := clients.FetchStockPrice(stockMetadata.Symbol, s.l)
	if err != nil {
		http.Error(rw, "Failed fetching stock value", http.StatusBadRequest)
	}
	var stock = data.Stock{
		Name:   stockMetadata.Name,
		Symbol: stockMetadata.Symbol,
		Value:  newStockValue,
	}

	data.AddStock(&stock)
}

func (s *Stocks) UpdateStockValue(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol, ok := vars["symbol"]
	if !ok {
		http.Error(rw, "Unable to find symbol in path", http.StatusBadRequest)
		return
	}

	stockValue, err := clients.FetchStockPrice(symbol, s.l)
	err = data.UpdateStockValue(symbol, stockValue)
	if err == data.ErrStockNotFound {
		http.Error(rw, "Stock not found", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "Failed updating stock value", http.StatusInternalServerError)
		return
	}
}

func (s Stocks) MiddlewareStockValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		stockMetadata := data.StockMetadata{}

		err := stockMetadata.FromJSON(r.Body)
		if err != nil {
			s.l.Println("[ERROR] failed to decode stock")
			http.Error(rw, "Error reading stock meta data", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyStock{}, stockMetadata)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}

func enableCors(rw *http.ResponseWriter) {
	(*rw).Header().Set("Access-Control-Allow-Origin", "*")
}
