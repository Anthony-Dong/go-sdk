package json_tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestPretty(t *testing.T) {
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(map[string]interface{}{
		"k1": "v1",
		"k2": map[string]interface{}{
			"k2_1": "v1",
		},
	}); err != nil {
		t.Fatal(err)
	}
	fmt.Println(buffer.String())
}
