package docker

import (
	"context"

	types "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Client struct {
	Client *client.Client
}

type Container = types.Container

func (c *Client) ExecInContainer(id string, cmd []string) error {
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

func (c *Client) GetContainers() ([]types.Container, error) {
	containers, err := c.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func CreateClient() (*Client, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client}, nil
}
