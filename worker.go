package worker

import (
	"fmt"
	"sync"

	"github.com/nouranatef/whale-watcher/models"
)

const WhaleThresholdSats = 10 * 100_000_000 // 10 BTC

type WorkerPool struct {
	JobQueue   chan models.Transaction
	NumWorkers int
	wg         sync.WaitGroup
	btcPrice   float64
	mu         sync.RWMutex
}

func NewWorkerPool(numWorkers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		JobQueue:   make(chan models.Transaction, queueSize),
		NumWorkers: numWorkers,
	}
}

func (wp *WorkerPool) SetBTCPrice(price float64) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.btcPrice = price
}

func (wp *WorkerPool) getBTCPrice() float64 {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.btcPrice
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.NumWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.JobQueue)
	wp.wg.Wait()
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for tx := range wp.JobQueue {
		wp.processTransaction(tx)
	}
}

func (wp *WorkerPool) processTransaction(tx models.Transaction) {
	var totalValue int64
	for _, out := range tx.Vout {
		totalValue += out.Value
	}

	if totalValue > WhaleThresholdSats {
		btcValue := float64(totalValue) / 100_000_000.0
		usdValue := btcValue * wp.getBTCPrice()

		fmt.Printf("🐳 Whale Alert! TxID: %s | Value: %.2f BTC ($%.2f)\n", tx.TxID, btcValue, usdValue)
	}
}
