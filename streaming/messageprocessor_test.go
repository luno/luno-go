package streaming

import (
	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"reflect"
	"testing"
)

type orderbookStatistics struct {
	Sequence  int64
	AskCount  int
	BidCount  int
	AskVolume decimal.Decimal
	BidVolume decimal.Decimal
}

func TestHandleMessageWithOrderbook(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage([]byte("\"\""))
	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	mp.HandleMessage([]byte("\"\""))
	mp.HandleMessage([]byte("\"\""))

	expected := orderbookStatistics{
		Sequence:  40413238,
		AskCount:  9214,
		BidCount:  3248,
		AskVolume: decimal.New(big.NewInt(784815424), 6),
		BidVolume: decimal.New(big.NewInt(2695234253), 6),
	}

	actual := calculateOrderbookStatistics(mp.OrderBookSnapshot())

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestHandleMessageWithInvalidOrderbook(t *testing.T) {
	mp := &messageProcessor{}

	err := mp.HandleMessage([]byte(`{"sequence": "40413238","asks": {"id": "BXEMZSYBRFYHSCF","price": "92655.00","volume": "0.495769"}, "bids": [{"id": "BXBAYA687URRT28","price": "92654.00","volume": "1.834379"}]}`))

	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestHandleMessageWithDelete(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	mp.HandleMessage([]byte(`{"sequence":"40413239","trade_updates":null,"create_update":null,"delete_update":{"order_id":"BXNC7TGBBJJ885S"},"timestamp":1530887350936}`))

	expected := orderbookStatistics{
		Sequence:  40413239,
		AskCount:  9214,
		BidCount:  3247,
		AskVolume: decimal.New(big.NewInt(784815424), 6),
		BidVolume: decimal.New(big.NewInt(2692184753), 6),
	}

	actual := calculateOrderbookStatistics(mp.OrderBookSnapshot())

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestHandleMessageWithInvalidDelete(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	err := mp.HandleMessage([]byte(`{"sequence":"40413239","trade_updates":null,"create_update":null,"delete_update":{"order_id":{"order_id":"BXNC7TGBBJJ885S"}},"timestamp":1530887350936}`))

	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestHandleMessageWithTrade(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	mp.HandleMessage([]byte(`{"sequence":"40413239","trade_updates":[{"base":"0.094976","counter":"8800.00128","maker_order_id":"BXEMZSYBRFYHSCF","taker_order_id":"BXGGSPFECZKFQ34","order_id":"BXEMZSYBRFYHSCF"}],"create_update":null,"delete_update":null,"timestamp":1530887351827}`))

	expected := orderbookStatistics{
		Sequence:  40413239,
		AskCount:  9214,
		BidCount:  3248,
		AskVolume: decimal.New(big.NewInt(784720448), 6),
		BidVolume: decimal.New(big.NewInt(2695234253), 6),
	}

	actual := calculateOrderbookStatistics(mp.OrderBookSnapshot())

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestHandleMessageWithNonpositiveTrade(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	err := mp.HandleMessage([]byte(`{"sequence":"40413239","trade_updates":[{"base":"-0.094976","counter":"8800.00128","maker_order_id":"BXEMZSYBRFYHSCF","taker_order_id":"BXGGSPFECZKFQ34","order_id":"BXEMZSYBRFYHSCF"}],"create_update":null,"delete_update":null,"timestamp":1530887351827}`))

	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestHandleMessageWithCreate(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	mp.HandleMessage([]byte(`{"sequence":"40413239","trade_updates":null,"create_update":{"order_id":"BXKQ7P9GK27486F","type":"BID","price":"88501.00","volume":"3.0485"},"delete_update":null,"timestamp":1530887351155}`))

	expected := orderbookStatistics{
		Sequence:  40413239,
		AskCount:  9214,
		BidCount:  3249,
		AskVolume: decimal.New(big.NewInt(784815424), 6),
		BidVolume: decimal.New(big.NewInt(2698282753), 6),
	}

	actual := calculateOrderbookStatistics(mp.OrderBookSnapshot())

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestHandleMessageWithUpdateBeforeOrderbook(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage([]byte(`{"sequence":"40413239","trade_updates":null,"create_update":{"order_id":"BXKQ7P9GK27486F","type":"BID","price":"88501.00","volume":"3.0485"},"delete_update":null,"timestamp":1530887351155}`))

	actualSeq, actualBids, actualAsks := mp.OrderBookSnapshot()

	if 0 != actualSeq {
		t.Errorf("Expected sequence to be 0, got %v", actualSeq)
	}
	if nil != actualBids {
		t.Errorf("Expected bids to be nil, got %v", actualBids)
	}
	if nil != actualAsks {
		t.Errorf("Expected asks to be nil, got %v", actualAsks)
	}
}

func TestHandleMessageWithPreviousUpdate(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	mp.HandleMessage([]byte(`{"sequence":"40413237","trade_updates":null,"create_update":{"order_id":"BXKQ7P9GK27486F","type":"BID","price":"88501.00","volume":"3.0485"},"delete_update":null,"timestamp":1530887351155}`))
	mp.HandleMessage([]byte(`{"sequence":"40413238","trade_updates":null,"create_update":{"order_id":"BXKQ7P9GK27486F","type":"BID","price":"88501.00","volume":"3.0485"},"delete_update":null,"timestamp":1530887351155}`))

	expected := orderbookStatistics{
		Sequence:  40413238,
		AskCount:  9214,
		BidCount:  3248,
		AskVolume: decimal.New(big.NewInt(784815424), 6),
		BidVolume: decimal.New(big.NewInt(2695234253), 6),
	}

	actual := calculateOrderbookStatistics(mp.OrderBookSnapshot())

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestHandleMessageWithOutOfSequenceUpdate(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	err := mp.HandleMessage([]byte(`{"sequence":"40413240","trade_updates":null,"create_update":{"order_id":"BXKQ7P9GK27486F","type":"BID","price":"88501.00","volume":"3.0485"},"delete_update":null,"timestamp":1530887351155}`))

	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestReset(t *testing.T) {
	mp := &messageProcessor{}

	mp.HandleMessage(loadFromFile(t, "fixture_orderbook.json"))
	mp.Reset()

	actualSeq, actualBids, actualAsks := mp.OrderBookSnapshot()

	if 0 != actualSeq {
		t.Errorf("Expected sequence to be 0, got %v", actualSeq)
	}
	if nil != actualBids {
		t.Errorf("Expected bids to be nil, got %v", actualBids)
	}
	if nil != actualAsks {
		t.Errorf("Expected asks to be nil, got %v", actualAsks)
	}
}

func calculateOrderbookStatistics(sequence int64, bids []luno.OrderBookEntry, asks []luno.OrderBookEntry) orderbookStatistics {
	var stats = orderbookStatistics{
		Sequence:  sequence,
		AskCount:  len(asks),
		BidCount:  len(bids),
		AskVolume: decimal.New(new(big.Int), 6),
		BidVolume: decimal.New(new(big.Int), 6),
	}

	for _, ask := range asks {
		stats.AskVolume = stats.AskVolume.Add(ask.Volume)
	}

	for _, bid := range bids {
		stats.BidVolume = stats.BidVolume.Add(bid.Volume)
	}

	return stats
}

func loadFromFile(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
