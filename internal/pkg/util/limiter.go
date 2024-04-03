package util

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

var (
	onceLimiter            sync.Once
	attributeSlidingWindow AttributeSlidingWindow
)

type AttributeSlidingWindow struct {
	mu  sync.Mutex // guards
	hub map[string]fiber.Handler
}

// New creates a new sliding window middleware handler
func (as *AttributeSlidingWindow) New(cfg limiter.Config) fiber.Handler {
	onceLimiter.Do(func() {
		attributeSlidingWindow = AttributeSlidingWindow{hub: make(map[string]fiber.Handler, 0)}
	})

	return func(ctx *fiber.Ctx) error {
		if itemLimiter, ok := attributeSlidingWindow.hub[cfg.KeyGenerator(ctx)]; ok {
			return itemLimiter(ctx)
		} else {
			as.mu.Lock()
			defer as.mu.Unlock()
			attributeSlidingWindow.hub[cfg.KeyGenerator(ctx)] = limiter.SlidingWindow{}.New(cfg)
			return attributeSlidingWindow.hub[cfg.KeyGenerator(ctx)](ctx)
		}
	}
}
