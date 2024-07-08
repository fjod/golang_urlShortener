package cache

import (
	"context"
	"time"
)

type Cache interface {
	SetLifetime(lifetime time.Duration)
	Exit(ctx context.Context)
}

type Kvp struct {
	key   string
	value string
}

type Record struct {
	value string
	added time.Time
}

type Service struct {
	data     map[string]Record
	lifetime time.Duration

	// Key в этот канал принять ключ и вернуть значение в канал Value
	Key chan string
	// Value в этот канал отдается значение, которое было запрошено из канала Key
	Value chan string
	// Kvp канал для записи новых данных в кэш
	Add chan Kvp
}

func NewCacheService(lifetime time.Duration) *Service {
	service := &Service{
		data:     make(map[string]Record),
		lifetime: lifetime,
		Key:      make(chan string),
		Value:    make(chan string),
		Add:      make(chan Kvp),
	}
	contChecking = true
	go KeepRunning(service)
	return service
}

func KeepRunning(cache *Service) {
	for contChecking {
		select {
		case values := <-(*cache).Add:
			(*cache).data[values.key] = Record{
				value: values.value,
				added: time.Now(),
			}
		case k := <-(*cache).Key:
			if value, ok := (*cache).data[k]; ok {
				(*cache).Value <- value.value
			}

		default:
			time.Sleep(time.Millisecond)
			for key, record := range (*cache).data {
				if !contChecking {
					return
				}
				if time.Since(record.added) > (*cache).lifetime {
					delete((*cache).data, key)
				}
			}
		}
	}
}

var contChecking = false

func (cache *Service) SetLifetime(lifetime time.Duration) {
	cache.lifetime = lifetime
}

func (cache *Service) Exit(ctx context.Context) {
	contChecking = false
	close(cache.Key)
	close(cache.Value)
	close(cache.Add)
	ctx.Done()
}
