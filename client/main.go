package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mcp_practice/mcp-core"
	"os"
	"os/exec"
)

func main() {
	// コマンドライン引数の解析
	toolName := flag.String("tool", "", "Tool name to call")
	toolArgs := flag.String("args", "{}", "Tool arguments in JSON")
	flag.Parse()

	// 1. MCPサーバーの起動 (stdio経由で通信)
	serverPath := "./mcp-date/date-server"
	if _, err := os.Stat(serverPath); os.IsNotExist(err) {
		serverPath = "../mcp-date/date-server"
	}

	// デモモードの場合のみログを出す
	if *toolName == "" {
		fmt.Printf("[*] MCPサーバーを起動中: %s\n", serverPath)
	}

	cmd := exec.Command(serverPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "StdinPipe error: %v\n", err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "StdoutPipe error: %v\n", err)
		return
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Start error: %v\n", err)
		return
	}
	defer func() {
		if *toolName == "" {
			fmt.Println("\n[*] プロセスを終了します。")
		}
		cmd.Process.Kill()
	}()

	reader := bufio.NewReader(stdout)

	// ヘルパー: メッセージ送信
	sendMessage := func(msg interface{}) error {
		sendData, err := mcpcore.MarshalDTO(msg)
		if err != nil {
			return fmt.Errorf("marshal error: %w", err)
		}
		if *toolName == "" {
			fmt.Printf("[送信 DTO] %+v\n", msg)
			fmt.Printf("[送信 JSON] %s\n", string(sendData))
		}
		if _, err := stdin.Write(append(sendData, '\n')); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		return nil
	}

	// ヘルパー: メッセージ受信
	receiveMessage := func() *mcpcore.JSONRPCResponse {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
			}
			return nil
		}
		msg, err := mcpcore.DecodeJSONRPCResponse([]byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Decode response error: %v\n", err)
			return nil
		}
		if *toolName == "" {
			if msg.Method != "" && msg.ID == nil {
				var params interface{}
				if len(msg.Params) > 0 {
					_ = mcpcore.UnmarshalDTO(msg.Params, &params)
				}
				fmt.Printf("[受信 Notification] {JSONRPC:%s Method:%s Params:%+v}\n", msg.JSONRPC, msg.Method, params)
				return &msg
			}

			var result interface{}
			if len(msg.Result) > 0 {
				_ = mcpcore.UnmarshalDTO(msg.Result, &result)
			}
			fmt.Printf("[受信 DTO] {JSONRPC:%s ID:%v Result:%+v Error:%+v}\n", msg.JSONRPC, msg.ID, result, msg.Error)
		}
		return &msg
	}

	waitForID := func(targetID string) *mcpcore.JSONRPCResponse {
		for {
			resp := receiveMessage()
			if resp == nil {
				return nil
			}
			if resp.ID == targetID {
				return resp
			}
		}
	}

	progressMeta := func(token string) map[string]interface{} {
		return map[string]interface{}{
			"progressToken": token,
		}
	}

	sendAndWait := func(id string, req interface{}, opName string) *mcpcore.JSONRPCResponse {
		if err := sendMessage(req); err != nil {
			fmt.Fprintf(os.Stderr, "Send %s error: %v\n", opName, err)
			return nil
		}
		resp := waitForID(id)
		if resp == nil {
			fmt.Fprintf(os.Stderr, "Receive %s response error: no response\n", opName)
			return nil
		}
		return resp
	}

	decodeResult := func(resp *mcpcore.JSONRPCResponse, dst interface{}, opName string) bool {
		if err := mcpcore.UnmarshalDTO(resp.Result, dst); err != nil {
			fmt.Fprintf(os.Stderr, "%s parse error: %v\n", opName, err)
			return false
		}
		return true
	}

	// --- ステップ1: initialize (ハンドシェイク) ---
	if *toolName == "" {
		fmt.Println("\n--- ステップ1: initialize ---")
	}
	initID := "init-go-1"

	// DTOを使用してリクエストを作成
	initReq := mcpcore.NewInitializeRequest(initID, "scratch-client-go", "1.0.0")
	if sendAndWait(initID, initReq, "initialize") == nil {
		return
	}

	// --- ステップ2: initialized (通知) ---
	if *toolName == "" {
		fmt.Println("\n--- ステップ2: initialized ---")
	}

	// DTOを使用して通知を作成
	initializedNotify := mcpcore.NewInitializedNotification()
	if err := sendMessage(initializedNotify); err != nil {
		fmt.Fprintf(os.Stderr, "Send initialized notification error: %v\n", err)
		return
	}

	// 引数が指定されている場合は、特定のツールを実行して終了
	if *toolName != "" {
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(*toolArgs), &args); err != nil {
			fmt.Fprintf(os.Stderr, "Args parse error: %v\n", err)
			return
		}

		callID := "oneshot-1"
		callReq := mcpcore.NewCallToolRequest(
			callID,
			*toolName,
			args,
		)
		callReq.Params.Meta = progressMeta("progress-oneshot-call")
		resp := sendAndWait(callID, callReq, "tools/call")
		if resp != nil {
			// 結果を出力 (Geminiがパースしやすい形式で)
			fmt.Println("MCP_RESULT_START")
			fmt.Println(string(resp.Result))
			fmt.Println("MCP_RESULT_END")
		}
		return
	}

	// --- ステップ3: prompts/list ---
	fmt.Println("\n--- ステップ3: prompts/list ---")
	listPromptsID := "prompts-list-1"
	listPromptsResp := sendAndWait(listPromptsID, mcpcore.NewListPromptsRequest(listPromptsID), "prompts/list")
	if listPromptsResp == nil {
		return
	}
	var promptsResult mcpcore.ListPromptsResult
	if !decodeResult(listPromptsResp, &promptsResult, "prompts/list") {
		return
	}
	fmt.Print("[*] 利用可能なプロンプト: ")
	for i, p := range promptsResult.Prompts {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(p.Name)
	}
	fmt.Println()

	// --- ステップ4: prompts/get ---
	if len(promptsResult.Prompts) > 0 {
		targetPrompt := promptsResult.Prompts[0].Name
		fmt.Printf("\n--- ステップ4: プロンプト取得 (%s) ---\n", targetPrompt)
		getPromptID := "prompts-get-1"
		getPromptReq := mcpcore.NewGetPromptRequest(getPromptID, targetPrompt, map[string]string{"city": "Tokyo"})
		getPromptResp := sendAndWait(getPromptID, getPromptReq, "prompts/get")
		if getPromptResp == nil {
			return
		}
		fmt.Printf("[*] prompts/get 結果:\n%s\n", string(getPromptResp.Result))
	}

	// --- ステップ5: resources/list ---
	fmt.Println("\n--- ステップ5: resources/list ---")
	listResourcesID := "resources-list-1"
	listResourcesResp := sendAndWait(listResourcesID, mcpcore.NewListResourcesRequest(listResourcesID), "resources/list")
	if listResourcesResp == nil {
		return
	}
	var resourcesResult mcpcore.ListResourcesResult
	if !decodeResult(listResourcesResp, &resourcesResult, "resources/list") {
		return
	}
	fmt.Print("[*] 利用可能なリソース: ")
	for i, r := range resourcesResult.Resources {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(r.URI)
	}
	fmt.Println()

	// --- ステップ6: resources/templates/list ---
	fmt.Println("\n--- ステップ6: resources/templates/list ---")
	listTemplatesID := "resources-templates-list-1"
	listTemplatesResp := sendAndWait(listTemplatesID, mcpcore.NewListResourceTemplatesRequest(listTemplatesID), "resources/templates/list")
	if listTemplatesResp == nil {
		return
	}
	var templatesResult mcpcore.ListResourceTemplatesResult
	if !decodeResult(listTemplatesResp, &templatesResult, "resources/templates/list") {
		return
	}
	fmt.Print("[*] 利用可能なリソーステンプレート: ")
	for i, t := range templatesResult.ResourceTemplates {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(t.URITemplate)
	}
	fmt.Println()

	// --- ステップ7: resources/read ---
	if len(resourcesResult.Resources) > 0 {
		targetURI := resourcesResult.Resources[0].URI
		fmt.Printf("\n--- ステップ7: リソース取得 (%s) ---\n", targetURI)
		readID := "resources-read-1"
		readReq := mcpcore.NewReadResourceRequest(readID, targetURI)
		readReq.Params.Meta = progressMeta("progress-resources-read-1")
		readResp := sendAndWait(readID, readReq, "resources/read")
		if readResp == nil {
			return
		}
		fmt.Printf("[*] resources/read 結果:\n%s\n", string(readResp.Result))
	}

	// --- ステップ8: completion/complete ---
	fmt.Println("\n--- ステップ8: completion/complete ---")
	completeID := "completion-1"
	ref := mcpcore.PromptReference{
		Type: "ref/prompt",
		Name: "city_time_prompt",
	}
	completeReq := mcpcore.NewCompleteRequest(
		completeID,
		ref,
		mcpcore.CompletionArgument{Name: "city", Value: "To"},
		nil,
	)
	completeReq.Params.Meta = progressMeta("progress-completion-1")
	completeResp := sendAndWait(completeID, completeReq, "completion/complete")
	if completeResp == nil {
		return
	}
	fmt.Printf("[*] completion/complete 結果:\n%s\n", string(completeResp.Result))

	// --- ステップ9: tools/list ---
	fmt.Println("\n--- ステップ9: tools/list ---")
	listID := "tools-list-1"
	listToolsResp := sendAndWait(listID, mcpcore.NewListToolsRequest(listID), "tools/list")
	if listToolsResp == nil {
		return
	}
	var toolsResult mcpcore.ListToolsResult
	if !decodeResult(listToolsResp, &toolsResult, "tools/list") {
		return
	}
	fmt.Print("[*] 利用可能なツール: ")
	for i, t := range toolsResult.Tools {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(t.Name)
	}
	fmt.Println()

	// --- ステップ10: tools/call ---
	if len(toolsResult.Tools) > 0 {
		targetTool := "get_city_date"
		fmt.Printf("\n--- ステップ10: ツールの実行 (%s) ---\n", targetTool)
		callID := "call-go-1"
		callReq := mcpcore.NewCallToolRequest(
			callID,
			targetTool,
			map[string]interface{}{"city": "Tokyo"},
		)
		callReq.Params.Meta = progressMeta("progress-tools-call-1")
		callResp := sendAndWait(callID, callReq, "tools/call")
		if callResp == nil {
			return
		}
		fmt.Printf("[*] 実行結果:\n%s\n", string(callResp.Result))
	}
}
