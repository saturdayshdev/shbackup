package lib

import (
	"context"

	types "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	Client *client.Client
}

func (c *DockerClient) ExecInContainer(id string, cmd []string) error {
	ctx := context.Background()
	exec, err := c.Client.ContainerExecCreate(ctx, id, types.ExecConfig{Cmd: cmd, Tty: false})
	if err != nil {
		return err
	}

	err = c.Client.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	return nil
}

func CreateDockerClient() (*DockerClient, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &DockerClient{Client: client}, nil
}
