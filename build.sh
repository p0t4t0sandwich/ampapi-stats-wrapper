#!/bin/bash

mkdir -p ./build
cd ./build

# Build
CGO_ENABLED=0 GOOS=linux go build -o ./build/ampapi-stats-wrapper
CGO_ENABLED=0 GOOS=windows go build -o ./build/ampapi-stats-wrapper.exe

# Copy files
cp ../index.html ./index.html
cp ../openapi.json ./openapi.json
cp ../settings.json ./settings.json

# Zip
zip -r ./ampapi-stats-wrapper.zip ./*
