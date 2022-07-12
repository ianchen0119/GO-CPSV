#!/usr/bin/env bash                                                                       

echo "(In $(pwd))"
echo docker build -t local/go-cpsv .
docker build -t local/go-cpsv .