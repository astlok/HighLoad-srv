#!/bin/bash
for (( number = 1; number <= $1 ; number++ ))
do
    curl -X GET 51.250.20.83:8080/api
    curl -X GET 51.250.27.26:8080/api
    curl -X GET 51.250.18.255:8080/api
done
