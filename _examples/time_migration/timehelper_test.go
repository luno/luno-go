package main

import (
	"encoding/json"
	"fmt"
)

// ExampleTimeCompatHelper demonstrates how to use TimeCompatHelper to handle both
// the old string format and new millisecond format for timestamps.
func ExampleTimeCompatHelper() {
	// Example with new format (milliseconds)
	newFormat := []byte(`{"timestamp":1621234567890}`)

	// Example with old format (string)
	oldFormat := []byte(`{"timestamp":"2021-05-17 12:34:56 +0000 UTC"}`)

	// Handle both formats
	for _, data := range [][]byte{newFormat, oldFormat} {
		var result map[string]json.RawMessage
		if err := json.Unmarshal(data, &result); err != nil {
			panic(err)
		}

		t, err := TimeCompatHelper(result["timestamp"])
		if err != nil {
			panic(err)
		}

		fmt.Printf("Parsed time: %v\n", t.UTC().Format("2006-01-02 15:04:05.999999999 -0700 MST"))
	}
	// Output:
	// Parsed time: 2021-05-17 06:56:07.89 +0000 UTC
	// Parsed time: 2021-05-17 12:34:56 +0000 UTC
}
