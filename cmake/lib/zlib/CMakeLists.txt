cmake_minimum_required(VERSION 3.14.2)

project(zlib)

set(ZLIB ${AOSP}/external/zlib)

set(ZLIB_SRC_FILES
    ${ZLIB}/adler32_simd.c
    ${ZLIB}/crc32_simd.c

    ${ZLIB}/crc_folding.c
    ${ZLIB}/fill_window_sse.c

    ${ZLIB}/adler32.c
    ${ZLIB}/compress.c
    ${ZLIB}/cpu_features.c
    ${ZLIB}/crc32.c
    ${ZLIB}/deflate.c
    ${ZLIB}/gzclose.c
    ${ZLIB}/gzlib.c
    ${ZLIB}/gzread.c
    ${ZLIB}/gzwrite.c
    ${ZLIB}/infback.c
    ${ZLIB}/inffast.c
    ${ZLIB}/inflate.c
    ${ZLIB}/inftrees.c
    ${ZLIB}/trees.c
    ${ZLIB}/uncompr.c
    ${ZLIB}/zutil.c
    )

include_directories(
    ${zlib_headers}
    ${cpu_features_headers}
    )


set(CMAKE_C_FLAGS "${GLOBAL_C_FLAGS} \
    -DX86_NOT_WINDOWS \
    -DCPU_NO_SIMD \
    -DINFLATE_CHUNK_READ_64LE \
    -DHAVE_HIDDEN \
    -DZLIB_CONST \
    -Wall \
    -Werror \
    -Wno-unused \
    -Wno-unused-parameter \
    ")

set(CMAKE_CXX_FLAGS "${GLOBAL_CXX_FLAGS}")


add_library(z STATIC ${ZLIB_SRC_FILES})

target_link_libraries(z)
