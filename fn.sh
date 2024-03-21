#!/bin/sh

# read from stdin into a variable
INPUT=$(cat)
chmod 755 umbilical-choir-proxy

# echo the result to stdout
echo -n $(./umbilical-choir-proxy "$INPUT")
