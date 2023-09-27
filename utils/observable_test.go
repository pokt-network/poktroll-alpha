package utils_test

import (
	"poktroll/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObservable(t *testing.T) {
	obs, controller := utils.NewControlledObservable[int](nil)

	subscription := obs.Subscribe()
	ch := subscription.Ch()
	go func() {
		controller <- 1
		close(controller)
	}()

	counter := 0
	for value := range ch {
		require.Equal(t, 1, value)
		counter++
	}

	require.Equal(t, 1, counter)
}
