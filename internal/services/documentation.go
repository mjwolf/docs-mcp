package services

import (
	"fmt"
	"strings"

	"elastic-integration-docs-mcp/internal/shared"
)

type DocumentationProvider struct {
	serviceDocs           map[string][]string
	troubleshootingGuides map[string]shared.TroubleshootingGuide
}

func NewDocumentationProvider() *DocumentationProvider {
	provider := &DocumentationProvider{
		serviceDocs:           make(map[string][]string),
		troubleshootingGuides: make(map[string]shared.TroubleshootingGuide),
	}
	provider.initializeServiceDocs()
	provider.initializeTroubleshootingGuides()
	return provider
}

func (d *DocumentationProvider) initializeServiceDocs() {
	d.serviceDocs["nginx"] = []string{
		"https://nginx.org/en/docs/",
		"https://nginx.org/en/docs/http/ngx_http_stub_status_module.html",
		"https://nginx.org/en/docs/http/log_module.html",
		"https://nginx.org/en/docs/http/ngx_http_access_module.html",
	}

	d.serviceDocs["mysql"] = []string{
		"https://dev.mysql.com/doc/",
		"https://dev.mysql.com/doc/refman/8.0/en/slow-query-log.html",
		"https://dev.mysql.com/doc/refman/8.0/en/replication.html",
		"https://dev.mysql.com/doc/refman/8.0/en/error-log.html",
	}

	d.serviceDocs["aws"] = []string{
		"https://docs.aws.amazon.com/",
		"https://docs.aws.amazon.com/cloudtrail/",
		"https://docs.aws.amazon.com/cloudwatch/",
		"https://docs.aws.amazon.com/guardduty/",
		"https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs.html",
	}

	d.serviceDocs["apache"] = []string{
		"https://httpd.apache.org/docs/",
		"https://httpd.apache.org/docs/2.4/mod/mod_status.html",
		"https://httpd.apache.org/docs/2.4/logs.html",
	}

	d.serviceDocs["postgresql"] = []string{
		"https://www.postgresql.org/docs/",
		"https://www.postgresql.org/docs/current/runtime-config-logging.html",
		"https://www.postgresql.org/docs/current/monitoring.html",
	}
}

