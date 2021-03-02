package order

import (
	"fmt"
	"time"
)

type Order struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Temperature  string  `json:"temp"`
	ShelfLife    int64   `json:"shelfLife"`
	DecayRate    float64 `json:"decayRate"`
	CreationTime time.Time
	index        int
}

// Deadline calculates the point in the future where the order
// becomes wasted.
func (o *Order) Deadline(decayModifier int) time.Time {
	totalDecayRate := o.DecayRate * float64(decayModifier)
	return o.CreationTime.Add(time.Duration(float64(o.ShelfLife)/totalDecayRate) * time.Second)
}

// Age returns the order age
func (o *Order) Age() time.Duration {
	return time.Now().Sub(o.CreationTime)
}

// Value calculates the order value, which is given by the formula:
//
// value = (shelfLife - decayRate * orderAge * shelfDecayModifier)
//         -------------------------------------------------------
//                          shelfLife
func (o *Order) Value(decayModifier int) float64 {
	totalDecayRate := o.DecayRate * float64(decayModifier)
	return 1.0 - totalDecayRate*o.Age().Seconds()/float64(o.ShelfLife)
}

// Index returns the heap key index of the order
func (o *Order) Index() int {
	return o.index
}

// IsWasted returns true the order is wasted
func (o *Order) IsWasted(decayModifier int) bool {
	return o.Value(decayModifier) <= 0.0
}

// String converts an order to a printable string
func (o *Order) String(decayModifier int) string {
	return fmt.Sprintf("Order{id=%s, value=%.3f}", o.Id, o.Value(decayModifier))
}
