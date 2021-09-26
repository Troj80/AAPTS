cmake_minimum_required(VERSION 3.14.2)

project(cmake)

set(CMAKE_MODULE_PATH ${CMAKE_SOURCE_DIR}/cmake/modules)

set(GLOBAL_C_FLAGS " \
    -target x86_64-linux-gnu \
    -Wall \
    -Werror \
    -Wno-unused-parameter \
    -m64 \
    -fcolor-diagnostics \
    -fPIE \
    -ffunction-sections \
    -fdata-sections \
    -fvisibility=hidden \
    -Oz \
    -DNDEBUG \
    -mtune=native \
    -flto=thin \
    ")

set(GLOBAL_CXX_FLAGS " \
    -target x86_64-linux-gnu \
    -Wall \
    -Werror \
    -Wno-unused-command-line-argument \
    -Wno-unused-parameter \
    -isystem ${libcxx_headers} \
    -isystem ${libcxxabi_headers} \
    -std=gnu++17 \
    -nostdinc++ \
    -fcolor-diagnostics \
    -fPIE \
    -ffunction-sections \
    -fdata-sections \
    -fvisibility=hidden \
    -Oz \
    -DNDEBUG \
    -mtune=native \
    -flto=thin \
    ")

set(CMAKE_EXE_LINKER_FLAGS " \
    -target x86_64-linux-gnu \
    -nostdlib++ \
    -static \
    -fuse-ld=lld \
    -Wl,--gc-sections \
    -flto=thin \
    ")

include_directories(
    ${libcxx_headers}
    ${libcxxabi_headers}	
    )

add_subdirectory(lib)
add_subdirectory(aapt)
add_subdirectory(aapt2)

