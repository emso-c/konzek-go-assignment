// This file is responsible for loading the configuration from the TOML file and setting the environment variables.
// The configuration is loaded from the config.toml file.
// The configuration values are then set as environment variables.
// The environment variables are then used by the other packages to access the configuration values.
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	toml "github.com/pelletier/go-toml"
)

// The default path to the configuration file
const DEFAULT_CONFIG_PATH = "config.toml"

// LoadEnv loads the configuration from the TOML file and sets the environment variables
//
// path: The path to the TOML file
//
// return: An error if one occurs
func LoadEnv(path string) error {
	// Read TOML data from file
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal TOML data into a map
	var config map[string]interface{}
	if err := toml.Unmarshal(data, &config); err != nil {
		return err
	}

	// Loop over each group of configuration values and set environment variables
	for group, values := range config {
		for key, val := range values.(map[string]interface{}) {
			// Make sure the environment variables are uppercase
			// Windows environment variables are case insensitive
			// but since we are using the environment variables in Docker,
			// we need to make sure they are uppercase
			group = strings.ToUpper(group)
			key = strings.ToUpper(key)
			if err := SetEnv(group, key, val); err != nil {
				return err
			}
		}
	}

	// raise error if SERVER_HOST or SERVER_PORT is not set
	if os.Getenv("SERVER_HOST") == "" || os.Getenv("SERVER_PORT") == "" {
		return fmt.Errorf("SERVER_HOST and SERVER_PORT must be set")
	}

	return nil
}

// SetEnv sets a TOML based environment variable.
// The environment variable is set in the format: <prefix>_<key>=<value>
// where <prefix> is the group name and <key> is the configuration key.
//
// Example:
//
// [database]
//
// host = "localhost"
//
// port = 5432
//
// The environment variable for the host would be: DATABASE_HOST=localhost
//
// The environment variable for the port would be: DATABASE_PORT=5432
//
// prefix: The group name
//
// key: The configuration key
//
// val: The configuration value
func SetEnv(prefix string, key string, val interface{}) error {
	switch v := val.(type) {
	case string:
		os.Setenv(fmt.Sprintf("%s_%s", prefix, key), v)
	case int:
		os.Setenv(fmt.Sprintf("%s_%s", prefix, key), strconv.Itoa(v))
	case int64:
		os.Setenv(fmt.Sprintf("%s_%s", prefix, key), strconv.FormatInt(v, 10))
	case bool:
		os.Setenv(fmt.Sprintf("%s_%s", prefix, key), strconv.FormatBool(v))
	case []interface{}:
		// Convert array to comma-separated string
		values := make([]string, len(v))
		for i, item := range v {
			values[i] = fmt.Sprintf("%v", item)
		}
		os.Setenv(fmt.Sprintf("%s_%s", prefix, key), strings.Join(values, ","))
	default:
		log.Printf("Ignoring unsupported type for key %s: %T", key, v)
	}
	return nil
}
