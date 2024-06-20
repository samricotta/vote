#!/usr/bin/env bash

set -e

echo "Generating gogo proto code"
cd x
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep -q "option go_package" "$file" && grep -H -o -c 'option go_package.*github.com/samricotta/vote/x/api' "$file" | grep -q ':0$'; then
      buf generate --template x/commit-reveal-scheme/proto/buf.gen.gogo.yaml $file
    fi
  done
done

echo "Generating pulsar proto code"
 cd $home
 buf generate --template x/commit-reveal-scheme/proto/buf.gen.gogo.yaml $file

cd ..

cp -r github.com/samricotta/vote/x/* ./
rm -rf api && mkdir api
mv samricotta/* ./api
rm -rf github.com samricotta