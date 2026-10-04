#!/usr/bin/env bash
set -euo pipefail

if (($# != 1)); then
  echo 'Uso: bash scripts/validate-commit-msg.sh <arquivo|->' >&2
  exit 2
fi

subject=''
if [[ "$1" == '-' ]]; then
  IFS= read -r subject || true
else
  IFS= read -r subject < "$1" || true
fi

if [[ ! "$subject" =~ ^(feat|fix|release|build|chore|ci|docs|perf|refactor|revert|style|test)(\([a-z0-9][a-z0-9._/-]*\))?(!)?:\ .+ ]]; then
  echo 'Commit inválido: use tipo(escopo opcional)!: descrição (ex.: feat(cli): adicionar versão).' >&2
  exit 1
fi
