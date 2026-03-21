# syntax=docker/dockerfile:1
# hadolint global ignore=DL3008
# =============================================================================
# Base Development Image
# =============================================================================
# Single base image with all tools + Claude Code (standalone)
# Languages (including Node.js) are added via devcontainer features at runtime
#
# Build:
#   docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/kodflow/devcontainer-template:latest .
#
# Layer optimization: ~15 real layers (consolidated from ~37)
# =============================================================================

# Base image uses tagged version for multi-arch support (amd64/arm64)
# Note: Digest pinning not used as it would require platform-specific digests
# The devcontainers base images are published by Microsoft with regular security updates
ARG BASE_IMAGE=mcr.microsoft.com/devcontainers/base:ubuntu-24.04
# hadolint ignore=DL3006
FROM ${BASE_IMAGE}

ARG TARGETARCH

ENV DEBIAN_FRONTEND=noninteractive \
    LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8

# =============================================================================
# System Dependencies (consolidated: core + network + VPN + desktop + python3-venv)
# =============================================================================
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        curl ca-certificates gnupg lsb-release \
        build-essential g++ gcc make cmake pkg-config \
        git jq zsh unzip zip tar xz-utils file \
        libssl-dev libpam0g-dev software-properties-common \
        libffi-dev libbz2-dev libreadline-dev libsqlite3-dev \
        libncurses5-dev libncursesw5-dev liblzma-dev zlib1g-dev \
        libgdbm-dev libyaml-dev \
        libsecret-1-0 libsecret-1-dev \
        shellcheck python3-venv \
        iputils-ping dnsutils net-tools nmap \
        traceroute mtr-tiny tcpdump \
        netcat-openbsd wget whois iperf3 \
        openvpn wireguard-tools strongswan pptp-linux \
        libwebkit2gtk-4.1-dev \
        libxdo-dev \
        libayatana-appindicator3-dev \
        librsvg2-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# CACHE_BUST_DYNAMIC: Forces cache invalidation for all dynamic installations.
# Set to current date (YYYY-MM-DD) during scheduled builds to pull latest versions.
# See: https://docs.docker.com/build/cache/invalidation/
ARG CACHE_BUST_DYNAMIC=stable
ARG GITHUB_TOKEN=

