package strategies

import (
	"errors"

	docker "github.com/saturdayshdev/shbackup/internal/docker"
)

type DumpConfig struct {
	Name      string
	User      string
	Password  string
	Container string
}

type Strategy interface {
	GetDump(docker *docker.Client, config DumpConfig) (*string, error)
}

func GetStrategy(strategy string) (Strategy, error) {
	if strategy == "postgres" {
		return &PostgresStrategy{}, nil
	}

	if strategy == "mysql" {
		return &MySQLStrategy{}, nil
	}

	return nil, errors.New("invalid backup strategy")
}
