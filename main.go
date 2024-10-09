package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

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
			&cli.BoolFlag{
				Name:  "set-default",
				Usage: "Set the distribution as the default",
			},
			&cli.BoolFlag{
				Name:  "start-menu",
				Usage: "Add the distribution to the Start Menu",
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

			if c.Bool("set-default") {
				err = setDefaultWsl(distroName)
				if err != nil {
					fmt.Printf("failed to set distro as default: %v\n", err)
				}
			}

			if c.Bool("start-menu") {
				err = addToStartMenu(distroName)
				if err != nil {
					fmt.Printf("failed to add to Start Menu: %v\n", err)
				}
			}

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
