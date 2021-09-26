cmake_minimum_required(VERSION 3.14.2)

project(libbase)

set(BASE ${AOSP}/system/libbase)

set(LIBBASE_SRC_FILES
    ${BASE}/abi_compatibility.cpp
    ${BASE}/chrono_utils.cpp
    ${BASE}/cmsg.cpp
    ${BASE}/file.cpp
    ${BASE}/hex.cpp
    ${BASE}/logging.cpp
    ${BASE}/mapped_file.cpp
    ${BASE}/parsebool.cpp
    ${BASE}/parsenetaddress.cpp
    ${BASE}/process.cpp
    ${BASE}/properties.cpp
    ${BASE}/stringprintf.cpp
    ${BASE}/strings.cpp
    ${BASE}/threads.cpp
    ${BASE}/test_utils.cpp
    ${BASE}/errors_unix.cpp
    )

include_directories(
    ${liblog_headers}
    ${libbase_headers}
    ${fmtlib_headers}
    )


set(CMAKE_CXX_FLAGS "${GLOBAL_CXX_FLAGS} \
    -Wall \
    -Werror \
    -Wextra \
    -Wexit-time-destructors \
    -Wno-c99-designator \
    ")


add_library(base STATIC ${LIBBASE_SRC_FILES})

target_link_libraries(base
    -Wl,--whole-archive
    fmt
    -Wl,--no-whole-archive
    )
