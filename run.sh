#!/bin/bash

if [ "$1" != "" ]; then
	go run -tags $1 .
else
	go run .
fi
