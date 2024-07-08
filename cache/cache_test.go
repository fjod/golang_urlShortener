package cache

import (
	"context"
	"testing"
	"time"
)

func TestGetFromCache(t *testing.T) {
	sut := NewCacheService(time.Second)
	defer (*sut).Exit(context.Background())
	(*sut).Add <- Kvp{"testKey", "testValue"}
	go func() {
		// wait for something to be received from the channel
		actual := <-(*sut).Value
		if actual != "testValue" {
			t.Errorf("actual = %v, want %v", actual, "testValue")
		}
	}()
	(*sut).Key <- "testKey"
}
