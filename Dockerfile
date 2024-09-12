FROM alpine:latest 
RUN apk update && apk add fish shadow
RUN chsh -s /usr/bin/fish
# Example add a env variable, Note: you cant use ENV
RUN fish -c "set -Ux key value"
# Example run a command on start up
RUN printf "[boot]\ncommand = /etc/entrypoint.sh" >> /etc/wsl.conf
RUN printf "#!/bin/sh\ntouch /root/booted" >> /etc/entrypoint.sh
RUN chmod +x /etc/entrypoint.sh
# Example add a file (only current directory and subdirectories are allowed)
ADD ./go.mod /test/go.mod