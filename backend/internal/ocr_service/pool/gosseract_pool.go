package pool

import (
	"github.com/otiai10/gosseract/v2"
	"log"
)

// GosseractPool represents a pool of reusable gosseract.Client instances for OCR operations.
// It helps manage the allocation and deallocation of gosseract.Client efficiently.
type GosseractPool struct {
	pool chan *gosseract.Client
}

// NewGosseractPool initializes a pool of gosseract.Client instances with a specified size and language setting.
func NewGosseractPool(size int, lang string) *GosseractPool {
	pool := make(chan *gosseract.Client, size)
	for i := 0; i < size; i++ {
		tempClient := gosseract.NewClient()
		err := tempClient.SetLanguage(lang)
		if err != nil {
			log.Fatalf("Failed to set language to %s", lang)
		}
		pool <- tempClient
	}
	return &GosseractPool{pool: pool}
}

// Get retrieves a *gosseract.Client instance from the pool for OCR operations. It blocks if the pool is empty until available.
func (cp *GosseractPool) Get() *gosseract.Client {
	return <-cp.pool
}

// Put adds a gosseract.Client instance back into the pool for reuse.
func (cp *GosseractPool) Put(client *gosseract.Client) {
	cp.pool <- client
}

// Close releases all clients in the pool and shuts down the pool, ensuring proper cleanup of resources.
func (cp *GosseractPool) Close() {
	close(cp.pool)
	for client := range cp.pool {
		client.Close()
	}
}
