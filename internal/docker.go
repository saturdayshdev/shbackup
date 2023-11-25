package lib

import (
	"context"

	types "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	Client *client.Client
}

func (c *DockerClient) ExecInContainer(id string, cmd []string) (*types.HijackedResponse, error) {
	exec, err := c.Client.ContainerExecCreate(context.Background(), id, types.ExecConfig{Cmd: cmd})
	if err != nil {
		return nil, err
	}

	attach, err := c.Client.ContainerExecAttach(context.Background(), exec.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}

	return &attach, nil
}

func CreateDockerClient() (*DockerClient, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &DockerClient{Client: client}, nil
}
