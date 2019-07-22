#!/usr/bin/env sh
set -ev

rm -r /temp-build || true
cp -a $PWD /temp-build
cd /temp-build
echo path $GOPATH
cp -a client $GOPATH/src/

# based on https://www.terraform.io/docs/extend/writing-custom-providers.html
go build -o ./terraform-provider-manageiq

terraform init
terraform plan

