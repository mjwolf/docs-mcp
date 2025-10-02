package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"elastic-integration-docs-mcp/internal/services"
	"elastic-integration-docs-mcp/internal/shared"
)

type Server struct {
	serviceInfo   *services.ServiceInfoProvider
	setupGuide    *services.SetupGuideProvider
	documentation *services.DocumentationProvider
	validation    *services.ValidationProvider
	integration   *services.IntegrationProvider
}

func NewServer() *Server {
	return &Server{
		serviceInfo:   services.NewServiceInfoProvider(),
		setupGuide:    services.NewSetupGuideProvider(),
		documentation: services.NewDocumentationProvider(),
		validation:    services.NewValidationProvider(),
		integration:   services.NewIntegrationProvider(),
	}
}

func (s *Server) Run() error {
	log.SetOutput(os.Stderr)
	log.Println("Elastic Integration Docs MCP server running on stdio")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var request JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			log.Printf("Error parsing request: %v", err)
			continue
		}

		response := s.handleRequest(request)
		responseData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			continue
		}

		fmt.Println(string(responseData))
	}

	return scanner.Err()
}

func (s *Server) handleRequest(request JSONRPCRequest) JSONRPCResponse {
	switch request.Method {
	case "initialize":
		return s.handleInitialize(request)
	case "tools/list":
		return s.handleListTools(request)
	case "tools/call":
		return s.handleCallTool(request)
	default:
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func (s *Server) handleInitialize(request JSONRPCRequest) JSONRPCResponse {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: true,
			},
		},
		ServerInfo: ServerInfo{
			Name:    "elastic-integration-docs",
			Version: "1.0.0",
		},
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func (s *Server) handleListTools(request JSONRPCRequest) JSONRPCResponse {
	tools := []Tool{
		{
			Name:        "get_service_info",
			Description: "Get comprehensive information about a service including requirements, capabilities, and supported versions",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"serviceName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service (e.g., nginx, mysql, aws)",
					},
				},
				"required": []string{"serviceName"},
			},
		},
		{
			Name:        "get_setup_instructions",
			Description: "Get step-by-step setup instructions for a service with platform-specific guidance",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"serviceName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
					"platform": map[string]interface{}{
						"type":        "string",
						"description": "Target platform (e.g., ubuntu, centos, docker)",
					},
					"version": map[string]interface{}{
						"type":        "string",
						"description": "Service version",
					},
				},
				"required": []string{"serviceName"},
			},
		},
		{
			Name:        "get_configuration_examples",
			Description: "Get configuration examples and templates for a service",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"serviceName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
					"configType": map[string]interface{}{
						"type":        "string",
						"description": "Type of configuration (e.g., logs, metrics, security)",
					},
				},
				"required": []string{"serviceName"},
			},
		},
		{
			Name:        "search_service_docs",
			Description: "Search for service-specific documentation and guides",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"serviceName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
					"docType": map[string]interface{}{
						"type":        "string",
						"description": "Type of documentation (official, community, troubleshooting)",
					},
				},
				"required": []string{"serviceName", "query"},
			},
		},
		{
			Name:        "validate_configuration",
			Description: "Validate service configuration and provide suggestions for improvements",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"serviceName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
					"configuration": map[string]interface{}{
						"type":        "string",
						"description": "Configuration content to validate",
					},
					"configType": map[string]interface{}{
						"type":        "string",
						"description": "Type of configuration (yaml, json, conf, etc.)",
					},
				},
				"required": []string{"serviceName", "configuration", "configType"},
			},
		},
		{
			Name:        "get_integration_details",
			Description: "Get details about Elastic integration including data streams and field mappings",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"integrationName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the Elastic integration",
					},
				},
				"required": []string{"integrationName"},
			},
		},
		{
			Name:        "get_troubleshooting_guide",
			Description: "Get troubleshooting guide for common issues with a service",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"serviceName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
					"issue": map[string]interface{}{
						"type":        "string",
						"description": "Specific issue or error message",
					},
				},
				"required": []string{"serviceName"},
			},
		},
		{
			Name:        "get_service_categories",
			Description: "Get list of available service categories and services within each category",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        "get_latest_docs",
			Description: "Get the latest official documentation for a service",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"serviceName": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
					"docType": map[string]interface{}{
						"type":        "string",
						"description": "Type of documentation (installation, configuration, api)",
					},
				},
				"required": []string{"serviceName"},
			},
		},
	}

	result := ListToolsResult{
		Tools: tools,
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func (s *Server) handleCallTool(request JSONRPCRequest) JSONRPCResponse {
	var callRequest CallToolRequest
	if err := json.Unmarshal(request.Params, &callRequest); err != nil {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	var result shared.CallToolResult
	var err error

	switch callRequest.Name {
	case "get_service_info":
		serviceName, ok := callRequest.Arguments["serviceName"].(string)
		if !ok {
			err = fmt.Errorf("serviceName is required")
			break
		}
		result, err = s.serviceInfo.GetServiceInfo(serviceName)

	case "get_setup_instructions":
		serviceName, ok := callRequest.Arguments["serviceName"].(string)
		if !ok {
			err = fmt.Errorf("serviceName is required")
			break
		}
		platform, _ := callRequest.Arguments["platform"].(string)
		version, _ := callRequest.Arguments["version"].(string)
		result, err = s.setupGuide.GetSetupInstructions(serviceName, platform, version)

	case "get_configuration_examples":
		serviceName, ok := callRequest.Arguments["serviceName"].(string)
		if !ok {
			err = fmt.Errorf("serviceName is required")
			break
		}
		configType, _ := callRequest.Arguments["configType"].(string)
		result, err = s.setupGuide.GetConfigurationExamples(serviceName, configType)

	case "search_service_docs":
		serviceName, ok := callRequest.Arguments["serviceName"].(string)
		if !ok {
			err = fmt.Errorf("serviceName is required")
			break
		}
		query, ok := callRequest.Arguments["query"].(string)
		if !ok {
			err = fmt.Errorf("query is required")
			break
		}
		docType, _ := callRequest.Arguments["docType"].(string)
		result, err = s.documentation.SearchServiceDocs(serviceName, query, docType)

	case "validate_configuration":
		serviceName, ok := callRequest.Arguments["serviceName"].(string)
		if !ok {
			err = fmt.Errorf("serviceName is required")
			break
		}
		configuration, ok := callRequest.Arguments["configuration"].(string)
		if !ok {
			err = fmt.Errorf("configuration is required")
			break
		}
		configType, ok := callRequest.Arguments["configType"].(string)
		if !ok {
			err = fmt.Errorf("configType is required")
			break
		}
		result, err = s.validation.ValidateConfiguration(serviceName, configuration, configType)

	case "get_integration_details":
		integrationName, ok := callRequest.Arguments["integrationName"].(string)
		if !ok {
			err = fmt.Errorf("integrationName is required")
			break
		}
		result, err = s.integration.GetIntegrationDetails(integrationName)

	case "get_troubleshooting_guide":
		serviceName, ok := callRequest.Arguments["serviceName"].(string)
		if !ok {
			err = fmt.Errorf("serviceName is required")
			break
		}
		issue, _ := callRequest.Arguments["issue"].(string)
		result, err = s.documentation.GetTroubleshootingGuide(serviceName, issue)

	case "get_service_categories":
		result, err = s.serviceInfo.GetServiceCategories()

	case "get_latest_docs":
		serviceName, ok := callRequest.Arguments["serviceName"].(string)
		if !ok {
			err = fmt.Errorf("serviceName is required")
			break
		}
		docType, _ := callRequest.Arguments["docType"].(string)
		result, err = s.documentation.GetLatestDocs(serviceName, docType)

	default:
		err = fmt.Errorf("unknown tool: %s", callRequest.Name)
	}

	if err != nil {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: shared.CallToolResult{
				Content: []shared.ToolContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Error: %v", err),
					},
				},
				IsError: true,
			},
		}
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}
