package services

import (
	"fmt"
	"strings"

	"elastic-integration-docs-mcp/internal/config"
	"elastic-integration-docs-mcp/internal/shared"
)

type ValidationProvider struct {
	configLoader *config.ConfigLoader
}

func NewValidationProvider(configDir string) *ValidationProvider {
	configLoader := config.NewConfigLoader(configDir)
	if err := configLoader.LoadAllServices(); err != nil {
		// In a real implementation, you might want to handle this error differently
		// For now, we'll create an empty loader
		configLoader = config.NewConfigLoader(configDir)
	}

	return &ValidationProvider{
		configLoader: configLoader,
	}
}

func (v *ValidationProvider) GetValidationSteps(serviceName string) (shared.CallToolResult, error) {
	serviceConfig, err := v.configLoader.GetServiceConfig(serviceName)
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

	validationSteps := fmt.Sprintf(`# %s Integration Validation Steps

%s

## Summary
These validation steps will help you verify that the %s integration is running properly and collecting data as expected.`,
		strings.ToUpper(serviceConfig.ServiceName),
		formatValidationSteps(serviceConfig.ValidationSteps.Steps),
		serviceConfig.ServiceName)

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: validationSteps,
			},
		},
	}, nil
}

func formatValidationSteps(steps []config.ValidationStep) string {
	var result strings.Builder
	for _, step := range steps {
		result.WriteString(fmt.Sprintf("\n## Step %d: %s\n%s\n\n",
			step.Step, step.Title, step.Description))

		if len(step.Commands) > 0 {
			result.WriteString("**Commands:**\n```bash\n")
			for _, cmd := range step.Commands {
				result.WriteString(cmd + "\n")
			}
			result.WriteString("```\n\n")
		}

		if step.ExpectedOutput != "" {
			result.WriteString(fmt.Sprintf("**Expected Output:**\n%s\n\n", step.ExpectedOutput))
		}
	}
	return result.String()
}
