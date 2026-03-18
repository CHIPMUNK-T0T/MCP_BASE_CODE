package mcpcore

import (
	"encoding/json"
)

// --- JSON-RPC 2.0 Base Structures ---

type Request struct {
	JsonRpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
}

type Response struct {
	JsonRpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// --- MCP Protocol Objects ---

type EmptyObject struct{}

type Implementation struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	WebsiteURL  string `json:"websiteUrl,omitempty"`
	Icons       []Icon `json:"icons,omitempty"`
}

type Icon struct {
	Src      string   `json:"src"`
	MimeType string   `json:"mimeType,omitempty"`
	Sizes    []string `json:"sizes,omitempty"`
	Theme    string   `json:"theme,omitempty"`
}

type RootsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type SamplingCapability struct {
	Context *EmptyObject `json:"context,omitempty"`
	Tools   *EmptyObject `json:"tools,omitempty"`
}

type ClientCapabilities struct {
	Experimental map[string]map[string]interface{} `json:"experimental,omitempty"`
	Roots        *RootsCapability                  `json:"roots,omitempty"`
	Sampling     *SamplingCapability               `json:"sampling,omitempty"`
	Elicitation  *EmptyObject                      `json:"elicitation,omitempty"`
	Tools        *EmptyObject                      `json:"tools,omitempty"`
}

type ToolsServerCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type PromptsServerCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ResourcesServerCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

type ServerCapabilities struct {
	Experimental map[string]map[string]interface{} `json:"experimental,omitempty"`
	Completions  *EmptyObject                      `json:"completions,omitempty"`
	Logging      *EmptyObject                      `json:"logging,omitempty"`
	Prompts      *PromptsServerCapability          `json:"prompts,omitempty"`
	Resources    *ResourcesServerCapability        `json:"resources,omitempty"`
	Tools        *ToolsServerCapability            `json:"tools,omitempty"`
}

type RequestMeta struct {
	ProgressToken interface{} `json:"progressToken,omitempty"`
}

type ProgressNotificationParams struct {
	ProgressToken interface{} `json:"progressToken"`
	Progress      float64     `json:"progress"`
	Total         float64     `json:"total,omitempty"`
	Message       string      `json:"message,omitempty"`
}

type InitializeParams struct {
	Meta            map[string]interface{} `json:"_meta,omitempty"`
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ClientCapabilities     `json:"capabilities"`
	ClientInfo      Implementation         `json:"clientInfo"`
}

type InitializeResult struct {
	Meta            map[string]interface{} `json:"_meta,omitempty"`
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ServerCapabilities     `json:"capabilities"`
	ServerInfo      Implementation         `json:"serverInfo"`
	Instructions    string                 `json:"instructions,omitempty"`
}

type ToolAnnotations struct {
	Title           string   `json:"title,omitempty"`
	ReadOnlyHint    bool     `json:"readOnlyHint,omitempty"`
	DestructiveHint bool     `json:"destructiveHint,omitempty"`
	IdempotentHint  bool     `json:"idempotentHint,omitempty"`
	OpenWorldHint   bool     `json:"openWorldHint,omitempty"`
	Audience        []string `json:"audience,omitempty"`
}

type Tool struct {
	Meta         map[string]interface{} `json:"_meta,omitempty"`
	Name         string                 `json:"name"`
	Title        string                 `json:"title,omitempty"`
	Description  string                 `json:"description,omitempty"`
	InputSchema  map[string]interface{} `json:"inputSchema"`
	OutputSchema map[string]interface{} `json:"outputSchema,omitempty"`
	Annotations  *ToolAnnotations       `json:"annotations,omitempty"`
	Icons        []Icon                 `json:"icons,omitempty"`
}

type ListToolsResult struct {
	Meta       map[string]interface{} `json:"_meta,omitempty"`
	Tools      []Tool                 `json:"tools"`
	NextCursor string                 `json:"nextCursor,omitempty"`
}

type TaskContext struct {
	ID string `json:"id,omitempty"`
}

type CallToolParams struct {
	Meta      map[string]interface{} `json:"_meta,omitempty"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Task      *TaskContext           `json:"task,omitempty"`
}

type CallToolResult struct {
	Meta              map[string]interface{} `json:"_meta,omitempty"`
	Content           []Content              `json:"content"`
	StructuredContent map[string]interface{} `json:"structuredContent,omitempty"`
	IsError           bool                   `json:"isError,omitempty"`
}

type Annotations struct {
	Audience     []string `json:"audience,omitempty"`
	Priority     float64  `json:"priority,omitempty"`
	LastModified string   `json:"lastModified,omitempty"`
}

type Content struct {
	Meta        map[string]interface{} `json:"_meta,omitempty"`
	Type        string                 `json:"type"`
	Text        string                 `json:"text,omitempty"`
	Data        string                 `json:"data,omitempty"`
	MimeType    string                 `json:"mimeType,omitempty"`
	URI         string                 `json:"uri,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Size        int64                  `json:"size,omitempty"`
	Resource    json.RawMessage        `json:"resource,omitempty"`
	Annotations *Annotations           `json:"annotations,omitempty"`
}

func NewTextContent(text string) Content {
	return Content{
		Type: "text",
		Text: text,
	}
}

