package postgres

import "time"

type Option func(*PostgreDB)

func MaxPoolSize(size int) Option {
	return func(c *PostgreDB) {
		c.maxPoolSize = size
	}
}

func ConnectionAttempts(attempts int) Option {
	return func(c *PostgreDB) {
		c.connectionAttempts = attempts
	}
}

func ConnectionTimeout(timeout time.Duration) Option {
	return func(c *PostgreDB) {
		c.connectionTimeout = timeout
	}
}
