#!/bin/bash

pr_id=${PR_ID}

if [ -f ".changelog/${pr_id}.txt" ]; then
    make changelogtest

    if [ $? -ne 0 ]; then
        printf "COMMIT FAILED\n"s
        exit 1
    fi
    exit 0
fi

make deltatest
if [ $? -ne 0 ]; then
    printf "COMMIT FAILED\n"
    exit 1
fi