package services

import (
	"fmt"
	"strings"

	"elastic-integration-docs-mcp/internal/shared"
)

type ValidationProvider struct {
	serviceValidators map[string]func(string, string) shared.ValidationResult
}

func NewValidationProvider() *ValidationProvider {
	provider := &ValidationProvider{
		serviceValidators: make(map[string]func(string, string) shared.ValidationResult),
	}
	provider.initializeValidators()
	return provider
}

func (v *ValidationProvider) initializeValidators() {
	// Nginx configuration validator
	v.serviceValidators["nginx"] = func(config, configType string) shared.ValidationResult {
		result := shared.ValidationResult{
			IsValid:     true,
			Errors:      []shared.ValidationError{},
			Warnings:    []shared.ValidationWarning{},
			Suggestions: []shared.ValidationSuggestion{},
		}

		// Basic syntax validation
		if configType == "conf" || configType == "nginx" {
			v.validateNginxSyntax(config, &result)
			v.validateNginxSecurity(config, &result)
			v.validateNginxPerformance(config, &result)
		}

		return result
	}

	// MySQL configuration validator
	v.serviceValidators["mysql"] = func(config, configType string) shared.ValidationResult {
		result := shared.ValidationResult{
			IsValid:     true,
			Errors:      []shared.ValidationError{},
			Warnings:    []shared.ValidationWarning{},
			Suggestions: []shared.ValidationSuggestion{},
		}

		if configType == "cnf" || configType == "ini" {
			v.validateMySQLSyntax(config, &result)
			v.validateMySQLSecurity(config, &result)
			v.validateMySQLPerformance(config, &result)
		}

		return result
	}

	// AWS configuration validator
	v.serviceValidators["aws"] = func(config, configType string) shared.ValidationResult {
		result := shared.ValidationResult{
			IsValid:     true,
			Errors:      []shared.ValidationError{},
			Warnings:    []shared.ValidationWarning{},
			Suggestions: []shared.ValidationSuggestion{},
		}

		if configType == "yaml" || configType == "yml" {
			v.validateAWSYaml(config, &result)
		} else if configType == "json" {
			v.validateAWSJson(config, &result)
		}

		return result
	}

	// Apache configuration validator
	v.serviceValidators["apache"] = func(config, configType string) shared.ValidationResult {
		result := shared.ValidationResult{
			IsValid:     true,
			Errors:      []shared.ValidationError{},
			Warnings:    []shared.ValidationWarning{},
			Suggestions: []shared.ValidationSuggestion{},
		}

		if configType == "conf" || configType == "apache" {
			v.validateApacheSyntax(config, &result)
			v.validateApacheSecurity(config, &result)
			v.validateApachePerformance(config, &result)
		}

		return result
	}

	// PostgreSQL configuration validator
	v.serviceValidators["postgresql"] = func(config, configType string) shared.ValidationResult {
		result := shared.ValidationResult{
			IsValid:     true,
			Errors:      []shared.ValidationError{},
			Warnings:    []shared.ValidationWarning{},
			Suggestions: []shared.ValidationSuggestion{},
		}

		if configType == "conf" || configType == "postgresql" {
			v.validatePostgreSQLSyntax(config, &result)
			v.validatePostgreSQLSecurity(config, &result)
			v.validatePostgreSQLPerformance(config, &result)
		}

		return result
	}
}

func (v *ValidationProvider) validateNginxSyntax(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1

		// Check for basic syntax issues
		if strings.Contains(line, "{") && !strings.Contains(line, "}") {
			// Look for matching closing brace
			braceCount := 1
			for j := i + 1; j < len(lines); j++ {
				nextLine := lines[j]
				braceCount += strings.Count(nextLine, "{")
				braceCount -= strings.Count(nextLine, "}")
				if braceCount == 0 {
					break
				}
			}
			if braceCount > 0 {
				result.Errors = append(result.Errors, shared.ValidationError{
					Type:     "syntax",
					Message:  "Unclosed brace",
					Line:     &lineNum,
					Severity: "error",
				})
				result.IsValid = false
			}
		}

		// Check for common syntax errors
		if strings.Contains(line, "server_name") && !strings.Contains(line, ";") {
			result.Errors = append(result.Errors, shared.ValidationError{
				Type:     "syntax",
				Message:  "Missing semicolon after server_name directive",
				Line:     &lineNum,
				Severity: "error",
			})
			result.IsValid = false
		}
	}
}

