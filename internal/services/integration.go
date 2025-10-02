package services

import (
	"fmt"
	"strings"

	"elastic-integration-docs-mcp/internal/shared"
)

type IntegrationProvider struct {
	integrations map[string]shared.IntegrationDetails
}

func NewIntegrationProvider() *IntegrationProvider {
	provider := &IntegrationProvider{
		integrations: make(map[string]shared.IntegrationDetails),
	}
	provider.initializeIntegrations()
	return provider
}

func (i *IntegrationProvider) initializeIntegrations() {
	// Nginx integration
	i.integrations["nginx"] = shared.IntegrationDetails{
		Name:        "nginx",
		Title:       "Nginx",
		Description: "Collect logs and metrics from Nginx HTTP servers with Elastic Agent.",
		Version:     "2.3.2",
		Categories:  []string{"web", "observability"},
		Requirements: shared.Requirements{
			Kibana:       "^8.13.0 || ^9.0.0",
			Subscription: "basic",
		},
		DataStreams: []shared.IntegrationDataStream{
			{
				Name:        "access",
				Type:        "logs",
				Description: "Nginx access logs containing HTTP request information",
				Fields: []shared.Field{
					{
						Name:        "nginx.access.remote_ip",
						Type:        "ip",
						Description: "Client IP address",
						Required:    true,
						Example:     "192.168.1.100",
					},
					{
						Name:        "nginx.access.method",
						Type:        "keyword",
						Description: "HTTP method",
						Required:    true,
						Example:     "GET",
					},
					{
						Name:        "nginx.access.url",
						Type:        "keyword",
						Description: "Requested URL",
						Required:    true,
						Example:     "/index.html",
					},
					{
						Name:        "nginx.access.response_code",
						Type:        "long",
						Description: "HTTP response code",
						Required:    true,
						Example:     "200",
					},
					{
						Name:        "nginx.access.body_sent.bytes",
						Type:        "long",
						Description: "Bytes sent in response body",
						Required:    false,
						Example:     "1024",
					},
				},
			},
			{
				Name:        "error",
				Type:        "logs",
				Description: "Nginx error logs containing server errors and warnings",
				Fields: []shared.Field{
					{
						Name:        "nginx.error.level",
						Type:        "keyword",
						Description: "Error level",
						Required:    true,
						Example:     "error",
					},
					{
						Name:        "nginx.error.message",
						Type:        "text",
						Description: "Error message",
						Required:    true,
						Example:     "connect() failed (111: Connection refused)",
					},
					{
						Name:        "nginx.error.pid",
						Type:        "long",
						Description: "Process ID",
						Required:    false,
						Example:     "12345",
					},
				},
			},
			{
				Name:        "stubstatus",
				Type:        "metrics",
				Description: "Nginx stub status metrics for monitoring server performance",
				Fields: []shared.Field{
					{
						Name:        "nginx.stubstatus.connections.active",
						Type:        "long",
						Description: "Active connections",
						Required:    true,
						Example:     "10",
					},
					{
						Name:        "nginx.stubstatus.requests.total",
						Type:        "long",
						Description: "Total requests",
						Required:    true,
						Example:     "1000",
					},
					{
						Name:        "nginx.stubstatus.connections.reading",
						Type:        "long",
						Description: "Connections reading",
						Required:    true,
						Example:     "2",
					},
					{
						Name:        "nginx.stubstatus.connections.writing",
						Type:        "long",
						Description: "Connections writing",
						Required:    true,
						Example:     "3",
					},
					{
						Name:        "nginx.stubstatus.connections.waiting",
						Type:        "long",
						Description: "Connections waiting",
						Required:    true,
						Example:     "5",
					},
				},
			},
		},
		PolicyTemplates: []shared.PolicyTemplate{
			{
				Name:        "nginx",
				Title:       "Nginx logs and metrics",
				Description: "Collect logs and metrics from Nginx instances",
				DataStreams: []string{"access", "error", "stubstatus"},
				Categories:  []string{"web", "observability"},
				Inputs: []shared.Input{
					{
						Type:        "logfile",
						Title:       "Collect logs from Nginx instances",
						Description: "Collecting Nginx access and error logs",
					},
					{
						Type:        "nginx/metrics",
						Title:       "Collect metrics from Nginx instances",
						Description: "Collecting Nginx stub status metrics",
						Vars: []shared.Variable{
							{
								Name:        "hosts",
								Type:        "text",
								Title:       "Hosts",
								Description: "Nginx status endpoint URLs",
								Required:    true,
								Multi:       true,
								Default:     []string{"http://127.0.0.1:80"},
							},
						},
					},
				},
			},
		},
		Screenshots: []string{
			"/img/nginx-metrics-overview.png",
			"/img/nginx-logs-access-error.png",
			"/img/nginx-logs-overview.png",
		},
		Icons: []string{"/img/logo_nginx.svg"},
		Owner: shared.Owner{
			GitHub: "elastic/obs-infraobs-integrations",
			Type:   "elastic",
		},
	}

	// MySQL integration
	i.integrations["mysql"] = shared.IntegrationDetails{
		Name:        "mysql",
		Title:       "MySQL",
		Description: "Collect logs and metrics from MySQL servers with Elastic Agent.",
		Version:     "1.28.1",
		Categories:  []string{"datastore", "observability", "database_security", "security"},
		Requirements: shared.Requirements{
			Kibana:       "^8.15.0 || ^9.0.0",
			Subscription: "basic",
		},
		DataStreams: []shared.IntegrationDataStream{
			{
				Name:        "error",
				Type:        "logs",
				Description: "MySQL error logs containing database errors and warnings",
				Fields: []shared.Field{
					{
						Name:        "mysql.error.level",
						Type:        "keyword",
						Description: "Error level",
						Required:    true,
						Example:     "ERROR",
					},
					{
						Name:        "mysql.error.message",
						Type:        "text",
						Description: "Error message",
						Required:    true,
						Example:     "Access denied for user",
					},
					{
						Name:        "mysql.error.thread_id",
						Type:        "long",
						Description: "Thread ID",
						Required:    false,
						Example:     "12345",
					},
				},
			},
			{
				Name:        "slowlog",
				Type:        "logs",
				Description: "MySQL slow query logs for performance analysis",
				Fields: []shared.Field{
					{
						Name:        "mysql.slowlog.query_time.sec",
						Type:        "float",
						Description: "Query execution time in seconds",
						Required:    true,
						Example:     "2.5",
					},
					{
						Name:        "mysql.slowlog.lock_time.sec",
						Type:        "float",
						Description: "Lock time in seconds",
						Required:    false,
						Example:     "0.1",
					},
					{
						Name:        "mysql.slowlog.rows_sent",
						Type:        "long",
						Description: "Rows sent",
						Required:    false,
						Example:     "100",
					},
					{
						Name:        "mysql.slowlog.rows_examined",
						Type:        "long",
						Description: "Rows examined",
						Required:    false,
						Example:     "1000",
					},
					{
						Name:        "mysql.slowlog.sql_text",
						Type:        "text",
						Description: "SQL query text",
						Required:    true,
						Example:     "SELECT * FROM users WHERE id = 1",
					},
				},
			},
			{
				Name:        "status",
				Type:        "metrics",
				Description: "MySQL status metrics for monitoring database performance",
				Fields: []shared.Field{
					{
						Name:        "mysql.status.connections",
						Type:        "long",
						Description: "Current connections",
						Required:    true,
						Example:     "50",
					},
					{
						Name:        "mysql.status.threads_connected",
						Type:        "long",
						Description: "Connected threads",
						Required:    true,
						Example:     "25",
					},
					{
						Name:        "mysql.status.threads_running",
						Type:        "long",
						Description: "Running threads",
						Required:    true,
						Example:     "5",
					},
					{
						Name:        "mysql.status.queries",
						Type:        "long",
						Description: "Total queries",
						Required:    true,
						Example:     "10000",
					},
					{
						Name:        "mysql.status.uptime",
						Type:        "long",
						Description: "Server uptime in seconds",
						Required:    true,
						Example:     "86400",
					},
				},
			},
			{
				Name:        "replica_status",
				Type:        "metrics",
				Description: "MySQL replication status metrics",
				Fields: []shared.Field{
					{
						Name:        "mysql.replica_status.slave_io_running",
						Type:        "keyword",
						Description: "Slave IO thread status",
						Required:    true,
						Example:     "Yes",
					},
					{
						Name:        "mysql.replica_status.slave_sql_running",
						Type:        "keyword",
						Description: "Slave SQL thread status",
						Required:    true,
						Example:     "Yes",
					},
					{
						Name:        "mysql.replica_status.seconds_behind_master",
						Type:        "long",
						Description: "Seconds behind master",
						Required:    true,
						Example:     "0",
					},
					{
						Name:        "mysql.replica_status.master_host",
						Type:        "keyword",
						Description: "Master host",
						Required:    true,
						Example:     "192.168.1.100",
					},
				},
			},
		},
		PolicyTemplates: []shared.PolicyTemplate{
			{
				Name:        "mysql",
				Title:       "MySQL logs and metrics",
				Description: "Collect logs and metrics from MySQL instances",
				DataStreams: []string{"error", "slowlog", "status", "replica_status"},
				Categories:  []string{"datastore", "observability"},
				Inputs: []shared.Input{
					{
						Type:        "logfile",
						Title:       "Collect logs from MySQL hosts",
						Description: "Collecting MySQL error and slowlog logs",
					},
					{
						Type:        "mysql/metrics",
						Title:       "Collect metrics from MySQL hosts",
						Description: "Collecting MySQL status and galera_status metrics",
						Vars: []shared.Variable{
							{
								Name:        "hosts",
								Type:        "text",
								Title:       "MySQL DSN",
								Description: "MySQL data source name",
								Required:    true,
								Multi:       true,
								Default:     []string{"tcp(127.0.0.1:3306)/"},
							},
							{
								Name:        "username",
								Type:        "text",
								Title:       "Username",
								Description: "MySQL username",
								Required:    false,
								Default:     "root",
							},
							{
								Name:        "password",
								Type:        "password",
								Title:       "Password",
								Description: "MySQL password",
								Required:    false,
								Secret:      true,
								Default:     "test",
							},
						},
					},
				},
			},
		},
		Screenshots: []string{
			"/img/kibana-mysql.png",
			"/img/metricbeat-mysql.png",
			"/img/mysql-replica_status-dashboard.png",
		},
		Icons: []string{"/img/logo_mysql.svg"},
		Owner: shared.Owner{
			GitHub: "elastic/obs-infraobs-integrations",
			Type:   "elastic",
		},
	}

	// AWS integration
	i.integrations["aws"] = shared.IntegrationDetails{
		Name:        "aws",
		Title:       "AWS",
		Description: "Collect logs and metrics from Amazon Web Services (AWS) with Elastic Agent.",
		Version:     "3.17.0",
		Categories:  []string{"aws", "cloud", "observability", "security"},
		Requirements: shared.Requirements{
			Kibana:       "~8.16.6 || ~8.17.4 || ^8.18.0 || ^9.0.0",
			Subscription: "basic",
		},
		DataStreams: []shared.IntegrationDataStream{
			{
				Name:        "cloudtrail",
				Type:        "logs",
				Description: "AWS CloudTrail logs for API activity monitoring",
				Fields: []shared.Field{
					{
						Name:        "aws.cloudtrail.event_name",
						Type:        "keyword",
						Description: "Event name",
						Required:    true,
						Example:     "CreateBucket",
					},
					{
						Name:        "aws.cloudtrail.event_source",
						Type:        "keyword",
						Description: "Event source",
						Required:    true,
						Example:     "s3.amazonaws.com",
					},
					{
						Name:        "aws.cloudtrail.user_identity.type",
						Type:        "keyword",
						Description: "User identity type",
						Required:    true,
						Example:     "IAMUser",
					},
					{
						Name:        "aws.cloudtrail.user_identity.user_name",
						Type:        "keyword",
						Description: "Username",
						Required:    false,
						Example:     "john.doe",
					},
					{
						Name:        "aws.cloudtrail.source_ip_address",
						Type:        "ip",
						Description: "Source IP address",
						Required:    false,
						Example:     "192.168.1.100",
					},
					{
						Name:        "aws.cloudtrail.response_elements",
						Type:        "object",
						Description: "Response elements",
						Required:    false,
					},
				},
			},
			{
				Name:        "cloudwatch_logs",
				Type:        "logs",
				Description: "AWS CloudWatch logs from various services",
				Fields: []shared.Field{
					{
						Name:        "aws.cloudwatch.log_group",
						Type:        "keyword",
						Description: "Log group name",
						Required:    true,
						Example:     "/aws/lambda/my-function",
					},
					{
						Name:        "aws.cloudwatch.log_stream",
						Type:        "keyword",
						Description: "Log stream name",
						Required:    true,
						Example:     "2023/01/01/[$LATEST]abc123",
					},
					{
						Name:        "aws.cloudwatch.message",
						Type:        "text",
						Description: "Log message",
						Required:    true,
						Example:     "Function execution started",
					},
					{
						Name:        "aws.cloudwatch.timestamp",
						Type:        "date",
						Description: "Log timestamp",
						Required:    true,
						Example:     "2023-01-01T00:00:00Z",
					},
				},
			},
			{
				Name:        "cloudwatch_metrics",
				Type:        "metrics",
				Description: "AWS CloudWatch metrics from various services",
				Fields: []shared.Field{
					{
						Name:        "aws.cloudwatch.namespace",
						Type:        "keyword",
						Description: "Metric namespace",
						Required:    true,
						Example:     "AWS/EC2",
					},
					{
						Name:        "aws.cloudwatch.metric_name",
						Type:        "keyword",
						Description: "Metric name",
						Required:    true,
						Example:     "CPUUtilization",
					},
					{
						Name:        "aws.cloudwatch.value",
						Type:        "float",
						Description: "Metric value",
						Required:    true,
						Example:     "75.5",
					},
					{
						Name:        "aws.cloudwatch.unit",
						Type:        "keyword",
						Description: "Metric unit",
						Required:    false,
						Example:     "Percent",
					},
					{
						Name:        "aws.cloudwatch.dimensions",
						Type:        "object",
						Description: "Metric dimensions",
						Required:    false,
					},
				},
			},
			{
				Name:        "vpcflow",
				Type:        "logs",
				Description: "AWS VPC Flow Logs for network traffic analysis",
				Fields: []shared.Field{
					{
						Name:        "aws.vpcflow.version",
						Type:        "long",
						Description: "Flow log version",
						Required:    true,
						Example:     "2",
					},
					{
						Name:        "aws.vpcflow.account_id",
						Type:        "keyword",
						Description: "AWS account ID",
						Required:    true,
						Example:     "123456789012",
					},
					{
						Name:        "aws.vpcflow.interface_id",
						Type:        "keyword",
						Description: "Network interface ID",
						Required:    true,
						Example:     "eni-12345678",
					},
					{
						Name:        "aws.vpcflow.srcaddr",
						Type:        "ip",
						Description: "Source IP address",
						Required:    true,
						Example:     "192.168.1.100",
					},
					{
						Name:        "aws.vpcflow.dstaddr",
						Type:        "ip",
						Description: "Destination IP address",
						Required:    true,
						Example:     "10.0.0.1",
					},
					{
						Name:        "aws.vpcflow.srcport",
						Type:        "long",
						Description: "Source port",
						Required:    true,
						Example:     "80",
					},
					{
						Name:        "aws.vpcflow.dstport",
						Type:        "long",
						Description: "Destination port",
						Required:    true,
						Example:     "443",
					},
					{
						Name:        "aws.vpcflow.protocol",
						Type:        "long",
						Description: "Protocol number",
						Required:    true,
						Example:     "6",
					},
					{
						Name:        "aws.vpcflow.packets",
						Type:        "long",
						Description: "Number of packets",
						Required:    true,
						Example:     "100",
					},
					{
						Name:        "aws.vpcflow.bytes",
						Type:        "long",
						Description: "Number of bytes",
						Required:    true,
						Example:     "1024",
					},
					{
						Name:        "aws.vpcflow.action",
						Type:        "keyword",
						Description: "Action taken",
						Required:    true,
						Example:     "ACCEPT",
					},
				},
			},
		},
		PolicyTemplates: []shared.PolicyTemplate{
			{
				Name:        "cloudtrail",
				Title:       "AWS CloudTrail",
				Description: "Collect AWS CloudTrail logs with Elastic Agent",
				DataStreams: []string{"cloudtrail"},
				Categories:  []string{"security"},
				Inputs: []shared.Input{
					{
						Type:        "aws-s3",
						Title:       "Collect CloudTrail logs from S3",
						Description: "Collecting logs from CloudTrail using aws-s3 input",
					},
					{
						Type:        "aws-cloudwatch",
						Title:       "Collect CloudTrail logs from CloudWatch",
						Description: "Collecting logs from CloudTrail using aws-cloudwatch input",
					},
				},
			},
			{
				Name:        "cloudwatch",
				Title:       "AWS CloudWatch",
				Description: "Use this integration to collect logs and metrics from Amazon CloudWatch with Elastic Agent",
				DataStreams: []string{"cloudwatch_logs", "cloudwatch_metrics"},
				Categories:  []string{"observability", "monitoring"},
				Inputs: []shared.Input{
					{
						Type:        "aws-cloudwatch",
						Title:       "Collect logs from CloudWatch",
						Description: "Collecting logs using aws-cloudwatch input",
					},
					{
						Type:        "aws/metrics",
						Title:       "Collect metrics from CloudWatch",
						Description: "Collecting metrics using AWS CloudWatch",
					},
				},
			},
			{
				Name:        "vpcflow",
				Title:       "Amazon VPC",
				Description: "Collect Amazon VPC flow logs with Elastic Agent",
				DataStreams: []string{"vpcflow"},
				Categories:  []string{"observability", "network"},
				Inputs: []shared.Input{
					{
						Type:        "aws-s3",
						Title:       "Collect VPC flow logs from S3",
						Description: "Collecting VPC Flow logs using aws-s3 input",
					},
					{
						Type:        "aws-cloudwatch",
						Title:       "Collect VPC flow logs from CloudWatch",
						Description: "Collecting VPC Flow logs using aws-cloudwatch input",
					},
				},
			},
		},
		Screenshots: []string{"/img/metricbeat-aws-overview.png"},
		Icons:       []string{"/img/logo_aws.svg"},
		Owner: shared.Owner{
			GitHub: "elastic/obs-ds-hosted-services",
			Type:   "elastic",
		},
	}
}

