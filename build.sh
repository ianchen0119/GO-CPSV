#!/usr/bin/env bash                                                                       

echo "(In $(pwd))"
echo docker build -t local/go-cpsv . --no-cache
docker build -t local/go-cpsv . # --no-cache
