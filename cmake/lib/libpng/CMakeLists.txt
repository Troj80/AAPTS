cmake_minimum_required(VERSION 3.14.2)

project(libpng)

set(LIBPNG ${AOSP}/external/libpng)

set(LIBPNG_SRC_FILES
    ${LIBPNG}/png.c
    ${LIBPNG}/pngerror.c
    ${LIBPNG}/pngget.c
    ${LIBPNG}/pngmem.c
    ${LIBPNG}/pngpread.c
    ${LIBPNG}/pngread.c
    ${LIBPNG}/pngrio.c
    ${LIBPNG}/pngrtran.c
    ${LIBPNG}/pngrutil.c
    ${LIBPNG}/pngset.c
    ${LIBPNG}/pngtrans.c
    ${LIBPNG}/pngwio.c
    ${LIBPNG}/pngwrite.c
    ${LIBPNG}/pngwtran.c
    ${LIBPNG}/pngwutil.c
    ${LIBPNG}/intel/filter_sse2_intrinsics.c
    ${LIBPNG}/intel/intel_init.c
    )

include_directories(
    ${libpng_headers}
    ${zlib_headers}
    )


set(CMAKE_C_FLAGS "${GLOBAL_C_FLAGS} \
    -std=gnu89 \
    -Wall \
    -Werror \
    -Wno-unused-parameter \
    -DPNG_INTEL_SSE_OPT=1 \
    -Wno-unused-but-set-variable \
    ")

set(CMAKE_CXX_FLAGS "${GLOBAL_CXX_FLAGS}")


add_library(png STATIC ${LIBPNG_SRC_FILES})

target_link_libraries(png)
