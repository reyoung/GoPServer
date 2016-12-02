#!/bin/bash
rm -rf GoPServer
flatc ParameterService.fbs  --go
cp -r GoPServer/* ../
rm -rf GoPServer