func (d *DocumentationProvider) initializeTroubleshootingGuides() {
	d.troubleshootingGuides["nginx"] = shared.TroubleshootingGuide{
		ServiceName: "nginx",
		CommonIssues: []shared.TroubleshootingIssue{
			{
				Issue: "Nginx fails to start",
				Symptoms: []string{
					"Service fails to start with systemctl",
					"Error messages in systemd journal",
					"Configuration test fails",
				},
				Causes: []string{
					"Syntax errors in configuration files",
					"Port conflicts with other services",
					"Missing dependencies or modules",
					"Permission issues with log directories",
				},
				Solutions: []string{
					"Check configuration syntax: `nginx -t`",
					"Review error logs: `journalctl -u nginx`",
					"Verify port availability: `netstat -tlnp | grep :80`",
					"Check file permissions on log directories",
					"Ensure required modules are installed",
				},
				Prevention: []string{
					"Always test configuration before reloading",
					"Use version control for configuration files",
					"Monitor disk space for log directories",
					"Regularly update Nginx and dependencies",
				},
			},
			{
				Issue: "Stub status module not accessible",
				Symptoms: []string{
					"404 error when accessing /nginx_status",
					"Connection refused errors",
					"Empty response from status endpoint",
				},
				Causes: []string{
					"Stub status module not enabled",
					"Incorrect location block configuration",
					"Access restrictions blocking requests",
					"Nginx not reloaded after configuration changes",
				},
				Solutions: []string{
					"Verify stub status module is compiled: `nginx -V 2>&1 | grep -o with-http_stub_status_module`",
					"Check location block configuration",
					"Ensure allow/deny directives are correct",
					"Reload Nginx: `systemctl reload nginx`",
					"Test from localhost: `curl http://localhost/nginx_status`",
				},
				Prevention: []string{
					"Document configuration changes",
					"Test status endpoint after changes",
					"Use monitoring to detect issues early",
				},
			},
		},
		DiagnosticCommands: []string{
			"nginx -t",
			"systemctl status nginx",
			"journalctl -u nginx -f",
			"ps aux | grep nginx",
			"netstat -tlnp | grep nginx",
			"curl -I http://localhost/nginx_status",
		},
		LogLocations: []string{
			"/var/log/nginx/access.log",
			"/var/log/nginx/error.log",
			"/var/log/syslog",
			"/var/log/messages",
		},
		SupportResources: []string{
			"https://nginx.org/en/docs/",
			"https://nginx.org/en/support.html",
			"https://serverfault.com/questions/tagged/nginx",
			"https://stackoverflow.com/questions/tagged/nginx",
		},
	}

	d.troubleshootingGuides["mysql"] = shared.TroubleshootingGuide{
		ServiceName: "mysql",
		CommonIssues: []shared.TroubleshootingIssue{
			{
				Issue: "MySQL service fails to start",
				Symptoms: []string{
					"Service fails to start with systemctl",
					"Error messages in MySQL error log",
					"Port 3306 not listening",
				},
				Causes: []string{
					"Corrupted data directory",
					"Insufficient disk space",
					"Configuration errors",
					"Permission issues",
					"Port conflicts",
				},
				Solutions: []string{
					"Check MySQL error log: `tail -f /var/log/mysql/error.log`",
					"Verify disk space: `df -h`",
					"Check configuration: `mysqld --help --verbose`",
					"Fix permissions: `chown -R mysql:mysql /var/lib/mysql`",
					"Check for port conflicts: `netstat -tlnp | grep 3306`",
				},
				Prevention: []string{
					"Regular backups of data directory",
					"Monitor disk space",
					"Test configuration changes",
					"Use proper file permissions",
				},
			},
		},
		DiagnosticCommands: []string{
			"systemctl status mysql",
			"mysql -u root -p -e \"SHOW PROCESSLIST;\"",
			"mysql -u root -p -e \"SHOW VARIABLES LIKE 'max_connections';\"",
			"mysql -u root -p -e \"SHOW ENGINE INNODB STATUS;\"",
			"tail -f /var/log/mysql/error.log",
			"netstat -tlnp | grep 3306",
		},
		LogLocations: []string{
			"/var/log/mysql/error.log",
			"/var/log/mysql/slow.log",
			"/var/log/mysql/general.log",
			"/var/log/syslog",
		},
		SupportResources: []string{
			"https://dev.mysql.com/doc/",
			"https://dev.mysql.com/doc/refman/8.0/en/error-handling.html",
			"https://stackoverflow.com/questions/tagged/mysql",
			"https://dba.stackexchange.com/questions/tagged/mysql",
		},
	}

	d.troubleshootingGuides["aws"] = shared.TroubleshootingGuide{
		ServiceName: "aws",
		CommonIssues: []shared.TroubleshootingIssue{
			{
				Issue: "Access denied errors",
				Symptoms: []string{
					"403 Forbidden errors",
					"Access denied messages",
					"Permission denied errors",
				},
				Causes: []string{
					"Insufficient IAM permissions",
					"Incorrect credentials",
					"Expired access keys",
					"Resource-based policies",
					"Service-specific permissions",
				},
				Solutions: []string{
					"Verify IAM user permissions",
					"Check resource-based policies",
					"Rotate access keys",
					"Use IAM roles when possible",
					"Review CloudTrail logs for denied actions",
				},
				Prevention: []string{
					"Use least privilege principle",
					"Regular access reviews",
					"Use IAM roles instead of access keys",
					"Enable MFA",
					"Monitor access patterns",
				},
			},
		},
		DiagnosticCommands: []string{
			"aws sts get-caller-identity",
			"aws iam list-attached-user-policies --user-name <username>",
			"aws cloudtrail describe-trails",
			"aws logs describe-log-groups",
			"aws s3 ls s3://your-log-bucket",
			"aws cloudwatch get-metric-statistics --namespace AWS/CloudTrail",
		},
		LogLocations: []string{
			"CloudTrail logs in S3",
			"CloudWatch Logs",
			"VPC Flow Logs",
			"Service-specific logs",
		},
		SupportResources: []string{
			"https://docs.aws.amazon.com/",
			"https://aws.amazon.com/support/",
			"https://repost.aws/",
			"https://stackoverflow.com/questions/tagged/amazon-web-services",
		},
	}
}

func (d *DocumentationProvider) SearchServiceDocs(serviceName, query, docType string) (shared.CallToolResult, error) {
	service := strings.ToLower(serviceName)
	docs, exists := d.serviceDocs[service]
	if !exists {
		availableServices := make([]string, 0, len(d.serviceDocs))
		for name := range d.serviceDocs {
			availableServices = append(availableServices, name)
		}
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Documentation for '%s' not found. Available services: %s", serviceName, strings.Join(availableServices, ", ")),
				},
			},
			IsError: true,
		}, nil
	}

	// For now, return a simple response with available documentation
	// In a real implementation, this would perform web scraping or API calls
	searchResults := fmt.Sprintf(`# Search Results for "%s" in %s Documentation

## Available Documentation
%s

## Search Query: "%s"
%s

## Additional Resources
%s`,
		query,
		strings.ToUpper(serviceName),
		formatLinks(docs),
		query,
		"Note: This is a simplified response. In a production environment, this would perform actual web search and content extraction.",
		formatLinks(docs))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: searchResults,
			},
		},
	}, nil
}

