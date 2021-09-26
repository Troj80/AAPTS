cmake_minimum_required(VERSION 3.14.2)

project(lib)

set(CMAKE_MODULE_PATH ${CMAKE_SOURCE_DIR}/cmake/modules)

# check host arch
execute_process(COMMAND uname -m 
    COMMAND tr -d '\n' 
    OUTPUT_VARIABLE 
    HOST_ARCH
    )

add_subdirectory(androidfw)
add_subdirectory(expat)
add_subdirectory(fmtlib)
add_subdirectory(libbase)
add_subdirectory(libbuildversion)
add_subdirectory(libcutils)
add_subdirectory(libcxx)
add_subdirectory(libcxxabi)
add_subdirectory(libincfs-utils)
add_subdirectory(liblog)
add_subdirectory(libpng)
add_subdirectory(libutils)
add_subdirectory(libziparchive)
add_subdirectory(protobuf)
add_subdirectory(zlib)
