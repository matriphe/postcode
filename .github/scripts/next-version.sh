#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 3 ]]; then
	printf 'usage: %s <branch-name> <latest-version-or-empty> <module-path>\n' "$0" >&2
	exit 2
fi

branch=${1,,}
latest_version=$2
module_path=$3

if [[ -z "$latest_version" ]]; then
	major=0
	minor=1
	patch=0
else
	if [[ ! "$latest_version" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
		printf 'latest version is not a stable semantic version: %s\n' "$latest_version" >&2
		exit 2
	fi
	major=${BASH_REMATCH[1]}
	minor=${BASH_REMATCH[2]}
	patch=${BASH_REMATCH[3]}

	if [[ "$branch" == *major* ]]; then
		((major += 1))
		minor=0
		patch=0
	elif [[ "$branch" == *minor* ]]; then
		((minor += 1))
		patch=0
	else
		((patch += 1))
	fi
fi

if [[ "$major" -ge 2 ]]; then
	if [[ "$module_path" =~ ^(.+)/v[0-9]+$ ]]; then
		module_base=${BASH_REMATCH[1]}
	else
		module_base=$module_path
	fi
	expected_module_path="$module_base/v$major"
	if [[ "$module_path" != "$expected_module_path" ]]; then
		printf 'version v%s requires module path %s; go.mod has %s\n' "$major" "$expected_module_path" "$module_path" >&2
		exit 1
	fi
elif [[ "$module_path" =~ /v[0-9]+$ ]]; then
	printf 'module path %s has a major-version suffix but version v%s does not use one\n' "$module_path" "$major" >&2
	exit 1
fi

printf 'v%s.%s.%s\n' "$major" "$minor" "$patch"
