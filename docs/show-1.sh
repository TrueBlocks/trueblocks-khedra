#!/usr/bin/env bash

input=$1
name=${input%.md}

mmdc -w 1600 -H 1200 -i $name.md -o img-$name.png --theme forest
open img-$name-1.png

