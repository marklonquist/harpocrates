package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BESTSELLER/harpocrates/config"
	"gopkg.in/yaml.v3"
)

// SecretJSON holds the information about which secrets to fetch and how to save them again
type SecretJSON struct {
	Append        *bool         `json:"append,omitempty"      yaml:"append,omitempty"`
	Format        string        `json:"format,omitempty"      yaml:"format,omitempty"`
	Output        string        `json:"output,omitempty"      yaml:"output,omitempty"`
	Owner         *int          `json:"owner,omitempty"       yaml:"owner,omitempty"`
	Prefix        string        `json:"prefix,omitempty"      yaml:"prefix,omitempty"`
	UpperCase     *bool         `json:"uppercase,omitempty"   yaml:"uppercase,omitempty"`
	Secrets       []interface{} `json:"secrets,omitempty"     yaml:"secrets,omitempty"`
	GcpWorkloadID bool          `json:"gcpWorkloadID,omitempty"     yaml:"gcpWorkloadID,omitempty"`
}

type Secret struct {
	Append    *bool         `json:"append,omitempty"      yaml:"append,omitempty"`
	Prefix    string        `json:"prefix,omitempty"      yaml:"prefix,omitempty"`
	Format    string        `json:"format,omitempty"      yaml:"format,omitempty"      mapstructure:"format,omitempty"`
	FileName  string        `json:"filename,omitempty"    yaml:"filename,omitempty"    mapstructure:"filename,omitempty"`
	UpperCase *bool         `json:"uppercase,omitempty"   yaml:"uppercase,omitempty"`
	Keys      []interface{} `json:"keys,omitempty"        yaml:"keys,omitempty"`
	Owner     *int          `json:"owner,omitempty"       yaml:"owner,omitempty"`
}

type SecretKeys struct {
	Append     *bool  `json:"append,omitempty"         yaml:"append,omitempty"`
	Prefix     string `json:"prefix,omitempty"         yaml:"prefix,omitempty"      mapstructure:"prefix,omitempty"`
	UpperCase  *bool  `json:"uppercase,omitempty"      yaml:"uppercase,omitempty"   mapstructure:"uppercase,omitempty"`
	SaveAsFile *bool  `json:"saveAsFile,omitempty"     yaml:"saveAsFile,omitempty"`
}

// ReadInput will read the input given to Harpocrates and try to parse it to SecretJSON
// Will also set some default values
func ReadInput(input string) SecretJSON {
	secretJSON := SecretJSON{}
	err := json.Unmarshal([]byte(input), &secretJSON)
	if err == nil {
		goto MoveOn
	}
	err = yaml.Unmarshal([]byte(input), &secretJSON)
	if err != nil {
		fmt.Printf("Your secret file contains an error, please refer to the documentation\n%v\n", err)
		os.Exit(1)
	}

MoveOn:
	if secretJSON.Format != "" {
		if secretJSON.Format != "json" && secretJSON.Format != "env" && secretJSON.Format != "secret" {
			fmt.Println("An invalid format was provided, only these formats are allowed at the moment:\njson\nenv\nsecret")
			os.Exit(1)
		}

		config.Config.Format = secretJSON.Format
	}

	if secretJSON.Output == "" {
		secretJSON.Output = "/secrets"
	}
	config.Config.Output = secretJSON.Output

	if secretJSON.Owner == nil {
		value := -1
		secretJSON.Owner = &value
	}
	config.Config.Owner = *secretJSON.Owner

	if len(secretJSON.Secrets) == 0 {
		fmt.Println("No secrets provided")
		os.Exit(1)
	}

	config.Config.Prefix = secretJSON.Prefix

	if secretJSON.UpperCase != nil {
		config.Config.UpperCase = *secretJSON.UpperCase
	}

	if secretJSON.Append != nil {
		config.Config.Append = *secretJSON.Append
	}

	if secretJSON.GcpWorkloadID {
		config.Config.GcpWorkloadID = secretJSON.GcpWorkloadID
	}

	return secretJSON
}
