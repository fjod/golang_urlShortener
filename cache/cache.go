package cache

import (
	"context"
	"time"
)

type Kvp struct {
	Key   string
	Value string
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

var service *Service = nil

func NewCacheService(lifetime time.Duration) *Service {
	if service != nil {
		return service
	}
	s := &Service{
		data:     make(map[string]Record),
		lifetime: lifetime,
		Key:      make(chan string),
		Value:    make(chan string),
		Add:      make(chan Kvp),
	}
	contChecking = true
	go KeepRunning(s)
	service = s
	return service
}

func KeepRunning(cache *Service) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for contChecking {
		select {
		case values, x := <-(*cache).Add:
			if !x {
				return
			}
			(*cache).data[values.Key] = Record{
				value: values.Value,
				added: time.Now(),
			}
		case k, x := <-(*cache).Key:
			//After the last Value has been received from a closed channel c,
			//any receive from c will succeed without blocking,
			//returning the zero Value for the channel element.
			if !x {
				return
			}
			value, ok := (*cache).data[k]
			if ok {
				(*cache).Value <- value.value
			}
			if !ok {
				(*cache).Value <- ""
			}

		case <-ticker.C:
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
	select {
	case <-ctx.Done():
		panic(ctx.Err())
	default:
		contChecking = false
		close(cache.Key)
		close(cache.Value)
		close(cache.Add)
	}
}
