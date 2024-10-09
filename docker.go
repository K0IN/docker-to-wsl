package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type dockerLogEntry struct {
	Stream      string `json:"stream"`
	ErrorDetail *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Error   string `json:"error"`
	} `json:"errorDetail"`
}

func buildDocker(dockerFilePath string) (imageName *string, err error) {
	imgName := "tmp-image:latest"
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}
	defer apiClient.Close()

	fullPath, err := filepath.Abs(filepath.Dir(dockerFilePath))

	if err != nil {
		return nil, fmt.Errorf("failed to get dockerfile path: %v", err)
	}

	ctx, err := archive.TarWithOptions(fullPath, &archive.TarOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create tar: %v", err)
	}

	dockerFileName := filepath.Base(dockerFilePath)
	res, err := apiClient.ImageBuild(context.Background(), ctx, types.ImageBuildOptions{Tags: []string{imgName}, Remove: false, Dockerfile: dockerFileName})
	if err != nil {
		return nil, fmt.Errorf("imageBuild failed: %v", err)
	}

	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		var entry dockerLogEntry
		err := json.Unmarshal([]byte(line), &entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
			continue
		}

		if entry.Stream != "" {
			fmt.Printf("> %s", entry.Stream)
		}

		if entry.ErrorDetail != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", entry.ErrorDetail.Message)
			return nil, fmt.Errorf("failed to build Docker image: %v", entry.ErrorDetail.Message)
		}
	}

	return &imgName, nil
}

func pullDockerImage(imageName string) (imgName *string, err error) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	defer apiClient.Close()
	res, err := apiClient.ImagePull(context.Background(), imageName, image.PullOptions{})
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(os.Stdout, res); err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	defer res.Close()
	return &imageName, nil
}

func exportDockerImage(imageName string) (err error) {
	containerName := "container"
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}
	defer apiClient.Close()

	createdContainer, err := apiClient.ContainerCreate(context.Background(),
		&container.Config{Image: imageName},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		&v1.Platform{},
		containerName)

	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	reader, exportError := apiClient.ContainerExport(context.Background(), createdContainer.ID)

	if err = apiClient.ContainerRemove(context.Background(), createdContainer.ID, container.RemoveOptions{Force: true, RemoveVolumes: true}); err != nil {
		return fmt.Errorf("failed to remove container: %v", err)
	}

	if exportError != nil {
		return fmt.Errorf("failed to export container: %v", exportError)
	}

	defer reader.Close()

	imageFile, err := os.Create("image.tar")
	if err != nil {
		return fmt.Errorf("failed to create image file: %v", err)
	}

	defer imageFile.Close()

	if _, err = io.Copy(imageFile, reader); err != nil {
		return fmt.Errorf("failed to copy image file: %v", err)
	}

	return nil
}
