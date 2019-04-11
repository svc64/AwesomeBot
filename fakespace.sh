#!/bin/sh
# Stolen from https://github.com/golang/dep/issues/1441

set -e

if [ ! -f "fakespace.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

project="AwesomeBot"
repositoryRoot="git.asafniv.me/blzit420"

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
dir="$workspace/src/$repositoryRoot"
if [ ! -L "$dir/$project" ]; then
    mkdir -p "$dir"
    cd "$dir"
    ln -s ../../../../../. $project
    cd "$root"
fi

# Set up the environment to use the workspace.
# Also add Godeps workspace so we build using canned dependencies.
GOPATH="$workspace"
GOBIN="$PWD/build/bin"
export GOPATH GOBIN

# Run the command inside the workspace.
cd "$dir/$project"
PWD="$dir/$project"

# Launch the arguments with the configured environment.
exec "$@"