type PromptArgument struct {
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

type Prompt struct {
	Meta        map[string]interface{} `json:"_meta,omitempty"`
	Icons       []Icon                 `json:"icons,omitempty"`
	Name        string                 `json:"name"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Arguments   []PromptArgument       `json:"arguments,omitempty"`
}

type PromptMessage struct {
	Role    string  `json:"role"`
	Content Content `json:"content"`
}

type ListPromptsResult struct {
	Meta       map[string]interface{} `json:"_meta,omitempty"`
	NextCursor string                 `json:"nextCursor,omitempty"`
	Prompts    []Prompt               `json:"prompts"`
}

type GetPromptParams struct {
	Meta      map[string]interface{} `json:"_meta,omitempty"`
	Name      string                 `json:"name"`
	Arguments map[string]string      `json:"arguments,omitempty"`
}

type GetPromptResult struct {
	Meta        map[string]interface{} `json:"_meta,omitempty"`
	Description string                 `json:"description,omitempty"`
	Messages    []PromptMessage        `json:"messages"`
}

type Resource struct {
	Meta        map[string]interface{} `json:"_meta,omitempty"`
	Icons       []Icon                 `json:"icons,omitempty"`
	Name        string                 `json:"name"`
	Title       string                 `json:"title,omitempty"`
	URI         string                 `json:"uri"`
	Description string                 `json:"description,omitempty"`
	MimeType    string                 `json:"mimeType,omitempty"`
	Annotations *Annotations           `json:"annotations,omitempty"`
	Size        int64                  `json:"size,omitempty"`
}

type ResourceTemplate struct {
	Meta        map[string]interface{} `json:"_meta,omitempty"`
	Icons       []Icon                 `json:"icons,omitempty"`
	Name        string                 `json:"name"`
	Title       string                 `json:"title,omitempty"`
	URITemplate string                 `json:"uriTemplate"`
	Description string                 `json:"description,omitempty"`
	MimeType    string                 `json:"mimeType,omitempty"`
	Annotations *Annotations           `json:"annotations,omitempty"`
}

type ListResourcesResult struct {
	Meta       map[string]interface{} `json:"_meta,omitempty"`
	NextCursor string                 `json:"nextCursor,omitempty"`
	Resources  []Resource             `json:"resources"`
}

type ListResourceTemplatesResult struct {
	Meta              map[string]interface{} `json:"_meta,omitempty"`
	NextCursor        string                 `json:"nextCursor,omitempty"`
	ResourceTemplates []ResourceTemplate     `json:"resourceTemplates"`
}

type ReadResourceParams struct {
	Meta map[string]interface{} `json:"_meta,omitempty"`
	URI  string                 `json:"uri"`
}

type ReadResourceResult struct {
	Meta     map[string]interface{} `json:"_meta,omitempty"`
	Contents []ResourceContent      `json:"contents"`
}

type ResourceContent struct {
	Meta     map[string]interface{} `json:"_meta,omitempty"`
	URI      string                 `json:"uri"`
	MimeType string                 `json:"mimeType,omitempty"`
	Text     string                 `json:"text,omitempty"`
	Blob     string                 `json:"blob,omitempty"`
}

func NewTextResourceContent(uri, mimeType, text string) ResourceContent {
	return ResourceContent{
		URI:      uri,
		MimeType: mimeType,
		Text:     text,
	}
}

type PromptReference struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Title string `json:"title,omitempty"`
}

type ResourceTemplateReference struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type CompletionArgument struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CompletionContext struct {
	Arguments map[string]string `json:"arguments,omitempty"`
}

type CompleteParams struct {
	Meta     map[string]interface{} `json:"_meta,omitempty"`
	Ref      interface{}            `json:"ref"`
	Argument CompletionArgument     `json:"argument"`
	Context  *CompletionContext     `json:"context,omitempty"`
}

type CompleteResult struct {
	Meta       map[string]interface{} `json:"_meta,omitempty"`
	Completion CompletionValues       `json:"completion"`
}

type CompletionValues struct {
	Values  []string `json:"values"`
	Total   int      `json:"total,omitempty"`
	HasMore bool     `json:"hasMore,omitempty"`
}

// --- Abstract Interfaces ---

// MCPServer はMCPサーバーが実装すべき抽象関数を定義します
type MCPServer interface {
	HandleInitialize(params InitializeParams) (InitializeResult, error)
	HandleListPrompts() (ListPromptsResult, error)
	HandleGetPrompt(params GetPromptParams) (GetPromptResult, error)
	HandleListResources() (ListResourcesResult, error)
	HandleListResourceTemplates() (ListResourceTemplatesResult, error)
	HandleReadResource(params ReadResourceParams) (ReadResourceResult, error)
	HandleComplete(params CompleteParams) (CompleteResult, error)
	HandleListTools() (ListToolsResult, error)
	HandleCallTool(params CallToolParams) (CallToolResult, error)
}

// MCPClient はMCPクライアントが提供すべき抽象関数を定義します
type MCPClient interface {
	Initialize() (InitializeResult, error)
	ListPrompts() (ListPromptsResult, error)
	GetPrompt(name string, args map[string]string) (GetPromptResult, error)
	ListResources() (ListResourcesResult, error)
	ListResourceTemplates() (ListResourceTemplatesResult, error)
	ReadResource(uri string) (ReadResourceResult, error)
	Complete(ref interface{}, argument CompletionArgument, context *CompletionContext) (CompleteResult, error)
	ListTools() (ListToolsResult, error)
	CallTool(name string, args map[string]interface{}) (CallToolResult, error)
	Close() error
}
