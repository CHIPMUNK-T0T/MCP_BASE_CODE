package mcpcore

import (
	"encoding/json"
	"fmt"
)

const (
	JSONRPCVersion               = "2.0"
	MCPProtocolVersion           = "2025-11-25"
	MethodInitialize             = "initialize"
	MethodInitialized            = "notifications/initialized"
	MethodProgress               = "notifications/progress"
	MethodPromptsListChanged     = "notifications/prompts/list_changed"
	MethodResourcesListChanged   = "notifications/resources/list_changed"
	MethodToolsListChanged       = "notifications/tools/list_changed"
	MethodPromptsList            = "prompts/list"
	MethodPromptsGet             = "prompts/get"
	MethodResourcesList          = "resources/list"
	MethodResourcesRead          = "resources/read"
	MethodResourcesTemplatesList = "resources/templates/list"
	MethodCompletionComplete     = "completion/complete"
	MethodToolsList              = "tools/list"
	MethodToolsCall              = "tools/call"
)

// --- Request DTOs ---

type InitializeRequest struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      interface{}      `json:"id"`
	Method  string           `json:"method"`
	Params  InitializeParams `json:"params"`
}

func NewInitializeRequest(id interface{}, clientName, clientVersion string) InitializeRequest {
	return InitializeRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  MethodInitialize,
		Params: InitializeParams{
			ProtocolVersion: MCPProtocolVersion,
			Capabilities: ClientCapabilities{
				Roots: &RootsCapability{
					ListChanged: true,
				},
				Sampling: &SamplingCapability{
					Context: &EmptyObject{},
					Tools:   &EmptyObject{},
				},
			},
			ClientInfo: Implementation{
				Name:    clientName,
				Version: clientVersion,
			},
		},
	}
}

type InitializedNotification struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

func NewInitializedNotification() InitializedNotification {
	return InitializedNotification{
		JSONRPC: JSONRPCVersion,
		Method:  MethodInitialized,
	}
}

type ListChangedNotification struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

func NewToolsListChangedNotification() ListChangedNotification {
	return ListChangedNotification{
		JSONRPC: JSONRPCVersion,
		Method:  MethodToolsListChanged,
		Params:  map[string]interface{}{},
	}
}

func NewPromptsListChangedNotification() ListChangedNotification {
	return ListChangedNotification{
		JSONRPC: JSONRPCVersion,
		Method:  MethodPromptsListChanged,
		Params:  map[string]interface{}{},
	}
}

func NewResourcesListChangedNotification() ListChangedNotification {
	return ListChangedNotification{
		JSONRPC: JSONRPCVersion,
		Method:  MethodResourcesListChanged,
		Params:  map[string]interface{}{},
	}
}

type ProgressNotification struct {
	JSONRPC string                     `json:"jsonrpc"`
	Method  string                     `json:"method"`
	Params  ProgressNotificationParams `json:"params"`
}

func NewProgressNotification(params ProgressNotificationParams) ProgressNotification {
	return ProgressNotification{
		JSONRPC: JSONRPCVersion,
		Method:  MethodProgress,
		Params:  params,
	}
}

type ListToolsParams struct {
	Meta   map[string]interface{} `json:"_meta,omitempty"`
	Cursor string                 `json:"cursor,omitempty"`
}

type ListPromptsParams struct {
	Meta   map[string]interface{} `json:"_meta,omitempty"`
	Cursor string                 `json:"cursor,omitempty"`
}

type ListResourcesParams struct {
	Meta   map[string]interface{} `json:"_meta,omitempty"`
	Cursor string                 `json:"cursor,omitempty"`
}

type ListResourceTemplatesParams struct {
	Meta   map[string]interface{} `json:"_meta,omitempty"`
	Cursor string                 `json:"cursor,omitempty"`
}

type ListToolsRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  ListToolsParams `json:"params"`
}

func NewListToolsRequest(id interface{}) ListToolsRequest {
	return ListToolsRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  MethodToolsList,
		Params:  ListToolsParams{},
	}
}

type ListPromptsRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      interface{}       `json:"id"`
	Method  string            `json:"method"`
	Params  ListPromptsParams `json:"params"`
}

func NewListPromptsRequest(id interface{}) ListPromptsRequest {
	return ListPromptsRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  MethodPromptsList,
		Params:  ListPromptsParams{},
	}
}

type GetPromptRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  GetPromptParams `json:"params"`
}

func NewGetPromptRequest(id interface{}, name string, args map[string]string) GetPromptRequest {
	return GetPromptRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  MethodPromptsGet,
		Params: GetPromptParams{
			Name:      name,
			Arguments: args,
		},
	}
}

type ListResourcesRequest struct {
	JSONRPC string              `json:"jsonrpc"`
	ID      interface{}         `json:"id"`
	Method  string              `json:"method"`
	Params  ListResourcesParams `json:"params"`
}

func NewListResourcesRequest(id interface{}) ListResourcesRequest {
	return ListResourcesRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  MethodResourcesList,
		Params:  ListResourcesParams{},
	}
}

type ListResourceTemplatesRequest struct {
	JSONRPC string                      `json:"jsonrpc"`
	ID      interface{}                 `json:"id"`
	Method  string                      `json:"method"`
	Params  ListResourceTemplatesParams `json:"params"`
}

func NewListResourceTemplatesRequest(id interface{}) ListResourceTemplatesRequest {
	return ListResourceTemplatesRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  MethodResourcesTemplatesList,
		Params:  ListResourceTemplatesParams{},
	}
}

type ReadResourceRequest struct {
	JSONRPC string             `json:"jsonrpc"`
	ID      interface{}        `json:"id"`
	Method  string             `json:"method"`
	Params  ReadResourceParams `json:"params"`
}

func NewReadResourceRequest(id interface{}, uri string) ReadResourceRequest {
	return ReadResourceRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  MethodResourcesRead,
		Params: ReadResourceParams{
			URI: uri,
		},
	}
}

type CompleteRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      interface{}    `json:"id"`
	Method  string         `json:"method"`
	Params  CompleteParams `json:"params"`
}

func NewCompleteRequest(id interface{}, ref interface{}, argument CompletionArgument, context *CompletionContext) CompleteRequest {
	return CompleteRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  MethodCompletionComplete,
		Params: CompleteParams{
			Ref:      ref,
			Argument: argument,
			Context:  context,
		},
	}
}

type CallToolRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      interface{}    `json:"id"`
	Method  string         `json:"method"`
	Params  CallToolParams `json:"params"`
}

func NewCallToolRequest(id interface{}, toolName string, args map[string]interface{}) CallToolRequest {
	return CallToolRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  MethodToolsCall,
		Params: CallToolParams{
			Name:      toolName,
			Arguments: args,
		},
	}
}

// --- Response DTOs ---

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// --- DTO <-> JSON helpers ---

func MarshalDTO(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalDTO(data []byte, dst interface{}) error {
	return json.Unmarshal(data, dst)
}

func DecodeJSONRPCResponse(data []byte) (JSONRPCResponse, error) {
	var resp JSONRPCResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return JSONRPCResponse{}, err
	}
	return resp, nil
}

func DecodeJSONRPCResult[T any](resp JSONRPCResponse) (T, error) {
	var zero T
	if resp.Error != nil {
		return zero, fmt.Errorf("json-rpc error (%d): %s", resp.Error.Code, resp.Error.Message)
	}
	if len(resp.Result) == 0 {
		return zero, fmt.Errorf("json-rpc result is empty")
	}

	var out T
	if err := json.Unmarshal(resp.Result, &out); err != nil {
		return zero, err
	}
	return out, nil
}
