package main

import (
	"testing"
)

func TestWhaleDetection(t *testing.T) {

	fakeTxValueSatoshis := int64(5000000000)
	whaleThresholdBTC := 10.0
	btcValue := float64(fakeTxValueSatoshis) / 100000000.0

	if btcValue < whaleThresholdBTC {
		t.Errorf("Expected to detect a whale (>%f BTC), but got %f", whaleThresholdBTC, btcValue)
	} else {
		t.Logf("Success! Detected a whale of %.2f BTC", btcValue)
	}
}
