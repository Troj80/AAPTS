cmake_minimum_required(VERSION 3.14.2)

project(aapt2)

set(AAPT2 ${AOSP}/frameworks/base/tools/aapt2)


set(AAPT2_PROTO_SRC)  # proto source files
set(AAPT2_PROTO_HDRS) # proto head files

file(GLOB_RECURSE PROTO_FILES ${AAPT2}/*.proto)
foreach(proto ${PROTO_FILES})
    get_filename_component(FIL_WE ${proto} NAME_WE)
    set(TARGET_CPP_FILE "${AAPT2}/${FIL_WE}.pb.cc")
    set(TARGET_HEAD_FILE "${AAPT2}/${FIL_WE}.pb.h")
    if(EXISTS ${TARGET_CPP_FILE} AND EXISTS ${TARGET_HEAD_FILE})
        list(APPEND AAPT2_PROTO_SRC ${TARGET_CPP_FILE})
        list(APPEND AAPT2_PROTO_HDRS ${TARGET_HEAD_FILE})
    else()
        execute_process(
            COMMAND protoc ${proto}
            --proto_path=${AAPT2}
            --cpp_out=${AAPT2}
            WORKING_DIRECTORY ${AAPT2}
        )
        message(STATUS "generate cpp file ${TARGET_CPP_FILE}")
        message(STATUS "generate head file ${TARGET_HEAD_FILE}")
    endif()
endforeach()

set_source_files_properties(${AAPT2_PROTO_SRC} ${AAPT2_PROTO_HDRS} PROPERTIES GENERATED TRUE)

set(AAPT2_SRC_FILES
    ${AAPT2}/Main.cpp

    ${AAPT2}/cmd/Command.cpp
    ${AAPT2}/cmd/Compile.cpp
    ${AAPT2}/cmd/Convert.cpp
    ${AAPT2}/cmd/Diff.cpp
    ${AAPT2}/cmd/Dump.cpp
    ${AAPT2}/cmd/Link.cpp
    ${AAPT2}/cmd/Optimize.cpp
    ${AAPT2}/cmd/Util.cpp

    ${AAPT2}/compile/IdAssigner.cpp
    ${AAPT2}/compile/InlineXmlFormatParser.cpp
    ${AAPT2}/compile/NinePatch.cpp
    ${AAPT2}/compile/Png.cpp
    ${AAPT2}/compile/PngChunkFilter.cpp
    ${AAPT2}/compile/PngCrunch.cpp
    ${AAPT2}/compile/PseudolocaleGenerator.cpp
    ${AAPT2}/compile/Pseudolocalizer.cpp
    ${AAPT2}/compile/XmlIdCollector.cpp
    ${AAPT2}/configuration/ConfigurationParser.cpp
    ${AAPT2}/dump/DumpManifest.cpp
    ${AAPT2}/filter/AbiFilter.cpp
    ${AAPT2}/filter/ConfigFilter.cpp
    ${AAPT2}/format/Archive.cpp
    ${AAPT2}/format/Container.cpp
    ${AAPT2}/format/binary/BinaryResourceParser.cpp
    ${AAPT2}/format/binary/ResChunkPullParser.cpp
    ${AAPT2}/format/binary/TableFlattener.cpp
    ${AAPT2}/format/binary/XmlFlattener.cpp
    ${AAPT2}/format/proto/ProtoDeserialize.cpp
    ${AAPT2}/format/proto/ProtoSerialize.cpp
    ${AAPT2}/io/BigBufferStream.cpp
    ${AAPT2}/io/File.cpp
    ${AAPT2}/io/FileStream.cpp
    ${AAPT2}/io/FileSystem.cpp
    ${AAPT2}/io/StringStream.cpp
    ${AAPT2}/io/Util.cpp
    ${AAPT2}/io/ZipArchive.cpp
    ${AAPT2}/link/AutoVersioner.cpp
    ${AAPT2}/link/ManifestFixer.cpp
    ${AAPT2}/link/NoDefaultResourceRemover.cpp
    ${AAPT2}/link/ProductFilter.cpp
    ${AAPT2}/link/PrivateAttributeMover.cpp
    ${AAPT2}/link/ReferenceLinker.cpp
    ${AAPT2}/link/ResourceExcluder.cpp
    ${AAPT2}/link/TableMerger.cpp
    ${AAPT2}/link/XmlCompatVersioner.cpp
    ${AAPT2}/link/XmlNamespaceRemover.cpp
    ${AAPT2}/link/XmlReferenceLinker.cpp
    ${AAPT2}/optimize/MultiApkGenerator.cpp
    ${AAPT2}/optimize/ResourceDeduper.cpp
    ${AAPT2}/optimize/ResourceFilter.cpp
    ${AAPT2}/optimize/ResourcePathShortener.cpp
    ${AAPT2}/optimize/VersionCollapser.cpp
    ${AAPT2}/process/SymbolTable.cpp
    ${AAPT2}/split/TableSplitter.cpp
    ${AAPT2}/text/Printer.cpp
    ${AAPT2}/text/Unicode.cpp
    ${AAPT2}/text/Utf8Iterator.cpp
    ${AAPT2}/util/BigBuffer.cpp
    ${AAPT2}/util/Files.cpp
    ${AAPT2}/util/Util.cpp
    ${AAPT2}/Debug.cpp
    ${AAPT2}/DominatorTree.cpp
    ${AAPT2}/java/AnnotationProcessor.cpp
    ${AAPT2}/java/ClassDefinition.cpp
    ${AAPT2}/java/JavaClassGenerator.cpp
    ${AAPT2}/java/ManifestClassGenerator.cpp
    ${AAPT2}/java/ProguardRules.cpp
    ${AAPT2}/LoadedApk.cpp
    ${AAPT2}/Resource.cpp
    ${AAPT2}/ResourceParser.cpp
    ${AAPT2}/ResourceTable.cpp
    ${AAPT2}/ResourceUtils.cpp
    ${AAPT2}/ResourceValues.cpp
    ${AAPT2}/SdkConstants.cpp
    ${AAPT2}/StringPool.cpp
    ${AAPT2}/trace/TraceBuffer.cpp
    ${AAPT2}/xml/XmlActionExecutor.cpp
    ${AAPT2}/xml/XmlDom.cpp
    ${AAPT2}/xml/XmlPullParser.cpp
    ${AAPT2}/xml/XmlUtil.cpp
    ${AAPT2}/Configuration.pb.cc
    ${AAPT2}/Resources.pb.cc
    ${AAPT2}/ResourcesInternal.pb.cc
    )

include_directories(
    ${aapt2_headers}
    ${protobuf_headers}
    ${androidfw_headers}
    ${libbase_headers}
    ${libziparchive_headers}
    ${liblog_headers}
    ${libcutils_headers}
    ${libutils_headers}
    ${libsystem_headers}
    ${zlib_headers}
    ${libpng_headers}
    ${expat_headers}
    ${libbuildversion_headers}
    ${libidmap2_policies_headers}
    ${fmtlib_headers}
    ${native_headers}
    ${libincfs-utils_headers}
    ${gtest_headers}
    )

set(EXECUTABLE_OUTPUT_PATH "../../bin")

message(STATUS "CMAKE_BINARY_DIR: ${CMAKE_BINARY_DIR}")

set(CMAKE_CXX_FLAGS "${GLOBAL_CXX_FLAGS} \
    -Wall \
    -Werror \
    -Wno-unused-parameter \
    -Wno-missing-field-initializers \
    -Wno-defaulted-function-deleted \
    -fno-exceptions \
    -fno-rtti \
    ")

set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -pthread -s")


add_executable(aapt2 ${AAPT2_SRC_FILES})

target_link_libraries(aapt2
    androidfw
    utils
    cutils
    expat
    ziparchive
    png
    base
    protobuf-cpp-full
    -Wl,--whole-archive
    buildversion
    -Wl,--no-whole-archive
    log
    z
    c++
    m
    dl
    rt
    pthread
    )
