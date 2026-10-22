#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 2 ]]; then
	printf 'usage: %s <base-sha> <head-sha>\n' "$0" >&2
	exit 2
fi

base=$1
head=$2
if [[ "$base" =~ ^0+$ ]]; then
	base=$(git hash-object -t tree /dev/null)
fi

changed=false
while IFS= read -r -d '' path; do
	if [[ "$path" == 'go.mod' || "$path" == 'go.sum' || ( "$path" == *.go && "$path" != .github/go/* ) ]]; then
		changed=true
		break
	fi
done < <(git diff --no-renames --name-only -z "$base" "$head")

printf 'package_changed=%s\n' "$changed"
