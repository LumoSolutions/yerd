#!/bin/bash




set -e


VERSION=${VERSION:-"1.0.0"}
OUTPUT_DIR="dist"
BINARY_NAME="yerd"


RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}


declare -a platforms=(
    "linux/amd64"
    "linux/arm64"
    "linux/386"
    "darwin/amd64"
    "darwin/arm64"
)


prepare_output() {
    print_status "Preparing output directory..."
    rm -rf "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR"
}


build_platform() {
    local platform=$1
    local os=${platform%/*}
    local arch=${platform#*/}
    local output_name="$BINARY_NAME"
    
    local output_path="${OUTPUT_DIR}/${output_name}"
    local archive_name="${BINARY_NAME}_${VERSION}_${os}_${arch}"
    
    print_status "Building for ${os}/${arch}..."
    

    GOOS=$os GOARCH=$arch go build -o "$output_path" \
        -ldflags="-s -w -X github.com/LumoSolutions/yerd/internal/version.Version=${VERSION}" \
        .
    
    if [ ! -f "$output_path" ]; then
        print_error "Build failed for ${os}/${arch}"
        return 1
    fi
    

    print_status "Creating archive for ${os}/${arch}..."
    
    (cd "$OUTPUT_DIR" && tar -czf "${archive_name}.tar.gz" "$output_name")
    

    rm "$output_path"
    
    print_success "Built ${archive_name}"
}


generate_checksums() {
    print_status "Generating checksums..."
    
    cd "$OUTPUT_DIR"
    

    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum * > checksums.txt
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 * > checksums.txt
    else
        print_error "No checksum utility found"
        return 1
    fi
    
    cd ..
    print_success "Checksums generated"
}


main() {

    while [[ $# -gt 0 ]]; do
        case $1 in
            --version)
                VERSION="$2"
                shift 2
                ;;
            --help|-h)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --version VERSION    Build specific version (default: 1.0.0)"
                echo "  --help, -h          Show this help message"
                echo ""
                echo "Environment variables:"
                echo "  VERSION             Version to build (default: 1.0.0)"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                           YERD Release Builder                              ║"
    echo "║                          Building v${VERSION}                                   ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo
    
    prepare_output
    

    local failed_builds=0
    for platform in "${platforms[@]}"; do
        if ! build_platform "$platform"; then
            ((failed_builds++))
        fi
    done
    
    if [ $failed_builds -gt 0 ]; then
        print_error "$failed_builds builds failed"
        exit 1
    fi
    
    generate_checksums
    
    print_success "All builds completed successfully!"
    print_status "Artifacts created in: $OUTPUT_DIR"
    
    echo
    echo "Built packages:"
    ls -la "$OUTPUT_DIR"
}


if ! command -v go >/dev/null 2>&1; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

main "$@"