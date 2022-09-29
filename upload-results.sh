#!/bin/bash

file=$1

tar -cvf "$file.tar" $file
gzip "$file.tar"

aws --profile personal s3 cp "$file.tar.gz" "s3://eigentrust-poc-files/$file.tar.gz"
