package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileName = ".statusline.config"

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".claude", ConfigFileName), nil
}

// Load reads the config from the default location
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	return LoadFromPath(path)
}

// LoadFromPath reads the config from a specific path
func LoadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return DefaultConfig(), nil
		}
		return nil, err
	}

	// Start with defaults, then overlay loaded config
	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save writes the config to the default location
func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	return SaveToPath(cfg, path)
}

// SaveToPath writes the config to a specific path
func SaveToPath(cfg *Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetStatuslineScriptPath returns the path where statusline.sh should be installed
func GetStatuslineScriptPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".claude", "statusline.sh"), nil
}

// InstallStatuslineScript copies the statusline.sh script to ~/.claude/statusline.sh
func InstallStatuslineScript() error {
	// Get the executable's directory to find statusline.sh
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	execDir := filepath.Dir(execPath)

	// Try multiple locations for statusline.sh
	possiblePaths := []string{
		filepath.Join(execDir, "statusline.sh"),
		filepath.Join(execDir, "..", "statusline.sh"),
		filepath.Join(execDir, "..", "..", "statusline.sh"),
	}

	var scriptContent []byte
	var found bool
	for _, srcPath := range possiblePaths {
		content, err := os.ReadFile(srcPath)
		if err == nil {
			scriptContent = content
			found = true
			break
		}
	}

	if !found {
		return os.ErrNotExist
	}

	destPath, err := GetStatuslineScriptPath()
	if err != nil {
		return err
	}

	// Ensure ~/.claude directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Write the script with executable permissions
	return os.WriteFile(destPath, scriptContent, 0755)
}

// SaveAndInstall saves the config and installs the statusline script globally
func SaveAndInstall(cfg *Config) error {
	// First save the config
	if err := Save(cfg); err != nil {
		return err
	}

	// Then install the script
	return InstallStatuslineScript()
}
