#FROM golang:alpine
FROM alpine

ENV GOPATH=/root/go/
ENV TF_LOG=DEBUG
ENV TF_LOG_PATH=/dev/stdout
ENTRYPOINT ["/usr/bin/terraform"]
CMD ["init"]

RUN apk add --no-cache \
    git \
    go \
    musl-dev \
    terraform; \
  mkdir -p $GOPATH/src; \
  mkdir -p ~/.terraform.d/plugins; \
  git clone https://github.com/hashicorp/terraform $GOPATH/src/github.com/hashicorp/terraform

COPY . /terraform-provider-manageiq
WORKDIR /terraform-provider-manageiq
RUN cp -a client $GOPATH/src/; \
  go build -o ~/.terraform.d/plugins/terraform-provider-manageiq; \
  mkdir -p /data; \
  cp example.tf /data/;

# mount the following dir
WORKDIR /data
