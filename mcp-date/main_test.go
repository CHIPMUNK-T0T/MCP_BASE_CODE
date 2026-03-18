package main

import (
	"mcp_practice/mcp-core"
	"strings"
	"testing"
)

func TestGetCityDate(t *testing.T) {
	// 既知の都市
	tokyo := getCityDate("Tokyo")
	if !strings.Contains(tokyo, "Time in Tokyo") {
		t.Errorf("Expected result for Tokyo, got: %s", tokyo)
	}

	// 未知の都市
	unknown := getCityDate("Atlantis")
	if !strings.Contains(unknown, "Timezone not found") {
		t.Errorf("Expected timezone error for Atlantis, got: %s", unknown)
	}
}

func TestHandleRequest_Initialize(t *testing.T) {
	s := DateServer{}
	result, err := s.HandleInitialize(mcpcore.InitializeParams{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.ServerInfo.Name != "date-server" {
		t.Errorf("Expected server name 'date-server', got %s", result.ServerInfo.Name)
	}
}

func TestHandleRequest_CallTool(t *testing.T) {
	s := DateServer{}
	result, err := s.HandleCallTool(mcpcore.CallToolParams{
		Name:      "get_city_date",
		Arguments: map[string]interface{}{"city": "London"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Content) == 0 {
		t.Error("Content is empty")
	}

	if !strings.Contains(result.Content[0].Text, "Time in London") {
		t.Errorf("Expected 'Time in London', got %s", result.Content[0].Text)
	}
}

func TestHandleListPromptsAndGetPrompt(t *testing.T) {
	s := DateServer{}

	listRes, err := s.HandleListPrompts()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(listRes.Prompts) == 0 {
		t.Fatal("Expected prompts, got empty")
	}

	getRes, err := s.HandleGetPrompt(mcpcore.GetPromptParams{
		Name:      "city_time_prompt",
		Arguments: map[string]string{"city": "Tokyo"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(getRes.Messages) == 0 {
		t.Fatal("Expected prompt messages, got empty")
	}
}

func TestHandleResourcesAndComplete(t *testing.T) {
	s := DateServer{}

	listRes, err := s.HandleListResources()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(listRes.Resources) == 0 {
		t.Fatal("Expected resources, got empty")
	}

	readRes, err := s.HandleReadResource(mcpcore.ReadResourceParams{URI: "resource://cities"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(readRes.Contents) == 0 {
		t.Fatal("Expected resource contents, got empty")
	}

	completeRes, err := s.HandleComplete(mcpcore.CompleteParams{
		Ref: mcpcore.PromptReference{
			Type: "ref/prompt",
			Name: "city_time_prompt",
		},
		Argument: mcpcore.CompletionArgument{
			Name:  "city",
			Value: "To",
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(completeRes.Completion.Values) == 0 {
		t.Fatal("Expected completion values, got empty")
	}
}