func (v *ValidationProvider) validateNginxSecurity(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1

		// Check for security issues
		if strings.Contains(line, "server_tokens") && !strings.Contains(line, "off") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "security",
				Message:    "Consider hiding server version with server_tokens off",
				Line:       &lineNum,
				Suggestion: "Add \"server_tokens off;\" to hide server version",
			})
		}

		if strings.Contains(line, "location /") && !strings.Contains(line, "deny all") && !strings.Contains(line, "allow") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "security",
				Message:    "Consider adding access controls to location blocks",
				Line:       &lineNum,
				Suggestion: "Add appropriate allow/deny directives",
			})
		}

		if strings.Contains(line, "ssl_protocols") && strings.Contains(line, "SSLv3") {
			result.Errors = append(result.Errors, shared.ValidationError{
				Type:     "security",
				Message:  "SSLv3 is insecure and should not be used",
				Line:     &lineNum,
				Severity: "error",
			})
			result.IsValid = false
		}
	}
}

func (v *ValidationProvider) validateNginxPerformance(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for performance issues
		if strings.Contains(line, "gzip") && !strings.Contains(line, "gzip on") {
			result.Suggestions = append(result.Suggestions, shared.ValidationSuggestion{
				Type:       "performance",
				Message:    "Consider enabling gzip compression",
				Suggestion: "Add \"gzip on;\" to enable compression",
				Impact:     "medium",
			})
		}

		if strings.Contains(line, "expires") && !strings.Contains(line, "expires") {
			result.Suggestions = append(result.Suggestions, shared.ValidationSuggestion{
				Type:       "performance",
				Message:    "Consider adding cache headers for static content",
				Suggestion: "Add \"expires 1y;\" for static assets",
				Impact:     "high",
			})
		}
	}
}

func (v *ValidationProvider) validateMySQLSyntax(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1

		// Check for basic syntax issues
		if strings.Contains(line, "=") && !strings.Contains(line, "=") {
			result.Errors = append(result.Errors, shared.ValidationError{
				Type:     "syntax",
				Message:  "Invalid configuration syntax",
				Line:     &lineNum,
				Severity: "error",
			})
			result.IsValid = false
		}

		// Check for common MySQL configuration issues
		if strings.Contains(line, "innodb_buffer_pool_size") && strings.Contains(line, "G") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "performance",
				Message:    "Very large buffer pool size may cause memory issues",
				Line:       &lineNum,
				Suggestion: "Consider using 70-80% of available RAM",
			})
		}
	}
}

func (v *ValidationProvider) validateMySQLSecurity(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1

		// Check for security issues
		if strings.Contains(line, "skip-grant-tables") {
			result.Errors = append(result.Errors, shared.ValidationError{
				Type:     "security",
				Message:  "skip-grant-tables is extremely dangerous and should not be used in production",
				Line:     &lineNum,
				Severity: "error",
			})
			result.IsValid = false
		}

		if strings.Contains(line, "local-infile") && strings.Contains(line, "1") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "security",
				Message:    "local-infile=1 allows loading local files and poses security risks",
				Line:       &lineNum,
				Suggestion: "Set local-infile=0 unless specifically needed",
			})
		}
	}
}

func (v *ValidationProvider) validateMySQLPerformance(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for performance issues
		if strings.Contains(line, "query_cache") && strings.Contains(line, "0") {
			result.Suggestions = append(result.Suggestions, shared.ValidationSuggestion{
				Type:       "performance",
				Message:    "Query cache is disabled",
				Suggestion: "Consider enabling query cache for read-heavy workloads",
				Impact:     "medium",
			})
		}

		if strings.Contains(line, "innodb_flush_log_at_trx_commit") && strings.Contains(line, "1") {
			result.Suggestions = append(result.Suggestions, shared.ValidationSuggestion{
				Type:       "performance",
				Message:    "innodb_flush_log_at_trx_commit=1 provides maximum durability but may impact performance",
				Suggestion: "Consider setting to 2 for better performance with some durability trade-off",
				Impact:     "high",
			})
		}
	}
}

func (v *ValidationProvider) validateAWSYaml(config string, result *shared.ValidationResult) {
	// Basic YAML validation - in a real implementation, this would use a YAML parser
	if strings.Contains(config, "aws:") && strings.Contains(config, "region:") {
		// Check for common AWS configuration issues
		if strings.Contains(config, "access_key_id:") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "security",
				Message:    "Access keys in configuration files are not recommended",
				Suggestion: "Use IAM roles or environment variables instead",
			})
		}
	}
}

