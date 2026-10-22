#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT
repo="$test_dir/repo"
mkdir -p "$repo"
git -C "$repo" init -q
git -C "$repo" config user.name 'CI Test'
git -C "$repo" config user.email 'ci@example.invalid'
printf 'module example.com/test\n' > "$repo/go.mod"
printf 'package example\n' > "$repo/example.go"
git -C "$repo" add .
git -C "$repo" commit -qm baseline
base=$(git -C "$repo" rev-parse HEAD)

assert_change() {
	local path=$1 expected=$2
	mkdir -p "$repo/$(dirname "$path")"
	printf '// change\n' >> "$repo/$path"
	git -C "$repo" add .
	git -C "$repo" commit -qm "change $path"
	local head actual
	head=$(git -C "$repo" rev-parse HEAD)
	actual=$(cd "$repo" && bash "$repo_root/ci/scripts/package-changed.sh" "$base" "$head")
	if [[ "$actual" != "$expected" ]]; then
		printf 'expected %s for %s; got %s\n' "$expected" "$path" "$actual" >&2
		exit 1
	fi
}

assert_change README.md package_changed=false
assert_change .github/workflows/ci.yml package_changed=false
assert_change ci/go/main.go package_changed=false
assert_change subpackage/package.go package_changed=true

git -C "$repo" checkout -q "$base" --detach
printf 'dependency change\n' >> "$repo/go.sum"
git -C "$repo" add go.sum
git -C "$repo" commit -qm 'change go.sum'
head=$(git -C "$repo" rev-parse HEAD)
actual=$(cd "$repo" && bash "$repo_root/ci/scripts/package-changed.sh" "$base" "$head")
if [[ "$actual" != package_changed=true ]]; then
	printf 'expected true for go.sum; got %s\n' "$actual" >&2
	exit 1
fi

printf 'package change detection tests passed\n'
