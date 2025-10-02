package services

import (
	"fmt"
	"strings"

	"elastic-integration-docs-mcp/internal/shared"
)

type SetupGuideProvider struct {
	setupGuides map[string]shared.SetupInstructions
}

func NewSetupGuideProvider() *SetupGuideProvider {
	provider := &SetupGuideProvider{
		setupGuides: make(map[string]shared.SetupInstructions),
	}
	provider.initializeSetupGuides()
	return provider
}

func (s *SetupGuideProvider) initializeSetupGuides() {
	// Nginx setup guide
	s.setupGuides["nginx"] = shared.SetupInstructions{
		ServiceName: "nginx",
		Prerequisites: []string{
			"Root or sudo access to the server",
			"Package manager (apt, yum, dnf, etc.)",
			"Basic understanding of web server configuration",
		},
		InstallationSteps: []shared.InstallationStep{
			{
				Step:        1,
				Title:       "Install Nginx",
				Description: "Install Nginx using your system package manager",
				Commands: []string{
					"# Ubuntu/Debian",
					"sudo apt update",
					"sudo apt install nginx",
					"",
					"# CentOS/RHEL/Rocky Linux",
					"sudo yum install nginx",
					"# or",
					"sudo dnf install nginx",
					"",
					"# Verify installation",
					"nginx -v",
				},
				Verification: "Check if Nginx is running: `systemctl status nginx`",
			},
			{
				Step:        2,
				Title:       "Enable and Start Nginx",
				Description: "Enable Nginx to start automatically on boot and start the service",
				Commands: []string{
					"sudo systemctl enable nginx",
					"sudo systemctl start nginx",
					"sudo systemctl status nginx",
				},
			},
			{
				Step:        3,
				Title:       "Configure Nginx for Logging",
				Description: "Configure Nginx to generate access and error logs",
				ConfigSnippets: []shared.ConfigSnippet{
					{
						Filename: "/etc/nginx/nginx.conf",
						Content: `http {
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log warn;
    
    # ... rest of configuration
}`,
					},
				},
			},
			{
				Step:        4,
				Title:       "Enable Stub Status Module",
				Description: "Configure Nginx to expose metrics via stub status module",
				ConfigSnippets: []shared.ConfigSnippet{
					{
						Filename: "/etc/nginx/sites-available/default",
						Content: `server {
    listen 80 default_server;
    listen [::]:80 default_server;
    
    # ... other configuration
    
    location /nginx_status {
        stub_status;
        allow 127.0.0.1;
        deny all;
    }
}`,
					},
				},
				Verification: "Test stub status: `curl http://localhost/nginx_status`",
			},
			{
				Step:        5,
				Title:       "Test Configuration and Reload",
				Description: "Test the Nginx configuration and reload the service",
				Commands: []string{
					"sudo nginx -t",
					"sudo systemctl reload nginx",
				},
			},
		},
		PostInstallation: []shared.PostInstallation{
			{
				Title:       "Firewall Configuration",
				Description: "Configure firewall to allow HTTP and HTTPS traffic",
				Steps: []string{
					"Open ports 80 and 443 in your firewall",
					"For UFW: `sudo ufw allow 80` and `sudo ufw allow 443`",
					"For firewalld: `sudo firewall-cmd --permanent --add-service=http` and `sudo firewall-cmd --permanent --add-service=https`",
				},
			},
			{
				Title:       "SSL/TLS Configuration",
				Description: "Configure SSL/TLS for secure connections",
				Steps: []string{
					"Obtain SSL certificates (Let's Encrypt recommended)",
					"Configure SSL in Nginx virtual host",
					"Redirect HTTP to HTTPS",
				},
			},
		},
		Troubleshooting: []shared.TroubleshootingItem{
			{
				Issue:    "Nginx fails to start",
				Solution: "Check configuration syntax with `nginx -t` and review error logs in `/var/log/nginx/error.log`",
			},
			{
				Issue:    "Stub status not accessible",
				Solution: "Verify the location block is properly configured and accessible from localhost only",
			},
			{
				Issue:    "Logs not being generated",
				Solution: "Check file permissions on log directories and ensure Nginx has write access",
			},
		},
	}

	// MySQL setup guide
	s.setupGuides["mysql"] = shared.SetupInstructions{
		ServiceName: "mysql",
		Prerequisites: []string{
			"Root or sudo access to the server",
			"Package manager (apt, yum, dnf, etc.)",
			"Basic understanding of database administration",
		},
		InstallationSteps: []shared.InstallationStep{
			{
				Step:        1,
				Title:       "Install MySQL Server",
				Description: "Install MySQL server using your system package manager",
				Commands: []string{
					"# Ubuntu/Debian",
					"sudo apt update",
					"sudo apt install mysql-server",
					"",
					"# CentOS/RHEL/Rocky Linux",
					"sudo yum install mysql-server",
					"# or",
					"sudo dnf install mysql-server",
					"",
					"# Verify installation",
					"mysql --version",
				},
			},
			{
				Step:        2,
				Title:       "Secure MySQL Installation",
				Description: "Run the MySQL security script to set up basic security",
				Commands: []string{
					"sudo mysql_secure_installation",
				},
				Verification: "Test MySQL connection: `mysql -u root -p`",
			},
			{
				Step:        3,
				Title:       "Enable and Start MySQL",
				Description: "Enable MySQL to start automatically on boot and start the service",
				Commands: []string{
					"sudo systemctl enable mysql",
					"sudo systemctl start mysql",
					"sudo systemctl status mysql",
				},
			},
			{
				Step:        4,
				Title:       "Configure MySQL Logging",
				Description: "Configure MySQL to generate error logs and slow query logs",
				ConfigSnippets: []shared.ConfigSnippet{
					{
						Filename: "/etc/mysql/mysql.conf.d/mysqld.cnf",
						Content: `[mysqld]
# Error logging
log-error = /var/log/mysql/error.log

# Slow query logging
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 2

# General query logging (optional, for debugging)
# general_log = 1
# general_log_file = /var/log/mysql/general.log

# Binary logging (for replication)
log-bin = mysql-bin
binlog_format = ROW`,
					},
				},
			},
			{
				Step:        5,
				Title:       "Create Monitoring User",
				Description: "Create a dedicated user for monitoring with appropriate permissions",
				Commands: []string{
					"mysql -u root -p",
					"CREATE USER 'monitoring'@'localhost' IDENTIFIED BY 'secure_password';",
					"GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'monitoring'@'localhost';",
					"FLUSH PRIVILEGES;",
					"EXIT;",
				},
			},
			{
				Step:        6,
				Title:       "Restart MySQL",
				Description: "Restart MySQL to apply configuration changes",
				Commands: []string{
					"sudo systemctl restart mysql",
					"sudo systemctl status mysql",
				},
			},
		},
		PostInstallation: []shared.PostInstallation{
			{
				Title:       "Performance Tuning",
				Description: "Optimize MySQL performance for your workload",
				Steps: []string{
					"Adjust buffer pool size based on available RAM",
					"Configure query cache if using MySQL 5.7 or earlier",
					"Set appropriate connection limits",
					"Enable query optimization features",
				},
			},
			{
				Title:       "Backup Configuration",
				Description: "Set up automated backups",
				Steps: []string{
					"Configure mysqldump for logical backups",
					"Set up binary log rotation",
					"Implement backup retention policies",
					"Test backup and restore procedures",
				},
			},
		},
		Troubleshooting: []shared.TroubleshootingItem{
			{
				Issue:    "MySQL fails to start",
				Solution: "Check error logs in `/var/log/mysql/error.log` and verify configuration syntax",
			},
			{
				Issue:    "Cannot connect to MySQL",
				Solution: "Verify MySQL service is running and check firewall settings",
			},
			{
				Issue:    "Slow query log not working",
				Solution: "Ensure slow_query_log is enabled and the log file has proper permissions",
			},
		},
	}

	// AWS setup guide
	s.setupGuides["aws"] = shared.SetupInstructions{
		ServiceName: "aws",
		Prerequisites: []string{
			"AWS account with appropriate permissions",
			"AWS CLI installed and configured",
			"Understanding of AWS services and IAM",
		},
		InstallationSteps: []shared.InstallationStep{
			{
				Step:        1,
				Title:       "Install AWS CLI",
				Description: "Install AWS CLI for managing AWS resources",
				Commands: []string{
					"# Using pip",
					"pip install awscli",
					"",
					"# Using package manager (Ubuntu/Debian)",
					"sudo apt install awscli",
					"",
					"# Using package manager (CentOS/RHEL)",
					"sudo yum install awscli",
					"",
					"# Verify installation",
					"aws --version",
				},
			},
			{
				Step:        2,
				Title:       "Configure AWS Credentials",
				Description: "Set up AWS credentials for authentication",
				Commands: []string{
					"aws configure",
					"# Enter your Access Key ID, Secret Access Key, default region, and output format",
				},
				Verification: "Test AWS connection: `aws sts get-caller-identity`",
			},
			{
				Step:        3,
				Title:       "Create IAM User and Policies",
				Description: "Create dedicated IAM user with necessary permissions for monitoring",
				Commands: []string{
					"# Create IAM user",
					"aws iam create-user --user-name elastic-monitoring",
					"",
					"# Attach necessary policies",
					"aws iam attach-user-policy --user-name elastic-monitoring --policy-arn arn:aws:iam::aws:policy/CloudWatchReadOnlyAccess",
					"aws iam attach-user-policy --user-name elastic-monitoring --policy-arn arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
				},
			},
			{
				Step:        4,
				Title:       "Enable CloudTrail",
				Description: "Enable CloudTrail for API activity monitoring",
				Commands: []string{
					"# Create CloudTrail trail",
					"aws cloudtrail create-trail --name elastic-monitoring-trail --s3-bucket-name your-log-bucket",
					"",
					"# Start logging",
					"aws cloudtrail start-logging --name elastic-monitoring-trail",
				},
			},
			{
				Step:        5,
				Title:       "Configure VPC Flow Logs",
				Description: "Enable VPC Flow Logs for network traffic monitoring",
				Commands: []string{
					"# Get VPC ID",
					"aws ec2 describe-vpcs --query \"Vpcs[0].VpcId\" --output text",
					"",
					"# Create flow log",
					"aws ec2 create-flow-logs --resource-type VPC --resource-ids vpc-12345678 --traffic-type ALL --log-destination-type s3 --log-destination arn:aws:s3:::your-log-bucket/vpc-flow-logs/",
				},
			},
		},
		PostInstallation: []shared.PostInstallation{
			{
				Title:       "Cost Optimization",
				Description: "Optimize AWS costs while maintaining monitoring coverage",
				Steps: []string{
					"Set up billing alerts and budgets",
					"Configure log retention policies",
					"Use appropriate log levels to reduce data volume",
					"Monitor CloudWatch costs and usage",
				},
			},
			{
				Title:       "Security Hardening",
				Description: "Enhance security of your AWS monitoring setup",
				Steps: []string{
					"Enable MFA for IAM users",
					"Use least privilege principle for IAM policies",
					"Encrypt log data at rest and in transit",
					"Regularly rotate access keys",
				},
			},
		},
		Troubleshooting: []shared.TroubleshootingItem{
			{
				Issue:    "Access denied errors",
				Solution: "Verify IAM permissions and ensure the user has necessary policies attached",
			},
			{
				Issue:    "CloudTrail not logging",
				Solution: "Check if CloudTrail is enabled and verify S3 bucket permissions",
			},
			{
				Issue:    "High AWS costs",
				Solution: "Review log retention policies and consider using log filters to reduce data volume",
			},
		},
	}
}

