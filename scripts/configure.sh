#!/usr/bin/sh

. ./colors.sh

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