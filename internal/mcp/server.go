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
}

func NewServer() *Server {
	// Try to find config directory relative to the executable
	configDir := "config"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// If not found, try relative to project root (go up one level from cmd/server)
		configDir = "../config"
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			// If still not found, try absolute path from project root
			configDir = "/Users/mwolf/git/docs-mcp/config"
		}
	}

	return &Server{
		serviceInfo:   services.NewServiceInfoProvider(configDir),
		setupGuide:    services.NewSetupGuideProvider(configDir),
		documentation: services.NewDocumentationProvider(configDir),
		validation:    services.NewValidationProvider(configDir),
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
			Name:        "search_documentation",
			Description: "Perform a web search of the search term, restricted to documentation sites for that service",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"search_term": map[string]interface{}{
						"type":        "string",
						"description": "Search term to look for in documentation",
					},
					"service_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service (e.g., apache, nginx, mysql)",
					},
				},
				"required": []string{"search_term", "service_name"},
			},
		},
		{
			Name:        "get_service_info",
			Description: "Get curated info on the service including common use cases, data types collected, compatibility, and scaling information",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"service_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service (e.g., nginx, mysql, aws)",
					},
				},
				"required": []string{"service_name"},
			},
		},
		{
			Name:        "get_service_setup_instructions",
			Description: "Return a list of known good, working steps to set up a service to prepare it to send data to the integration",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"service_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
					"version": map[string]interface{}{
						"type":        "string",
						"description": "Service version (optional)",
					},
				},
				"required": []string{"service_name"},
			},
		},
		{
			Name:        "get_kibana_setup_instructions",
			Description: "Return the steps to configure the service in Kibana",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"service_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
					"input_type": map[string]interface{}{
						"type":        "string",
						"description": "Input type (optional, e.g., tcp, udp)",
					},
					"version": map[string]interface{}{
						"type":        "string",
						"description": "Service version (optional)",
					},
				},
				"required": []string{"service_name"},
			},
		},
		{
			Name:        "get_troubleshooting_help",
			Description: "Return a list of common problems and solutions for the service",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"service_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
				},
				"required": []string{"service_name"},
			},
		},
		{
			Name:        "get_validation_steps",
			Description: "Return a list of steps for how to validate that the integration is running properly",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"service_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the service",
					},
				},
				"required": []string{"service_name"},
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
	case "search_documentation":
		searchTerm, ok := callRequest.Arguments["search_term"].(string)
		if !ok {
			err = fmt.Errorf("search_term is required")
			break
		}
		serviceName, ok := callRequest.Arguments["service_name"].(string)
		if !ok {
			err = fmt.Errorf("service_name is required")
			break
		}
		result, err = s.documentation.SearchDocumentation(searchTerm, serviceName)

	case "get_service_info":
		serviceName, ok := callRequest.Arguments["service_name"].(string)
		if !ok {
			err = fmt.Errorf("service_name is required")
			break
		}
		result, err = s.serviceInfo.GetServiceInfo(serviceName)

	case "get_service_setup_instructions":
		serviceName, ok := callRequest.Arguments["service_name"].(string)
		if !ok {
			err = fmt.Errorf("service_name is required")
			break
		}
		version, _ := callRequest.Arguments["version"].(string)
		result, err = s.setupGuide.GetServiceSetupInstructions(serviceName, version)

	case "get_kibana_setup_instructions":
		serviceName, ok := callRequest.Arguments["service_name"].(string)
		if !ok {
			err = fmt.Errorf("service_name is required")
			break
		}
		inputType, _ := callRequest.Arguments["input_type"].(string)
		version, _ := callRequest.Arguments["version"].(string)
		result, err = s.setupGuide.GetKibanaSetupInstructions(serviceName, inputType, version)

	case "get_troubleshooting_help":
		serviceName, ok := callRequest.Arguments["service_name"].(string)
		if !ok {
			err = fmt.Errorf("service_name is required")
			break
		}
		result, err = s.documentation.GetTroubleshootingHelp(serviceName)

	case "get_validation_steps":
		serviceName, ok := callRequest.Arguments["service_name"].(string)
		if !ok {
			err = fmt.Errorf("service_name is required")
			break
		}
		result, err = s.validation.GetValidationSteps(serviceName)

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
