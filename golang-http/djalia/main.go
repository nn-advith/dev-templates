package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nbeans/djalia/server"
)

func main() {
	fmt.Println("Server")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "You have hit root")
	})

	cfg := server.Config{
		Address:      ":4000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	app, err := server.NewServer(cfg, mux)
	if err != nil {
		fmt.Println("err :", err)
	}

	go func() {
		if err := app.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Println("start error: ", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		fmt.Println("Shutdown error:", err)
	}

	fmt.Println("Server gracefully stopped")

}
