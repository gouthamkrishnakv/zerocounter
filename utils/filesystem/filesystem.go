package utils_fs

import (
	// stdlib imports

	"os"
	"path/filepath"

	// third-party imports
	"github.com/adrg/xdg"

	// internal imports
	"github.com/gouthamkrishnakv/zerocounter/constants"
)

// -- functions --

// generatePath is an internal function which provides a "common" method
// which will then generate the paths required
func generatePath(relativePath string, xdgPathFunction func(string) (string, error)) (string, error) {
	// This would be "<PATH_TO_APP>/zerocounter/<relativePath>"
	return xdgPathFunction(filepath.Join(constants.DefaultDir, relativePath))
}

// ConfigFile creates the standard file-path for configuration and then
// generates absolute path for config, provided its relative path is given.
func ConfigFile(relativePath string) (string, error) {
	return generatePath(relativePath, xdg.ConfigFile)
}

// DataFile creates the standrad file-path for data-directory and then
// generates absolute path for data-files, provided a relative path.
func DataFile(relativePath string) (string, error) {
	return generatePath(relativePath, xdg.DataFile)
}

// SearchConfigFile searches for a file in config dirsin xdg-standard provided
// it's path and if exists, returns the path to file
func SearchConfigFile(relativePath string) (string, error) {
	configPath, configErr := ConfigFile(relativePath)
	if configErr != nil {
		return "", configErr
	}
	_, statErr := os.Stat(configPath)
	if statErr != nil {
		return configPath, statErr
	}
	return configPath, nil
}

// SearchDataFile searches for a file in data dirs in xdg-standard provided
// it's path and if exists, returns path to file
func SearchDataFile(relativePath string) (string, error) {
	return xdg.SearchDataFile(filepath.Join(constants.DefaultDir, relativePath))
}
