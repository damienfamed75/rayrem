#!/bin/bash

set -e

mkdir bin || true
mkdir bin/config || true

pkger

go build .

mv ./rayrem.exe ./bin/rayrem.exe

go run ./cmd/encrypter/

mv ./config/game.config ./bin/config/game.config
mv ./config/settings.config ./bin/config/settings.config