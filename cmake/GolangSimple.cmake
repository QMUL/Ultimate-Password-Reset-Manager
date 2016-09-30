set(GOPATH "${CMAKE_CURRENT_BINARY_DIR}/go")
file(MAKE_DIRECTORY ${GOPATH})

function(GO_GET TARG)
  add_custom_target(${TARG} env GOPATH=${GOPATH} go get ${ARGN})
endfunction(GO_GET)


function(GO_COPY MODULE_NAME MODULE_SRC)
  get_filename_component(MODULE_SRC_ABS "../../${MODULE_SRC}" ABSOLUTE)

  message(STATUS "Copying Local Module ${MODULE_SRC_ABS}")
  add_custom_target(${MODULE_NAME}
                    COMMAND ${CMAKE_COMMAND} -E copy_directory
                    ${MODULE_SRC_ABS} ${GOPATH}/src/${MODULE_SRC})

endfunction(GO_COPY)


function(ADD_GO_INSTALLABLE_PROGRAM NAME MAIN_SRC)
  get_filename_component(MAIN_SRC_ABS ${MAIN_SRC} ABSOLUTE)
  add_custom_target(${NAME} ALL DEPENDS ${NAME})
  add_custom_command(TARGET ${NAME}
                    COMMAND env GOPATH=${GOPATH} go build 
                    -o "${CMAKE_CURRENT_BINARY_DIR}/${NAME}"
                    ${CMAKE_GO_FLAGS} ${MAIN_SRC}
                    WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR}
                    DEPENDS ${MAIN_SRC_ABS})
  foreach(DEP ${ARGN})
    add_dependencies(${NAME} ${DEP})
  endforeach()
  
endfunction(ADD_GO_INSTALLABLE_PROGRAM)

