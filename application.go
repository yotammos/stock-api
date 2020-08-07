package main

import (
	"context"
	"github.com/gorilla/mux"
	"gopkg.in/robfig/cron.v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stock-api/clients"
	"stock-api/data"
	"stock-api/handlers"
	"time"
)

func main() {
	l := log.New(os.Stdout, "stock-api", log.LstdFlags)
	sh := handlers.NewStocks(l)
	sm := mux.NewRouter()

	setupStockUpdate()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", sh.GetStocks)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{symbol:[A-Z]+}", sh.UpdateStockValue)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", sh.AddStocks)
	postRouter.Use(sh.MiddlewareStockValidation)

	s := http.Server{
		Addr: ":5000",
		Handler: sm,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <- c
	log.Println("Got signal: ", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	s.Shutdown(ctx)
}

func setupStockUpdate() {
	l := log.New(os.Stdout, "update-stocks-job", log.LstdFlags)
	c := cron.New()
	_, err := c.AddFunc("* * * * *", func() {
		stocks := data.GetStocks()
		for _, stock := range stocks {
			newValue, err := clients.FetchStockPrice(stock.Symbol, l)
			if err != nil {
				l.Println("failed fetching stock price")
			} else {
				err = data.UpdateStockValue(stock.Symbol, newValue)
				if err != nil {
					l.Println("failed updating stock price")
				} else {
					l.Printf("updated %v price to %v", stock.Symbol, newValue)
				}
			}
		}
	})

	if err != nil {
		l.Println("failed creating cron job")
		return
	}

	l.Println("starting cron job")
	c.Start()
}
