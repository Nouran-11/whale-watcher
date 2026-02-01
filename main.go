package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nouranatef/whale-watcher/models"
	"github.com/nouranatef/whale-watcher/worker"
)

const (
	MempoolWS    = "wss://mempool.space/api/v1/ws"
	MempoolAPI   = "https://mempool.space/api"
	WorkerCount  = 5
	JobQueueSize = 100
)

type WSMessage struct {
	Block models.Block `json:"block"`
}

type PriceResponse struct {
	USD float64 `json:"USD"` // Check actual response structure
}

func main() {
	// 1. Setup Worker Pool
	wp := worker.NewWorkerPool(WorkerCount, JobQueueSize)
	wp.Start()
	defer wp.Stop()

	// 2. Fetch Initial Price
	price := fetchBTCPrice()
	wp.SetBTCPrice(price)
	fmt.Printf("Current BTC Price: $%.2f\n", price)

	// 3. Connect to WebSocket
	c, _, err := websocket.DefaultDialer.Dial(MempoolWS, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// 4. Subscribe to blocks
	// The prompt specifically asked for "mempool-blocks" event.
	subMsg := `{"action":"want","data":["mempool-blocks"]}`
	err = c.WriteMessage(websocket.TextMessage, []byte(subMsg))
	if err != nil {
		log.Println("write:", err)
		return
	}
	fmt.Println("Connected to Mempool.space WebSocket. Listening for blocks...")

	// 5. Handle Interrupts
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			handleMessage(message, wp)
		}
	}()

	<-interrupt
	fmt.Println("Shutting down...")
}

func handleMessage(msg []byte, wp *worker.WorkerPool) {
	// The mempool-blocks event structure is slightly different depending on if it's the `init` or update.
	// But usually it has a `block` field.
	// Actually, the `want` `blocks` returns an initial array of blocks, then updates.
	// Let's generic parse.

	var response struct {
		Block models.Block `json:"block"`
	}
	// Try parsing as single block notification
	if err := json.Unmarshal(msg, &response); err == nil && response.Block.ID != "" {
		fmt.Printf("📦 New Block Found: %d (Hash: %s)\n", response.Block.Height, response.Block.ID)
		go processBlock(response.Block.ID, wp)
		return
	}

	// If it's an array (initial state), we might ignore or process the tip.
	// For now, let's just focus on updates which usually come as object.
}

func processBlock(blockHash string, wp *worker.WorkerPool) {
	// Fetch transactions from REST API
	// Note: /block/:hash/txs is paginated (25 per page).
	// To be thorough we should fetch all. For demo, we'll fetch the first 50 (2 pages).

	client := &http.Client{Timeout: 10 * time.Second}

	// Fetch first page
	fetchAndEnqueue(client, blockHash, "", wp)
}

func fetchAndEnqueue(client *http.Client, blockHash string, lastTxID string, wp *worker.WorkerPool) {
	// Construct URL. Pagination with /txs usually uses a start index or last txid?
	// Mempool.space API docs say: GET /block/:hash/txs[/:start_index]
	// start_index is an integer offset.

	// Let's implement a simple loop for the first few pages.
	for i := 0; i < 100; i += 25 { // Fetch first 100 txs
		url := fmt.Sprintf("%s/block/%s/txs/%d", MempoolAPI, blockHash, i)
		if i == 0 {
			url = fmt.Sprintf("%s/block/%s/txs", MempoolAPI, blockHash)
		}

		resp, err := client.Get(url)
		if err != nil {
			log.Printf("Failed to fetch txs: %v", err)
			return
		}
		defer resp.Body.Close()

		var txs []models.Transaction
		if err := json.NewDecoder(resp.Body).Decode(&txs); err != nil {
			log.Printf("Failed to decode txs: %v", err)
			return
		}

		if len(txs) == 0 {
			break
		}

		for _, tx := range txs {
			wp.JobQueue <- tx
		}
	}
}

func fetchBTCPrice() float64 {
	// Use mempool.space price endpoint? Or coindesk?
	// Mempool.space doesn't document a public price API easily found in v1/ws docs, but `mempool.space/api/v1/prices` exists?
	// Let's check `https://mempool.space/api/v1/prices`.
	// Response: {"USD": 12345, ...}

	resp, err := http.Get("https://mempool.space/api/v1/prices")
	if err != nil {
		log.Println("Error fetching price:", err)
		return 60000.0 // Fallback
	}
	defer resp.Body.Close()

	var p models.PriceInfo
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return 60000.0
	}
	return p.USD
}
