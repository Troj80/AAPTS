cmake_minimum_required(VERSION 3.14.2)

project(libincfs-utils)

set(INCFS-UTILS ${AOSP}/system/incremental_delivery/util)

set(LIBINCFS-UTILS_SRC_FILES
    ${INCFS-UTILS}/map_ptr.cpp
    )

include_directories(
    ${libincfs-utils_headers}
    ${libbase_headers}
    ${libutils_headers}
    )


set(CMAKE_CXX_FLAGS "${GLOBAL_CXX_FLAGS} \
    -Werror \
    -Wall \
    -Wextra \
    -std=c++2a \
    ")


add_library(incfs-utils STATIC ${LIBINCFS-UTILS_SRC_FILES})

target_link_libraries(incfs-utils)
