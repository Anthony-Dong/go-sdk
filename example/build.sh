#!/usr/bin/env bash
mod_name="github.com/anthony-dong/go-tool"
pwd="internal/example"

function build() {
    if [ -z "$1" ]; then echo "usage: build.sh [path]
help: 帮助构建子mod"; exit 1; fi

if [ ! -d "$1" ]; then mkdir -p "$1" || exit 1;  fi
sub_mod_name="${mod_name}/${pwd}/$1"
echo "sub_mod_name: ${sub_mod_name}"

# build
if [ ! -e "$1/go.mod" ]; then cd $1 || exit 1; echo "init sub mod ${sub_mod_name}"; go mod init ${sub_mod_name} || exit 1; cd - || exit 1; fi

# go mod tidy
if [ `find "$1" -name "*.go" | wc -l` -gt 0 ]; then echo "已经存在项目了"; cd $1 || exit 1; go mod tidy || exit 1; cd - || exit 1; fi

if [ `go list -m all | grep -E 'github.com/anthony-dong/go-tool\s'` ]; then echo "replace github.com/anthony-dong/go-too"; \
  cd $1 || exit 1; go mod edit -replace github.com/anthony-dong/go-tool=../../../ || exit 1; go mod tidy || exit 1; cd -; \
fi

}

for elem in `ls | grep -v -E '.\.(sh|md)'`; do build $elem; done
