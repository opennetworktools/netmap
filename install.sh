#!/bin/sh

# This script is for installing the latest version of Netmap CLI on your machine.

set -e

BINARY_NAME="netmap"
VERSION="v0.1.3"

probe_arch() {
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="x86_64"  ;;
        aarch64) ARCH="arm64" ;;
        arm64) ARCH="arm64" ;;
        *) printf "Architecture ${ARCH} is not supported by this installation script\n"; exit 1 ;;
    esac
}

probe_os() {
    OS=$(uname -s)
    case $OS in
        Darwin) OS="Darwin" ;;
        Linux) OS="Linux" ;;
        *) printf "Operating system ${OS} is not supported by this installation script\n"; exit 1 ;;
    esac
}

install_netmap_cli() {
    printf "\nInstalling Netmap!\n"

    case $ARCH in
        x86_64) ARCH_TARGET="amd64" ;;
        aarch64) ARCH_TARGET="aarch64" ;;
        arm64) ARCH_TARGET="arm64" ;;
        *)
            printf "Architecture ${ARCH} is not supported for netmap\n"
            return 1
            ;;
    esac

    case $OS in
        Darwin) OS_TARGET="darwin" ;;
        Linux) OS_TARGET="linux" ;;
        *)
            printf "Operating system ${OS} is not supported for netmap\n"
            return 1
            ;;
    esac

    NETMAP_TARGET="${BINARY_NAME}-${OS_TARGET}-${ARCH_TARGET}-${VERSION}"

    # URL for MacOS AMD64: https://github.com/opennetworktools/netmap/releases/download/v0.1.3/netmap-darwin-amd64-v0.1.3
    NETMAP_URL_PREFIX="https://github.com/opennetworktools/netmap/releases/download/${VERSION}"
    NETMAP_URL="${NETMAP_URL_PREFIX}/${NETMAP_TARGET}"

    # printf "\n$NETMAP_URL"

    curl --progress-bar -L -o $BINARY_NAME $NETMAP_URL
    chmod +x $BINARY_NAME
    sudo mv $BINARY_NAME /usr/local/bin/$BINARY_NAME
}

main() {
  printf "\nWelcome to the Netmap CLI installer!\n"

  probe_arch
  probe_os

  install_netmap_cli

  printf "\nNetmap CLI installed!\n" 
  Printf "\nRun \"netmap version\" to verify the installation."
  printf "\nRun \"netmap help\" to get started!\n\n"
}

main