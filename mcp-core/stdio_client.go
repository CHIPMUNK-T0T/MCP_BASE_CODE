package mcpcore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

// StdioClient implements MCPClient over a subprocess using stdin/stdout.
type StdioClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
	id     int
}

var _ MCPClient = (*StdioClient)(nil)

func NewStdioClient(command string, args ...string) (*StdioClient, error) {
	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &StdioClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewScanner(stdout),
		id:     1,
	}, nil
}

func (c *StdioClient) call(method string, params interface{}) (JSONRPCResponse, error) {
	var paramsRaw json.RawMessage
	if params != nil {
		b, _ := json.Marshal(params)
		paramsRaw = b
	}

	reqID := c.id
	req := Request{
		JsonRpc: JSONRPCVersion,
		Method:  method,
		Params:  paramsRaw,
		ID:      reqID,
	}
	c.id++

	reqBytes, _ := json.Marshal(req)
	if _, err := fmt.Fprintln(c.stdin, string(reqBytes)); err != nil {
		return JSONRPCResponse{}, err
	}

	for c.stdout.Scan() {
		line := c.stdout.Bytes()
		resp, err := DecodeJSONRPCResponse(line)
		if err != nil {
			return JSONRPCResponse{}, err
		}
		// ignore server notifications
		if resp.Method != "" && resp.ID == nil {
			continue
		}
		if fmt.Sprint(resp.ID) != fmt.Sprint(reqID) {
			continue
		}
		if resp.Error != nil {
			return JSONRPCResponse{}, fmt.Errorf("json-rpc error (%d): %s", resp.Error.Code, resp.Error.Message)
		}
		return resp, nil
	}
	return JSONRPCResponse{}, fmt.Errorf("failed to read response for id=%v", reqID)
}

func (c *StdioClient) notify(method string, params interface{}) error {
	var paramsRaw json.RawMessage
	if params != nil {
		b, _ := json.Marshal(params)
		paramsRaw = b
	}

	req := Request{
		JsonRpc: JSONRPCVersion,
		Method:  method,
		Params:  paramsRaw,
	}
	reqBytes, _ := json.Marshal(req)
	_, err := fmt.Fprintln(c.stdin, string(reqBytes))
	return err
}

func (c *StdioClient) Initialize() (InitializeResult, error) {
	params := InitializeParams{
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
			Name:    "demo-client",
			Version: "1.0.0",
		},
	}
	resp, err := c.call(MethodInitialize, params)
	if err != nil {
		return InitializeResult{}, err
	}
	result, err := DecodeJSONRPCResult[InitializeResult](resp)
	if err != nil {
		return InitializeResult{}, err
	}
	if err := c.notify(MethodInitialized, map[string]interface{}{}); err != nil {
		return InitializeResult{}, err
	}
	return result, nil
}

func (c *StdioClient) ListPrompts() (ListPromptsResult, error) {
	resp, err := c.call(MethodPromptsList, ListPromptsParams{})
	if err != nil {
		return ListPromptsResult{}, err
	}
	return DecodeJSONRPCResult[ListPromptsResult](resp)
}

func (c *StdioClient) GetPrompt(name string, args map[string]string) (GetPromptResult, error) {
	params := GetPromptParams{
		Name:      name,
		Arguments: args,
	}
	resp, err := c.call(MethodPromptsGet, params)
	if err != nil {
		return GetPromptResult{}, err
	}
	return DecodeJSONRPCResult[GetPromptResult](resp)
}

func (c *StdioClient) ListResources() (ListResourcesResult, error) {
	resp, err := c.call(MethodResourcesList, ListResourcesParams{})
	if err != nil {
		return ListResourcesResult{}, err
	}
	return DecodeJSONRPCResult[ListResourcesResult](resp)
}

func (c *StdioClient) ListResourceTemplates() (ListResourceTemplatesResult, error) {
	resp, err := c.call(MethodResourcesTemplatesList, ListResourceTemplatesParams{})
	if err != nil {
		return ListResourceTemplatesResult{}, err
	}
	return DecodeJSONRPCResult[ListResourceTemplatesResult](resp)
}

func (c *StdioClient) ReadResource(uri string) (ReadResourceResult, error) {
	params := ReadResourceParams{URI: uri}
	resp, err := c.call(MethodResourcesRead, params)
	if err != nil {
		return ReadResourceResult{}, err
	}
	return DecodeJSONRPCResult[ReadResourceResult](resp)
}

func (c *StdioClient) Complete(ref interface{}, argument CompletionArgument, context *CompletionContext) (CompleteResult, error) {
	params := CompleteParams{
		Ref:      ref,
		Argument: argument,
		Context:  context,
	}
	resp, err := c.call(MethodCompletionComplete, params)
	if err != nil {
		return CompleteResult{}, err
	}
	return DecodeJSONRPCResult[CompleteResult](resp)
}

func (c *StdioClient) ListTools() (ListToolsResult, error) {
	resp, err := c.call(MethodToolsList, nil)
	if err != nil {
		return ListToolsResult{}, err
	}
	return DecodeJSONRPCResult[ListToolsResult](resp)
}

func (c *StdioClient) CallTool(name string, args map[string]interface{}) (CallToolResult, error) {
	params := CallToolParams{
		Name:      name,
		Arguments: args,
	}
	resp, err := c.call(MethodToolsCall, params)
	if err != nil {
		return CallToolResult{}, err
	}
	return DecodeJSONRPCResult[CallToolResult](resp)
}

func (c *StdioClient) Close() error {
	c.stdin.Close()
	return c.cmd.Wait()
}
