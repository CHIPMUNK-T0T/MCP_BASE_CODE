package main

import (
	"fmt"
	"mcp_practice/mcp-core"
	"os"
	"sort"
	"strings"
	"time"
)

type DateServer struct{}

var cityTimezones = map[string]string{
	"new york":    "America/New_York",
	"london":      "Europe/London",
	"paris":       "Europe/Paris",
	"tokyo":       "Asia/Tokyo",
	"sydney":      "Australia/Sydney",
	"los angeles": "America/Los_Angeles",
	"dubai":       "Asia/Dubai",
	"seoul":       "Asia/Seoul",
}

func (DateServer) HandleInitialize(params mcpcore.InitializeParams) (mcpcore.InitializeResult, error) {
	return mcpcore.InitializeResult{
		ProtocolVersion: mcpcore.MCPProtocolVersion,
		Capabilities: mcpcore.ServerCapabilities{
			Tools:       &mcpcore.ToolsServerCapability{ListChanged: true},
			Prompts:     &mcpcore.PromptsServerCapability{ListChanged: true},
			Resources:   &mcpcore.ResourcesServerCapability{Subscribe: true, ListChanged: true},
			Completions: &mcpcore.EmptyObject{},
		},
		ServerInfo: mcpcore.Implementation{
			Name:    "date-server",
			Version: "0.1.0",
		},
	}, nil
}

func (DateServer) HandleListPrompts() (mcpcore.ListPromptsResult, error) {
	return mcpcore.ListPromptsResult{
		Prompts: []mcpcore.Prompt{
			{
				Name:        "city_time_prompt",
				Title:       "City Time Prompt",
				Description: "Builds a user prompt to ask the current time in a city.",
				Arguments: []mcpcore.PromptArgument{
					{
						Name:        "city",
						Description: "City name (e.g., Tokyo, Paris)",
						Required:    true,
					},
				},
			},
		},
	}, nil
}

func (DateServer) HandleGetPrompt(params mcpcore.GetPromptParams) (mcpcore.GetPromptResult, error) {
	if params.Name != "city_time_prompt" {
		return mcpcore.GetPromptResult{}, fmt.Errorf("prompt not found: %s", params.Name)
	}

	city := params.Arguments["city"]
	if city == "" {
		return mcpcore.GetPromptResult{}, fmt.Errorf("city argument is required")
	}

	return mcpcore.GetPromptResult{
		Description: "Prompt for asking city current time.",
		Messages: []mcpcore.PromptMessage{
			{
				Role:    "user",
				Content: mcpcore.NewTextContent(fmt.Sprintf("What is the current local time in %s?", city)),
			},
		},
	}, nil
}

func (DateServer) HandleListResources() (mcpcore.ListResourcesResult, error) {
	return mcpcore.ListResourcesResult{
		Resources: []mcpcore.Resource{
			{
				Name:        "supported-cities",
				Title:       "Supported Cities",
				URI:         "resource://cities",
				Description: "List of cities supported by get_city_date",
				MimeType:    "text/plain",
			},
		},
	}, nil
}

func (DateServer) HandleListResourceTemplates() (mcpcore.ListResourceTemplatesResult, error) {
	return mcpcore.ListResourceTemplatesResult{
		ResourceTemplates: []mcpcore.ResourceTemplate{
			{
				Name:        "city-timezone-template",
				Title:       "City Timezone",
				URITemplate: "resource://city/{city}",
				Description: "Timezone information for a city",
				MimeType:    "text/plain",
			},
		},
	}, nil
}

