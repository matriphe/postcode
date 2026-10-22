#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
version_script="$repo_root/.github/ci/scripts/next-version.sh"
module_path='github.com/matriphe/postcode'

assert_version() {
	local branch=$1
	local latest=$2
	local module=$3
	local expected=$4
	local actual
	actual=$(bash "$version_script" "$branch" "$latest" "$module")
	if [[ "$actual" != "$expected" ]]; then
		printf 'expected %s for branch=%s latest=%s; got %s\n' "$expected" "$branch" "$latest" "$actual" >&2
		exit 1
	fi
}

assert_failure() {
	local branch=$1
	local latest=$2
	local module=$3
	if bash "$version_script" "$branch" "$latest" "$module" >/dev/null 2>&1; then
		printf 'expected failure for branch=%s latest=%s module=%s\n' "$branch" "$latest" "$module" >&2
		exit 1
	fi
}

assert_version 'feature/major-bootstrap' '' "$module_path" 'v0.1.0'
assert_version 'main' 'v0.1.0' "$module_path" 'v0.1.1'
assert_version 'feature/minor-rules' 'v0.1.0' "$module_path" 'v0.2.0'
assert_version 'feature/MAJOR-rules' 'v0.9.8' "$module_path" 'v1.0.0'
assert_version 'feature/major-minor-rules' 'v1.2.3' "$module_path/v2" 'v2.0.0'
assert_version 'maintenance' 'v2.0.0' "$module_path/v2" 'v2.0.1'
assert_failure 'feature/major-release' 'v1.2.3' "$module_path"