func (i *IntegrationProvider) GetIntegrationDetails(integrationName string) (shared.CallToolResult, error) {
	integration, exists := i.integrations[strings.ToLower(integrationName)]
	if !exists {
		availableIntegrations := make([]string, 0, len(i.integrations))
		for name := range i.integrations {
			availableIntegrations = append(availableIntegrations, name)
		}
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Integration '%s' not found. Available integrations: %s", integrationName, strings.Join(availableIntegrations, ", ")),
				},
			},
			IsError: true,
		}, nil
	}

	details := fmt.Sprintf(`# %s Integration Details

## Overview
%s

## Version
%s

## Categories
%s

## Requirements
- **Kibana**: %s
- **Elasticsearch**: %s
- **Subscription**: %s

## Data Streams

%s

## Policy Templates

%s

## Screenshots
%s

## Icons
%s

## Owner
- **GitHub**: %s
- **Type**: %s`,
		integration.Title,
		integration.Description,
		integration.Version,
		formatList(integration.Categories),
		integration.Requirements.Kibana,
		integration.Requirements.Elasticsearch,
		integration.Requirements.Subscription,
		formatIntegrationDataStreams(integration.DataStreams),
		formatPolicyTemplates(integration.PolicyTemplates),
		formatList(integration.Screenshots),
		formatList(integration.Icons),
		integration.Owner.GitHub,
		integration.Owner.Type)

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: details,
			},
		},
	}, nil
}