func (v *ValidationProvider) validateAWSJson(config string, result *shared.ValidationResult) {
	// Basic JSON validation - in a real implementation, this would use a JSON parser
	if strings.Contains(config, "\"aws\"") && strings.Contains(config, "\"region\"") {
		// Check for common AWS configuration issues
		if strings.Contains(config, "\"access_key_id\"") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "security",
				Message:    "Access keys in configuration files are not recommended",
				Suggestion: "Use IAM roles or environment variables instead",
			})
		}
	}
}

func (v *ValidationProvider) validateApacheSyntax(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1

		// Check for basic syntax issues
		if strings.HasPrefix(line, "<") && !strings.Contains(line, ">") {
			result.Errors = append(result.Errors, shared.ValidationError{
				Type:     "syntax",
				Message:  "Unclosed XML tag",
				Line:     &lineNum,
				Severity: "error",
			})
			result.IsValid = false
		}

		// Check for common Apache configuration issues
		if strings.Contains(line, "ServerRoot") && !strings.Contains(line, "/") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "best_practice",
				Message:    "ServerRoot should be an absolute path",
				Line:       &lineNum,
				Suggestion: "Use absolute path like /etc/apache2",
			})
		}
	}
}

func (v *ValidationProvider) validateApacheSecurity(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1

		// Check for security issues
		if strings.Contains(line, "ServerTokens") && !strings.Contains(line, "Prod") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "security",
				Message:    "Consider hiding server version with ServerTokens Prod",
				Line:       &lineNum,
				Suggestion: "Add \"ServerTokens Prod\" to hide server version",
			})
		}

		if strings.Contains(line, "Options") && strings.Contains(line, "Indexes") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "security",
				Message:    "Directory indexing may expose sensitive files",
				Line:       &lineNum,
				Suggestion: "Remove Indexes from Options directive",
			})
		}
	}
}

func (v *ValidationProvider) validateApachePerformance(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for performance issues
		if strings.Contains(line, "KeepAlive") && !strings.Contains(line, "On") {
			result.Suggestions = append(result.Suggestions, shared.ValidationSuggestion{
				Type:       "performance",
				Message:    "Consider enabling KeepAlive for better performance",
				Suggestion: "Add \"KeepAlive On\" to enable persistent connections",
				Impact:     "medium",
			})
		}

		if strings.Contains(line, "LoadModule") && strings.Contains(line, "mod_deflate") {
			result.Suggestions = append(result.Suggestions, shared.ValidationSuggestion{
				Type:       "performance",
				Message:    "Consider enabling compression with mod_deflate",
				Suggestion: "Add compression rules for text content",
				Impact:     "high",
			})
		}
	}
}

func (v *ValidationProvider) validatePostgreSQLSyntax(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1

		// Check for basic syntax issues
		if strings.Contains(line, "=") && !strings.Contains(line, "=") {
			result.Errors = append(result.Errors, shared.ValidationError{
				Type:     "syntax",
				Message:  "Invalid configuration syntax",
				Line:     &lineNum,
				Severity: "error",
			})
			result.IsValid = false
		}

		// Check for common PostgreSQL configuration issues
		if strings.Contains(line, "shared_buffers") && strings.Contains(line, "MB") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "performance",
				Message:    "shared_buffers is very small",
				Line:       &lineNum,
				Suggestion: "Consider increasing to 25% of available RAM",
			})
		}
	}
}

func (v *ValidationProvider) validatePostgreSQLSecurity(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1

		// Check for security issues
		if strings.Contains(line, "ssl") && strings.Contains(line, "off") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "security",
				Message:    "SSL is disabled",
				Line:       &lineNum,
				Suggestion: "Enable SSL for secure connections",
			})
		}

		if strings.Contains(line, "log_statement") && strings.Contains(line, "none") {
			result.Warnings = append(result.Warnings, shared.ValidationWarning{
				Type:       "security",
				Message:    "Statement logging is disabled",
				Line:       &lineNum,
				Suggestion: "Consider enabling logging for security monitoring",
			})
		}
	}
}

