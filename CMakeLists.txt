# CMakeList for Go borrowed from here and modified
# https://github.com/cpconduce/go_cmake

cmake_minimum_required(VERSION 3.0)
project(prm)

include(cmake/PrmVersion.cmake)
set_version(MAJOR ${VERSION_MAJOR} MINOR ${VERSION_MINOR} PATCH ${VERSION_PATCH} RELEASE ${VERSION_RELEASE})

include(cmake/GolangSimple.cmake)
add_subdirectory(buildinfo)
add_subdirectory(prmserver)

#Tests
include(cmake/GoTests.cmake)
set_target_for_tests("tests.xml" "coverage.xml")

#Packaging
include(cmake/PrmPackage.cmake)
