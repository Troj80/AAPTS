cmake_minimum_required(VERSION 3.14.2)

project(aapt)

set(AAPT ${AOSP}/frameworks/base/tools/aapt)

set(AAPT_SRC_FILES
    ${AAPT}/Main.cpp
    ${AAPT}/AaptAssets.cpp
    ${AAPT}/AaptConfig.cpp
    ${AAPT}/AaptUtil.cpp
    ${AAPT}/AaptXml.cpp
    ${AAPT}/ApkBuilder.cpp
    ${AAPT}/Command.cpp
    ${AAPT}/CrunchCache.cpp
    ${AAPT}/FileFinder.cpp
    ${AAPT}/Images.cpp
    ${AAPT}/Package.cpp
    ${AAPT}/pseudolocalize.cpp
    ${AAPT}/Resource.cpp
    ${AAPT}/ResourceFilter.cpp
    ${AAPT}/ResourceIdCache.cpp
    ${AAPT}/ResourceTable.cpp
    ${AAPT}/SourcePos.cpp
    ${AAPT}/StringPool.cpp
    ${AAPT}/WorkQueue.cpp
    ${AAPT}/XMLNode.cpp
    ${AAPT}/ZipEntry.cpp
    ${AAPT}/ZipFile.cpp
    )

include_directories(
    ${aapt_headers}
    ${androidfw_headers}
    ${libpng_headers}
    ${libutils_headers}
    ${liblog_headers}
    ${libcutils_headers}
    ${expat_headers}
    ${libziparchive_headers}
    ${libbase_headers}
    ${zlib_headers}
    ${libsystem_headers}
    ${fmtlib_headers}
    ${native_headers}
    ${libbuildversion_headers}
    ${libincfs-utils_headers}
    )

set(EXECUTABLE_OUTPUT_PATH "../../bin")

set(CMAKE_CXX_FLAGS "${GLOBAL_CXX_FLAGS} \
    -Wall \
    -Werror \
    -Wno-format-y2k \
    -DSTATIC_ANDROIDFW_FOR_TOOLS \
    -Wno-unused-but-set-variable \
    
    ")

set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -pthread -s")


add_executable(aapt ${AAPT_SRC_FILES})

target_link_libraries(aapt
    png
    androidfw
    utils
    log
    cutils
    expat
    ziparchive
    base
    -Wl,--whole-archive
    buildversion
    -Wl,--no-whole-archive
    z
    c++
    m
    dl
    rt
    pthread
    )
