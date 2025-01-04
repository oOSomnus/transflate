package pool

import (
	"github.com/otiai10/gosseract/v2"
	"log"
)

/*
GosseractPool is a structure that manages a pool of gosseract.Client instances.

Fields:
  - pool (chan *gosseract.Client): A buffered channel holding gosseract.Client instances.
*/
type GosseractPool struct {
	pool chan *gosseract.Client
}

/*
NewGosseractPool initializes a new pool of gosseract.Client instances.

Parameters:
  - size (int): The number of gosseract.Client instances to initialize in the pool.
  - lang (string): The language to set for the gosseract.Client instances.

Returns:
  - (*GosseractPool): A pointer to a new GosseractPool containing initialized gosseract.Client instances.
*/
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

/*
Get retrieves a gosseract.Client instance from the pool.

Returns:
  - (*gosseract.Client): A pointer to a gosseract.Client instance from the pool.
*/
func (cp *GosseractPool) Get() *gosseract.Client {
	return <-cp.pool
}

/*
Put returns a gosseract.Client instance back to the pool.

Parameters:
  - client (*gosseract.Client): The gosseract.Client instance to be returned to the pool.
*/
func (cp *GosseractPool) Put(client *gosseract.Client) {
	cp.pool <- client
}

/*
Close closes the GosseractPool and releases all gosseract.Client resources.
*/
func (cp *GosseractPool) Close() {
	close(cp.pool)
	for client := range cp.pool {
		client.Close()
	}
}
