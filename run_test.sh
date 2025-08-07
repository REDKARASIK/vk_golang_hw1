#!/bin/sh

modified_directories=$(git diff origin/master --name-only | grep 99_hw | awk -F/ '{print $1}' | uniq | sort -u)
if [ -z "$modified_directories" ]; then
    echo "No directories were modified"
    exit 0
fi

echo "Modified directories: $modified_directories"

is_tests_failed=false

for dir in $modified_directories; do
    test_path="./${dir}/99_hw/..."

    if [ ! -d "$test_path" ]; then
      echo "No 99_hw dir in $dir"
      continue
    else
      echo "Running tests in: $test_path"
    fi

    if [ "$dir" = "./04_net2/99_hw/..." ]; then
      go mod vendor

      go test --mod=vendor --race "$test_path"

      rm -rf vendor
    else
        go test --race "$test_path"
    fi

    if [ $? -ne 0 ]; then
        is_tests_failed=true
    fi
done

if $is_tests_failed; then
  exit 1
fi
