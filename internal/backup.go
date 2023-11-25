package lib

import (
	"errors"
	"time"

	types "github.com/docker/docker/api/types"
)

type PostgresStrategy struct{}
type MySQLStrategy struct{}

type BackupStrategy interface {
	GetDump(name string, client *DockerClient, container *types.Container) (*string, error)
}

func (s *PostgresStrategy) GetDump(name string, client *DockerClient, container *types.Container) (*string, error) {
	timestamp := time.Now().Unix()
	fileName := string(timestamp) + "_" + name + ".sql"

	err := client.ExecInContainer(container.ID, []string{"pg_dump", "-U", "postgres", "-f", fileName})
	if err != nil {
		return nil, err
	}

	return &fileName, nil
}

func (s *MySQLStrategy) GetDump(name string, client *DockerClient, container *types.Container) (*string, error) {
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

func (c *DockerClient) BackupDatabase(name string, container *types.Container, strategy BackupStrategy) (*string, error) {
	return strategy.GetDump(name, c, container)
}
