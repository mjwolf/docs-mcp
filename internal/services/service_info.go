package services

import (
	"fmt"
	"strings"

	"elastic-integration-docs-mcp/internal/shared"
)

type ServiceInfoProvider struct {
	services map[string]shared.ServiceInfo
}

func NewServiceInfoProvider() *ServiceInfoProvider {
	provider := &ServiceInfoProvider{
		services: make(map[string]shared.ServiceInfo),
	}
	provider.initializeServices()
	return provider
}

func (s *ServiceInfoProvider) initializeServices() {
	// Web servers
	s.services["nginx"] = shared.ServiceInfo{
		Name:        "nginx",
		Title:       "Nginx",
		Description: "Collect logs and metrics from Nginx HTTP servers with Elastic Agent.",
		Categories:  []string{"web", "observability"},
		Requirements: shared.Requirements{
			Kibana:       "^8.13.0 || ^9.0.0",
			Subscription: "basic",
		},
		SupportedVersions: []string{"1.19.5", "1.20.x", "1.21.x", "1.22.x", "1.23.x", "1.24.x", "1.25.x"},
		DataStreams: []shared.DataStream{
			{
				Name:        "access",
				Type:        "logs",
				Description: "Nginx access logs containing HTTP request information",
			},
			{
				Name:        "error",
				Type:        "logs",
				Description: "Nginx error logs containing server errors and warnings",
			},
			{
				Name:        "stubstatus",
				Type:        "metrics",
				Description: "Nginx stub status metrics for monitoring server performance",
			},
		},
		SetupComplexity: "low",
		CommonUseCases: []string{
			"Web server monitoring",
			"HTTP request analysis",
			"Performance monitoring",
			"Security monitoring",
			"Load balancer monitoring",
		},
		OfficialDocs: []string{
			"https://nginx.org/en/docs/",
			"https://nginx.org/en/docs/http/ngx_http_stub_status_module.html",
			"https://nginx.org/en/docs/http/log_module.html",
		},
		CommunityResources: []string{
			"https://github.com/nginx/nginx",
			"https://www.nginx.com/resources/wiki/",
		},
	}

	// Databases
	s.services["mysql"] = shared.ServiceInfo{
		Name:        "mysql",
		Title:       "MySQL",
		Description: "Collect logs and metrics from MySQL servers with Elastic Agent.",
		Categories:  []string{"datastore", "observability", "database_security", "security"},
		Requirements: shared.Requirements{
			Kibana:       "^8.15.0 || ^9.0.0",
			Subscription: "basic",
		},
		SupportedVersions: []string{"5.7.x", "8.0.x", "8.1.x", "8.2.x", "8.3.x", "8.4.x"},
		DataStreams: []shared.DataStream{
			{
				Name:        "error",
				Type:        "logs",
				Description: "MySQL error logs containing database errors and warnings",
			},
			{
				Name:        "slowlog",
				Type:        "logs",
				Description: "MySQL slow query logs for performance analysis",
			},
			{
				Name:        "status",
				Type:        "metrics",
				Description: "MySQL status metrics for monitoring database performance",
			},
			{
				Name:        "replica_status",
				Type:        "metrics",
				Description: "MySQL replication status metrics",
			},
			{
				Name:        "galera_status",
				Type:        "metrics",
				Description: "MySQL Galera cluster status metrics",
			},
			{
				Name:        "performance",
				Type:        "metrics",
				Description: "MySQL performance metrics",
			},
		},
		SetupComplexity: "medium",
		CommonUseCases: []string{
			"Database performance monitoring",
			"Query performance analysis",
			"Replication monitoring",
			"Security monitoring",
			"Capacity planning",
		},
		OfficialDocs: []string{
			"https://dev.mysql.com/doc/",
			"https://dev.mysql.com/doc/refman/8.0/en/slow-query-log.html",
			"https://dev.mysql.com/doc/refman/8.0/en/replication.html",
		},
		CommunityResources: []string{
			"https://github.com/mysql/mysql-server",
			"https://www.percona.com/resources/mysql",
		},
	}

	// Cloud services
	s.services["aws"] = shared.ServiceInfo{
		Name:        "aws",
		Title:       "AWS",
		Description: "Collect logs and metrics from Amazon Web Services (AWS) with Elastic Agent.",
		Categories:  []string{"aws", "cloud", "observability", "security"},
		Requirements: shared.Requirements{
			Kibana:       "~8.16.6 || ~8.17.4 || ^8.18.0 || ^9.0.0",
			Subscription: "basic",
		},
		SupportedVersions: []string{"All current AWS services"},
		DataStreams: []shared.DataStream{
			{
				Name:        "cloudtrail",
				Type:        "logs",
				Description: "AWS CloudTrail logs for API activity monitoring",
			},
			{
				Name:        "cloudwatch_logs",
				Type:        "logs",
				Description: "AWS CloudWatch logs from various services",
			},
			{
				Name:        "cloudwatch_metrics",
				Type:        "metrics",
				Description: "AWS CloudWatch metrics from various services",
			},
			{
				Name:        "vpcflow",
				Type:        "logs",
				Description: "AWS VPC Flow Logs for network traffic analysis",
			},
			{
				Name:        "s3access",
				Type:        "logs",
				Description: "AWS S3 access logs for object access monitoring",
			},
			{
				Name:        "guardduty",
				Type:        "logs",
				Description: "AWS GuardDuty security findings",
			},
			{
				Name:        "securityhub_findings",
				Type:        "logs",
				Description: "AWS Security Hub security findings",
			},
		},
		SetupComplexity: "high",
		CommonUseCases: []string{
			"Cloud security monitoring",
			"Cost optimization",
			"Performance monitoring",
			"Compliance monitoring",
			"Incident response",
		},
		OfficialDocs: []string{
			"https://docs.aws.amazon.com/",
			"https://docs.aws.amazon.com/cloudtrail/",
			"https://docs.aws.amazon.com/cloudwatch/",
			"https://docs.aws.amazon.com/guardduty/",
		},
		CommunityResources: []string{
			"https://github.com/aws",
			"https://aws.amazon.com/blogs/",
		},
	}

	// Add more services as needed
	s.services["apache"] = shared.ServiceInfo{
		Name:        "apache",
		Title:       "Apache HTTP Server",
		Description: "Collect logs and metrics from Apache HTTP servers with Elastic Agent.",
		Categories:  []string{"web", "observability"},
		Requirements: shared.Requirements{
			Kibana:       "^8.13.0 || ^9.0.0",
			Subscription: "basic",
		},
		SupportedVersions: []string{"2.4.x"},
		DataStreams: []shared.DataStream{
			{
				Name:        "access",
				Type:        "logs",
				Description: "Apache access logs containing HTTP request information",
			},
			{
				Name:        "error",
				Type:        "logs",
				Description: "Apache error logs containing server errors and warnings",
			},
			{
				Name:        "status",
				Type:        "metrics",
				Description: "Apache status metrics for monitoring server performance",
			},
		},
		SetupComplexity: "low",
		CommonUseCases: []string{
			"Web server monitoring",
			"HTTP request analysis",
			"Performance monitoring",
			"Security monitoring",
		},
		OfficialDocs: []string{
			"https://httpd.apache.org/docs/",
			"https://httpd.apache.org/docs/2.4/mod/mod_status.html",
		},
		CommunityResources: []string{
			"https://github.com/apache/httpd",
			"https://httpd.apache.org/docs/2.4/logs.html",
		},
	}

	s.services["postgresql"] = shared.ServiceInfo{
		Name:        "postgresql",
		Title:       "PostgreSQL",
		Description: "Collect logs and metrics from PostgreSQL servers with Elastic Agent.",
		Categories:  []string{"datastore", "observability", "database_security", "security"},
		Requirements: shared.Requirements{
			Kibana:       "^8.15.0 || ^9.0.0",
			Subscription: "basic",
		},
		SupportedVersions: []string{"12.x", "13.x", "14.x", "15.x", "16.x"},
		DataStreams: []shared.DataStream{
			{
				Name:        "log",
				Type:        "logs",
				Description: "PostgreSQL log files containing database events",
			},
			{
				Name:        "statement",
				Type:        "logs",
				Description: "PostgreSQL statement logs for query analysis",
			},
			{
				Name:        "activity",
				Type:        "metrics",
				Description: "PostgreSQL activity metrics for monitoring database performance",
			},
		},
		SetupComplexity: "medium",
		CommonUseCases: []string{
			"Database performance monitoring",
			"Query performance analysis",
			"Security monitoring",
			"Capacity planning",
		},
		OfficialDocs: []string{
			"https://www.postgresql.org/docs/",
			"https://www.postgresql.org/docs/current/runtime-config-logging.html",
		},
		CommunityResources: []string{
			"https://github.com/postgres/postgres",
			"https://www.postgresql.org/docs/current/monitoring.html",
		},
	}
}

