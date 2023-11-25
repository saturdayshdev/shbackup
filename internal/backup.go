package lib

import (
	"errors"
	"time"

	types "github.com/docker/docker/api/types"
)

type PostgresStrategy struct{}
type MySQLStrategy struct{}

type BackupConfig struct {
	Client    *DockerClient
	Name      string
	User      string
	Password  string
	Container *types.Container
	Strategy  BackupStrategy
}
type BackupStrategy interface {
	GetDump(config *BackupConfig) (*string, error)
}

func (s *PostgresStrategy) GetDump(config *BackupConfig) (*string, error) {
	timestamp := time.Now().Unix()
	file := string(rune(timestamp)) + "_" + config.Name + ".sql"

	err := config.Client.ExecInContainer(config.Container.ID, []string{"pg_dump", "-U", "postgres", "-f", file})
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (s *MySQLStrategy) GetDump(config *BackupConfig) (*string, error) {
	return nil, errors.New("not implemented")
}

func GetBackupStrategy(labels map[string]string) (BackupStrategy, error) {
	label := labels["shbackup.type"]

	if label == "postgres" {
		return &PostgresStrategy{}, nil
	}

	if label == "mysql" {
		return &MySQLStrategy{}, nil
	}

	return nil, errors.New("unknown backup strategy")
}

func (c *DockerClient) BackupDatabase(config *BackupConfig) (*string, error) {
	return config.Strategy.GetDump(config)
}
