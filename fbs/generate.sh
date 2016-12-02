#!/bin/bash
rm -rf GoPServer
flatc *.fbs  --go
cp -r GoPServer/* ../
rm -rf GoPServer
