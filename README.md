# Elastic Integration Documentation MCP Server (Go)

A Model Context Protocol (MCP) server implemented in Go that provides comprehensive documentation and setup guidance for Elastic integrations. This server helps generate accurate, service-specific documentation for Elastic integrations by providing detailed information about service setup, configuration, and troubleshooting.

## Features

- **Service Information**: Get comprehensive details about supported services including requirements, capabilities, and supported versions
- **Setup Instructions**: Step-by-step setup guides with platform-specific instructions
- **Configuration Examples**: Ready-to-use configuration templates and examples
- **Documentation Search**: Search through official and community documentation
- **Configuration Validation**: Validate service configurations and get improvement suggestions
- **Integration Details**: Get detailed information about Elastic integrations including data streams and field mappings
- **Troubleshooting Guides**: Comprehensive troubleshooting guides for common issues
- **Provider Agnostic**: Works with any MCP-compatible LLM provider (Bedrock, Google AI, etc.)

## Supported Services

- **Web Servers**: Nginx, Apache
- **Databases**: MySQL, PostgreSQL
- **Cloud Services**: AWS (CloudTrail, CloudWatch, VPC Flow Logs, etc.)
- **And more...**

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd docs-mcp
```

2. Install dependencies:
```bash
go mod tidy
```

## Usage

### Building the Server

```bash
go build -o elastic-integration-docs-mcp cmd/server/main.go
```

### Running the Server

```bash
./elastic-integration-docs-mcp
```

The server runs on stdio and can be connected to by MCP-compatible clients.

### Available Tools

#### `get_service_info`
Get comprehensive information about a service including requirements, capabilities, and supported versions.

**Parameters:**
- `serviceName` (string): Name of the service (e.g., nginx, mysql, aws)

#### `get_setup_instructions`
Get step-by-step setup instructions for a service with platform-specific guidance.

**Parameters:**
- `serviceName` (string): Name of the service
- `platform` (string, optional): Target platform (e.g., ubuntu, centos, docker)
- `version` (string, optional): Service version

#### `get_configuration_examples`
Get configuration examples and templates for a service.

**Parameters:**
- `serviceName` (string): Name of the service
- `configType` (string, optional): Type of configuration (e.g., logs, metrics, security)

#### `search_service_docs`
Search for service-specific documentation and guides.

**Parameters:**
- `serviceName` (string): Name of the service
- `query` (string): Search query
- `docType` (string, optional): Type of documentation (official, community, troubleshooting)

#### `validate_configuration`
Validate service configuration and provide suggestions for improvements.

**Parameters:**
- `serviceName` (string): Name of the service
- `configuration` (string): Configuration content to validate
- `configType` (string): Type of configuration (yaml, json, conf, etc.)

#### `get_integration_details`
Get details about Elastic integration including data streams and field mappings.

**Parameters:**
- `integrationName` (string): Name of the Elastic integration

#### `get_troubleshooting_guide`
Get troubleshooting guide for common issues with a service.

**Parameters:**
- `serviceName` (string): Name of the service
- `issue` (string, optional): Specific issue or error message

#### `get_service_categories`
Get list of available service categories and services within each category.

#### `get_latest_docs`
Get the latest official documentation for a service.

**Parameters:**
- `serviceName` (string): Name of the service
- `docType` (string, optional): Type of documentation (installation, configuration, api)

## Integration with Elastic Package

This MCP server is designed to work with the `elastic-package` LLM agent to help generate documentation for Elastic integrations. The agent can use this server to:

1. Get detailed service information
2. Retrieve setup instructions
3. Validate configurations
4. Search for troubleshooting information
5. Get integration-specific details

## Architecture

The server is built with Go and implements the Model Context Protocol. It consists of several service providers:

- **ServiceInfoProvider**: Manages service metadata and information
- **SetupGuideProvider**: Provides setup instructions and configuration examples
- **DocumentationProvider**: Handles documentation search and troubleshooting guides
- **ValidationProvider**: Validates configurations and provides suggestions
- **IntegrationProvider**: Manages Elastic integration details

## Development

### Prerequisites

- Go 1.21+

### Building

```bash
go build -o elastic-integration-docs-mcp cmd/server/main.go
```

### Testing

```bash
go test ./...
```

### Running Tests

```bash
go test -v ./...
```

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go          # Main server executable
├── internal/
│   ├── mcp/
│   │   ├── types.go         # MCP protocol types
│   │   └── server.go        # MCP server implementation
│   └── services/
│       ├── service_info.go  # Service information provider
│       ├── setup_guide.go   # Setup guide provider
│       ├── documentation.go # Documentation provider
│       ├── validation.go    # Configuration validation
│       └── integration.go   # Integration details provider
├── go.mod                   # Go module file
├── go.sum                   # Go module checksums
└── README.md               # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the Elastic License.

## Support

For issues and questions, please open an issue in the repository.
