#!/usr/bin/env bash

platforms=$(go tool dist list)

echo -ne "Please enter a version: "
read version

rm build/keygen_*

for platform in $platforms
do
  read -r os arch <<<$(echo "${platform}" | tr '/' ' ')
  tag=$(echo "${version}" | sed 's/[-.+]/_/g')
  out="build/keygen_${os}_${arch}_${tag}"

  if [ "${os}" = "windows" ]; then
    out="${out}.exe"
  fi

  env GOOS="${os}" GOARCH="${arch}" \
    go build -o "${out}"
done
