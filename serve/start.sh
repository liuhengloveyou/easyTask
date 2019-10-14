#!/bin/bash

mkdir ./logs &>/dev/null
nohup ./serv -log_dir="./logs" &>/dev/null &

# ./serv -logtostderr=true
