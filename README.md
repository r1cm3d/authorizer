# authorizer
This CLI application authorizes a transaction according predefined rules

## Table of Contents
* [TL;DR](#tldr)
* [Prerequisites](#prerequisites)
    * [To run](#to-run-in-docker-container)
    * [To develop](#to-run-locally)
* [About](#about)
* [Test](#test)
  * [Unit Test](#unit-test)
  * [Integration Test](#integration-test)
  * [Both of them](#both-of-them)

### TLDR
``` shell
./configure
```
To copy **INPUT_JSON_PATH** into test/ directory to be build inside Docker container.
``` shell
make && make install
```
To assemble Docker container and create executable file (just an alias for `docker run` command).
``` shell
./authorizer < operations
```

### Prerequisites
TODO: Add badges here :smile: 

#### To run in Docker container
You only need [Make](https://www.gnu.org/software/make/) (tested with GNU Make 4.2.1) and [Docker](https://www.docker.com/) (tested with Docker version 20.10.2, build 2291f61)  installed. It is important to give execution permission for `scripts/`:
``` shell
chmod +x scripts/
```
They are safe, but you could check it out if you don't believe me. :smile:
#### To run locally
[Go 1.16](https://golang.org/dl/) ecosystem installed. I'm pretty sure it might work in the recent old versions `1.1[0-6]`, but I actually did not test.

### About
[WIP]


### Test
#### Unit test
``` shell
make unit-test
```
#### Integration test
``` shell
make integration-test
```
#### Both of them
``` shell
make test
```
It is important to make sure [Go 1.16](https://golang.org/dl/) and [Make](https://www.gnu.org/software/make/) is installed.