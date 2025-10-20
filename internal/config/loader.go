package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ServiceConfig represents the complete configuration for a service
type ServiceConfig struct {
	ServiceName             string                  `yaml:"service_name"`
	Title                   string                  `yaml:"title"`
	Description             string                  `yaml:"description"`
	ServiceInfo             ServiceInfo             `yaml:"service_info"`
	SetupInstructions       SetupInstructions       `yaml:"setup_instructions"`
	KibanaSetupInstructions KibanaSetupInstructions `yaml:"kibana_setup_instructions"`
	Troubleshooting         Troubleshooting         `yaml:"troubleshooting"`
	ValidationSteps         ValidationSteps         `yaml:"validation_steps"`
	DocumentationSites      []string                `yaml:"documentation_sites"`
}

// ServiceInfo represents service information for get_service_info tool
type ServiceInfo struct {
	CommonUseCases        []string              `yaml:"common_use_cases"`
	DataTypesCollected    []string              `yaml:"data_types_collected"`
	Compatibility         Compatibility         `yaml:"compatibility"`
	ScalingAndPerformance ScalingAndPerformance `yaml:"scaling_and_performance"`
}

// Compatibility represents compatibility information
type Compatibility struct {
	ElasticStackVersions []string `yaml:"elastic_stack_versions"`
	ServiceVersions      []string `yaml:"service_versions"`
}

// ScalingAndPerformance represents scaling and performance information
type ScalingAndPerformance struct {
	Description             string   `yaml:"description"`
	PerformanceExpectations []string `yaml:"performance_expectations"`
	ScalingGuidance         []string `yaml:"scaling_guidance"`
}

// SetupInstructions represents setup instructions for get_service_setup_instructions tool
type SetupInstructions struct {
	Prerequisites     []string           `yaml:"prerequisites"`
	InstallationSteps []InstallationStep `yaml:"installation_steps"`
}

// InstallationStep represents a single installation step
type InstallationStep struct {
	Step           int             `yaml:"step"`
	Title          string          `yaml:"title"`
	Description    string          `yaml:"description"`
	Commands       []string        `yaml:"commands,omitempty"`
	ConfigSnippets []ConfigSnippet `yaml:"config_snippets,omitempty"`
	Verification   string          `yaml:"verification,omitempty"`
}

// ConfigSnippet represents a configuration snippet
type ConfigSnippet struct {
	Filename string `yaml:"filename"`
	Content  string `yaml:"content"`
}

// KibanaSetupInstructions represents Kibana setup instructions
type KibanaSetupInstructions struct {
	Default KibanaSetupSteps `yaml:"default"`
	TCP     KibanaSetupSteps `yaml:"tcp,omitempty"`
	UDP     KibanaSetupSteps `yaml:"udp,omitempty"`
}

// KibanaSetupSteps represents steps for Kibana setup
type KibanaSetupSteps struct {
	Steps []KibanaSetupStep `yaml:"steps"`
}

// KibanaSetupStep represents a single Kibana setup step
type KibanaSetupStep struct {
	Step        int    `yaml:"step"`
	Instruction string `yaml:"instruction"`
}

// Troubleshooting represents troubleshooting information
type Troubleshooting struct {
	CommonIssues []TroubleshootingIssue `yaml:"common_issues"`
}

// TroubleshootingIssue represents a troubleshooting issue
type TroubleshootingIssue struct {
	Issue    string `yaml:"issue"`
	Solution string `yaml:"solution"`
}

// ValidationSteps represents validation steps
type ValidationSteps struct {
	Steps []ValidationStep `yaml:"steps"`
}

// ValidationStep represents a single validation step
type ValidationStep struct {
	Step           int      `yaml:"step"`
	Title          string   `yaml:"title"`
	Description    string   `yaml:"description"`
	Commands       []string `yaml:"commands"`
	ExpectedOutput string   `yaml:"expected_output"`
}

// ConfigLoader handles loading service configurations from YAML files
type ConfigLoader struct {
	configDir string
	services  map[string]*ServiceConfig
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(configDir string) *ConfigLoader {
	return &ConfigLoader{
		configDir: configDir,
		services:  make(map[string]*ServiceConfig),
	}
}

// LoadAllServices loads all service configurations from the config directory
func (cl *ConfigLoader) LoadAllServices() error {
	servicesDir := filepath.Join(cl.configDir, "services")

	// Check if services directory exists
	if _, err := os.Stat(servicesDir); os.IsNotExist(err) {
		return fmt.Errorf("services directory does not exist: %s", servicesDir)
	}

	// Read all YAML files in the services directory
	files, err := ioutil.ReadDir(servicesDir)
	if err != nil {
		return fmt.Errorf("failed to read services directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			serviceName := strings.TrimSuffix(file.Name(), ".yaml")
			serviceName = strings.TrimSuffix(serviceName, ".yml")

			configPath := filepath.Join(servicesDir, file.Name())
			config, err := cl.LoadServiceConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config for %s: %v", serviceName, err)
			}

			cl.services[serviceName] = config
		}
	}

	return nil
}

// LoadServiceConfig loads a single service configuration from a YAML file
func (cl *ConfigLoader) LoadServiceConfig(configPath string) (*ServiceConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	var config ServiceConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config file %s: %v", configPath, err)
	}

	return &config, nil
}

// GetServiceConfig returns the configuration for a specific service
func (cl *ConfigLoader) GetServiceConfig(serviceName string) (*ServiceConfig, error) {
	config, exists := cl.services[strings.ToLower(serviceName)]
	if !exists {
		availableServices := make([]string, 0, len(cl.services))
		for name := range cl.services {
			availableServices = append(availableServices, name)
		}
		return nil, fmt.Errorf("service '%s' not found. Available services: %s", serviceName, strings.Join(availableServices, ", "))
	}
	return config, nil
}

// GetAllServiceNames returns all available service names
func (cl *ConfigLoader) GetAllServiceNames() []string {
	names := make([]string, 0, len(cl.services))
	for name := range cl.services {
		names = append(names, name)
	}
	return names
}

// GetServiceConfigs returns all service configurations
func (cl *ConfigLoader) GetServiceConfigs() map[string]*ServiceConfig {
	return cl.services
}
