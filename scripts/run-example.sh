#!/bin/bash

EXAMPLE=$1

if [ -z $1 ]; then
    echo "Please select one of the following examples to run:"
    select e in $(ls -1 examples/); do
        EXAMPLE=$e
        break
    done
fi

if [ -z $TOKEN ]; then
    echo "The TOKEN environment variable must be set to the Discord bots token."
    exit 1
fi

echo "Running example $EXAMPLE ..."
go run -v examples/$EXAMPLE/main.go