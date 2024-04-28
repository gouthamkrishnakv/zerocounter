package config

import (
	// stdlib imports
	"fmt"
	"os"

	// third-party imports
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"

	// internal imports
	"github.com/gouthamkrishnakv/zerocounter/logging"
	utils_fs "github.com/gouthamkrishnakv/zerocounter/utils/filesystem"
)

// -- constants --

// Configuration file path (relative)
const DefaultConfigFile = "zc.yaml"

// Default config file, this is added and setup during configuration time
var defaultConfig = &Config{
	DevMode: false,
}

// -- variables --

// logger
var logger *zerolog.Logger = nil

// Database Path
var databasePath string = ""

// Config Path
var configPath string = ""

// config
var config *Config = nil

// -- structs --

// Config stores "absolutely necessary" and user-configurable options
// in a single-file.
type Config struct {
	// TODO: decide on whether to even have a dev-mode here.
	DevMode bool `yaml:"dev_mode"`
}

// -- functions --

// Initialize sets up configuration from the "first-found" configuration file.
// This isn't expected to change, so we'll just stick with most accepted
// configuration paths.
func Initialize() error {
	// setup logging
	configLogger := logging.L().With().Str("module", "config").Timestamp().Stack().Logger()
	logger = &configLogger

	// search if the configuration *file* is there, if not, run setup
	filePath, searchErr := utils_fs.SearchConfigFile(DefaultConfigFile)
	if searchErr != nil {
		if os.IsNotExist(searchErr) {
			var setupErr error
			logger.Warn().Str("config_path", filePath).Err(setupErr).Msg("config file not found. running setup")
			filePath, setupErr = setup()
			if setupErr != nil {
				return fmt.Errorf("Initialize.setupErr: %v", setupErr)
			}
		}
	}

	// read the config file
	if readConfigFileErr := readConfigFile(filePath); readConfigFileErr != nil {
		return fmt.Errorf("Initialize.readConfigErr: %v", readConfigFileErr)
	}
	return nil
}

// GetConfig returns the configuration file
func GetConfig() *Config {
	return config
}

// setup sets up the directories, creates the config and the database files
func setup() (string, error) {
	var configFileErr error

	// generate the (xdg-style) config path that we require (if needed) and
	// return the path
	configPath, configFileErr = utils_fs.ConfigFile(DefaultConfigFile)
	if configFileErr != nil {
		return "", configFileErr
	}
	logger.Debug().Str("file_path", configPath).Msg("configuration file path generated")

	// serialize the configuration as bytes
	configContents, marshalErr := yaml.Marshal(defaultConfig)
	if marshalErr != nil {
		return "", marshalErr
	}

	// create the new config file
	configFile, fileCreateErr := os.Create(configPath)
	if fileCreateErr != nil {
		return "", fileCreateErr
	}
	// file to be closed
	defer configFile.Close()

	// write the marshaled data to the created config file
	if _, writeErr := configFile.Write(configContents); writeErr != nil {
		return "", writeErr
	}

	logger.Info().Str("configPath", configPath).Msg("default cofiguration written")
	return configPath, nil
}

func readConfigFile(filePath string) error {
	// read the new configuration file
	configFileContents, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return fmt.Errorf("readConfigFile: %v", readErr)
	}
	logger.Debug().Str("file_path", filePath).Msg("config file read")

	// create a config variable to unmarshal config to, and set the configuration
	newConfig := new(Config)
	if unmarshalErr := yaml.Unmarshal(configFileContents, newConfig); unmarshalErr != nil {
		return unmarshalErr
	}
	config = newConfig

	logging.L().Info().Str("file_path", filePath).Msg("configuration loaded")
	return nil
}
