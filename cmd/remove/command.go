package remove

import (
	"github.com/decentralized-hse/go-log-gossip/infra/config"
	"os"
)

func CommandRemoveConfig() (err error) {
	cfg, err := config.LoadFromFile()
	if err != nil {
		return err
	}

	err = removeFolders(cfg)
	return cfg.DeleteFile()
}

func removeFolders(c *config.Config) (err error) {
	for _, folder := range c.ListFolders() {
		err = os.RemoveAll(folder)
		if err != nil {
			return err
		}
	}

	return nil
}
