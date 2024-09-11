# Docker to WSL

Docker to WSL is a tool that converts Docker images into WSL distributions. This project allows you to build or pull Docker images and then import them into WSL for further use.

## Installation

### Dependencies

- Go 1.22 or later
- Docker
- WSL (Windows Subsystem for Linux)

### Steps

```
go install github.com/K0IN/docker-to-wsl/v2@main
```

OR 

1. Clone the repository:
    ```
    git clone https://github.com/K0IN/docker-to-wsl.git
    cd docker-to-wsl
    ```

2. Install the tool:
    ```
    go install
    ```

## Usage

### Building and Importing a Dockerfile

1. Create a Dockerfile in your current directory.
2. Run the tool:
    ```
    docker-to-wsl --distro-name myDistro
    ```

### Pulling and Importing a Docker Image

1. Run the tool:
    ```
    docker-to-wsl --image <docker-image-name> --distro-name myDistro
    ```

### Launching the WSL Distribution

1. Add the `--launch` flag to the command:
    ```
    docker-to-wsl --image <docker-image-name> --distro-name myDistro --launch
    ```

## Dependencies and Licensing

This project requires the following dependencies:
- [Docker](https://github.com/docker/docker)
- [Opencontainers Image Spec](https://github.com/opencontainers/image-spec)
- [urfave/cli](https://github.com/urfave/cli)
- [yuk7/wsllib-go](https://github.com/yuk7/wsllib-go)

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Overview of Main Functions

- `buildDocker(dockerFilePath string)`: Builds a Docker image from a Dockerfile.
- `pullDockerImage(imageName string)`: Pulls a Docker image from a registry.
- `exportDockerImage(imageName string)`: Exports a Docker image to a tar file.
- `importWsl(distroName string)`: Imports the tar file into WSL as a new distribution.
- `launchWsl()`: Launches the WSL distribution.

## Quickstart

Example dockerfile (you can change this as you like):

```Dockerfile
# example image
FROM alpine:latest 
RUN apk update && apk add fish shadow
RUN chsh -s /usr/bin/fish
# Example add a env variable, Note: you cant use ENV
RUN fish -c "set -Ux key value"
# Example run a command on start up
RUN printf "[boot]\ncommand = /etc/entrypoint.sh" >> /etc/wsl.conf
RUN printf "#!/bin/sh\ntouch /root/booted" >> /etc/entrypoint.sh
RUN chmod +x /etc/entrypoint.sh
```

then run:

```bash
docker-to-wsl --distro-name myDistro
wsl -d myDistro
```


