package utils

import (
	"github.com/otiai10/gosseract/v2"
	"log"
)

type GosseractPool struct {
	pool chan *gosseract.Client
}

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

func (cp *GosseractPool) Get() *gosseract.Client {
	return <-cp.pool
}

func (cp *GosseractPool) Put(client *gosseract.Client) {
	cp.pool <- client
}

func (cp *GosseractPool) Close() {
	close(cp.pool)
	for client := range cp.pool {
		client.Close()
	}
}
