package anypool

import (
	"github.com/google/uuid"
	"sync"
)

type AnyPool interface {
	Get() *Wrapper
	Size() int
	Close()
}

type Wrapper struct {
	Conn       any
	ReuseCount int
	ID         uuid.UUID
}

type anyPool struct {
	mu         sync.Mutex
	size       int
	reuseLimit int
	factory    func() any
	clients    chan *Wrapper
}

func New(size int, factory func() any, opts ...Option) AnyPool {
	pool := &anyPool{
		size:    size,
		factory: factory,
		clients: make(chan *Wrapper, size),
	}
	for _, opt := range opts {
		opt(pool)
	}

	for i := 0; i < pool.size; i++ {
		client := pool.factory()
		pool.clients <- &Wrapper{
			Conn:       client,
			ReuseCount: 0,
			ID:         uuid.New(),
		}
	}

	return pool
}

func (c *anyPool) Get() *Wrapper {
	for {
		select {
		case wrapper := <-c.clients:
			if wrapper.ReuseCount < c.reuseLimit {
				wrapper.ReuseCount++
				c.clients <- wrapper
				return wrapper
			}

			newClient := c.factory()
			newWrapper := &Wrapper{
				Conn:       newClient,
				ReuseCount: 0,
				ID:         uuid.New(),
			}
			c.clients <- newWrapper
		default:
			c.mu.Lock()
			newClient := c.factory()
			newWrapper := &Wrapper{
				Conn:       newClient,
				ReuseCount: 0,
				ID:         uuid.New(),
			}
			c.clients <- newWrapper
			c.mu.Unlock()
		}
	}
}

func (c *anyPool) Size() int {
	return len(c.clients)
}

func (c *anyPool) Close() {
	close(c.clients)
}
