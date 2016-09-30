#based on http://stackoverflow.com/questions/3780667/use-cmake-to-get-build-time-svn-revision
find_package(Git REQUIRED)
exec_program(${GIT_EXECUTABLE} 
    ${SOURCE_DIR}
    ARGS "describe --dirty"
    OUTPUT_VARIABLE GIT_REVISION)
set(PRM_REVISION "${GIT_REVISION}")

#get the build host
site_name(PRM_BUILD_HOST)

#get the build timestamp in the default format
string(TIMESTAMP BUILD_TIMESTAMP)

# write relevant information into template header 
file(WRITE ${CMAKE_CURRENT_BINARY_DIR}/buildinfo.go.txt
"package prm

func GetVersionString() string {
        return \"${PRM_REVISION} ${PRM_BUILD_HOST} ${BUILD_TIMESTAMP}\"
}
")

# copy the file to the final header only if the data changes (currently it will be every time due to the timestamp data) 
execute_process(COMMAND ${CMAKE_COMMAND} -E copy_if_different ${CMAKE_CURRENT_BINARY_DIR}/buildinfo.go.txt ${DEST_DIR}/version.go)
