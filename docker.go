package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

var ErrNoContainers = errors.New("no containers found")

type DockerService struct {
	client *client.Client
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerService{client: cli}, nil
}

func (s *DockerService) GetContainers() (map[string]string, error) {
	ctx := context.Background()
	containers, err := s.client.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	if len(containers.Items) == 0 {
		return nil, ErrNoContainers
	}
	containerResult := make(map[string]string)
	for _, container := range containers.Items {
		containerResult[container.ID] = strings.TrimPrefix(container.Names[0], "/")
	}
	return containerResult, nil
}

func (s *DockerService) WatchLogs(containerID, containerName string, lines string, color string) error {
	loggerConfig := Config{
		Output:        os.Stdout,
		ContainerName: containerName,
		Color:         color,
	}
	lw := NewLogWriter(loggerConfig, "stdout")

	ctx := context.Background()

	stdoutReader, err := s.client.ContainerLogs(ctx, containerID, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       lines,
	})
	if err != nil {
		return err
	}
	defer func(logReader client.ContainerLogsResult) {
		err := logReader.Close()
		if err != nil {
			panic(err)
		}
	}(stdoutReader)

	_, _ = stdcopy.StdCopy(
		lw,
		lw,
		stdoutReader,
	)
	return nil
}

func (s *DockerService) Close() error {
	return s.client.Close()
}
