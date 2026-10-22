#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 || ! "$1" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
	printf 'usage: %s <stable-version-tag>\n' "$0" >&2
	exit 2
fi

tag=$1
if gh release view "$tag" >/dev/null 2>&1; then
	printf 'GitHub Release %s already exists\n' "$tag"
else
	gh release create "$tag" --verify-tag --generate-notes
fi