func (d *DocumentationProvider) GetTroubleshootingGuide(serviceName, issue string) (shared.CallToolResult, error) {
	service := strings.ToLower(serviceName)
	guide, exists := d.troubleshootingGuides[service]
	if !exists {
		availableServices := make([]string, 0, len(d.troubleshootingGuides))
		for name := range d.troubleshootingGuides {
			availableServices = append(availableServices, name)
		}
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Troubleshooting guide for '%s' not found. Available services: %s", serviceName, strings.Join(availableServices, ", ")),
				},
			},
			IsError: true,
		}, nil
	}

	var relevantIssues []shared.TroubleshootingIssue
	if issue != "" {
		// Filter by specific issue if provided
		for _, item := range guide.CommonIssues {
			if strings.Contains(strings.ToLower(item.Issue), strings.ToLower(issue)) {
				relevantIssues = append(relevantIssues, item)
			} else {
				// Check symptoms
				for _, symptom := range item.Symptoms {
					if strings.Contains(strings.ToLower(symptom), strings.ToLower(issue)) {
						relevantIssues = append(relevantIssues, item)
						break
					}
				}
			}
		}
	} else {
		relevantIssues = guide.CommonIssues
	}

	troubleshootingGuide := fmt.Sprintf(`# %s Troubleshooting Guide

%s

## Common Issues

%s

## Diagnostic Commands

%s

## Log Locations

%s

## Support Resources

%s`,
		strings.ToUpper(serviceName),
		func() string {
			if issue != "" {
				return fmt.Sprintf("## Issue: %s", issue)
			}
			return ""
		}(),
		formatTroubleshootingIssues(relevantIssues),
		"```bash\n"+strings.Join(guide.DiagnosticCommands, "\n")+"\n```",
		formatList(guide.LogLocations),
		formatLinks(guide.SupportResources))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: troubleshootingGuide,
			},
		},
	}, nil
}

func (d *DocumentationProvider) GetLatestDocs(serviceName, docType string) (shared.CallToolResult, error) {
	service := strings.ToLower(serviceName)
	docs, exists := d.serviceDocs[service]
	if !exists {
		availableServices := make([]string, 0, len(d.serviceDocs))
		for name := range d.serviceDocs {
			availableServices = append(availableServices, name)
		}
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Documentation for '%s' not found. Available services: %s", serviceName, strings.Join(availableServices, ", ")),
				},
			},
			IsError: true,
		}, nil
	}

	latestDocs := fmt.Sprintf(`# Latest Documentation for %s

## Official Documentation
%s

## Quick Reference

%s`,
		strings.ToUpper(serviceName),
		formatLinks(docs),
		d.getQuickReference(service, docType))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: latestDocs,
			},
		},
	}, nil
}

func formatTroubleshootingIssues(issues []shared.TroubleshootingIssue) string {
	var result strings.Builder
	for _, item := range issues {
		result.WriteString(fmt.Sprintf("\n### %s\n\n#### Symptoms\n%s\n\n#### Possible Causes\n%s\n\n#### Solutions\n%s\n\n#### Prevention\n%s\n\n",
			item.Issue,
			formatList(item.Symptoms),
			formatList(item.Causes),
			formatList(item.Solutions),
			formatList(item.Prevention)))
	}
	return result.String()
}

func (d *DocumentationProvider) getQuickReference(service, docType string) string {
	switch service {
	case "nginx":
		return `### Configuration Files
- Main config: /etc/nginx/nginx.conf
- Site configs: /etc/nginx/sites-available/
- Logs: /var/log/nginx/

### Key Commands
- Test config: nginx -t
- Reload: systemctl reload nginx
- Status: systemctl status nginx`
	case "mysql":
		return `### Configuration Files
- Main config: /etc/mysql/mysql.conf.d/mysqld.cnf
- Logs: /var/log/mysql/

### Key Commands
- Connect: mysql -u root -p
- Status: systemctl status mysql
- Logs: tail -f /var/log/mysql/error.log`
	case "aws":
		return `### Key Commands
- Identity: aws sts get-caller-identity
- Regions: aws ec2 describe-regions
- CloudTrail: aws cloudtrail describe-trails

### Important Services
- CloudTrail: API activity logging
- CloudWatch: Metrics and logs
- VPC Flow Logs: Network traffic`
	default:
		return `### General Tips
- Check service status regularly
- Monitor logs for errors
- Keep configurations documented
- Test changes in development first`
	}
}
