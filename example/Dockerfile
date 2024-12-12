FROM ubuntu:24.10

# Set environment variables
ENV USERNAME=k0in
ENV GO_VERSION=1.23.2

WORKDIR /tmp
# Basic setup
RUN apt-get update && apt-get upgrade -y && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y \
  software-properties-common build-essential fish sudo curl wget stow make \
  gcc g++ git vim nano openssh-client clang gdb cmake python3 python3-pip net-tools \
  python3-venv unzip btop htop jq tmux npm && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

# Setup user
RUN useradd -m -s /usr/bin/fish -G sudo $USERNAME && \
  echo "$USERNAME:$USERNAME" | chpasswd $USERNAME && \
  echo "$USERNAME ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# change root shell to fish
RUN chsh -s /usr/bin/fish root

# Configure WSL
RUN printf "[boot]\ncommand = /etc/entrypoint.sh\n[user]\ndefault=${USERNAME}" >> /etc/wsl.conf

# SSH setup
# COPY --chown=$USERNAME:$USERNAME files/.ssh /home/$USERNAME/.ssh
# RUN chmod 700 /home/$USERNAME/.ssh && chmod 600 /home/$USERNAME/.ssh/*

# Setup entrypoint
COPY entrypoint.sh /etc/entrypoint.sh
RUN chmod +x /etc/entrypoint.sh

# Install Go
RUN wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz && \
  mkdir -p /usr/local/go && \
  tar -C /usr/local/ -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
  rm go${GO_VERSION}.linux-amd64.tar.gz && \
  fish -c "set -Ux fish_user_paths /usr/local/go/bin/ $fish_user_paths"

RUN wget https://github.com/earthly/earthly/releases/latest/download/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly && /usr/local/bin/earthly bootstrap --with-autocomplete

# Switch to user
USER $USERNAME
WORKDIR /home/$USERNAME

# Setup dotfiles
COPY --chown=$USERNAME:$USERNAME dotfiles/ /home/$USERNAME/dotfiles/
RUN cd dotfiles && make all

# Install oh-my-posh
RUN curl -s https://ohmyposh.dev/install.sh | bash -s && fish -c "fish_add_path ~/.local/bin"

# install fzf
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && ~/.fzf/install --key-bindings --completion --no-update-rc

# install neovim
RUN curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux64.tar.gz && \
  tar -C ~/.local/bin -xzf nvim-linux64.tar.gz &&\
  rm nvim-linux64.tar.gz && \
  fish -c "fish_add_path ~/.local/bin/nvim-linux64/bin"

RUN fish -c "fish_add_path /usr/local/go/bin/"

SHELL ["fish", "-c"]
RUN curl -sL https://raw.githubusercontent.com/jorgebucaran/fisher/main/functions/fisher.fish | source && fisher install jorgebucaran/fisher
RUN fisher install eth-p/fish-plugin-sudo && bind \cs '__ethp_commandline_toggle_sudo'
RUN fisher install realiserad/fish-ai




