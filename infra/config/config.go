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
	Nodes struct {
		BoostrapNodeIp   string `yaml:"boostrap_node_ip"`
		BoostrapNodePort string `yaml:"boostrap_node_port"`
	} `yaml:"nodes"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Paths: struct {
			PathToFolderWithRSAKeys string `yaml:"path_to_folder_with_rsa_keys"`
			PathToFolderWithLogs    string `yaml:"path_to_folder_with_logs"`
		}{
			PathToFolderWithRSAKeys: RSAKeysFolder,
			PathToFolderWithLogs:    LogsFolder,
		},
		Nodes: struct {
			BoostrapNodeIp   string `yaml:"boostrap_node_ip"`
			BoostrapNodePort string `yaml:"boostrap_node_port"`
		}{
			BoostrapNodeIp:   BoostrapNodeIp,
			BoostrapNodePort: BoostrapNodePort,
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
