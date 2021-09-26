cmake_minimum_required(VERSION 3.14.2)

project(libcutils)

set(LIBCUTILS ${AOSP}/system/core/libcutils)

set(LIBCUTILS_SRC_FILES
    ${LIBCUTILS}/config_utils.cpp
    ${LIBCUTILS}/canned_fs_config.cpp
    ${LIBCUTILS}/iosched_policy.cpp
    ${LIBCUTILS}/load_file.cpp
    ${LIBCUTILS}/native_handle.cpp
    ${LIBCUTILS}/properties.cpp
    ${LIBCUTILS}/record_stream.cpp
    ${LIBCUTILS}/strlcpy.c
    ${LIBCUTILS}/threads.cpp
    ${LIBCUTILS}/trace-host.cpp

    ${LIBCUTILS}/sockets.cpp
    ${LIBCUTILS}/socket_inaddr_any_server_windows.cpp
    ${LIBCUTILS}/socket_network_client_windows.cpp
    ${LIBCUTILS}/sockets_windows.cpp
    )

include_directories(
    ${libcutils_headers}
    ${liblog_headers}
    ${libbase_headers}
    ${libutils_headers}
    )


set(CMAKE_CXX_FLAGS "${GLOBAL_CXX_FLAGS} \
    -Wno-exit-time-destructors \
    -Werror \
    -Wall \
    -Wextra \
    -D_GNU_SOURCE \
    ")


add_library(cutils STATIC ${LIBCUTILS_SRC_FILES})

target_link_libraries(cutils
    ws2_32
    )