func (s *ServiceInfoProvider) GetServiceInfo(serviceName string) (shared.CallToolResult, error) {
	service, exists := s.services[strings.ToLower(serviceName)]
	if !exists {
		availableServices := make([]string, 0, len(s.services))
		for name := range s.services {
			availableServices = append(availableServices, name)
		}
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Service '%s' not found. Available services: %s", serviceName, strings.Join(availableServices, ", ")),
				},
			},
			IsError: true,
		}, nil
	}

	info := fmt.Sprintf(`# %s Service Information

## Overview
%s

## Categories
%s

## Requirements
- **Kibana**: %s
- **Subscription**: %s

## Supported Versions
%s

## Data Streams
%s

## Setup Complexity
**%s** - %s

## Common Use Cases
%s

## Official Documentation
%s

## Community Resources
%s`,
		service.Title,
		service.Description,
		formatList(service.Categories),
		service.Requirements.Kibana,
		service.Requirements.Subscription,
		formatList(service.SupportedVersions),
		formatDataStreams(service.DataStreams),
		strings.ToUpper(service.SetupComplexity),
		s.getComplexityDescription(service.SetupComplexity),
		formatList(service.CommonUseCases),
		formatLinks(service.OfficialDocs),
		formatLinks(service.CommunityResources))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: info,
			},
		},
	}, nil
}

