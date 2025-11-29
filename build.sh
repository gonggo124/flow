#!/bin/sh

set -xe

CFLAGS="-Wall -Wextra -Wshadow -Werror --pedantic"
LIBS=""

gcc $CFLAGS -o mnc mnc.c Tokenizer/tokenizer.c $LIBS
./mnc dtpk/src -o dtpk/data
