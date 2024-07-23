package anypool

type Option func(pool *anyPool)

func WithReuseLimit(count int) Option {
	return func(pool *anyPool) {
		pool.reuseLimit = count
	}
}
