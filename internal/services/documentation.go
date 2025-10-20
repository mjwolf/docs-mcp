package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"elastic-integration-docs-mcp/internal/config"
	"elastic-integration-docs-mcp/internal/shared"
)

type DocumentationProvider struct {
	configLoader *config.ConfigLoader
	httpClient   *http.Client
}

func NewDocumentationProvider(configDir string) *DocumentationProvider {
	configLoader := config.NewConfigLoader(configDir)
	if err := configLoader.LoadAllServices(); err != nil {
		// In a real implementation, you might want to handle this error differently
		// For now, we'll create an empty loader
		configLoader = config.NewConfigLoader(configDir)
	}

	return &DocumentationProvider{
		configLoader: configLoader,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (d *DocumentationProvider) SearchDocumentation(searchTerm, serviceName string) (shared.CallToolResult, error) {
	serviceConfig, err := d.configLoader.GetServiceConfig(serviceName)
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

	// For Apache, perform a Google search restricted to documentation sites
	if strings.ToLower(serviceName) == "apache" {
		return d.performWebSearch(searchTerm, "httpd.apache.org/docs/")
	}

	// For other services, search their documentation sites
	if len(serviceConfig.DocumentationSites) > 0 {
		// Use the first documentation site for the search
		site := serviceConfig.DocumentationSites[0]
		return d.performWebSearch(searchTerm, site)
	}

	// Fallback: return available documentation sites
	searchResults := fmt.Sprintf(`# Search Results for "%s" in %s Documentation

## Available Documentation Sites
%s

## Search Query: "%s"
Note: This is a simplified response. In a production environment, this would perform actual web search and content extraction.

## Manual Search Instructions
To search for "%s" in %s documentation, visit the following sites:
%s`,
		searchTerm,
		strings.ToUpper(serviceName),
		formatList(serviceConfig.DocumentationSites),
		searchTerm,
		searchTerm,
		strings.ToUpper(serviceName),
		formatList(serviceConfig.DocumentationSites))

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: searchResults,
			},
		},
	}, nil
}

func (d *DocumentationProvider) performWebSearch(searchTerm, site string) (shared.CallToolResult, error) {
	// Construct Google search URL with site restriction
	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s+site:%s",
		url.QueryEscape(searchTerm),
		url.QueryEscape(site))

	// In a real implementation, you would:
	// 1. Make HTTP request to Google search
	// 2. Parse the HTML response
	// 3. Extract search results
	// 4. Return formatted results

	// For now, return a placeholder response
	searchResults := fmt.Sprintf(`# Search Results for "%s" on %s

## Search URL
%s

## Note
This is a placeholder response. In a production environment, this would:
1. Perform actual web search using the URL above
2. Parse and extract search results
3. Return formatted documentation links and snippets

## Manual Search
You can manually search by visiting: %s`,
		searchTerm,
		site,
		searchURL,
		searchURL)

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: searchResults,
			},
		},
	}, nil
}

func (d *DocumentationProvider) GetTroubleshootingHelp(serviceName string) (shared.CallToolResult, error) {
	serviceConfig, err := d.configLoader.GetServiceConfig(serviceName)
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

	// Format troubleshooting issues as JSON-like structure as shown in requirements
	var issuesJSON strings.Builder
	issuesJSON.WriteString("[\n")

	for i, issue := range serviceConfig.Troubleshooting.CommonIssues {
		issuesJSON.WriteString(fmt.Sprintf("  {\n    \"Issue\":    \"%s\",\n    \"Solution\": \"%s\"\n  }",
			issue.Issue, issue.Solution))
		if i < len(serviceConfig.Troubleshooting.CommonIssues)-1 {
			issuesJSON.WriteString(",")
		}
		issuesJSON.WriteString("\n")
	}

	issuesJSON.WriteString("]")

	return shared.CallToolResult{
		Content: []shared.ToolContent{
			{
				Type: "text",
				Text: issuesJSON.String(),
			},
		},
	}, nil
}
