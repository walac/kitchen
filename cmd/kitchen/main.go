package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/walac/kitchen/courier"
	"github.com/walac/kitchen/order"
)

func main() {
	ordersPerSecond := flag.Int("r", 2, "Rate of orders in order per sec")
	flag.Parse()
	db := flag.Arg(0)

	if db == "" {
		fmt.Fprintf(os.Stderr, "The path to the order database file should be passed")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	source, err := NewSupplier(ctx, *ordersPerSecond, db)
	if err != nil {
		fmt.Println(err)
		return
	}

	sm := order.NewShelfManager()
	mediator := courier.NewLocalMediator(sm)
	courierDone := courier.New(ctx, mediator)
	kitchenDone := NewKitchen(ctx, mediator, source)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig

	fmt.Fprintln(os.Stderr, "Finish...")
	cancel()

	<-kitchenDone
	<-courierDone
	mediator.Close()
}
