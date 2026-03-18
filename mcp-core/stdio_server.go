package mcpcore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// ServeStdio runs a JSON-RPC 2.0 loop over stdin/stdout using the provided MCPServer.
func ServeStdio(impl MCPServer) error {
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
			continue
		}

		// Notifications (no ID) should not get a response.
		if req.Method == MethodInitialized && req.ID == nil {
			_ = emitNotification(NewToolsListChangedNotification())
			_ = emitNotification(NewPromptsListChangedNotification())
			_ = emitNotification(NewResourcesListChangedNotification())
			continue
		}

		progressToken := extractProgressToken(req.Params)
		if req.ID != nil && progressToken != nil {
			_ = emitNotification(NewProgressNotification(ProgressNotificationParams{
				ProgressToken: progressToken,
				Progress:      0,
				Total:         1,
				Message:       "started",
			}))
		}

		resp := Response{
			JsonRpc: JSONRPCVersion,
			ID:      req.ID,
		}

		switch req.Method {
		case MethodInitialize:
			var params InitializeParams
			if len(req.Params) > 0 {
				if err := json.Unmarshal(req.Params, &params); err != nil {
					resp.Error = &Error{Code: -32602, Message: "Invalid params"}
					break
				}
			}
			result, err := impl.HandleInitialize(params)
			if err != nil {
				resp.Error = &Error{Code: -32000, Message: err.Error()}
				break
			}
			resp.Result = result

		case MethodPromptsList:
			result, err := impl.HandleListPrompts()
			if err != nil {
				resp.Error = &Error{Code: -32000, Message: err.Error()}
				break
			}
			resp.Result = result

		case MethodPromptsGet:
			var params GetPromptParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				resp.Error = &Error{Code: -32602, Message: "Invalid params"}
				break
			}
			result, err := impl.HandleGetPrompt(params)
			if err != nil {
				resp.Error = &Error{Code: -32000, Message: err.Error()}
				break
			}
			resp.Result = result

		case MethodResourcesList:
			result, err := impl.HandleListResources()
			if err != nil {
				resp.Error = &Error{Code: -32000, Message: err.Error()}
				break
			}
			resp.Result = result

		case MethodResourcesTemplatesList:
			result, err := impl.HandleListResourceTemplates()
			if err != nil {
				resp.Error = &Error{Code: -32000, Message: err.Error()}
				break
			}
			resp.Result = result

		case MethodResourcesRead:
			var params ReadResourceParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				resp.Error = &Error{Code: -32602, Message: "Invalid params"}
				break
			}
			result, err := impl.HandleReadResource(params)
			if err != nil {
				resp.Error = &Error{Code: -32000, Message: err.Error()}
				break
			}
			resp.Result = result

		case MethodCompletionComplete:
			var params CompleteParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				resp.Error = &Error{Code: -32602, Message: "Invalid params"}
				break
			}
			result, err := impl.HandleComplete(params)
			if err != nil {
				resp.Error = &Error{Code: -32000, Message: err.Error()}
				break
			}
			resp.Result = result

		case MethodToolsList:
			result, err := impl.HandleListTools()
			if err != nil {
				resp.Error = &Error{Code: -32000, Message: err.Error()}
				break
			}
			resp.Result = result

		case MethodToolsCall:
			var params CallToolParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				resp.Error = &Error{Code: -32602, Message: "Invalid params"}
				break
			}
			result, err := impl.HandleCallTool(params)
			if err != nil {
				resp.Error = &Error{Code: -32000, Message: err.Error()}
				break
			}
			resp.Result = result

		default:
			if req.ID == nil {
				continue
			}
			resp.Error = &Error{Code: -32601, Message: "Method not found"}
		}

		if req.ID == nil {
			continue
		}

		if progressToken != nil {
			_ = emitNotification(NewProgressNotification(ProgressNotificationParams{
				ProgressToken: progressToken,
				Progress:      1,
				Total:         1,
				Message:       "completed",
			}))
		}

		bytes, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
			continue
		}
		fmt.Println(string(bytes))
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func emitNotification(msg interface{}) error {
	bytes, err := MarshalDTO(msg)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}

func extractProgressToken(params json.RawMessage) interface{} {
	if len(params) == 0 {
		return nil
	}

	var wrapper struct {
		Meta *RequestMeta `json:"_meta,omitempty"`
	}
	if err := json.Unmarshal(params, &wrapper); err != nil {
		return nil
	}
	if wrapper.Meta == nil {
		return nil
	}
	return wrapper.Meta.ProgressToken
}
