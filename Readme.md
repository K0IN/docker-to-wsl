# Docker to WSL

Docker to WSL is a tool that converts Docker images into WSL distributions. 
This project allows you to build or pull Docker images and then import them into WSL for further use.

## Why

- Easily manage and replicate development environments across multiple systems
- Avoid corrupting your main WSL distribution when testing or experimenting
- Quickly start over with a fresh environment when needed
- Leverage Docker's vast ecosystem of images for WSL use
- Simplify the process of creating custom WSL distributions

## Features

- Convert Docker images to WSL distributions
- Build custom Docker images and import them as WSL distributions
- Pull existing Docker images and convert them to WSL distributions
- Launch newly created WSL distributions directly
- Support for custom Dockerfiles and configurations

## Installation

### Dependencies

This is a Windows only app.

- Go 1.22 or later
- Docker
- WSL (Windows Subsystem for Linux)

### Steps

```
go install github.com/k0in/docker-to-wsl/v2@main
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

- `--distro-name`: Set the name for the new WSL distribution (required)
- `--image`: Specify a Docker image to pull and convert or if a local file is specified, it will be built
- `--launch`: Launch the new WSL distribution after creation
- `--help`: Show help information

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

## Quick start

Example Dockerfile (you can change this as you like):

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

## Add your distro to the start menu (optional)

Please fill in <>

```powershell
$WScriptShell = New-Object -ComObject WScript.Shell
$Shortcut = $WScriptShell.CreateShortcut("$env:APPDATA\Microsoft\Windows\Start Menu\Programs\<stat-menu-name>.lnk")
$Shortcut.TargetPath = "wsl.exe"
$Shortcut.Arguments = "-d <dist-name>"
$Shortcut.Save()
```

## Complex example

Here is a full Dockerfile for a more complex setup

```Dockerfile
FROM ubuntu:24.10

# basic setup
RUN apt-get update && apt-get upgrade -y && apt-get install -y software-properties-common
RUN apt update && apt install -y fish sudo curl
RUN chsh -s /usr/bin/fish

# setup user
RUN useradd -m -s /usr/bin/fish -G sudo k0in
# set password to 'k0in'
RUN echo "k0in:k0in" | chpasswd k0in
# confgure wsl
RUN printf "[user]\ndefault=k0in" >> /etc/wsl.conf

# ssh setup
COPY --chown=k0in:k0in files/.ssh /home/k0in/.ssh
RUN chmod 700 /home/k0in/.ssh
RUN chmod 600 /home/k0in/.ssh/*

# setup sudo
RUN echo "k0in ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# install packages
RUN apt-get install -y wget git vim nano openssh-client clang gcc g++ make cmake gdb python3 python3-pip python3-venv

# example you can use x11 apps :)
RUN apt-get install -y x11-apps
```
