FROM golang:alpine
RUN apk add --no-cache git terraform; \
  git clone https://github.com/hashicorp/terraform $GOPATH/src/github.com/hashicorp/terraform; \
  go get gopkg.in/yaml.v2
