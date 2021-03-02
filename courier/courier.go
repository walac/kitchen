package courier

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// New instanciates a new courier and return the channel that receives
// a value when it is done
func New(ctx context.Context, mediator CourierMediator) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		rand.Seed(time.Now().Unix())
		fmt.Fprintln(os.Stderr, "Starting courier service")
		for {
			select {
			case <-ctx.Done():
				fmt.Fprintln(os.Stderr, "Exiting courier service")
				done <- struct{}{}
				return

			case id := <-mediator.NotifyChan():
				go func(id string, delay int) {
					select {
					case <-ctx.Done():
						return
					case <-time.After(time.Duration(delay) * time.Second):
						o := mediator.RequestOrder(id)
						if o != nil {
							fmt.Printf("Order %s delivered\n", id)
						}
					}
				}(id, rand.Intn(5)+2)
			}
		}
	}()

	return done
}
