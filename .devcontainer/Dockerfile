# Use an official Go image as a parent image
FROM golang:1.20

# Install necessary packages
RUN apt-get update && apt-get install -y \
    libgtk-4-dev \
    libadwaita-1-dev \
    curl \
    git \
    neovim \
    flatpak \
    flatpak-builder \
    zsh && \
    apt-get clean

# Set Zsh as the default shell
RUN chsh -s /bin/zsh

# Install Go language server
RUN go install golang.org/x/tools/gopls@latest

# Install oh-my-zsh for a better zsh experience (optional)
RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended

# Set working directory
WORKDIR /workspace
