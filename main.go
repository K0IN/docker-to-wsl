package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/urfave/cli"
	"github.com/yuk7/wsllib-go"
)

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

	if _, err = io.Copy(os.Stdout, res.Body); err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if err = res.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %v", err)
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

func importWsl(distroName string) (err error) {
	_ = wsllib.WslUnregisterDistribution(distroName)
	path, err := filepath.Abs("image.tar")
	if err != nil {
		return err
	}

	var re = regexp.MustCompile(`~[\p{L}0-9\s]+`)
	escapedPath := re.ReplaceAllString(distroName, `-`)

	cmd := exec.Command("wsl", "--import", distroName, fmt.Sprintf("./%s", escapedPath), "image.tar", "--version", "2")
	_, err = cmd.Output()
	if err != nil {
		return wsllib.WslRegisterDistribution(distroName, path)
	}

	return nil
}

func launchWsl(distroName string) (err error) {
	_, err = wsllib.WslLaunchInteractive(distroName, "", true)
	return err
}

func main() {
	app := &cli.App{
		Name:  "Docker 2 WSL",
		Usage: "Convert a Docker image to a WSL distribution",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "distro-name",
				Usage: "the name of the WSL distribution",
				Value: "dist",
			},
			&cli.StringFlag{
				Name:  "image",
				Usage: "The name of the image or 'Dockerfile' if you want to build the Dockerfile in your current directory",
				Value: "Dockerfile",
			},
			&cli.BoolFlag{
				Name:  "launch",
				Usage: "Launch the WSL distribution after importing",
			},
		},
		Action: func(c *cli.Context) error {
			var ImageName = c.String("image")
			var distroName = c.String("distro-name")

			_, err := os.Stat(ImageName)
			hasDockerfile := err == nil
			var img *string

			if !hasDockerfile {
				println("Pulling Docker image...")
				img, err = pullDockerImage(ImageName)
				if err != nil {
					return fmt.Errorf("failed to pull Docker image: %v", err)
				}
			} else {
				println("Building Docker image...")
				img, err = buildDocker(ImageName)
				if err != nil {
					return fmt.Errorf("failed to build Docker image: %v", err)
				}
			}

			fmt.Printf("Successfully built Docker image '%s'\n", *img)

			err = exportDockerImage(*img)
			if err != nil {
				return fmt.Errorf("failed to export Docker image: %v", err)
			}

			defer os.Remove("image.tar")

			err = importWsl(distroName)
			if err != nil {
				return fmt.Errorf("failed to import WSL: %v", err)
			}

			fmt.Printf("Successfully imported %s to WSL\n", *img)
			fmt.Printf("Run `wsl -d %s` to launch the distribution\n", distroName)
			fmt.Printf("Run `wsl --unregister %s` to remove the distribution and files\n", distroName)

			if c.Bool("launch") {
				err = launchWsl(distroName)
				if err != nil {
					return fmt.Errorf("failed to launch WSL: %v", err)
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
