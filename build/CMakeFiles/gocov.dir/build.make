# CMAKE generated file: DO NOT EDIT!
# Generated by "Unix Makefiles" Generator, CMake Version 3.5

# Delete rule output on recipe failure.
.DELETE_ON_ERROR:


#=============================================================================
# Special targets provided by cmake.

# Disable implicit rules so canonical targets will work.
.SUFFIXES:


# Remove some rules from gmake that .SUFFIXES does not remove.
SUFFIXES =

.SUFFIXES: .hpux_make_needs_suffix_list


# Suppress display of executed commands.
$(VERBOSE).SILENT:


# A target that is always out of date.
cmake_force:

.PHONY : cmake_force

#=============================================================================
# Set environment variables for the build.

# The shell in which to execute make rules.
SHELL = /bin/sh

# The CMake executable.
CMAKE_COMMAND = /usr/bin/cmake

# The command to remove a file.
RM = /usr/bin/cmake -E remove -f

# Escaping for special characters.
EQUALS = =

# The top-level source directory on which CMake was run.
CMAKE_SOURCE_DIR = /home/oni/Projects/gocode/src/pass.hpc.qmul.ac.uk

# The top-level build directory on which CMake was run.
CMAKE_BINARY_DIR = /home/oni/Projects/gocode/src/pass.hpc.qmul.ac.uk/build

# Utility rule file for gocov.

# Include the progress variables for this target.
include CMakeFiles/gocov.dir/progress.make

CMakeFiles/gocov:
	env GOPATH=/home/oni/Projects/gocode/src/pass.hpc.qmul.ac.uk/build/go go get github.com/axw/gocov/gocov

gocov: CMakeFiles/gocov
gocov: CMakeFiles/gocov.dir/build.make

.PHONY : gocov

# Rule to build all files generated by this target.
CMakeFiles/gocov.dir/build: gocov

.PHONY : CMakeFiles/gocov.dir/build

CMakeFiles/gocov.dir/clean:
	$(CMAKE_COMMAND) -P CMakeFiles/gocov.dir/cmake_clean.cmake
.PHONY : CMakeFiles/gocov.dir/clean

CMakeFiles/gocov.dir/depend:
	cd /home/oni/Projects/gocode/src/pass.hpc.qmul.ac.uk/build && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /home/oni/Projects/gocode/src/pass.hpc.qmul.ac.uk /home/oni/Projects/gocode/src/pass.hpc.qmul.ac.uk /home/oni/Projects/gocode/src/pass.hpc.qmul.ac.uk/build /home/oni/Projects/gocode/src/pass.hpc.qmul.ac.uk/build /home/oni/Projects/gocode/src/pass.hpc.qmul.ac.uk/build/CMakeFiles/gocov.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : CMakeFiles/gocov.dir/depend

