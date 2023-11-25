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
	Strategy  BackupStrategy
	User      string
	Password  string
	Container *types.Container
}
type BackupStrategy interface {
	GetDump(config *BackupConfig) (*string, error)
}

func (s *PostgresStrategy) GetDump(config *BackupConfig) (*string, error) {
	timestamp := time.Now().Unix()
	file := string(rune(timestamp)) + "_" + config.Name + ".sql"

	cmd := []string{"pg_dump", "-U", config.User, ">", file}
	err := config.Client.ExecInContainer(config.Container.ID, cmd)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (s *MySQLStrategy) GetDump(config *BackupConfig) (*string, error) {
	timestamp := time.Now().Unix()
	file := string(rune(timestamp)) + "_" + config.Name + ".sql"

	cmd := []string{"mysqldump", "-u", config.User, "-p" + config.Password, ">", file}
	err := config.Client.ExecInContainer(config.Container.ID, cmd)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func GetBackupStrategy(labels map[string]string) (BackupStrategy, error) {
	strategy := labels["shbackup.strategy"]

	if strategy == "postgres" {
		return &PostgresStrategy{}, nil
	}

	if strategy == "mysql" {
		return &MySQLStrategy{}, nil
	}

	return nil, errors.New("unknown backup strategy")
}

func (c *DockerClient) BackupDatabase(config *BackupConfig) (*string, error) {
	return config.Strategy.GetDump(config)
}
