#!/bin/sh

pwd |grep -q 'github.com/ziutek/de$' || exit 1

rm -rf de32
mkdir de32
cp *.go de32
cd de32
sed -i 's/loat64/loat32/g' *.go
sed -i 's/ackage de/ackage de32/g' *.go
sed -i 's/ziutek\/matrix/ziutek\/matrix\/matrix32/g' *.go
sed -i 's/matrix\./matrix32./g' *.go