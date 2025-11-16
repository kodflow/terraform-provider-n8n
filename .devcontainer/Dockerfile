FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

# Set shell with pipefail for RUN commands with pipes
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install HashiCorp tools and 1Password CLI
# hadolint ignore=DL3008
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        curl \
        gnupg \
        lsb-release \
        cmake \
        libpam0g-dev \
        jq \
        yq \
        zsh \
        unzip \
        wget \
        openjdk-17-jre-headless \
    # Add HashiCorp repository
    && mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /etc/apt/keyrings/hashicorp-archive-keyring.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/hashicorp.list \
    # Add 1Password repository
    && curl -fsSL https://downloads.1password.com/linux/keys/1password.asc | gpg --dearmor -o /etc/apt/keyrings/1password-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/$(dpkg --print-architecture) stable main" | tee /etc/apt/sources.list.d/1password.list \
    # Install packages
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        vault \
        consul \
        nomad \
        1password-cli

# Switch to vscode user
USER vscode
WORKDIR /home/vscode

# Create directories for persistent volumes
RUN mkdir -p /home/vscode/.cache \
    /home/vscode/.config \
    /home/vscode/.config/op \
    /home/vscode/.zsh_history_dir \
    /home/vscode/.local/bin \
    /home/vscode/.local/share/npm-global \
    /home/vscode/.claude \
    /home/vscode/.local/share/@anthropic \
    /home/vscode/.git-hooks \
    # Set proper permissions
    && chmod 700 /home/vscode/.config/op \
    && chmod -R 755 /home/vscode/.local

# Install Oh My Zsh and themes (rarely changes - good for cache)
# Using single-branch clone for faster downloads
RUN --mount=type=cache,target=/home/vscode/.cache/git,uid=1000,gid=1000 \
    rm -rf "$HOME/.oh-my-zsh" && \
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended && \
    git clone --depth=1 --single-branch https://github.com/romkatv/powerlevel10k.git \
        "${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/powerlevel10k" & \
    git clone --depth=1 --single-branch https://github.com/dracula/zsh.git \
        "${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/dracula" & \
    wait

# Configure zshrc (changes rarely - done before copying files)
RUN sed -i 's/^ZSH_THEME=.*/ZSH_THEME="powerlevel10k\/powerlevel10k"/' "$HOME/.zshrc" && \
    { \
        echo ''; \
        echo '# To customize prompt, run "p10k configure" or edit ~/.p10k.zsh'; \
        echo '[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh'; \
        echo ''; \
        echo '# Persistent zsh history'; \
        echo 'export HISTFILE=~/.zsh_history_dir/.zsh_history'; \
        echo 'export HISTSIZE=10000'; \
        echo 'export SAVEHIST=10000'; \
        echo ''; \
        echo '# npm global packages path (using NPM_CONFIG_PREFIX to avoid nvm conflicts)'; \
        echo 'export NPM_CONFIG_PREFIX="/home/vscode/.local/share/npm-global"'; \
        echo 'export PATH="/home/vscode/.local/share/npm-global/bin:$PATH"'; \
        echo ''; \
        echo '# Go environment (features install to /usr/local/go)'; \
        echo 'export GOPATH="/home/vscode/.cache/go"'; \
        echo 'export GOBIN="/home/vscode/.cache/go/bin"'; \
        echo 'export PATH="/usr/local/go/bin:$GOBIN:$PATH"'; \
        echo ''; \
        echo '# Local binaries'; \
        echo 'export PATH="/home/vscode/.local/bin:$PATH"'; \
        echo ''; \
        echo '# Super Claude alias'; \
        echo 'alias super-claude="claude --dangerously-skip-permissions --mcp-config /workspace/.devcontainer/mcp.json"'; \
    } >> "$HOME/.zshrc"

# Copy configuration files (moved to end to improve cache hit rate)
# These files change more frequently than system setup
COPY --chown=vscode:vscode .devcontainer/p10k.sh /home/vscode/.p10k.zsh
COPY --chown=vscode:vscode --chmod=755 .github/hooks/ /home/vscode/.git-hooks/

# Set working directory
WORKDIR /workspace
