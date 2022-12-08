package init

import (
	"github.com/decentralized-hse/go-log-gossip/infra/config"
	"github.com/decentralized-hse/go-log-gossip/infra/keys"
	"os"
)

func CommandInitConfig() (err error) {
	newConfig := config.NewDefaultConfig()
	err = newConfig.SaveToFile()
	if err != nil {
		return err
	}

	err = createFoldersFromConfig(newConfig)
	if err != nil {
		return err
	}

	err = createRSAKeys(newConfig)
	if err != nil {
		return err
	}

	return nil
}

func createFoldersFromConfig(c *config.Config) (err error) {
	for _, folder := range c.ListFolders() {
		err = os.Mkdir(folder, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func createRSAKeys(c *config.Config) (err error) {
	newPair, err := keys.GenerateNewPair()
	if err != nil {
		return err
	}

	err = newPair.SaveToFiles(c.Paths.PathToFolderWithRSAKeys)
	return err
}