func (s *SetupGuideProvider) GetSetupInstructions(serviceName, platform, version string) (shared.CallToolResult, error) {
	guide, exists := s.setupGuides[strings.ToLower(serviceName)]
	if !exists {
		availableServices := make([]string, 0, len(s.setupGuides))
		for name := range s.setupGuides {
			availableServices = append(availableServices, name)
		}
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Setup guide for '%s' not found. Available services: %s", serviceName, strings.Join(availableServices, ", ")),
				},
			},
			IsError: true,
		}, nil
	}

	platformInfo := ""
	if platform != "" {
		platformInfo = fmt.Sprintf("\n**Platform**: %s", platform)
	}
	versionInfo := ""
	if version != "" {
		versionInfo = fmt.Sprintf("\n**Version**: %s", version)
	}

	instructions := fmt.Sprintf(`# %s Setup Instructions%s%s

## Prerequisites
%s

## Installation Steps

%s

## Post-Installation Configuration

%s

## Troubleshooting

%s`,
		strings.ToUpper(guide.ServiceName),
		platformInfo,
		versionInfo,
		formatList(guide.Prerequisites),
		formatInstallationSteps(guide.InstallationSteps),
		formatPostInstallation(guide.PostInstallation),
		formatTroubleshooting(guide.Troubleshooting))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: instructions,
			},
		},
	}, nil
}