# =============================================================================
# PPA-Based Tools (HashiCorp, 1Password, Ansible) - consolidated
# =============================================================================
RUN curl -fsSL https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" > /etc/apt/sources.list.d/hashicorp.list && \
    curl -fsSL https://downloads.1password.com/linux/keys/1password.asc | gpg --dearmor -o /usr/share/keyrings/1password-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/$(dpkg --print-architecture) stable main" > /etc/apt/sources.list.d/1password.list && \
    mkdir -p /etc/debsig/policies/AC2D62742012EA22/ && \
    curl -fsSL https://downloads.1password.com/linux/debian/debsig/1password.pol -o /etc/debsig/policies/AC2D62742012EA22/1password.pol && \
    mkdir -p /usr/share/debsig/keyrings/AC2D62742012EA22 && \
    curl -fsSL https://downloads.1password.com/linux/keys/1password.asc | gpg --dearmor -o /usr/share/debsig/keyrings/AC2D62742012EA22/debsig.gpg && \
    add-apt-repository --yes --update ppa:ansible/ansible && \
    apt-get install -y --no-install-recommends \
        terraform vault consul nomad packer terraform-ls \
        1password-cli \
        ansible && \
    rm -rf /var/lib/apt/lists/*

# =============================================================================
# Cloud CLIs (AWS, GCP, Azure, GitHub) - consolidated
# =============================================================================
RUN ARCH_AWS=$([ "$TARGETARCH" = "arm64" ] && echo "aarch64" || echo "x86_64") && \
    curl -fsSL "https://awscli.amazonaws.com/awscli-exe-linux-${ARCH_AWS}.zip" -o /tmp/awscliv2.zip && \
    unzip -q /tmp/awscliv2.zip -d /tmp && /tmp/aws/install && rm -rf /tmp/aws /tmp/awscliv2.zip && \
    curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" > /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && \
    chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" > /etc/apt/sources.list.d/github-cli.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends google-cloud-cli google-cloud-cli-gke-gcloud-auth-plugin gh && \
    rm -rf /var/lib/apt/lists/* && \
    # Azure CLI: use official Microsoft apt repository (auditable, consistent with GCP/GitHub pattern)
    curl -fsSL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor -o /usr/share/keyrings/microsoft.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends azure-cli && \
    rm -rf /var/lib/apt/lists/*

# =============================================================================
# Kubernetes + Binary Tools (kubectl, Helm, Bazelisk, yq, grepai, rtk) - consolidated
# =============================================================================
RUN ARCH=$([ "$TARGETARCH" = "arm64" ] && echo "arm64" || echo "amd64") && \
    KUBECTL_LATEST=$(curl -fsSL "https://dl.k8s.io/release/stable.txt") && \
    echo "Installing kubectl ${KUBECTL_LATEST}..." && \
    curl -fsSL "https://dl.k8s.io/release/${KUBECTL_LATEST}/bin/linux/${ARCH}/kubectl" -o /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    HELM_LATEST=$(curl -fsSL ${GITHUB_TOKEN:+-H "Authorization: token ${GITHUB_TOKEN}"} "https://api.github.com/repos/helm/helm/releases/latest" | jq -r '.tag_name') && \
    echo "Installing Helm ${HELM_LATEST}..." && \
    curl -fsSL "https://get.helm.sh/helm-${HELM_LATEST}-linux-${ARCH}.tar.gz" | tar xz -C /tmp && \
    mv /tmp/linux-${ARCH}/helm /usr/local/bin/helm && \
    rm -rf /tmp/linux-${ARCH} && \
    echo "Installing Bazelisk (latest)..." && \
    curl -fsSL "https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-${ARCH}" -o /usr/local/bin/bazel && \
    chmod +x /usr/local/bin/bazel && \
    echo "Installing yq (latest)..." && \
    curl -fsSL "https://github.com/mikefarah/yq/releases/latest/download/yq_linux_${ARCH}" -o /usr/bin/yq && \
    chmod +x /usr/bin/yq && \
    GREPAI_LATEST=$(curl -fsSL ${GITHUB_TOKEN:+-H "Authorization: token ${GITHUB_TOKEN}"} "https://api.github.com/repos/yoanbernabeu/grepai/releases/latest" | jq -r '.tag_name') && \
    echo "Installing grepai ${GREPAI_LATEST}..." && \
    curl -fsSL "https://github.com/yoanbernabeu/grepai/releases/download/${GREPAI_LATEST}/grepai_${GREPAI_LATEST#v}_linux_${ARCH}.tar.gz" | \
    tar xz -C /usr/local/bin grepai && \
    RTK_LATEST=$(curl -fsSL ${GITHUB_TOKEN:+-H "Authorization: token ${GITHUB_TOKEN}"} "https://api.github.com/repos/rtk-ai/rtk/releases/latest" | jq -r '.tag_name') && \
    echo "Installing rtk ${RTK_LATEST}..." && \
    RTK_ARCH=$([ "$TARGETARCH" = "arm64" ] && echo "aarch64-unknown-linux-gnu" || echo "x86_64-unknown-linux-musl") && \
    curl -fsSL "https://github.com/rtk-ai/rtk/releases/download/${RTK_LATEST}/rtk-${RTK_ARCH}.tar.gz" | \
    tar xz -C /usr/local/bin rtk

# =============================================================================
# MkDocs Material (Documentation Server)
# =============================================================================
# python3-venv installed in base dependencies above
# Uses virtual environment to avoid system package conflicts (Markdown 3.5.2)
# Note: mkdocs-material >= 8.2.0 has native Mermaid support (no plugin needed)
RUN python3 -m venv /opt/mkdocs && \
    /opt/mkdocs/bin/pip install --no-cache-dir mkdocs-material && \
    ln -s /opt/mkdocs/bin/mkdocs /usr/local/bin/mkdocs

# =============================================================================
# User Setup (directories + Oh My Zsh + shell config) - consolidated
# =============================================================================
USER vscode
WORKDIR /home/vscode

# Copy p10k configuration (referenced at runtime by .zshrc)
COPY --chown=vscode:vscode .p10k.zsh /home/vscode/.p10k.zsh

# Create cache/config directories + install Oh My Zsh + configure shells
RUN mkdir -p .cache .config .local/bin .local/share .zsh_history_dir \
    .cache/{terraform,helm,aws,bazel} \
    .config/{aws,gcloud,azure,op} .kube \
    .config/{openvpn,wireguard,strongswan,pptp} \
    .claude .config/@anthropic .cache/@anthropic .local/share/@anthropic && \
    rm -rf ~/.oh-my-zsh && \
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended && \
    git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ~/.oh-my-zsh/custom/themes/powerlevel10k && \
    git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions && \
    git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting && \
    sed -i 's/^ZSH_THEME=.*/ZSH_THEME="powerlevel10k\/powerlevel10k"/' ~/.zshrc && \
    sed -i 's/^plugins=.*/plugins=(git zsh-autosuggestions zsh-syntax-highlighting)/' ~/.zshrc && \
    printf '\n[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh\nexport HISTFILE=~/.zsh_history_dir/.zsh_history\nexport HISTSIZE=50000\nexport SAVEHIST=50000\n[[ -f ~/.devcontainer-env.sh ]] && source ~/.devcontainer-env.sh\n' >> ~/.zshrc && \
    printf '\n[[ -f ~/.devcontainer-env.sh ]] && source ~/.devcontainer-env.sh\n' >> ~/.bashrc

# =============================================================================
# Claude Code (standalone installation - no Node.js required)
# =============================================================================
# Using official installer: https://claude.ai/install.sh
# This avoids requiring Node.js/npm in the base image.
# If you need Node.js for your project, enable the "nodejs" feature.
# Note: The official installer creates the symlink at ~/.local/bin/claude automatically.
RUN curl -fsSL https://claude.ai/install.sh | bash && \
    mkdir -p /home/vscode/.claude/sessions

# =============================================================================
# Claude Configuration
# =============================================================================
COPY --chown=vscode:vscode .claude/ /home/vscode/.claude/

# =============================================================================
# System Templates (root-owned, copied at runtime via postStart.sh)
# =============================================================================
USER root

# Claude defaults backup (full directory for volume restoration)
COPY --chown=root:root .claude/ /etc/claude-defaults/

# VPN scripts (staged in /tmp for renaming in consolidated RUN below)
COPY --chown=root:root scripts/vpn/ /tmp/vpn/

# MCP template (secrets injected at runtime via postStart.sh)
COPY --chown=vscode:vscode mcp.json.tpl /etc/mcp/mcp.json.tpl

# GrepAI config template (optimized for 12 languages)
COPY --chown=vscode:vscode grepai.config.yaml /etc/grepai/config.yaml

# RTK config template (token savings proxy, root-owned system template)
COPY --chown=root:root rtk.config.toml /etc/rtk/config.toml

# Lifecycle hooks embedded in image (auto-updated on rebuild)
# Workspace stubs delegate to /etc/devcontainer-hooks/ at runtime
COPY --chown=root:root hooks/ /etc/devcontainer-hooks/

# Consolidated system setup: permissions, VPN scripts, sudoers, default shell
RUN chmod -R 755 /etc/claude-defaults/ /home/vscode/.claude/scripts/ /etc/devcontainer-hooks/ && \
    mkdir -p /etc/openvpn && \
    cp /tmp/vpn/update-dns.sh /etc/openvpn/update-dns && \
    cp /tmp/vpn/vpn-connect.sh /usr/local/bin/vpn-connect && \
    cp /tmp/vpn/vpn-disconnect.sh /usr/local/bin/vpn-disconnect && \
    cp /tmp/vpn/vpn-status.sh /usr/local/bin/vpn-status && \
    rm -rf /tmp/vpn && \
    chmod +x /etc/openvpn/update-dns /usr/local/bin/vpn-connect \
             /usr/local/bin/vpn-disconnect /usr/local/bin/vpn-status && \
    printf '%s\n' \
      'vscode ALL=(ALL) NOPASSWD: /usr/sbin/openvpn' \
      'vscode ALL=(ALL) NOPASSWD: /usr/bin/killall openvpn' \
      'vscode ALL=(ALL) NOPASSWD: /usr/bin/wg' \
      'vscode ALL=(ALL) NOPASSWD: /usr/bin/wg-quick' \
      'vscode ALL=(ALL) NOPASSWD: /usr/sbin/ipsec' \
      'vscode ALL=(ALL) NOPASSWD: /usr/bin/killall charon' \
      'vscode ALL=(ALL) NOPASSWD: /usr/sbin/pppd' \
      'vscode ALL=(ALL) NOPASSWD: /usr/sbin/pptp' \
      > /etc/sudoers.d/vpn && \
    chmod 440 /etc/sudoers.d/vpn && \
    chsh -s /bin/zsh vscode

USER vscode

# =============================================================================
# CodeRabbit CLI (AI-powered code review)
# =============================================================================
# Using official installer: https://cli.coderabbit.ai/install.sh
# Provides: coderabbit (cr) - local AI code reviews
RUN curl -fsSL https://cli.coderabbit.ai/install.sh | sh

# =============================================================================
# Qodo Command CLI (AI-powered code review + agents)
# =============================================================================
# Provides: qodo - local AI code reviews, chains, custom agents
# Docs: https://docs.qodo.ai/qodo-documentation/qodo-command
# Note: Node.js downloaded temporarily for npm install (not kept in image)
# Runtime requires Node.js via devcontainer feature
# hadolint ignore=DL3003,DL4006,SC2086
RUN NODE_ARCH=$([ "$TARGETARCH" = "arm64" ] && echo "arm64" || echo "x64") && \
    NODE_VERSION=$(curl -fsSL https://nodejs.org/dist/index.json | jq -r '[.[] | select(.lts != false)][0].version') && \
    curl -fsSL "https://nodejs.org/dist/${NODE_VERSION}/node-${NODE_VERSION}-linux-${NODE_ARCH}.tar.xz" | tar xJ -C /tmp && \
    export PATH="/tmp/node-${NODE_VERSION}-linux-${NODE_ARCH}/bin:$PATH" && \
    npm install -g --prefix="$HOME/.local" @qodo/command@0.36.0 && \
    rm -rf "/tmp/node-${NODE_VERSION}-linux-${NODE_ARCH}"

# =============================================================================
# User Tools (status-line, ktn-linter, xdg-open) - consolidated
# =============================================================================
RUN ARCH=$([ "$TARGETARCH" = "arm64" ] && echo "arm64" || echo "amd64") && \
    mkdir -p ~/.local/bin && \
    echo "Installing status-line (latest)..." && \
    curl -fsSL "https://github.com/kodflow/status-line/releases/latest/download/status-line-linux-${ARCH}" \
      -o ~/.local/bin/status-line && \
    curl -fsSL "https://github.com/kodflow/status-line/releases/latest/download/status-line-linux-${ARCH}.sha256" \
      -o /tmp/status-line.sha256 && \
    sed -i "s|status-line-linux-${ARCH}|$HOME/.local/bin/status-line|" /tmp/status-line.sha256 && \
    sha256sum -c /tmp/status-line.sha256 && rm -f /tmp/status-line.sha256 && \
    chmod +x ~/.local/bin/status-line && \
    echo "Installing ktn-linter (latest)..." && \
    curl -fsSL "https://github.com/kodflow/ktn-linter/releases/latest/download/ktn-linter-linux-${ARCH}" \
      -o ~/.local/bin/ktn-linter && \
    chmod +x ~/.local/bin/ktn-linter && \
    printf '%s\n' '#!/bin/sh' \
    '[ $# -eq 0 ] && { echo "Usage: xdg-open <url> [url...]" >&2; exit 2; }' \
    'for url in "$@"; do printf "\n  Open in browser: %s\n" "$url"; done' \
    'echo ""' \
    > ~/.local/bin/xdg-open && chmod +x ~/.local/bin/xdg-open

# =============================================================================
# Environment Variables
# =============================================================================
ENV TF_PLUGIN_CACHE_DIR=/home/vscode/.cache/terraform \
    HELM_CACHE_HOME=/home/vscode/.cache/helm \
    KUBECONFIG=/home/vscode/.kube/config \
    BAZEL_USER_ROOT=/home/vscode/.cache/bazel \
    CLAUDE_CONFIG_DIR=/home/vscode/.claude \
    ENABLE_TOOL_SEARCH=auto \
    PATH=/home/vscode/.local/bin:$PATH

WORKDIR /workspace

# =============================================================================
# Labels
# =============================================================================
LABEL org.opencontainers.image.source="https://github.com/kodflow/devcontainer-template" \
      org.opencontainers.image.description="Base Development Image - All tools + Claude Code included" \
      org.opencontainers.image.licenses="MIT"
