package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Paths struct {
		PathToFolderWithRSAKeys string `yaml:"path_to_folder_with_rsa_keys"`
		PathToFolderWithLogs    string `yaml:"path_to_folder_with_logs"`
	} `yaml:"paths"`
	Gossip struct {
		// SelfNodeName marks current node name (default is os.Hostname())
		SelfNodeName string `yaml:"self_node_name"`

		// SelfNodePort
		SelfNodePort int `yaml:"self_node_port"`

		// SecretKey AES key stored as base64
		SecretKey string `yaml:"secret_key"`

		// BoostrapNodeIp is used for boostraping cluster nodes
		BoostrapNodeAddr string `yaml:"boostrap_node_addr"`
		IsBoostrapNode   bool   `yaml:"is_boostrap_node"`
	}
	Api struct {
		Addr string `yaml:"addr"`
	}
}

func NewDefaultConfig() *Config {
	hostName, _ := os.Hostname()

	return &Config{
		Paths: struct {
			PathToFolderWithRSAKeys string `yaml:"path_to_folder_with_rsa_keys"`
			PathToFolderWithLogs    string `yaml:"path_to_folder_with_logs"`
		}{
			PathToFolderWithRSAKeys: RSAKeysFolder,
			PathToFolderWithLogs:    LogsFolder,
		},
		Gossip: struct {
			SelfNodeName     string `yaml:"self_node_name"`
			SelfNodePort     int    `yaml:"self_node_port"`
			SecretKey        string `yaml:"secret_key"`
			BoostrapNodeAddr string `yaml:"boostrap_node_addr"`
			IsBoostrapNode   bool   `yaml:"is_boostrap_node"`
		}{
			SelfNodeName:     hostName,
			SelfNodePort:     SelfPort,
			SecretKey:        SecretKey,
			BoostrapNodeAddr: BoostrapNodeIp,
			IsBoostrapNode:   IsBoostrapNode,
		},
		Api: struct {
			Addr string `yaml:"addr"`
		}{
			Addr: ApiAddr,
		},
	}
}

func (c *Config) SaveToFile() error {
	f, err := os.Create(ConfigFilename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	encoder := yaml.NewEncoder(f)
	defer func(encoder *yaml.Encoder) {
		_ = encoder.Close()
	}(encoder)

	err = encoder.Encode(c)
	return err
}

func (c *Config) DeleteFile() error {
	return os.Remove(ConfigFilename)
}

func (c *Config) ListFolders() (folders []string) {
	folders = append(folders, c.Paths.PathToFolderWithLogs, c.Paths.PathToFolderWithRSAKeys)
	return folders
}

func LoadFromFile() (config *Config, err error) {
	data, err := os.ReadFile(ConfigFilename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}
