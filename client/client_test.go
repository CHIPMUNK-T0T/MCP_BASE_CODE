package main

import (
	"fmt"
	"mcp_practice/mcp-core"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var serverPath string

func TestMain(m *testing.M) {
	// 1. テスト前にサーバーをビルド
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get working directory: %v\n", err)
		os.Exit(1)
	}

	// client/ から実行されるので、一つ上へ
	projectRoot := filepath.Dir(wd)
	serverSrc := filepath.Join(projectRoot, "mcp-date", "main.go")
	serverBin := filepath.Join(projectRoot, "mcp-date", "date-server-test-bin")

	// go build
	cmd := exec.Command("go", "build", "-o", serverBin, serverSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to build test server: %v\nOutput: %s\n", err, string(out))
		os.Exit(1)
	}

	serverPath = serverBin

	// 2. テスト実行
	code := m.Run()

	// 3. クリーンアップ
	os.Remove(serverBin)

	os.Exit(code)
}

func TestClientIntegration(t *testing.T) {
	// サーバープロセスを起動してクライアントを作成
	client, err := mcpcore.NewStdioClient(serverPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// 1. Initialize
	initRes, err := client.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	if initRes.ServerInfo.Name != "date-server" {
		t.Errorf("Expected server name 'date-server', got %s", initRes.ServerInfo.Name)
	}

	// 2. List Tools
	promptsRes, err := client.ListPrompts()
	if err != nil {
		t.Fatalf("ListPrompts failed: %v", err)
	}
	if len(promptsRes.Prompts) == 0 {
		t.Error("Expected prompts list, got empty")
	}

	getPromptRes, err := client.GetPrompt("city_time_prompt", map[string]string{"city": "Tokyo"})
	if err != nil {
		t.Fatalf("GetPrompt failed: %v", err)
	}
	if len(getPromptRes.Messages) == 0 {
		t.Error("Expected prompt messages, got empty")
	}

	resourcesRes, err := client.ListResources()
	if err != nil {
		t.Fatalf("ListResources failed: %v", err)
	}
	if len(resourcesRes.Resources) == 0 {
		t.Error("Expected resources list, got empty")
	}

	resourceReadRes, err := client.ReadResource("resource://cities")
	if err != nil {
		t.Fatalf("ReadResource failed: %v", err)
	}
	if len(resourceReadRes.Contents) == 0 {
		t.Error("Expected resource contents, got empty")
	}

	templatesRes, err := client.ListResourceTemplates()
	if err != nil {
		t.Fatalf("ListResourceTemplates failed: %v", err)
	}
	if len(templatesRes.ResourceTemplates) == 0 {
		t.Error("Expected resource template list, got empty")
	}

	completeRes, err := client.Complete(
		mcpcore.PromptReference{Type: "ref/prompt", Name: "city_time_prompt"},
		mcpcore.CompletionArgument{Name: "city", Value: "To"},
		nil,
	)
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}
	if len(completeRes.Completion.Values) == 0 {
		t.Error("Expected completion values, got empty")
	}

	toolsRes, err := client.ListTools()
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
	if len(toolsRes.Tools) == 0 {
		t.Error("Expected tools list, got empty")
	}

	// 3. Call Tool
	args := map[string]interface{}{"city": "Paris"}
	callRes, err := client.CallTool("get_city_date", args)
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	if len(callRes.Content) == 0 {
		t.Error("Content is empty")
	} else {
		text := callRes.Content[0].Text
		if !strings.Contains(text, "Time in Paris") {
			t.Errorf("Unexpected tool output: %s", text)
		}
	}
}
