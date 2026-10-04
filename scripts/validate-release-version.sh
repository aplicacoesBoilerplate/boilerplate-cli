#!/usr/bin/env bash
set -euo pipefail

if (($# != 2)); then
  echo 'Uso: bash scripts/validate-release-version.sh <tag-aprovada> <SemVer-calculada>' >&2
  exit 2
fi

tag=$1
version=$2
if [[ ! "$tag" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]] || [[ "$version" != "${tag#v}" ]]; then
  echo "Versão calculada '$version' difere da tag estável aprovada '$tag'." >&2
  exit 1
fi

if git show-ref --verify --quiet "refs/tags/$tag"; then
  echo "A tag '$tag' já existe; nenhuma release deve ser sobrescrita." >&2
  exit 1
fi
