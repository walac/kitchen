package order

import (
	"container/heap"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestShelf(t *testing.T) {
	shelf := NewShelf(1)
	heap.Init(shelf)

	order1 := Order{
		Id:           "1",
		Name:         "MyOrder",
		Temperature:  "hot",
		ShelfLife:    300,
		DecayRate:    1.0,
		CreationTime: time.Now(),
	}

	heap.Push(shelf, &order1)
	assert.Equal(t, shelf.Len(), 1)
	assert.Equal(t, order1.Index(), 0)

	order2 := Order{
		Id:           "2",
		Name:         "MyOrder",
		Temperature:  "hot",
		ShelfLife:    10,
		DecayRate:    1.0,
		CreationTime: time.Now(),
	}

	heap.Push(shelf, &order2)
	assert.Equal(t, shelf.Len(), 2)
	assert.Equal(t, order2.Index(), 0)
	assert.Equal(t, order1.Index(), 1)
}
