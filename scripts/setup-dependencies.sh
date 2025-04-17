#!/bin/bash

set -e

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if Docker is running
echo "Checking for Docker..."
if ! command_exists docker; then
    echo "Error: Docker is not installed. Please install Docker Desktop or Rancher Desktop."
    exit 1
fi
if ! docker ps >/dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker Desktop or Rancher Desktop."
    exit 1
fi
echo "Docker is installed and running."

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
echo "Detected OS: $OS"

# Install dependencies based on OS
case "$OS" in
    linux)
        # Update package manager
        if command_exists apt-get; then
            sudo apt-get update
        elif command_exists yum; then
            sudo yum update
        else
            echo "Warning: No supported package manager (apt/yum) found. Attempting direct downloads."
        fi

        # Install kind
        if ! command_exists kind; then
            echo "Installing kind..."
            curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.23.0/kind-linux-amd64
            chmod +x ./kind
            sudo mv ./kind /usr/local/bin/kind
        else
            echo "kind is already installed."
        fi

        # Install kubectl
        if ! command_exists kubectl; then
            echo "Installing kubectl..."
            curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
            chmod +x ./kubectl
            sudo mv ./kubectl /usr/local/bin/kubectl
        else
            echo "kubectl is already installed."
        fi

        # Install helm
        if ! command_exists helm; then
            echo "Installing helm..."
            curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
            chmod 700 get_helm.sh
            ./get_helm.sh
            rm get_helm.sh
        else
            echo "helm is already installed."
        fi
        ;;

    darwin)
        # Check for Homebrew
        if ! command_exists brew; then
            echo "Installing Homebrew..."
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        fi

        # Install kind
        if ! command_exists kind; then
            echo "Installing kind..."
            brew install kind
        else
            echo "kind is already installed."
        fi

        # Install kubectl
        if ! command_exists kubectl; then
            echo "Installing kubectl..."
            brew install kubectl
        else
            echo "kubectl is already installed."
        fi

        # Install helm
        if ! command_exists helm; then
            echo "Installing helm..."
            brew install helm
        else
            echo "helm is already installed."
        fi
        ;;

    *)
        echo "Error: Unsupported OS: $OS"
        exit 1
        ;;
esac

# Verify installations
echo "Verifying installations..."
for cmd in kind kubectl helm; do
    if command_exists "$cmd"; then
        echo "$cmd installed: $($cmd version --client 2>/dev/null || $cmd version)"
    else
        echo "Error: $cmd installation failed."
        exit 1
    fi
done

echo "All dependencies installed successfully! You can now use kindctl."