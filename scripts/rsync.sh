#!/bin/bash

cd "$(dirname "$0")"/../ || exit 1
dir="$(basename "$(pwd)")"

# define Host pi at ~/.ssh/config to specify copy target
fswatch -o ../ | xargs -I{} rsync -av --delete . "$1":work/"${dir}"
