cmake_minimum_required(VERSION 3.14.2)

project(libcxx)

set(LIBCXX ${AOSP}/external/libcxx)

set(LIBCXX_SRC_FILES
    ${LIBCXX}/src/algorithm.cpp
    ${LIBCXX}/src/any.cpp
    ${LIBCXX}/src/bind.cpp
    ${LIBCXX}/src/charconv.cpp
    ${LIBCXX}/src/chrono.cpp
    ${LIBCXX}/src/condition_variable.cpp
    ${LIBCXX}/src/debug.cpp
    ${LIBCXX}/src/exception.cpp
    ${LIBCXX}/src/future.cpp
    ${LIBCXX}/src/hash.cpp
    ${LIBCXX}/src/ios.cpp
    ${LIBCXX}/src/iostream.cpp
    ${LIBCXX}/src/locale.cpp
    ${LIBCXX}/src/memory.cpp
    ${LIBCXX}/src/mutex.cpp
    ${LIBCXX}/src/new.cpp
    ${LIBCXX}/src/optional.cpp
    ${LIBCXX}/src/random.cpp
    ${LIBCXX}/src/regex.cpp
    ${LIBCXX}/src/shared_mutex.cpp
    ${LIBCXX}/src/stdexcept.cpp
    ${LIBCXX}/src/string.cpp
    ${LIBCXX}/src/strstream.cpp
    ${LIBCXX}/src/system_error.cpp
    ${LIBCXX}/src/thread.cpp
    ${LIBCXX}/src/typeinfo.cpp
    ${LIBCXX}/src/utility.cpp
    ${LIBCXX}/src/valarray.cpp
    ${LIBCXX}/src/variant.cpp
    ${LIBCXX}/src/vector.cpp
    ${LIBCXX}/src/support/win32/locale_win32.cpp
    ${LIBCXX}/src/support/win32/support.cpp
    ${LIBCXX}/src/support/win32/thread_win32.cpp
    )


set(CMAKE_CXX_FLAGS "${GLOBAL_CXX_FLAGS} \
    -Wall \
    -Werror \
    -Wno-unused-parameter \
    -std=c++14 \
    -frtti \
    -fexceptions \
    -DLIBCXX_BUILDING_LIBCXXABI \
    -D_LIBCPP_BUILDING_LIBRARY \
    -D_LIBCPP_HAS_THREAD_API_WIN32 \
    -D_LIBCXXABI_BUILDING_LIBRARY \
    -D_LIBCPP_DISABLE_VISIBILITY_ANNOTATIONS \
    -D_LIBCXXABI_DISABLE_VISIBILITY_ANNOTATIONS \
    -UWIN32_LEAN_AND_MEAN \
    ")


add_library(c++ STATIC ${LIBCXX_SRC_FILES})

target_link_libraries(c++
    c++abi
    )