func (v *ValidationProvider) validatePostgreSQLPerformance(config string, result *shared.ValidationResult) {
	lines := strings.Split(config, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for performance issues
		if strings.Contains(line, "checkpoint_completion_target") && strings.Contains(line, "0.5") {
			result.Suggestions = append(result.Suggestions, shared.ValidationSuggestion{
				Type:       "performance",
				Message:    "checkpoint_completion_target is at default value",
				Suggestion: "Consider increasing to 0.9 for better performance",
				Impact:     "medium",
			})
		}

		if strings.Contains(line, "random_page_cost") && strings.Contains(line, "4.0") {
			result.Suggestions = append(result.Suggestions, shared.ValidationSuggestion{
				Type:       "performance",
				Message:    "random_page_cost is high for SSD storage",
				Suggestion: "Consider reducing to 1.1 for SSD storage",
				Impact:     "high",
			})
		}
	}
}

func (v *ValidationProvider) ValidateConfiguration(serviceName, configuration, configType string) (shared.CallToolResult, error) {
	service := strings.ToLower(serviceName)
	validator, exists := v.serviceValidators[service]
	if !exists {
		availableServices := make([]string, 0, len(v.serviceValidators))
		for name := range v.serviceValidators {
			availableServices = append(availableServices, name)
		}
		return shared.CallToolResult{
			Content: []shared.ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Validation for '%s' not supported. Available services: %s", serviceName, strings.Join(availableServices, ", ")),
				},
			},
			IsError: true,
		}, nil
	}

	result := validator(configuration, configType)

	validationReport := fmt.Sprintf(`# Configuration Validation Report for %s

## Overall Status
%s

## Errors
%s

## Warnings
%s

## Suggestions
%s

## Summary
- **Total Errors**: %d
- **Total Warnings**: %d
- **Total Suggestions**: %d
- **Configuration Status**: %s

%s`,
		strings.ToUpper(serviceName),
		func() string {
			if result.IsValid {
				return "✅ **VALID**"
			}
			return "❌ **INVALID**"
		}(),
		formatValidationErrors(result.Errors),
		formatValidationWarnings(result.Warnings),
		formatValidationSuggestions(result.Suggestions),
		len(result.Errors),
		len(result.Warnings),
		len(result.Suggestions),
		func() string {
			if result.IsValid {
				return "Valid"
			}
			return "Invalid"
		}(),
		func() string {
			if result.IsValid {
				return "✅ Your configuration is valid and ready to use."
			}
			return "❌ Please fix the errors before using this configuration."
		}())

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: validationReport,
			},
		},
	}, nil
}

func formatValidationErrors(errors []shared.ValidationError) string {
	if len(errors) == 0 {
		return "No errors found."
	}

	var result strings.Builder
	for _, error := range errors {
		result.WriteString(fmt.Sprintf("\n### %s: %s\n- **Type**: %s\n- **Line**: %s\n- **Column**: %s\n",
			strings.ToUpper(error.Severity),
			error.Message,
			error.Type,
			func() string {
				if error.Line != nil {
					return fmt.Sprintf("%d", *error.Line)
				}
				return "N/A"
			}(),
			func() string {
				if error.Column != nil {
					return fmt.Sprintf("%d", *error.Column)
				}
				return "N/A"
			}()))
	}
	return result.String()
}

func formatValidationWarnings(warnings []shared.ValidationWarning) string {
	if len(warnings) == 0 {
		return "No warnings found."
	}

	var result strings.Builder
	for _, warning := range warnings {
		result.WriteString(fmt.Sprintf("\n### %s: %s\n- **Line**: %s\n- **Column**: %s\n",
			strings.ToUpper(warning.Type),
			warning.Message,
			func() string {
				if warning.Line != nil {
					return fmt.Sprintf("%d", *warning.Line)
				}
				return "N/A"
			}(),
			func() string {
				if warning.Column != nil {
					return fmt.Sprintf("%d", *warning.Column)
				}
				return "N/A"
			}()))
		if warning.Suggestion != "" {
			result.WriteString(fmt.Sprintf("- **Suggestion**: %s\n", warning.Suggestion))
		}
	}
	return result.String()
}

func formatValidationSuggestions(suggestions []shared.ValidationSuggestion) string {
	if len(suggestions) == 0 {
		return "No suggestions available."
	}

	var result strings.Builder
	for _, suggestion := range suggestions {
		result.WriteString(fmt.Sprintf("\n### %s: %s\n- **Suggestion**: %s\n- **Impact**: %s\n",
			strings.ToUpper(suggestion.Type),
			suggestion.Message,
			suggestion.Suggestion,
			strings.ToUpper(suggestion.Impact)))
	}
	return result.String()
}
