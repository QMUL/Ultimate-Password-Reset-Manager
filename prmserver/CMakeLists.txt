project(prm_server)

GO_GET(go-ldap github.com/go-ldap/ldap)
GO_GET(yaml.v2 gopkg.in/yaml.v2)
GO_GET(ldap.v2 gopkg.in/ldap.v2)
GO_COPY(prm pass.hpc.qmul.ac.uk/prm)

ADD_GO_INSTALLABLE_PROGRAM(prm_server # executable name
  prm_server.go # `package main` source file
  prm
  go-ldap
  yaml.v2
  ldap.v2)

install(PROGRAMS ${CMAKE_CURRENT_BINARY_DIR}/prm_server DESTINATION passwordmanager RENAME prm_server.fcgi)
install(DIRECTORY ${CMAKE_SOURCE_DIR}/static/ DESTINATION static)
install(DIRECTORY ${CMAKE_SOURCE_DIR}/templates/ DESTINATION templates)

