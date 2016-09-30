#utility function to generate a shell script from the input list containing commands
function(generate_shell_script SCRIPTFILENAME INPUTLIST)
    file(WRITE ${SCRIPTFILENAME} "#!/bin/sh\n")
    FOREACH(command ${INPUTLIST})
        file(APPEND ${SCRIPTFILENAME} "${command}\n")
    ENDFOREACH(command)
endfunction(generate_shell_script)

#define the generator and other required settings
set(CPACK_GENERATOR "RPM")
set(CPACK_PACKAGE_CONTACT "ITS Research <its-research-support@qmul.ac.uk>")
set(CPACK_PACKAGE_DESCRIPTION_SUMMARY "ITSR password reset manager application")
set(CPACK_RPM_PACKAGE_DESCRIPTION "ITSR password reset manager application")
set(CPACK_PACKAGE_VENDOR "ITS Research")
SET(CPACK_PACKAGE_NAME "prm")
set(CPACK_PACKAGE_VERSION_MAJOR "${PRM_PACKAGE_VERSION_MAJOR}")
set(CPACK_PACKAGE_VERSION_MINOR "${PRM_PACKAGE_VERSION_MINOR}")
set(CPACK_PACKAGE_VERSION_PATCH "${PRM_PACKAGE_VERSION_PATCH}")
set(CPACK_PACKAGING_INSTALL_PREFIX "/var/www/html")

#RPM specific options
set(CPACK_RPM_PACKAGE_RELEASE "${PRM_RPM_PACKAGE_RELEASE}")
set(CPACK_RPM_PACKAGE_ARCHITECTURE "x86_64")
set(CPACK_RPM_PACKAGE_REQUIRES "policycoreutils >= 2.2.5")
set(CPACK_RPM_EXCLUDE_FROM_AUTO_FILELIST_ADDITION "/var/www/html")

#set the os release
string (REGEX MATCH "\\el[1-9]" OS_VERSION ${CMAKE_SYSTEM})

#set the output file name
set(CPACK_PACKAGE_FILE_NAME "${CPACK_PACKAGE_NAME}-${PRM_PACKAGE_VERSION}.${OS_VERSION}.${CPACK_RPM_PACKAGE_ARCHITECTURE}")

set(PRM_POST_INSTALL_CMDS "restorecon -R ${CPACK_PACKAGING_INSTALL_PREFIX}/{passwordmanager,static,templates}")
set(PRM_POST_INSTALL ${CMAKE_CURRENT_BINARY_DIR}/rpm/postinstall)
generate_shell_script(${PRM_POST_INSTALL} "${PRM_POST_INSTALL_CMDS}")
set(CPACK_RPM_POST_INSTALL_SCRIPT_FILE "${PRM_POST_INSTALL}")
#set(CPACK_RPM_POST_INSTALL_SCRIPT_FILE "")

include(CPack)


