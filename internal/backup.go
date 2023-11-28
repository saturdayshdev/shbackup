package lib

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	types "github.com/docker/docker/api/types"
)

type PostgresStrategy struct{}
type MySQLStrategy struct{}

type BackupConfig struct {
	Name      string
	Strategy  BackupStrategy
	User      string
	Password  string
	Container *types.Container
}
type BackupStrategy interface {
	GetDump(docker *DockerClient, config *BackupConfig) (*string, error)
}

func (s *PostgresStrategy) GetDump(docker *DockerClient, config *BackupConfig) (*string, error) {
	file := fmt.Sprint(time.Now().Unix()) + "_" + config.Name + ".sql"

	cmd := []string{"pg_dump", "-U", config.User, "-W", config.Password, "-f", file}
	err := docker.ExecInContainer(config.Container.ID, cmd)
	if err != nil {
		return nil, err
	}

	var stream io.ReadCloser
	for i := 0; i < 10; i++ {
		stream, _, err = docker.Client.CopyFromContainer(context.Background(), config.Container.ID, "/"+file)
		if err == nil {
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	tr := tar.NewReader(stream)
	if _, err := tr.Next(); err != nil {
		return nil, err
	}

	dest, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	defer dest.Close()

	_, err = io.Copy(dest, tr)
	if err != nil {
		return nil, err
	}

	cmd = []string{"rm", file}
	err = docker.ExecInContainer(config.Container.ID, cmd)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (s *MySQLStrategy) GetDump(docker *DockerClient, config *BackupConfig) (*string, error) {
	file := fmt.Sprint(time.Now().Unix()) + "_" + config.Name + ".sql"

	cmd := []string{"mysqldump", "-u", config.User, "-p" + config.Password, "-f", file}
	err := docker.ExecInContainer(config.Container.ID, cmd)
	if err != nil {
		return nil, err
	}

	var stream io.ReadCloser
	for i := 0; i < 10; i++ {
		stream, _, err = docker.Client.CopyFromContainer(context.Background(), config.Container.ID, "/"+file)
		if err == nil {
			break
		}

		time.Sleep(1000 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	tr := tar.NewReader(stream)
	if _, err := tr.Next(); err != nil {
		return nil, err
	}

	dest, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	defer dest.Close()

	_, err = io.Copy(dest, tr)
	if err != nil {
		return nil, err
	}

	cmd = []string{"rm", file}
	err = docker.ExecInContainer(config.Container.ID, cmd)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func GetBackupStrategy(strategy string) (BackupStrategy, error) {
	if strategy == "postgres" {
		return &PostgresStrategy{}, nil
	}

	if strategy == "mysql" {
		return &MySQLStrategy{}, nil
	}

	return nil, errors.New("invalid backup strategy")
}

func GetBackupConfig(labels map[string]string, container *types.Container) (*BackupConfig, error) {
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

	backupStrategy, err := GetBackupStrategy(strategy)
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

func BackupDatabase(docker *DockerClient, storage *StorageClient, config *BackupConfig) error {
	log.Printf("Backing up %s database\n", config.Name)

	file, err := config.Strategy.GetDump(docker, config)
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
