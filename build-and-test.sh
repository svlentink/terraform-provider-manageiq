#!/usr/bin/env sh
set -ev

if [ -z "$GOPATH" ]; then
  echo "GOPATH not set"
  exit 1
fi
cp -a client $GOPATH/src/

mkdir -p ~/.terraform.d/plugins
go build -o ~/.terraform.d/plugins/terraform-provider-manageiq

terraform init
terraform plan

