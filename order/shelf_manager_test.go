package order

import (
	"fmt"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func fillTemperaturesShelfs(sm *ShelfManager) int {
	temps := []string{"hot", "cold", "frozen"}
	id := 1

	// fill all the temperature shelfs
	for _, temp := range temps {
		for i := 0; i < shelfCapacity; i += 1 {
			sm.Add(&Order{
				Id:           fmt.Sprint(id),
				Name:         "My Order",
				Temperature:  temp,
				ShelfLife:    300,
				DecayRate:    0.5,
				CreationTime: time.Now(),
			})

			id += 1
		}
	}

	return id
}

func TestInsert(t *testing.T) {
	hotOrder := Order{
		Id:           "1",
		Name:         "MyOrder",
		Temperature:  "hot",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	}
	coldOrder := Order{
		Id:           "2",
		Name:         "MyOrder",
		Temperature:  "cold",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	}
	frozenOrder := Order{
		Id:           "3",
		Name:         "MyOrder",
		Temperature:  "frozen",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	}

	sm := NewShelfManager()
	assert.Equal(t, len(sm.temperatureShelfs), 3)
	assert.Equal(t, sm.overflowShelf.Len(), 0)
	assert.Equal(t, len(sm.ordersMap), 0)
	assert.Equal(t, len(sm.overflowOrderMap), 3)

	assert.Equal(t, sm.temperatureShelfs[hotOrder.Temperature].Len(), 0)
	assert.Equal(t, sm.temperatureShelfs[coldOrder.Temperature].Len(), 0)
	assert.Equal(t, sm.temperatureShelfs[frozenOrder.Temperature].Len(), 0)

	// Insert 1 element in each temperature shelf
	sm.Add(&hotOrder)
	sm.Add(&coldOrder)
	sm.Add(&frozenOrder)

	assert.Equal(t, sm.temperatureShelfs[hotOrder.Temperature].Len(), 1)
	assert.Equal(t, sm.temperatureShelfs[coldOrder.Temperature].Len(), 1)
	assert.Equal(t, sm.temperatureShelfs[frozenOrder.Temperature].Len(), 1)
	assert.Equal(t, len(sm.ordersMap), 3)
	assert.Equal(t, sm.overflowShelf.Len(), 0)
}

func TestInsertOverflow(t *testing.T) {
	sm := NewShelfManager()
	id := fillTemperaturesShelfs(sm)

	assert.Equal(t, len(sm.ordersMap), len(sm.temperatureShelfs)*shelfCapacity)
	assert.Equal(t, sm.overflowShelf.Len(), 0)
	for _, v := range sm.temperatureShelfs {
		assert.Equal(t, v.Len(), shelfCapacity)
	}

	// this order should be in the overflow shelf
	sm.Add(&Order{
		Id:           fmt.Sprint(id),
		Name:         "MyOrder",
		Temperature:  "cold",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	})

	assert.Equal(t, sm.overflowShelf.Len(), 1)
	assert.Equal(t, sm.overflowShelf.queue[0].Id, fmt.Sprint(id))
	id += 2
	for _, v := range sm.temperatureShelfs {
		assert.Equal(t, v.Len(), shelfCapacity)
	}

	// make the overflow shelf full
	for i := 0; i < overflowShelfCapacity-1; i += 1 {
		sm.Add(&Order{
			Id:           fmt.Sprint(id),
			Name:         "MyOrder",
			Temperature:  "hot",
			ShelfLife:    300,
			DecayRate:    0.5,
			CreationTime: time.Now(),
		})

		id += 1
	}

	assert.Equal(t, sm.overflowShelf.Len(), overflowShelfCapacity)
	overflowTop := sm.overflowShelf.queue[0].Id
	sm.Add(&Order{
		Id:           fmt.Sprint(id),
		Name:         "MyOrder",
		Temperature:  "frozen",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	})
	id += 1

	// in case the overflow shelf is full, we discard the top of the heap
	assert.Assert(t, overflowTop != sm.overflowShelf.queue[0].Id)

	// open room in the hot shelf
	assert.Assert(t, sm.Remove(sm.temperatureShelfs["hot"].queue[0].Id) != nil)
	assert.Equal(t, sm.temperatureShelfs["hot"].Len(), shelfCapacity-1)

	// this operation should move one "hot" order to the hot shelf and then
	// insert our "cold" order to the overflow shelf
	sm.Add(&Order{
		Id:           fmt.Sprint(id),
		Name:         "MyOrder",
		Temperature:  "cold",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	})
	assert.Equal(t, sm.temperatureShelfs["hot"].Len(), shelfCapacity)
	assert.Equal(t, sm.overflowShelf.queue[sm.ordersMap[fmt.Sprint(id)].index].Id, fmt.Sprint(id))
}

func TestRemove(t *testing.T) {
	hotOrder := Order{
		Id:           "1",
		Name:         "MyOrder",
		Temperature:  "hot",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	}
	coldOrder := Order{
		Id:           "2",
		Name:         "MyOrder",
		Temperature:  "cold",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	}
	frozenOrder := Order{
		Id:           "3",
		Name:         "MyOrder",
		Temperature:  "frozen",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	}

	sm := NewShelfManager()

	sm.Add(&hotOrder)
	sm.Add(&coldOrder)
	sm.Add(&frozenOrder)

	assert.Equal(t, sm.Remove(hotOrder.Id).Id, hotOrder.Id)
	assert.Equal(t, sm.temperatureShelfs[hotOrder.Temperature].Len(), 0)
	assert.Equal(t, sm.Remove(coldOrder.Id).Id, coldOrder.Id)
	assert.Equal(t, sm.temperatureShelfs[coldOrder.Temperature].Len(), 0)
	assert.Equal(t, sm.Remove(frozenOrder.Id).Id, frozenOrder.Id)
	assert.Equal(t, sm.temperatureShelfs[frozenOrder.Temperature].Len(), 0)
}

func TestRemoveOverflow(t *testing.T) {
	sm := NewShelfManager()
	id := fillTemperaturesShelfs(sm)

	testOrder := Order{
		Id:           fmt.Sprint(id),
		Name:         "MyOrder",
		Temperature:  "cold",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: time.Now(),
	}

	// this should go in the overflow shelf
	sm.Add(&testOrder)

	assert.Equal(t, sm.overflowShelf.Len(), 1)
	assert.Equal(t, sm.Remove(testOrder.Id).Id, testOrder.Id)
	assert.Equal(t, sm.overflowShelf.Len(), 0)
}
