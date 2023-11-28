package strategies

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	docker "github.com/saturdayshdev/shbackup/internal/docker"
)

type PostgresStrategy struct{}

func (s *PostgresStrategy) GetDump(docker *docker.Client, config DumpConfig) (*string, error) {
	file := fmt.Sprint(time.Now().Unix()) + "_" + config.Name + ".sql"

	cmd := []string{"pg_dump", "-U", config.User, "-W", config.Password, "-f", file}
	err := docker.ExecInContainer(config.Container, cmd)
	if err != nil {
		return nil, err
	}

	var stream io.ReadCloser
	for i := 0; i < 10; i++ {
		stream, _, err = docker.Client.CopyFromContainer(context.Background(), config.Container, "/"+file)
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
	err = docker.ExecInContainer(config.Container, cmd)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
