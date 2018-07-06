package streaming

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestHandleMessage(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage([]byte("\"\""))
	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	mp.HandleMessage([]byte("\"\""))
	mp.HandleMessage([]byte("\"\""))

	actualSequence, actualBids, actualAsks := mp.OrderBookSnapshot()

	var expectedSequence int64 = 40413238
	if actualSequence != expectedSequence {
		t.Errorf("Expected sequence to be %v, got %v", expectedSequence, actualSequence)
	}

	var expectedBidCount = 3248
	var actualBidCount = len(actualBids)
	if expectedBidCount != actualBidCount {
		t.Errorf("Expected bid count to be %v, got %v", expectedBidCount, actualBidCount)
	}

	var expectedAskCount = 9214
	var actualAskCount = len(actualAsks)
	if expectedAskCount != actualAskCount {
		t.Errorf("Expected ask count to be %v, got %v", expectedAskCount, actualAskCount)
	}
}

func loadFromFile(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
