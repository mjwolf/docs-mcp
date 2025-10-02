package shared

// TroubleshootingGuide represents a troubleshooting guide for a service
type TroubleshootingGuide struct {
	ServiceName        string                 `json:"serviceName"`
	CommonIssues       []TroubleshootingIssue `json:"commonIssues"`
	DiagnosticCommands []string               `json:"diagnosticCommands"`
	LogLocations       []string               `json:"logLocations"`
	SupportResources   []string               `json:"supportResources"`
}

// TroubleshootingIssue represents a specific troubleshooting issue
type TroubleshootingIssue struct {
	Issue      string   `json:"issue"`
	Symptoms   []string `json:"symptoms"`
	Causes     []string `json:"causes"`
	Solutions  []string `json:"solutions"`
	Prevention []string `json:"prevention"`
}

// CallToolResult represents the result of a tool call
type CallToolResult struct {
	Content []ToolContent          `json:"content"`
	IsError bool                   `json:"isError,omitempty"`
	Meta    map[string]interface{} `json:"_meta,omitempty"`
}

// ToolContent represents content returned by a tool
type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// IntegrationDetails represents details about an Elastic integration
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

// IntegrationDataStream represents a data stream in an integration
type IntegrationDataStream struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Fields      []Field `json:"fields"`
}

// Field represents a field in a data stream
type Field struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Example     string `json:"example,omitempty"`
}

// PolicyTemplate represents a policy template for an integration
type PolicyTemplate struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Inputs      []Input  `json:"inputs"`
	DataStreams []string `json:"dataStreams"`
	Categories  []string `json:"categories"`
}

// Input represents an input for a policy template
type Input struct {
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Vars        []Variable `json:"vars,omitempty"`
}

// Variable represents a variable in an input
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

// Owner represents the owner of an integration
type Owner struct {
	GitHub string `json:"github"`
	Type   string `json:"type"`
}

// Requirements represents requirements for a service
type Requirements struct {
	Elasticsearch string `json:"elasticsearch,omitempty"`
	Kibana        string `json:"kibana,omitempty"`
	Subscription  string `json:"subscription,omitempty"`
}

// ServiceInfo represents information about a service
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

// DataStream represents a data stream for a service
type DataStream struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// SetupInstructions represents setup instructions for a service
type SetupInstructions struct {
	ServiceName       string                `json:"serviceName"`
	Platform          string                `json:"platform,omitempty"`
	Version           string                `json:"version,omitempty"`
	Prerequisites     []string              `json:"prerequisites"`
	InstallationSteps []InstallationStep    `json:"installationSteps"`
	PostInstallation  []PostInstallation    `json:"postInstallation"`
	Troubleshooting   []TroubleshootingItem `json:"troubleshooting"`
}

// InstallationStep represents a step in the installation process
type InstallationStep struct {
	Step           int             `json:"step"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Commands       []string        `json:"commands,omitempty"`
	ConfigSnippets []ConfigSnippet `json:"configSnippets,omitempty"`
	Verification   string          `json:"verification,omitempty"`
}

// ConfigSnippet represents a configuration snippet
type ConfigSnippet struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

// PostInstallation represents post-installation configuration
type PostInstallation struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
}

// TroubleshootingItem represents a troubleshooting item
type TroubleshootingItem struct {
	Issue    string `json:"issue"`
	Solution string `json:"solution"`
}

// ValidationResult represents the result of configuration validation
type ValidationResult struct {
	IsValid     bool                   `json:"isValid"`
	Errors      []ValidationError      `json:"errors"`
	Warnings    []ValidationWarning    `json:"warnings"`
	Suggestions []ValidationSuggestion `json:"suggestions"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Line     *int   `json:"line,omitempty"`
	Column   *int   `json:"column,omitempty"`
	Severity string `json:"severity"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Line       *int   `json:"line,omitempty"`
	Column     *int   `json:"column,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

// ValidationSuggestion represents a validation suggestion
type ValidationSuggestion struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
	Impact     string `json:"impact"`
}
