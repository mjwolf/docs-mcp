package services

import (
	"fmt"
	"strings"

	"elastic-integration-docs-mcp/internal/config"
	"elastic-integration-docs-mcp/internal/shared"
)

type SetupGuideProvider struct {
	configLoader *config.ConfigLoader
}

func NewSetupGuideProvider(configDir string) *SetupGuideProvider {
	configLoader := config.NewConfigLoader(configDir)
	if err := configLoader.LoadAllServices(); err != nil {
		// In a real implementation, you might want to handle this error differently
		// For now, we'll create an empty loader
		configLoader = config.NewConfigLoader(configDir)
	}

	return &SetupGuideProvider{
		configLoader: configLoader,
	}
}

func (s *SetupGuideProvider) GetServiceSetupInstructions(serviceName, version string) (shared.CallToolResult, error) {
	serviceConfig, err := s.configLoader.GetServiceConfig(serviceName)
	if err != nil {
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: err.Error(),
				},
			},
			IsError: true,
		}, nil
	}

	versionInfo := ""
	if version != "" {
		versionInfo = fmt.Sprintf("\n**Version**: %s", version)
	}

	instructions := fmt.Sprintf(`# %s Setup Instructions%s

## Prerequisites
%s

## Installation Steps

%s`,
		strings.ToUpper(serviceConfig.ServiceName),
		versionInfo,
		formatList(serviceConfig.SetupInstructions.Prerequisites),
		formatInstallationSteps(serviceConfig.SetupInstructions.InstallationSteps))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: instructions,
			},
		},
	}, nil
}

func (s *SetupGuideProvider) GetKibanaSetupInstructions(serviceName, inputType, version string) (shared.CallToolResult, error) {
	serviceConfig, err := s.configLoader.GetServiceConfig(serviceName)
	if err != nil {
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: err.Error(),
				},
			},
			IsError: true,
		}, nil
	}

	// Select the appropriate setup instructions based on input type
	var steps []config.KibanaSetupStep
	if inputType != "" {
		switch strings.ToLower(inputType) {
		case "tcp":
			if serviceConfig.KibanaSetupInstructions.TCP.Steps != nil {
				steps = serviceConfig.KibanaSetupInstructions.TCP.Steps
			}
		case "udp":
			if serviceConfig.KibanaSetupInstructions.UDP.Steps != nil {
				steps = serviceConfig.KibanaSetupInstructions.UDP.Steps
			}
		}
	}

	// Fall back to default if no specific input type or if not found
	if steps == nil {
		steps = serviceConfig.KibanaSetupInstructions.Default.Steps
	}

	// Format the steps as JSON-like structure as shown in requirements
	var stepsJSON strings.Builder
	stepsJSON.WriteString("{\n  \"steps\": [\n")

	for i, step := range steps {
		stepsJSON.WriteString(fmt.Sprintf("    {\n      \"step\": %d,\n      \"instruction\": \"%s\"\n    }",
			step.Step, step.Instruction))
		if i < len(steps)-1 {
			stepsJSON.WriteString(",")
		}
		stepsJSON.WriteString("\n")
	}

	stepsJSON.WriteString("  ]\n}")

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: stepsJSON.String(),
			},
		},
	}, nil
}

func formatInstallationSteps(steps []config.InstallationStep) string {
	var result strings.Builder
	for _, step := range steps {
		result.WriteString(fmt.Sprintf("\n### Step %d: %s\n%s\n\n", step.Step, step.Title, step.Description))

		if len(step.Commands) > 0 {
			result.WriteString("**Commands:**\n```bash\n")
			for _, cmd := range step.Commands {
				result.WriteString(cmd + "\n")
			}
			result.WriteString("```\n\n")
		}

		if len(step.ConfigSnippets) > 0 {
			for _, snippet := range step.ConfigSnippets {
				result.WriteString(fmt.Sprintf("**Configuration File: %s**\n```%s\n%s\n```\n\n",
					snippet.Filename, getFileExtension(snippet.Filename), snippet.Content))
			}
		}

		if step.Verification != "" {
			result.WriteString(fmt.Sprintf("**Verification:**\n%s\n\n", step.Verification))
		}
	}
	return result.String()
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return "text"
	}
	ext := strings.ToLower(parts[len(parts)-1])
	switch ext {
	case "conf", "cnf":
		return "ini"
	case "yml", "yaml":
		return "yaml"
	case "json":
		return "json"
	default:
		return "text"
	}
}
