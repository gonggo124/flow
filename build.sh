#!/bin/sh

set -xe

CFLAGS="-Wall -Wextra -Wshadow -Werror --pedantic"
LIBS=""

gcc $CFLAGS -o mnc mnc.c $LIBS
./mnc dtpk/src -o dtpk/data
