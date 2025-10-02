package mcp

import (
	"encoding/json"
)

// JSON-RPC 2.0 structures
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP Protocol structures
type InitializeRequest struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ClientCapabilities     `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
	Meta            map[string]interface{} `json:"_meta,omitempty"`
}

type ClientCapabilities struct {
	Roots    *RootsCapability    `json:"roots,omitempty"`
	Sampling *SamplingCapability `json:"sampling,omitempty"`
}

type RootsCapability struct {
	ListChanged bool `json:"listChanged"`
}

type SamplingCapability struct{}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ServerCapabilities     `json:"capabilities"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
	Meta            map[string]interface{} `json:"_meta,omitempty"`
}

type ServerCapabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Tools structures
type ListToolsRequest struct {
	Meta map[string]interface{} `json:"_meta,omitempty"`
}

type ListToolsResult struct {
	Tools []Tool                 `json:"tools"`
	Meta  map[string]interface{} `json:"_meta,omitempty"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
	Meta      map[string]interface{} `json:"_meta,omitempty"`
}

type CallToolResult struct {
	Content []ToolContent          `json:"content"`
	IsError bool                   `json:"isError,omitempty"`
	Meta    map[string]interface{} `json:"_meta,omitempty"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Service-specific structures
type ServiceInfo struct {
	Name               string       `json:"name"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	Categories         []string     `json:"categories"`
	Requirements       Requirements `json:"requirements"`
	SupportedVersions  []string     `json:"supportedVersions"`
	DataStreams        []DataStream `json:"dataStreams"`
	SetupComplexity    string       `json:"setupComplexity"`
	CommonUseCases     []string     `json:"commonUseCases"`
	OfficialDocs       []string     `json:"officialDocs"`
	CommunityResources []string     `json:"communityResources"`
}

type Requirements struct {
	Elasticsearch string `json:"elasticsearch,omitempty"`
	Kibana        string `json:"kibana,omitempty"`
	Subscription  string `json:"subscription,omitempty"`
}

type DataStream struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type SetupInstructions struct {
	ServiceName       string                `json:"serviceName"`
	Platform          string                `json:"platform,omitempty"`
	Version           string                `json:"version,omitempty"`
	Prerequisites     []string              `json:"prerequisites"`
	InstallationSteps []InstallationStep    `json:"installationSteps"`
	PostInstallation  []PostInstallation    `json:"postInstallation"`
	Troubleshooting   []TroubleshootingItem `json:"troubleshooting"`
}

type InstallationStep struct {
	Step           int             `json:"step"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Commands       []string        `json:"commands,omitempty"`
	ConfigSnippets []ConfigSnippet `json:"configSnippets,omitempty"`
	Verification   string          `json:"verification,omitempty"`
}

type ConfigSnippet struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type PostInstallation struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
}

type TroubleshootingItem struct {
	Issue    string `json:"issue"`
	Solution string `json:"solution"`
}

type ValidationResult struct {
	IsValid     bool                   `json:"isValid"`
	Errors      []ValidationError      `json:"errors"`
	Warnings    []ValidationWarning    `json:"warnings"`
	Suggestions []ValidationSuggestion `json:"suggestions"`
}

type ValidationError struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Line     *int   `json:"line,omitempty"`
	Column   *int   `json:"column,omitempty"`
	Severity string `json:"severity"`
}

type ValidationWarning struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Line       *int   `json:"line,omitempty"`
	Column     *int   `json:"column,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

type ValidationSuggestion struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
	Impact     string `json:"impact"`
}

type IntegrationDetails struct {
	Name            string                  `json:"name"`
	Title           string                  `json:"title"`
	Description     string                  `json:"description"`
	Version         string                  `json:"version"`
	Categories      []string                `json:"categories"`
	DataStreams     []IntegrationDataStream `json:"dataStreams"`
	PolicyTemplates []PolicyTemplate        `json:"policyTemplates"`
	Requirements    Requirements            `json:"requirements"`
	Screenshots     []string                `json:"screenshots"`
	Icons           []string                `json:"icons"`
	Owner           Owner                   `json:"owner"`
}

type IntegrationDataStream struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Fields      []Field `json:"fields"`
}

type Field struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Example     string `json:"example,omitempty"`
}

type PolicyTemplate struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Inputs      []Input  `json:"inputs"`
	DataStreams []string `json:"dataStreams"`
	Categories  []string `json:"categories"`
}

type Input struct {
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Vars        []Variable `json:"vars,omitempty"`
}

type Variable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Multi       bool        `json:"multi,omitempty"`
	Secret      bool        `json:"secret,omitempty"`
}

type Owner struct {
	GitHub string `json:"github"`
	Type   string `json:"type"`
}
