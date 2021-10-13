#!/bin/sh

set -e

for ex in examples/*; do
    go build -v -o "bin/${ex##*/}" "$ex/main.go"
done