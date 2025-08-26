#!/bin/bash




set -e


REPO="LumoSolutions/yerd"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="yerd"
YERD_DIRS="/opt/yerd"


RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color


print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}


detect_arch() {
    local arch
    arch=$(uname -m)
    case $arch in
        x86_64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        i386|i686)
            echo "386"
            ;;
        armv7l)
            echo "arm"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}


detect_os() {
    local os
    os=$(uname -s)
    case $os in
        Linux)
            echo "linux"
            ;;
        Darwin)
            echo "darwin"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
}


detect_distribution() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo "$ID"
    elif [ -f /etc/arch-release ]; then
        echo "arch"
    elif [ -f /etc/debian_version ]; then
        echo "debian"
    elif [ -f /etc/redhat-release ]; then
        echo "rhel"
    else
        echo "unknown"
    fi
}


command_exists() {
    command -v "$1" >/dev/null 2>&1
}


check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if ! command_exists curl && ! command_exists wget; then
        print_error "Neither curl nor wget found. Please install one of them."
        

        local distro
        distro=$(detect_distribution)
        case $distro in
            ubuntu|debian)
                print_status "Try: sudo apt update && sudo apt install curl wget"
                ;;
            arch|manjaro)
                print_status "Try: sudo pacman -S curl wget"
                ;;
            fedora)
                print_status "Try: sudo dnf install curl wget"
                ;;
            *)
                print_status "Please install curl or wget using your system's package manager"
                ;;
        esac
        exit 1
    fi
    
    if ! command_exists tar; then
        print_error "tar command not found. Please install tar."
        exit 1
    fi
    
    print_success "Prerequisites check passed"
}


get_latest_version() {
    local api_url="https://api.github.com/repos/${REPO}/releases/latest"
    
    if command_exists curl; then
        curl -s "$api_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/'
    elif command_exists wget; then
        wget -qO- "$api_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/'
    else
        print_error "Failed to get latest version"
        exit 1
    fi
}


install_yerd() {
    local version="$1"
    local os="$2"
    local arch="$3"
    

    version=${version#v}
    
    local filename="${BINARY_NAME}_${version}_${os}_${arch}.tar.gz"
    local download_url="https://github.com/${REPO}/releases/download/v${version}/${filename}"
    local temp_dir
    temp_dir=$(mktemp -d)
    
    print_status "Downloading YERD v${version} for ${os}/${arch}..."
    
    if command_exists curl; then
        if ! curl -sL "$download_url" -o "${temp_dir}/${filename}"; then
            print_error "Failed to download using curl"
            exit 1
        fi
    elif command_exists wget; then
        if ! wget -q "$download_url" -O "${temp_dir}/${filename}"; then
            print_error "Failed to download using wget"
            exit 1
        fi
    fi
    
    if [ ! -f "${temp_dir}/${filename}" ]; then
        print_error "Failed to download ${filename}"
        print_error "URL: ${download_url}"
        exit 1
    fi
    
    print_status "Extracting archive..."
    if ! tar -xzf "${temp_dir}/${filename}" -C "$temp_dir"; then
        print_error "Failed to extract archive"
        print_error "Archive might be corrupted or incompatible"
        exit 1
    fi
    
    if [ ! -f "${temp_dir}/${BINARY_NAME}" ]; then
        print_error "Binary not found in archive"
        exit 1
    fi
    
    print_status "Installing YERD to ${INSTALL_DIR}..."
    

    if [ -w "$INSTALL_DIR" ]; then
        cp "${temp_dir}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        print_warning "Installing to ${INSTALL_DIR} requires sudo privileges"
        sudo cp "${temp_dir}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    

    rm -rf "$temp_dir"
    
    print_success "YERD binary installed successfully"
}


setup_directories() {
    print_status "Setting up YERD directories..."
    
    if [ -w "$(dirname "$YERD_DIRS")" ]; then
        mkdir -p "${YERD_DIRS}"/{bin,php,etc}
    else
        print_warning "Setting up directories in ${YERD_DIRS} requires sudo privileges"
        sudo mkdir -p "${YERD_DIRS}"/{bin,php,etc}
        

        if [ "$EUID" -ne 0 ]; then
            read -p "Set user ownership for ${YERD_DIRS}? This allows installation without sudo. (y/N): " -r
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                sudo chown -R "$USER:$USER" "$YERD_DIRS"
                print_success "User ownership set for ${YERD_DIRS}"
            fi
        fi
    fi
    
    print_success "YERD directories created"
}


verify_installation() {
    print_status "Verifying installation..."
    
    if command_exists "$BINARY_NAME"; then
        local installed_version

        installed_version=$($BINARY_NAME --version 2>/dev/null | head -n1 | grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+' | sed 's/v//' || echo "unknown")
        if [ "$installed_version" = "unknown" ]; then
            print_success "YERD installed successfully!"
        else
            print_success "YERD v${installed_version} installed successfully!"
        fi
        print_status "Try: ${BINARY_NAME} --help"
    else
        print_error "Installation failed. ${BINARY_NAME} command not found."
        print_warning "Make sure ${INSTALL_DIR} is in your PATH"
        

        if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
            print_warning "${INSTALL_DIR} is not in your PATH"
            print_status "Add this line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
            print_status "export PATH=\"${INSTALL_DIR}:\$PATH\""
            print_status "Then run: source ~/.bashrc (or restart your shell)"
        fi
        exit 1
    fi
}


main() {
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                           YERD Installation Script                           ║"
    echo "║      A powerful, developer-friendly tool for managing PHP versions           ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo
    

    local force_version=""
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version)
                force_version="$2"
                shift 2
                ;;
            --help|-h)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --version VERSION    Install specific version"
                echo "  --help, -h          Show this help message"
                echo ""
                echo "Examples:"
                echo "  $0                  Install latest version"
                echo "  $0 --version 1.0.0  Install version 1.0.0"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    check_prerequisites
    
    local os arch version
    os=$(detect_os)
    arch=$(detect_arch)
    
    if [ -n "$force_version" ]; then
        version="$force_version"
        print_status "Installing specified version: v${version}"
    else
        print_status "Getting latest version from GitHub..."
        version=$(get_latest_version)
        if [ -z "$version" ]; then
            print_error "Failed to get latest version"
            exit 1
        fi
        print_status "Latest version: ${version}"
    fi
    
    install_yerd "$version" "$os" "$arch"
    setup_directories
    verify_installation
    
    echo
    print_success "🎉 YERD installation completed!"
    echo
    echo "Next steps:"
    echo "1. Run: yerd --help"
    echo "2. Check system status: yerd status"  
    echo "3. List PHP versions: yerd php list"
    echo "4. Install PHP: sudo yerd php add 8.4"
    echo
}

# Call main function with all arguments
main "$@"