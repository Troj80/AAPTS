cmake_minimum_required(VERSION 3.14.2)

project(fmtlib)

set(FMTLIB ${AOSP}/external/fmtlib)

set(FMTLIB_SRC_FILES
    ${FMTLIB}/src/format.cc
    )

include_directories(
    ${fmtlib_headers}
    )


set(CMAKE_CXX_FLAGS "${GLOBAL_CXX_FLAGS} \
    -Wall \
    -Werror \
    -Wextra \
    -fno-exceptions \
    -UNDEBUG \
    ")


add_library(fmt STATIC ${FMTLIB_SRC_FILES})

target_link_libraries(fmt)
