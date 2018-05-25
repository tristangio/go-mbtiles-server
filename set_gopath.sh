#!/bin/bash

## This file show you what to add to your GOPATH to do a simple build
## Just run it in the instalation directory

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

CMD='export GOPATH=$GOPATH:"'$DIR'"'
echo "to do a simple go build, add to your GOPATH the current dir with :"
echo "$CMD"
