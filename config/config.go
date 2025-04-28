package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

type Option struct {
	Prefix       string
	Delimiter    string
	Separator    string
	YamlFilePath string
	CallBackEnv  func(string) string
}

// defaultCallbackEnv process environment variable keys based on prefix and separator.
func defaultCallbackEnv(source, prefix, separator string) string {
	base := strings.ToLower(strings.TrimPrefix(source, prefix))

	return strings.ReplaceAll(base, separator, ".")
}

// Load function loads configuration from YAML file and environment variables based on provided options.
func Load(options Option, config, defaultConfig interface{}) {
	k := koanf.New(options.Delimiter)

	// Load default configuration from Default function.
	if defaultConfig != nil {
		if err := k.Load(structs.Provider(defaultConfig, "koanf"), nil); err != nil {
			log.Fatalf("error loading default config: %s", err.Error())
		}
	}

	// Load configuration from YAML file if provided.
	if options.YamlFilePath != "" {
		if err := k.Load(file.Provider(options.YamlFilePath), yaml.Parser()); err != nil {
			log.Fatalf("error when trying to loading yml config file: %s", err.Error())
		}
	}

	// Define callback function for environment variables.
	callback := options.CallBackEnv
	if callback == nil {
		// Set default callback using the prefix and separator from options.
		callback = func(source string) string {
			return defaultCallbackEnv(source, options.Prefix, options.Separator)
		}
	}

	// Load environment variables with the specified prefix and callback.
	if err := k.Load(env.Provider(options.Prefix, options.Delimiter, callback), nil); err != nil {
		log.Fatalf("error when loading environment variables: %s", err.Error())
	}

	// Unmarshal into provided config structure (passing address).
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("error unmarshalling config: %s", err.Error())
	}
}
