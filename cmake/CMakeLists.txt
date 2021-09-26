cmake_minimum_required(VERSION 3.14.2)

project(cmake)

set(CMAKE_MODULE_PATH ${CMAKE_SOURCE_DIR}/cmake/modules)

set(MINGW $MINGW)

set(GLOBAL_C_FLAGS " \
    -target x86_64-pc-windows-gnu \
    -Wall \
    -Werror \
    -Wno-unused-parameter \
    -DUSE_MINGW \
    -DWIN32_LEAN_AND_MEAN \
    -D__STDC_FORMAT_MACROS \
    -D__STDC_CONSTANT_MACROS \
    -D__USE_MINGW_ANSI_STDIO=1 \
    -D_WIN32_WINNT=0x0601 \
    -DWINVER=0x0601 \
    -D_FILE_OFFSET_BITS=64 \
    -mno-ms-bitfields \
    -m64 \
    -ffunction-sections \
    -fdata-sections \
    -fvisibility=hidden \
    -Oz \
    -DNDEBUG \
    -mtune=native \
    -flto=thin \
    ")
if(UNIX)
    set(GLOBAL_C_FLAGS "${GLOBAL_C_FLAGS} \
        --sysroot $ENV{MINGW}/x86_64-w64-mingw32 \
        -isystem $ENV{MINGW}/x86_64-w64-mingw32/include \
        ")
endif()

set(GLOBAL_CXX_FLAGS " \
    -target x86_64-pc-windows-gnu \
    -Wall \
    -Werror \
    -Wno-unused-command-line-argument \
    -Wno-unused-parameter \
    -isystem ${libcxx_headers} \
    -isystem ${libcxxabi_headers} \
    -std=gnu++17 \
    -nostdinc++ \
    -D_LIBCPP_DISABLE_VISIBILITY_ANNOTATIONS \
    -D_LIBCXXABI_DISABLE_VISIBILITY_ANNOTATIONS \
    -D__STDC_FORMAT_MACROS \
    -D__STDC_CONSTANT_MACROS \
    -D_LIBCPP_HAS_THREAD_API_WIN32 \
    -DUSE_MINGW -DWIN32_LEAN_AND_MEAN \
    -D__USE_MINGW_ANSI_STDIO=1 \
    -D_WIN32_WINNT=0x0601 \
    -DWINVER=0x0601 \
    -D_FILE_OFFSET_BITS=64 \
    -mno-ms-bitfields \
    -ffunction-sections \
    -fdata-sections \
    -fvisibility=hidden \
    -Oz \
    -DNDEBUG \
    -mtune=native \
    -flto=thin \
    ")
if(UNIX)
    set(GLOBAL_CXX_FLAGS "${GLOBAL_CXX_FLAGS} \
        --sysroot $ENV{MINGW}/x86_64-w64-mingw32 \
        -isystem $ENV{MINGW}/x86_64-w64-mingw32/include \
        ")
endif()

set(CMAKE_EXE_LINKER_FLAGS " \
    -target x86_64-pc-windows-gnu \
    -nostdlib++ \
    -static \
    -static-libgcc \
    -Wl,--dynamicbase \
    -Wl,--nxcompat \
    -Wl,--Xlink=-Brepro \
    -Wl,--high-entropy-va \
    -fuse-ld=lld \
    -Wl,--gc-sections \
    -flto=thin \
    ")
if(UNIX)
    set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} \
        -L$ENV{MINGW}/x86_64-w64-mingw32/lib \
        -B$ENV{MINGW}/x86_64-w64-mingw32/lib \
        -B$ENV{MINGW}/x86_64-w64-mingw32/bin \
        -B$ENV{MINGW}/lib/gcc/x86_64-w64-mingw32/11.2.0 \
        -L$ENV{MINGW}/lib/gcc/x86_64-w64-mingw32/11.2.0 \
        ")
endif()

include_directories(
    ${libcxx_headers}
    ${libcxxabi_headers}	
    )

add_subdirectory(lib)
add_subdirectory(aapt)
add_subdirectory(aapt2)

