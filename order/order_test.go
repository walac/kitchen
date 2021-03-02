package order

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestDeadline(t *testing.T) {
	now := time.Now()
	testOrder := Order{
		Id:           "1",
		Name:         "MyOrder",
		Temperature:  "hot",
		ShelfLife:    300,
		DecayRate:    0.5,
		CreationTime: now,
	}

	assert.Equal(t, testOrder.CreationTime.Add(600*time.Second), testOrder.Deadline(1))
	assert.Equal(t, testOrder.CreationTime, now)
}

func TestIsWasted(t *testing.T) {
	testOrder := Order{
		Id:           "1",
		Name:         "My order",
		Temperature:  "hot",
		ShelfLife:    1,
		DecayRate:    1.0,
		CreationTime: time.Now(),
	}

	time.Sleep(1100 * time.Millisecond)
	assert.Assert(t, testOrder.IsWasted(1))
}
