package bus

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestEvent struct{}

func TestSubscribe(t *testing.T) {
	topic := fmt.Sprintf("%T", *new(TestEvent))
	closeSub1 := Subscribe("", func(ctx context.Context, event TestEvent) error { return nil })
	closeSub2 := Subscribe("", func(ctx context.Context, event TestEvent) error { return nil })

	assert.Len(t, subHandlers[topic], 2)
	assert.Equal(t, 0, subHandlers[topic][0].ID)
	assert.Equal(t, 1, subHandlers[topic][1].ID)

	closeSub1()

	assert.Len(t, subHandlers[topic], 1)
	assert.Equal(t, 2, cap(subHandlers[topic]))
	assert.Equal(t, 1, subHandlers[topic][0].ID)

	closeSub2()
	assert.Len(t, subHandlers[topic], 0)
}
