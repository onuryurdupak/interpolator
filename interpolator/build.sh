#!/bin/bash
git diff --quiet HEAD

if [ "$?" != "0" ]; then
    echo "Warning: Can not build: repository is dirty."
    exit 1
fi

DATE=$(date +'%Y.%m.%d')
COMMIT_HASH=$(git rev-parse --short HEAD)

interpolator ./main.go ':=' 'stamp_build_date\s+=\s+"\${build_date}":=stamp_build_date = '\"$DATE\"
interpolator ./main.go ':=' 'stamp_commit_hash\s+=\s+"\${commit_hash}":=stamp_commit_hash = '\"$COMMIT_HASH\"

go build
git reset --hard
