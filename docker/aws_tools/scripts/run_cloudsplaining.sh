#!/usr/bin/env bash

OUT_DIR="/work/cloudsplaining"
PWD="$(pwd)"

mkdir -p "${OUT_DIR}"

cd "${OUT_DIR}"

if [ -r exclusions.yml ]
then
    EXCLUSIONS="-e exclusions.yml"
fi

cloudsplaining download
cloudsplaining scan -i "default.json" ${EXCLUSIONS} -s

cd "${PWD}"
