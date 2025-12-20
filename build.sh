#!/bin/sh

set -xe

CFLAGS="-Wall -Wextra -Wshadow -Werror --pedantic -std=c23"
LIBS=""

gcc $CFLAGS -o mnc mnc.c Tokenizer/*.c Parser/*.c $LIBS
./mnc dtpk/src -o dtpk/data