func formatIntegrationDataStreams(streams []shared.IntegrationDataStream) string {
	var result strings.Builder
	for _, stream := range streams {
		result.WriteString(fmt.Sprintf("\n### %s (%s)\n%s\n\n#### Key Fields\n%s\n\n",
			stream.Name,
			stream.Type,
			stream.Description,
			formatFields(stream.Fields)))
	}
	return result.String()
}

func formatFields(fields []shared.Field) string {
	var result strings.Builder
	for _, field := range fields {
		result.WriteString(fmt.Sprintf("- **%s** (%s)%s\n  - %s\n",
			field.Name,
			field.Type,
			func() string {
				if field.Required {
					return " *required*"
				}
				return ""
			}(),
			field.Description))
		if field.Example != "" {
			result.WriteString(fmt.Sprintf("  - Example: %s\n", field.Example))
		}
	}
	return result.String()
}

func formatPolicyTemplates(templates []shared.PolicyTemplate) string {
	var result strings.Builder
	for _, template := range templates {
		result.WriteString(fmt.Sprintf("\n### %s\n%s\n\n#### Data Streams\n%s\n\n#### Categories\n%s\n\n#### Inputs\n%s\n\n",
			template.Title,
			template.Description,
			formatList(template.DataStreams),
			formatList(template.Categories),
			formatInputs(template.Inputs)))
	}
	return result.String()
}

func formatInputs(inputs []shared.Input) string {
	var result strings.Builder
	for _, input := range inputs {
		result.WriteString(fmt.Sprintf("- **%s** (%s)\n  - %s\n",
			input.Title,
			input.Type,
			input.Description))
		if len(input.Vars) > 0 {
			result.WriteString("  **Variables:**\n")
			for _, variable := range input.Vars {
				result.WriteString(fmt.Sprintf("  - **%s** (%s)%s\n    - %s\n",
					variable.Name,
					variable.Type,
					func() string {
						if variable.Required {
							return " *required*"
						}
						return ""
					}(),
					variable.Title))
				if variable.Description != "" {
					result.WriteString(fmt.Sprintf("    - %s\n", variable.Description))
				}
				if variable.Default != nil {
					result.WriteString(fmt.Sprintf("    - Default: %v\n", variable.Default))
				}
				if variable.Multi {
					result.WriteString("    - Multi-value supported\n")
				}
				if variable.Secret {
					result.WriteString("    - Secret value\n")
				}
			}
		}
	}
	return result.String()
}
