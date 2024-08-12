#!/bin/bash

GIT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

echo "\"$CURRENT_BRANCH\"" >"$GIT_ROOT/branch-name.nix"

echo "Branch name \"$CURRENT_BRANCH\" written to $GIT_ROOT/branch-name.nix"
