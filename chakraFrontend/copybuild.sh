#!/usr/bin/sh

rm -r ../modules/cis --force;

mkdir ../modules/cis;

cd ./build/client;

for i in *;
do mv   -v  $i ../../../modules/cis/$i ;
done;