func (s *ServiceInfoProvider) GetServiceCategories() (shared.CallToolResult, error) {
	categories := make(map[string][]string)

	for _, service := range s.services {
		for _, category := range service.Categories {
			categories[category] = append(categories[category], service.Name)
		}
	}

	var result strings.Builder
	result.WriteString("# Available Service Categories\n\n")

	for category, services := range categories {
		result.WriteString(fmt.Sprintf("## %s\n", category))
		for _, service := range services {
			result.WriteString(fmt.Sprintf("- %s\n", service))
		}
		result.WriteString("\n")
	}

	result.WriteString("## All Services\n")
	allServices := make([]string, 0, len(s.services))
	for name := range s.services {
		allServices = append(allServices, name)
	}
	result.WriteString(formatList(allServices))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: result.String(),
			},
		},
	}, nil
}

func (s *ServiceInfoProvider) getComplexityDescription(complexity string) string {
	switch complexity {
	case "low":
		return "Simple setup with minimal configuration required"
	case "medium":
		return "Moderate setup complexity with some configuration and dependencies"
	case "high":
		return "Complex setup requiring significant configuration and understanding of the service"
	default:
		return "Unknown complexity level"
	}
}

func formatList(items []string) string {
	var result strings.Builder
	for _, item := range items {
		result.WriteString(fmt.Sprintf("- %s\n", item))
	}
	return result.String()
}

func formatDataStreams(streams []shared.DataStream) string {
	var result strings.Builder
	for _, stream := range streams {
		result.WriteString(fmt.Sprintf("### %s (%s)\n%s\n\n", stream.Name, stream.Type, stream.Description))
	}
	return result.String()
}

func formatLinks(urls []string) string {
	var result strings.Builder
	for _, url := range urls {
		result.WriteString(fmt.Sprintf("- [%s](%s)\n", url, url))
	}
	return result.String()
}
