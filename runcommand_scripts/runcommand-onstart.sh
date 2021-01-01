#!/bin/sh

rom="\"`basename $3`\""
line="START $1 $rom"
timeout 1 bash -c "echo $line > /tmp/gamestats"
