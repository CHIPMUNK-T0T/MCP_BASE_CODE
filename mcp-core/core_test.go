package mcpcore

import (
	"encoding/json"
	"testing"
)

func TestRequestMarshal(t *testing.T) {
	params := map[string]string{"key": "value"}
	paramsBytes, _ := json.Marshal(params)

	req := Request{
		JsonRpc: "2.0",
		Method:  "test_method",
		Params:  paramsBytes,
		ID:      1,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	expected := `{"jsonrpc":"2.0","method":"test_method","params":{"key":"value"},"id":1}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestResponseUnmarshal(t *testing.T) {
	jsonStr := `{"jsonrpc": "2.0", "result": {"success": true}, "id": 1}`

	var resp Response
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.JsonRpc != "2.0" {
		t.Errorf("Expected jsonrpc 2.0, got %s", resp.JsonRpc)
	}

	// idは数値だが、json.Unmarshalではinterface{}としてfloat64になることが多い
	id, ok := resp.ID.(float64)
	if !ok || id != 1 {
		t.Errorf("Expected id 1, got %v", resp.ID)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok || resultMap["success"] != true {
		t.Errorf("Expected result success=true, got %v", resp.Result)
	}
}
