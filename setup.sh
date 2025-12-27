#!/bin/bash
# Setup script for TMNT Pinball - Downloads SDL2 locally

set -e

echo "🐢 TMNT Pinball - Local SDL2 Setup 🍕"
echo "========================================"
echo ""

# Create lib directory
mkdir -p lib
cd lib

# Detect architecture
ARCH=$(uname -m)
OS=$(uname -s)

if [ "$OS" != "Darwin" ]; then
    echo "❌ This script is for macOS only."
    echo "For other platforms, please install SDL2 manually."
    exit 1
fi

echo "📦 Downloading SDL2 for macOS ($ARCH)..."

# SDL2 version
SDL2_VERSION="2.28.5"

# Download SDL2 framework for macOS
if [ ! -d "SDL2.framework" ]; then
    echo "Downloading SDL2 ${SDL2_VERSION}..."
    curl -L "https://github.com/libsdl-org/SDL/releases/download/release-${SDL2_VERSION}/SDL2-${SDL2_VERSION}.dmg" -o SDL2.dmg

    echo "Mounting SDL2 disk image..."
    hdiutil attach SDL2.dmg -quiet

    echo "Copying SDL2.framework to local lib directory..."
    cp -R "/Volumes/SDL2/SDL2.framework" .

    echo "Unmounting disk image..."
    hdiutil detach "/Volumes/SDL2" -quiet

    echo "Cleaning up..."
    rm SDL2.dmg

    echo "✅ SDL2.framework installed locally!"
else
    echo "✅ SDL2.framework already exists locally!"
fi

cd ..

# Update go.mod to use local SDL2
echo ""
echo "📝 Setting up Go environment..."

# Create a build script
cat > build.sh << 'BUILDSCRIPT'
#!/bin/bash
# Build script for TMNT Pinball with local SDL2

set -e

# Set environment variables for local SDL2
export CGO_CFLAGS="-F$(pwd)/lib -I$(pwd)/lib/SDL2.framework/Headers"
export CGO_LDFLAGS="-F$(pwd)/lib -framework SDL2 -Wl,-rpath,@executable_path/../lib"

echo "🔨 Building TMNT Pinball..."
go build -o tmnt-pinball

echo "✅ Build complete! Run with: ./tmnt-pinball"
BUILDSCRIPT

chmod +x build.sh

# Create a run script
cat > run.sh << 'RUNSCRIPT'
#!/bin/bash
# Run script for TMNT Pinball with local SDL2

set -e

# Set DYLD_FRAMEWORK_PATH to find local SDL2 framework
export DYLD_FRAMEWORK_PATH="$(pwd)/lib:${DYLD_FRAMEWORK_PATH}"

# Build if binary doesn't exist
if [ ! -f "tmnt-pinball" ]; then
    echo "🔨 Building first..."
    ./build.sh
fi

echo "🎮 Starting TMNT Pinball..."
./tmnt-pinball
RUNSCRIPT

chmod +x run.sh

echo ""
echo "✨ Setup complete! ✨"
echo ""
echo "To build and run the game:"
echo "  ./run.sh"
echo ""
echo "Or build separately and run:"
echo "  ./build.sh"
echo "  ./tmnt-pinball"
echo ""
echo "Cowabunga! 🐢🍕"
