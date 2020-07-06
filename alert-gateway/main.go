package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/models"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/routers"
)

func main() {
	defer models.Ormer.Close()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routers.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
