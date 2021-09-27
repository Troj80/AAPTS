### This is a fork to build aapt/aapt2 for Linux and Windows
https://github.com/Lzhiyong/sdk-tools
https://github.com/MrIkso/sdk-tools/tree/apktool
____
Instrument preparation
```sh
mkir toolchains
cd toolchains
wget https://archlinux.org/packages/community/x86_64/mingw-w64-gcc/download
wget https://archlinux.org/packages/community/x86_64/mingw-w64-binutils/download
wget https://archlinux.org/packages/community/any/mingw-w64-crt/download
wget https://archlinux.org/packages/community/any/mingw-w64-headers/download
wget https://archlinux.org/packages/community/any/mingw-w64-winpthreads/download
for file in download*; do tar -I zstd -xvpf "${file}"; done
mv usr mingw-w64-11.2-linux-x86_64
rm .BUILDINFO .MTREE .PKGINFO download download.1 download.2 download.3 download.4
wget https://android.googlesource.com/platform/prebuilts/clang/host/linux-x86/+archive/refs/heads/master/clang-r433403.tar.gz
mkdir llvm-clang-13.0.2-linux-x86_64
tar -C llvm-clang-13.0.2-linux-x86_64 -xvpf clang-r433403.tar.gz
rm clang-r433403.tar.gz
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-linux-x86_64.zip
unzip protoc-3.9.1-linux-x86_64.zip -d protoc-3.9.1-linux-x86_64
rm protoc-3.9.1-linux-x86_64.zip
wget https://github.com/Kitware/CMake/releases/download/v3.21.3/cmake-3.21.3-linux-x86_64.tar.gz
tar -xvpf cmake-3.21.3-linux-x86_64.tar.gz
rm cmake-3.21.3-linux-x86_64.tar.gz
```
Build for Linux
```
git clone
cd 
export CC=/home/pasha/programs/toolchains/llvm-clang-13.0.2-linux-x86_64/bin/clang
export CXX=/home/pasha/programs/toolchains/llvm-clang-13.0.2-linux-x86_64/bin/clang++
export PATH=/home/pasha/programs/toolchains/protoc-3.9.1-linux-x86_64/bin:$PATH
export PATH=/home/pasha/programs/toolchains/cmake-3.21.3-linux-x86_64/bin:$PATH
cmake -G Ninja -S. -Bbuild && cmake --build build
```
Build for Windows
```
git clone
cd 
export CC=/home/pasha/programs/toolchains/llvm-clang-13.0.2-linux-x86_64/bin/clang
export CXX=/home/pasha/programs/toolchains/llvm-clang-13.0.2-linux-x86_64/bin/clang++
export MINGW=/home/pasha/programs/toolchains/mingw-w64-11.2-linux-x86_64
export PATH=/home/pasha/programs/toolchains/protoc-3.9.1-linux-x86_64/bin:$PATH
export PATH=/home/pasha/programs/toolchains/cmake-3.21.3-linux-x86_64/bin:$PATH
cmake -G Ninja -S. -Bbuild && cmake --build build
```
It is also possible to build in Windows aapt/aapt2 for Windows using MSYS2.
