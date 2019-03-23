package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// BornToDie blocks until being interrupted or cancelled
func BornToDie(ctx context.Context) {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-ctx.Done():
		log.Println("parent context done")
	case s := <-signals:
		log.Printf("signal received: %s\n", s)
	}
	signal.Stop(signals)
}
