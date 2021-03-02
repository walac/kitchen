package courier

import "github.com/walac/kitchen/order"

// CourierMediator defines the interface between the kitchen and the courier
type CourierMediator interface {
	// NotifyChanR returns a read channel that receives the ID of
	// orders that are ready to pickup
	NotifyChan() <-chan string

	// Notify notifies the courier about a new order
	Notify(id string)

	// GetOrder returns the order correspoding to the
	// given id. It returns nil if the order is wasted or if it not found.
	RequestOrder(id string) *order.Order

	// AddOrder add a new order to the shelf
	AddOrder(o *order.Order)

	// Close closes the communication channels
	Close()
}

// LocalMediator implements the CourierMediator interface for
// local communication in the same process
type LocalMediator struct {
	notifyChan   chan string
	shelfManager *order.ShelfManager
}

// NewLocalMediator instantiates a new mediator object
func NewLocalMediator(shelfManager *order.ShelfManager) CourierMediator {
	return &LocalMediator{
		notifyChan:   make(chan string),
		shelfManager: shelfManager,
	}
}

func (m *LocalMediator) NotifyChan() <-chan string {
	return m.notifyChan
}

func (m *LocalMediator) Notify(id string) {
	m.notifyChan <- id
}

func (m *LocalMediator) RequestOrder(id string) *order.Order {
	return m.shelfManager.Remove(id)
}

func (m *LocalMediator) AddOrder(o *order.Order) {
	m.shelfManager.Add(o)
}

func (m *LocalMediator) Close() {
	close(m.notifyChan)
}
