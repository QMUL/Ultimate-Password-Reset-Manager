GO_GET(go-junit-report github.com/jstemmer/go-junit-report)
GO_GET(gocov github.com/axw/gocov/gocov)
GO_GET(gocov-xml github.com/AlekSi/gocov-xml)

FUNCTION(SET_TARGET_FOR_TESTS output_test_file output_coverage_file)

  ADD_CUSTOM_TARGET(tests

    DEPENDS go-junit-report gocov gocov-xml

    #Run the tests
    COMMAND env GOPATH=${GOPATH} go test -v pass.hpc.qmul.ac.uk/prm > tests.out 
    
    #Run the conversion to xUnit format
    COMMAND cat ${CMAKE_BINARY_DIR}/tests.out | ${GOPATH}/bin/go-junit-report > ${CMAKE_BINARY_DIR}/${output_test_file}

    #Run the coverage tests
    COMMAND env GOPATH=${GOPATH} ${GOPATH}/bin/gocov test pass.hpc.qmul.ac.uk/prm > gocov.txt

    #Run the conversion to xml
    COMMAND cat ${CMAKE_BINARY_DIR}/gocov.txt | ${GOPATH}/bin/gocov-xml > ${CMAKE_BINARY_DIR}/${output_coverage_file}
    
    WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
    COMMENT "Running tests."
  )

  # Show info where to find the report
  ADD_CUSTOM_COMMAND(TARGET tests POST_BUILD
    COMMAND ;
    COMMENT "Test report saved in ${output_test_file} and coverage report saved in ${output_coverage_file}"
  )
ENDFUNCTION() # SET_TARGET_FOR_TESTS
