package services

import (
	"fmt"
	"strings"

	"elastic-integration-docs-mcp/internal/config"
	"elastic-integration-docs-mcp/internal/shared"
)

type ServiceInfoProvider struct {
	configLoader *config.ConfigLoader
}

func NewServiceInfoProvider(configDir string) *ServiceInfoProvider {
	configLoader := config.NewConfigLoader(configDir)
	if err := configLoader.LoadAllServices(); err != nil {
		// In a real implementation, you might want to handle this error differently
		// For now, we'll create an empty loader
		configLoader = config.NewConfigLoader(configDir)
	}

	return &ServiceInfoProvider{
		configLoader: configLoader,
	}
}

func (s *ServiceInfoProvider) GetServiceInfo(serviceName string) (shared.CallToolResult, error) {
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

	// Format the service info according to the requirements
	info := fmt.Sprintf(`# %s Service Information

## Common Use Cases
%s

## Data Types Collected
%s

## Compatibility
- **Elastic Stack Versions**: %s
- **Service Versions**: %s

## Scaling and Performance
%s

### Performance Expectations
%s

### Scaling Guidance
%s`,
		serviceConfig.Title,
		formatList(serviceConfig.ServiceInfo.CommonUseCases),
		formatList(serviceConfig.ServiceInfo.DataTypesCollected),
		strings.Join(serviceConfig.ServiceInfo.Compatibility.ElasticStackVersions, ", "),
		strings.Join(serviceConfig.ServiceInfo.Compatibility.ServiceVersions, ", "),
		serviceConfig.ServiceInfo.ScalingAndPerformance.Description,
		formatList(serviceConfig.ServiceInfo.ScalingAndPerformance.PerformanceExpectations),
		formatList(serviceConfig.ServiceInfo.ScalingAndPerformance.ScalingGuidance))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: info,
			},
		},
	}, nil
}
