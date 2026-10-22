#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT
mkdir "$test_dir/bin"

cat > "$test_dir/bin/gh" <<'STUB'
#!/usr/bin/env bash
set -euo pipefail

if [[ "$1 $2" == 'release view' ]]; then
	[[ "${GH_STUB_RELEASE_EXISTS:-0}" == '1' ]]
elif [[ "$1 $2" == 'release create' ]]; then
	printf '%s\n' "$@" > "$GH_STUB_LOG"
else
	printf 'unexpected gh invocation: %s\n' "$*" >&2
	exit 2
fi
STUB
chmod +x "$test_dir/bin/gh"

export PATH="$test_dir/bin:$PATH"
export GH_STUB_LOG="$test_dir/gh.log"
: > "$GH_STUB_LOG"

GH_STUB_RELEASE_EXISTS=1 bash "$repo_root/.github/ci/scripts/ensure-github-release.sh" v0.1.0
if [[ -s "$GH_STUB_LOG" ]]; then
	printf 'existing release should not be recreated\n' >&2
	exit 1
fi

GH_STUB_RELEASE_EXISTS=0 bash "$repo_root/.github/ci/scripts/ensure-github-release.sh" v0.1.1
expected=$(printf 'release\ncreate\nv0.1.1\n--verify-tag\n--generate-notes')
actual=$(cat "$GH_STUB_LOG")
if [[ "$actual" != "$expected" ]]; then
	printf 'unexpected release create arguments:\n%s\n' "$actual" >&2
	exit 1
fi
