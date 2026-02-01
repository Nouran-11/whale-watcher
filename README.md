# Whale Watcher 🐳

A high-performance Go backend service that monitors the Bitcoin blockchain for "Whale" transactions (> 10 BTC).

## Features
- **Real-time Monitoring**: Connects to Mempool.space WebSocket (`mempool-blocks` event).
- **Concurrency**: Uses a Worker Pool pattern to process transactions efficiently.
- **Whale Alert**: Logs transactions > 10 BTC with USD value.

## Technology Stack
- **Language**: Go (Golang)
- **Libs**: `github.com/gorilla/websocket`
- **API**: Mempool.space (WebSocket & REST)

## Setup & Run

1. **Clone the repository**
   ```bash
   git clone <repo_url>
   cd whale-watcher
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the Service**
   ```bash
   go run main.go
   ```

## How to Test

### 1. Unit Tests
Run the included unit tests to verify the whale detection logic:
```bash
go test -v
```
You should see output indicating that the logic correctly identifies transactions over the threshold.

### 2. Manual Verification
Run the application and wait for a new block (approx. every 10 minutes), or verify the initial connection log:
```bash
go run main.go
```
Expected Output:
```
Current BTC Price: $XXXXX.XX
Connected to Mempool.space WebSocket. Listening for blocks...
```

## Architecture
- **Main**: Manages WS connection and dispatches Block IDs.
- **Worker Pool**: A pool of workers (`worker/worker.go`) processes transactions concurrently.
- **Models**: internal data structures.

## Logic
1. Service starts and fetches current BTC price.
2. Listens for `mempool-blocks` WebSocket event.
3. When a block is found, it fetches the first 100 transactions (mocking full block processing for demo) from Mempool.space REST API.
4. Transactions are sent to the Worker Pool.
5. Workers check if `Value > 10 BTC`. If yes, logs: `🐳 Whale Alert! ...`

## Note
- The service fetches the first 100 transactions of a block to avoid hitting API rate limits during the demo. In a production environment, full block pagination would be implemented.
