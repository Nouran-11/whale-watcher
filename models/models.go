package models

// WebSocketResponse represents the structure of the message received from the WebSocket.
type WebSocketResponse struct {
	Block Block `json:"block"`
}

// Block represents a Bitcoin block.
type Block struct {
	ID     string `json:"id"`
	Height int    `json:"height"`
	Size   int    `json:"size"`
}

// Transaction represents a Bitcoin transaction.
type Transaction struct {
	TxID string `json:"txid"`
	Vout []Vout `json:"vout"`
}

// Vout represents a transaction output.
type Vout struct {
	ScriptPubKey string `json:"scriptpubkey"`
	ScriptType   string `json:"scriptpubkey_type"`
	Value        int64  `json:"value"` // Value in Satoshis
	Address      string `json:"scriptpubkey_address,omitempty"`
}

// WhaleTransaction holds the info we want to log.
type WhaleTransaction struct {
	TxID     string
	ValueBTC float64
	ValueUSD float64
}

type PriceInfo struct {
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
}
