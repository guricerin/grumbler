#!/bin/bash

readonly SCRIPT_ROOT=$(cd $(dirname ${0}); pwd)

main() {
    cp $SCRIPT_ROOT/pre-commit $SCRIPT_ROOT/.git/hooks
}

main
