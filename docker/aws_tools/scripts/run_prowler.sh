#!/usr/bin/env bash

OUT_DIR="/work/prowler"
PWD="$(pwd)"

mkdir -p "${OUT_DIR}"

prowler -M text,html,json,csv -o "${OUT_DIR}" $*
