package main

import (
	"context"
	"fmt"
	"os"

	"github.com/walac/kitchen/courier"
	"github.com/walac/kitchen/order"
)

// NewKitchen starts a new kitchen service
//
// New orders are received through the recv channel
func NewKitchen(ctx context.Context, mediator courier.CourierMediator, recv <-chan *order.Order) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		fmt.Fprintln(os.Stderr, "Starting kitchen service")

		for {
			select {
			case <-ctx.Done():
				fmt.Fprintln(os.Stderr, "Exiting kitchen service")
				done <- struct{}{}
				return

			case o := <-recv:
				mediator.AddOrder(o)
				mediator.Notify(o.Id)
			}
		}
	}()

	return done
}
