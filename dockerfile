FROM golang:1.20.1-bullseye

ENV KEYRINGS /usr/local/share/keyrings
ENV CARGO_ROOT /root/.cargo/bin

RUN set -eux; \
    mkdir -p $KEYRINGS; \
    apt-get update && apt-get install -y gpg curl; \
    # clang/lld/llvm
    curl --fail https://apt.llvm.org/llvm-snapshot.gpg.key | gpg --dearmor > $KEYRINGS/llvm.gpg; \
    # wine
    curl --fail https://dl.winehq.org/wine-builds/winehq.key | gpg --dearmor > $KEYRINGS/winehq.gpg; \
    echo "deb [signed-by=$KEYRINGS/llvm.gpg] http://apt.llvm.org/bullseye/ llvm-toolchain-bullseye-13 main" > /etc/apt/sources.list.d/llvm.list; \
    echo "deb [signed-by=$KEYRINGS/winehq.gpg] https://dl.winehq.org/wine-builds/debian/ bullseye main" > /etc/apt/sources.list.d/winehq.list;

RUN set -eux; \
    dpkg --add-architecture i386; \
    # Skipping all of the "recommended" cruft reduces total images size by ~300MiB
    apt-get update && apt-get install --no-install-recommends -y \
    clang-13 \
    # llvm-ar
    llvm-13 \
    lld-13 \
    # get a recent wine so we can run tests
    winehq-staging \
    # Unpack xwin
    tar; \
    # ensure that clang/clang++ are callable directly
    ln -s clang-13 /usr/bin/clang && ln -s clang /usr/bin/clang++ && ln -s lld-13 /usr/bin/ld.lld; \
    # We also need to setup symlinks ourselves for the MSVC shims because they aren't in the debian packages
    ln -s clang-13 /usr/bin/clang-cl && ln -s llvm-ar-13 /usr/bin/llvm-lib && ln -s lld-link-13 /usr/bin/lld-link; \
    # Verify the symlinks are correct
    clang++ -v; \
    ld.lld -v; \
    # Doesn't have an actual -v/--version flag, but it still exits with 0
    llvm-lib -v; \
    clang-cl -v; \
    lld-link --version; \
    # Use clang instead of gcc when compiling binaries targeting the host (eg proc macros, build files)
    update-alternatives --install /usr/bin/cc cc /usr/bin/clang 100; \
    update-alternatives --install /usr/bin/c++ c++ /usr/bin/clang++ 100; \
    apt-get remove -y --auto-remove; \
    rm -rf /var/lib/apt/lists/*;

# Install rustup


RUN set -eux; \
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --default-toolchain stable; \
    xwin_version="0.1.1"; \
    xwin_prefix="xwin-$xwin_version-x86_64-unknown-linux-musl"; \
    # Install xwin to cargo/bin via github release. Note you could also just use `cargo install xwin`.
    ${CARGO_ROOT}/cargo install xwin && \
    # Splat the CRT and SDK files to /xwin/crt and /xwin/sdk respectively
    ${CARGO_ROOT}/xwin --accept-license splat --output /xwin; \
    # Remove unneeded files to reduce image size
    rm -rf .xwin-cache /usr/local/cargo/bin/xwin;

# Note that we're using the full target triple for each variable instead of the
# simple CC/CXX/AR shorthands to avoid issues when compiling any C/C++ code for
# build dependencies that need to compile and execute in the host environment
ENV CC_x86_64_pc_windows_msvc="clang-cl" \
    CXX_x86_64_pc_windows_msvc="clang-cl" \
    AR_x86_64_pc_windows_msvc="llvm-lib" \
    GOCC="clang-cl" \
    GOOS="windows" \
    GOARCH="amd64" \
    CGO_ENABLED="1" \
    FILENAME="app.exe" \

    # wine can be quite spammy with log messages and they're generally uninteresting
    WINEDEBUG="-all" \
    # Use wine to run test executables
    CARGO_TARGET_X86_64_PC_WINDOWS_MSVC_RUNNER="wine" \
    # Note that we only disable unused-command-line-argument here since clang-cl
    # doesn't implement all of the options supported by cl, but the ones it doesn't
    # are _generally_ not interesting.
    CL_FLAGS="-Wno-error=unused-command-line-argument -fuse-ld=lld-link /imsvc/xwin/crt/include /imsvc/xwin/sdk/include/ucrt /imsvc/xwin/sdk/include/um /imsvc/xwin/sdk/include/shared"

# These are separate since docker/podman won't transform environment variables defined in the same ENV block
ENV CFLAGS_x86_64_pc_windows_msvc="$CL_FLAGS" \
    CXXFLAGS_x86_64_pc_windows_msvc="$CL_FLAGS" \
    CFLAGS="$CL_FLAGS" \
    CXXFLAGS="$CL_FLAGS" \
    CC="clang-cl" \
    CXX="clang-cl" \
    CGO_CFLAGS="$CL_FLAGS" \
    CGO_CXXFLAGS="$CL_FLAGS" 

# Run wineboot just to setup the default WINEPREFIX so we don't do it every
# container run
RUN wine wineboot --init

WORKDIR /usr/src/app

CMD ["go", "build", "-v", "-o", "${FILENAME}", "-ldflags", "-extld=$CC", "main.go"]