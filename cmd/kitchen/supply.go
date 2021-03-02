package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/walac/kitchen/order"
)

// NewSupplier starts the order supplier service
func NewSupplier(ctx context.Context, ordersPerSec int, orderDb string) (<-chan *order.Order, error) {
	contents, err := ioutil.ReadFile(orderDb)
	if err != nil {
		return nil, err
	}

	var orders []order.Order
	if err := json.Unmarshal(contents, &orders); err != nil {
		return nil, err
	}

	sender := make(chan *order.Order, ordersPerSec)

	go func() {
		fmt.Fprintln(os.Stderr, "Starting supply service")
		i := 0
		timer := time.Tick(time.Second)
		for {
			select {
			case <-ctx.Done():
				fmt.Fprintln(os.Stderr, "Exiting supply service")
				close(sender)
				return

			case <-timer:
				for j := 0; j < ordersPerSec && i < len(orders); j += 1 {
					o := &orders[i]
					o.CreationTime = time.Now()
					sender <- o
					i += 1
				}
			}
		}
	}()

	return sender, nil
}
