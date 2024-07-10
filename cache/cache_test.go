package cache

import (
	"context"
	"testing"
	"time"
)

func TestGetFromCache(t *testing.T) {
	sut := NewCacheService(time.Second)
	(*sut).Add <- Kvp{"testKey", "testValue"}
	go func() {
		actual := <-(*sut).Value
		if actual != "testValue" {
			t.Errorf("actual = %v, want %v", actual, "testValue")
		}
		sut.Exit(context.Background())
	}()
	(*sut).Key <- "testKey"
}

func TestCacheEvicting(t *testing.T) {
	sut := NewCacheService(time.Second)

	(*sut).Add <- Kvp{"testKey", "testValue"}
	time.Sleep(time.Second * 2)

	(*sut).Key <- "testKey"
	actual := <-(*sut).Value
	if actual != "" {
		t.Errorf("actual = %v, want %v", actual, "")
	}
	(*sut).Exit(context.Background())
}
