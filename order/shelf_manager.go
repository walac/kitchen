package order

import (
	"container/heap"
	"container/list"
	"fmt"
	"os"
	"sync"
)

const (
	shelfCapacity         = 10
	overflowShelfCapacity = 15
)

// ShelfManager manages the insertion and deletion of orders from the shelfs
type ShelfManager struct {
	temperatureShelfs map[string]*Shelf
	overflowShelf     *Shelf
	// map order Id to order object
	ordersMap map[string]*Order
	// map order temperature to the list of orders of that
	// temperature in the overflow shelf
	overflowOrderMap map[string]*list.List

	mutex sync.Mutex
}

// NewShelfManager returns a initialized ShelfManager object
func NewShelfManager() *ShelfManager {
	temps := []string{"hot", "cold", "frozen"}
	temperatureShelfs := make(map[string]*Shelf)
	overflowOrderMap := make(map[string]*list.List)
	for _, temp := range temps {
		temperatureShelfs[temp] = NewShelf(1)
		overflowOrderMap[temp] = list.New()
	}

	return &ShelfManager{
		temperatureShelfs: temperatureShelfs,
		overflowShelf:     NewShelf(2),
		ordersMap:         make(map[string]*Order),
		overflowOrderMap:  overflowOrderMap,
		mutex:             sync.Mutex{},
	}
}

// Add inserts the order to the corresponding shelf
func (m *ShelfManager) Add(o *Order) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	fmt.Fprintf(os.Stderr, "Adding order %s\n", o.Id)

	entry, ok := m.temperatureShelfs[o.Temperature]
	if !ok {
		fmt.Fprintf(os.Stderr, "Temperature %s not valid\n", o.Temperature)
		return
	}

	// remove all wasted orders from the temperature shelf
	for entry.IsTopWasted() {
		removed := heap.Pop(entry).(*Order)
		fmt.Fprintf(os.Stderr, "Removing %s from the %s shelf because it is wasted\n",
			removed.String(entry.decayModifier), removed.Temperature)
	}

	if entry.Len() == shelfCapacity {
		// remove all wasted orders from the overflow shelf
		for m.overflowShelf.IsTopWasted() {
			removed := m.popOverflow()
			fmt.Fprintf(os.Stderr, "Removing %s from the overflow shelf because it is wasted\n",
				removed.String(m.overflowShelf.decayModifier))
		}

		if m.overflowShelf.Len() == overflowShelfCapacity {
			if !m.flushOverflow() {
				// if the flush operation couldn't make room
				// in the overflow shelf, remove one order
				removed := m.popOverflow()
				fmt.Fprintf(os.Stderr, "Removing %s from the overflow to make room in the shelf\n",
					removed.String(m.overflowShelf.decayModifier))
			}
		}

		m.insertOverflow(o)
		fmt.Fprintf(os.Stderr, "Inserted %s in the overflow shelf\n",
			o.String(m.overflowShelf.decayModifier))
	} else {
		heap.Push(entry, o)
		fmt.Fprintf(os.Stderr, "Inserted %s in the %s shelf\n",
			o.String(entry.decayModifier), o.Temperature)
	}

	m.ordersMap[o.Id] = o
}

// Remove removes from the corresponding shelf and returns the order
// referenced by the given id. It returns nil if the order is wasted
// or not found.
func (m *ShelfManager) Remove(id string) *Order {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	o, ok := m.ordersMap[id]
	if !ok {
		return nil
	}

	entry := m.temperatureShelfs[o.Temperature]
	var temp string

	// if the order isn't in the temperature shelf it must be in the overflow shelf
	if o.index >= entry.Len() || o.Id != entry.queue[o.index].Id {
		entry = m.overflowShelf
		orders := m.overflowOrderMap[o.Temperature]
		element := orders.Front()
		for element != nil {
			if element.Value.(*Order).Id == o.Id {
				orders.Remove(element)
				break
			}

			element = element.Next()
		}

		temp = "overflow"
	} else {
		temp = o.Temperature
	}

	fmt.Fprintf(os.Stderr, "Removing %s from the %s shelf\n",
		o.String(entry.decayModifier), temp)
	heap.Remove(entry, o.index)
	delete(m.ordersMap, id)

	if o.IsWasted(entry.decayModifier) {
		fmt.Fprintf(os.Stderr, "Order %s removed because it is wasted\n", o.String(entry.decayModifier))
		o = nil
	}
	return o
}

// popOverflow removes the element from the top of the overflow shelf
func (m *ShelfManager) popOverflow() *Order {
	o := heap.Pop(m.overflowShelf).(*Order)
	orderList := m.overflowOrderMap[o.Temperature]
	element := orderList.Front()

	// Since the bucket size is bouned, this is effectively an O(1)
	// operation. In practice, dependending on the shelfCapacity
	// value, we may need to to switch orderList to a more search
	// efficient data structure
	for element != nil {
		if element.Value.(*Order).Id == o.Id {
			orderList.Remove(element)
			break
		}

		element = element.Next()
	}

	delete(m.ordersMap, o.Id)
	return o
}

// insertOverflow adds an order to the overflow shelf
func (m *ShelfManager) insertOverflow(o *Order) {
	heap.Push(m.overflowShelf, o)
	m.overflowOrderMap[o.Temperature].PushBack(o)
}

// flushOverflow flushes all orders in the overflow shelf to the non-full
// corresponding temperature shelf
func (m *ShelfManager) flushOverflow() bool {
	found := false
	for k, v := range m.temperatureShelfs {
		orders := m.overflowOrderMap[k]
		element := orders.Front()
		for v.Len() < shelfCapacity && element != nil {
			next := element.Next()
			o := orders.Remove(element).(*Order)
			fmt.Fprintf(os.Stderr, "Moving %s from the overflow to the %s shelf\n",
				o.String(v.decayModifier), k)
			element = next
			heap.Remove(m.overflowShelf, o.index)
			heap.Push(v, o)
			found = true
		}
	}

	return found
}
