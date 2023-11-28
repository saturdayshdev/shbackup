package lib

import (
	"errors"
	"log"
	"os"

	strategies "github.com/saturdayshdev/shbackup/internal/backup/strategies"
	docker "github.com/saturdayshdev/shbackup/internal/docker"
)

type BackupConfig struct {
	Name      string
	Strategy  strategies.Strategy
	User      string
	Password  string
	Container *docker.Container
}

func GetBackupConfig(container *docker.Container) (*BackupConfig, error) {
	labels := container.Labels
	if labels["shbackup.enabled"] != "true" {
		return nil, errors.New("shbackup.enabled label not found")
	}

	name, ok := labels["shbackup.name"]
	if !ok {
		return nil, errors.New("shbackup.name label not found")
	}

	strategy, ok := labels["shbackup.strategy"]
	if !ok {
		return nil, errors.New("shbackup.strategy label not found")
	}

	user, ok := labels["shbackup.user"]
	if !ok {
		return nil, errors.New("shbackup.user label not found")
	}

	password, ok := labels["shbackup.password"]
	if !ok {
		return nil, errors.New("shbackup.password label not found")
	}

	backupStrategy, err := strategies.GetStrategy(strategy)
	if err != nil {
		return nil, err
	}

	return &BackupConfig{
		Name:      name,
		Strategy:  backupStrategy,
		User:      user,
		Password:  password,
		Container: container,
	}, nil
}

func BackupDatabase(docker *docker.Client, storage *StorageClient, config *BackupConfig) error {
	log.Printf("Backing up %s database\n", config.Name)

	file, err := config.Strategy.GetDump(docker, strategies.DumpConfig{
		Name:      config.Name,
		User:      config.User,
		Password:  config.Password,
		Container: config.Container.ID,
	})
	if err != nil {
		return err
	}

	errCh := make(chan error)
	go storage.UploadFile(*file, *file, errCh)

	err = <-errCh
	if err != nil {
		return err
	}

	err = os.Remove(*file)
	if err != nil {
		return err
	}

	log.Printf("Backup of %s database completed\n", config.Name)

	return nil
}