func (DateServer) HandleReadResource(params mcpcore.ReadResourceParams) (mcpcore.ReadResourceResult, error) {
	switch {
	case params.URI == "resource://cities":
		cities := sortedCityNames()
		return mcpcore.ReadResourceResult{
			Contents: []mcpcore.ResourceContent{
				mcpcore.NewTextResourceContent(params.URI, "text/plain", strings.Join(cities, "\n")),
			},
		}, nil

	case strings.HasPrefix(params.URI, "resource://city/"):
		city := strings.TrimPrefix(params.URI, "resource://city/")
		tzName, ok := cityTimezones[strings.ToLower(city)]
		if !ok {
			return mcpcore.ReadResourceResult{}, fmt.Errorf("city not found: %s", city)
		}
		text := fmt.Sprintf("city=%s\ntimezone=%s", city, tzName)
		return mcpcore.ReadResourceResult{
			Contents: []mcpcore.ResourceContent{
				mcpcore.NewTextResourceContent(params.URI, "text/plain", text),
			},
		}, nil
	}

	return mcpcore.ReadResourceResult{}, fmt.Errorf("resource not found: %s", params.URI)
}

func (DateServer) HandleComplete(params mcpcore.CompleteParams) (mcpcore.CompleteResult, error) {
	refType, err := completionRefType(params.Ref)
	if err != nil {
		return mcpcore.CompleteResult{}, err
	}
	if refType != "ref/prompt" && refType != "ref/resource" {
		return mcpcore.CompleteResult{}, fmt.Errorf("unsupported ref type: %s", refType)
	}

	prefix := strings.ToLower(strings.TrimSpace(params.Argument.Value))
	values := make([]string, 0)
	for _, city := range sortedCityNames() {
		if strings.HasPrefix(strings.ToLower(city), prefix) {
			values = append(values, city)
		}
	}

	return mcpcore.CompleteResult{
		Completion: mcpcore.CompletionValues{
			Values: values,
			Total:  len(values),
		},
	}, nil
}

func (DateServer) HandleListTools() (mcpcore.ListToolsResult, error) {
	return mcpcore.ListToolsResult{
		Tools: []mcpcore.Tool{
			{
				Name:        "get_city_date",
				Description: "Get the current date and time for a specific city",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"city": map[string]interface{}{
							"type":        "string",
							"description": "The name of the city (e.g., 'New York', 'Tokyo')",
						},
					},
					"required": []string{"city"},
				},
			},
		},
	}, nil
}

func (DateServer) HandleCallTool(params mcpcore.CallToolParams) (mcpcore.CallToolResult, error) {
	if params.Name != "get_city_date" {
		return mcpcore.CallToolResult{
			Content: []mcpcore.Content{mcpcore.NewTextContent("Error: tool not found")},
			IsError: true,
		}, nil
	}

	city, ok := params.Arguments["city"].(string)
	if !ok {
		return mcpcore.CallToolResult{
			Content: []mcpcore.Content{mcpcore.NewTextContent("Error: city argument required")},
			IsError: true,
		}, nil
	}

	return mcpcore.CallToolResult{
		Content: []mcpcore.Content{mcpcore.NewTextContent(getCityDate(city))},
	}, nil
}

func main() {
	if err := mcpcore.ServeStdio(DateServer{}); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}

func getCityDate(city string) string {
	tzName, found := cityTimezones[strings.ToLower(city)]
	if !found {
		return fmt.Sprintf("Timezone not found for city: %s. UTC: %s", city, time.Now().UTC().Format(time.RFC1123))
	}

	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return fmt.Sprintf("Error loading timezone %s: %v", tzName, err)
	}

	return fmt.Sprintf("Time in %s: %s", city, time.Now().In(loc).Format(time.RFC1123))
}

func sortedCityNames() []string {
	cities := make([]string, 0, len(cityTimezones))
	for city := range cityTimezones {
		cities = append(cities, displayCityName(city))
	}
	sort.Strings(cities)
	return cities
}

func displayCityName(city string) string {
	parts := strings.Split(city, " ")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

func completionRefType(ref interface{}) (string, error) {
	switch v := ref.(type) {
	case mcpcore.PromptReference:
		return v.Type, nil
	case mcpcore.ResourceTemplateReference:
		return v.Type, nil
	}

	refMap, ok := ref.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid completion ref")
	}
	refType, ok := refMap["type"].(string)
	if !ok || refType == "" {
		return "", fmt.Errorf("completion ref type is required")
	}
	return refType, nil
}
