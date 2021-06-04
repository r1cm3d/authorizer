#!/usr/bin/sh

#                                           DISCLAIMER IMPORTANT
#
# All commands in this script are POSIX standard compliant. I do not know if Darwin (MacOS) is compliance with all of them, so,
# some of them may fail in this OS. Unfortunately, I do not have any Darwin computer around to test it (or fortunately ;] ).
# I googled about some MacOS that runs on Docker container, but I just found this one (https://github.com/sickcodes/Docker-OSX)
# which I've thought "too much" for this toy project.
#
# So if this script fails and you have a MacOS, please install Make, Docker and run 'make && make install'"

red='\e[1;31m'
grn='\e[1;32m'
cyn='\e[1;36m'
yel='\e[1;33m'
end='\e[0m'

check() {
  (printf "\nChecking if %s is installed\n" "$1" && \
                  command -v "$1" 1>/dev/null && \
                  printf "\n%b%s is installed%b\n" "$grn" "$1" "$end" && \
                  return 0)                               || \
   (printf "\n%b%s is not installed%b\n" "$red" "$1" "$end" && return 1)
}

check "make" || printf "\nPlease install %bMake%b without %bMake%b is impossible to build this application\n" "$cyn" "$end" "$cyn" "$end"
check "docker" || printf "\nPlease install %bDocker%b to containerize this application\n" "$cyn" "$end"
check "go" || printf "\nPlease install %bGo%b to run tests and build locally\n" "$cyn" "$end"
printf "\nAll setup. Run %bmake%b to assemble %bDocker%b container\n" "$yel" "$end" "$cyn" "$end"