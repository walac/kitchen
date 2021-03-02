package order

// Shelf implements the container/heap.Interface interface
// to allow an list of Order objects to be ordered by Deadline
type Shelf struct {
	queue         []*Order
	decayModifier int
}

// NewShelf creates a new shelf instance
func NewShelf(decayModifier int) *Shelf {
	return &Shelf{
		queue:         []*Order{},
		decayModifier: decayModifier,
	}
}

// Len return the length of the shelf
func (shelf *Shelf) Len() int {
	return len(shelf.queue)
}

// Less implements the < comparator based on the order deadline.
//
// It is implemented in such a way that the earliest deadline in on the
// top of the heap.
func (shelf *Shelf) Less(i, j int) bool {
	return shelf.queue[i].Deadline(shelf.decayModifier).Before(shelf.queue[j].Deadline(shelf.decayModifier))
}

// Swap swaps two orders in the shelf
func (shelf *Shelf) Swap(i, j int) {
	shelf.queue[i], shelf.queue[j] = shelf.queue[j], shelf.queue[i]
	shelf.queue[i].index = i
	shelf.queue[j].index = j
}

// Push pushes a new order to the end of the shelf
func (shelf *Shelf) Push(x interface{}) {
	o := x.(*Order)
	o.index = shelf.Len()
	shelf.queue = append(shelf.queue, o)
}

// Pop removes an order from the end of the shelf
func (shelf *Shelf) Pop() interface{} {
	x := shelf.queue[shelf.Len()-1]
	shelf.queue = shelf.queue[:shelf.Len()-1]
	return x
}

// IsTopWasted checks if the top item is wasted
func (shelf *Shelf) IsTopWasted() bool {
	if shelf.Len() == 0 {
		return false
	}

	return shelf.queue[0].IsWasted(shelf.decayModifier)
}
