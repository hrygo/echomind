package registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/pkg/ai"
)

// FactoryFunc is a function that creates an AIProvider.
type FactoryFunc func(ctx context.Context, settings configs.ProviderSettings, prompts map[string]string) (ai.AIProvider, error)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]FactoryFunc)
)

// Register registers a factory function for a given provider protocol.
// This function is typically called in the init() function of a provider package.
func Register(protocol string, factory FactoryFunc) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if factory == nil {
		panic("ai/registry: Register factory is nil")
	}
	if _, dup := providers[protocol]; dup {
		panic("ai/registry: Register called twice for provider " + protocol)
	}
	providers[protocol] = factory
}

// Get returns the factory function for a given provider protocol.
func Get(protocol string) (FactoryFunc, error) {
	providersMu.RLock()
	defer providersMu.RUnlock()
	factory, ok := providers[protocol]
	if !ok {
		return nil, fmt.Errorf("ai/registry: unknown provider protocol %q (forgot to import?)", protocol)
	}
	return factory, nil
}
