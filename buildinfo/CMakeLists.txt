cmake_minimum_required(VERSION 3.0)
project(buildinfo)

#reqire either subversion or git here so that we can find out
#before the custom command is executed
find_package(Git REQUIRED)

# a custom target that is always built
add_custom_target(genversion ALL DEPENDS buildinfo)

#run autogeneration of the header
add_custom_command(OUTPUT buildinfo
                   COMMAND ${CMAKE_COMMAND}
                           -DSOURCE_DIR=${CMAKE_SOURCE_DIR}
                           -DDEST_DIR=${GOPATH}/src/pass.hpc.qmul.ac.uk/prm
                           -P ${CMAKE_SOURCE_DIR}/cmake/PrmBuildInfo.cmake)