func (s *SetupGuideProvider) GetConfigurationExamples(serviceName, configType string) (shared.CallToolResult, error) {
	guide, exists := s.setupGuides[strings.ToLower(serviceName)]
	if !exists {
		availableServices := make([]string, 0, len(s.setupGuides))
		for name := range s.setupGuides {
			availableServices = append(availableServices, name)
		}
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Configuration examples for '%s' not found. Available services: %s", serviceName, strings.Join(availableServices, ", ")),
				},
			},
			IsError: true,
		}, nil
	}

	// Extract configuration snippets from the setup guide
	var configSnippets []shared.ConfigSnippet
	for _, step := range guide.InstallationSteps {
		configSnippets = append(configSnippets, step.ConfigSnippets...)
	}

	if len(configSnippets) == 0 {
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("No configuration examples available for '%s'.", serviceName),
				},
			},
		}, nil
	}

	examples := fmt.Sprintf(`# %s Configuration Examples

%s

## Additional Configuration Tips

%s`,
		strings.ToUpper(serviceName),
		formatConfigSnippets(configSnippets),
		s.getConfigurationTips(serviceName))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: examples,
			},
		},
	}, nil
}

func formatInstallationSteps(steps []shared.InstallationStep) string {
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

func formatPostInstallation(configs []shared.PostInstallation) string {
	var result strings.Builder
	for _, config := range configs {
		result.WriteString(fmt.Sprintf("\n### %s\n%s\n\n%s\n\n", config.Title, config.Description, formatList(config.Steps)))
	}
	return result.String()
}

func formatTroubleshooting(items []shared.TroubleshootingItem) string {
	var result strings.Builder
	for _, item := range items {
		result.WriteString(fmt.Sprintf("\n### %s\n%s\n\n", item.Issue, item.Solution))
	}
	return result.String()
}

func formatConfigSnippets(snippets []shared.ConfigSnippet) string {
	var result strings.Builder
	for _, snippet := range snippets {
		result.WriteString(fmt.Sprintf("\n## %s\n```%s\n%s\n```\n\n",
			snippet.Filename, getFileExtension(snippet.Filename), snippet.Content))
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

func (s *SetupGuideProvider) getConfigurationTips(serviceName string) string {
	switch strings.ToLower(serviceName) {
	case "nginx":
		return `- Always test configuration changes with nginx -t before reloading
- Use separate log files for different virtual hosts
- Enable gzip compression for better performance
- Set appropriate cache headers for static content
- Consider using rate limiting for security`
	case "mysql":
		return `- Monitor slow query log regularly for performance issues
- Set appropriate buffer pool size (70-80% of available RAM)
- Enable binary logging for point-in-time recovery
- Use connection pooling for high-traffic applications
- Regularly analyze and optimize database schema`
	case "aws":
		return `- Use IAM roles instead of access keys when possible
- Enable CloudTrail in all regions you use
- Set up log retention policies to control costs
- Use VPC Flow Logs for network security monitoring
- Enable GuardDuty for threat detection`
	default:
		return `- Always backup configuration files before making changes
- Test configuration changes in a development environment first
- Monitor logs after configuration changes
- Document any custom configurations
- Keep configurations under version control`
	}
}
