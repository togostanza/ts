#!/bin/bash

set -eu

DESTDIR='provider/assets/components/webcomponentsjs'

rm -rf $DESTDIR
mkdir -p $DESTDIR
cp node_modules/@webcomponents/webcomponentsjs/webcomponents-*{.js,.js.map} $DESTDIR
