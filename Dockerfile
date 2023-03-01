
# Base image
FROM golang:1.19.5-alpine3.17

# Specify work directory on the image.
# All commands will refer to this work directory from now on below.
WORKDIR /app

# Copy local go.mod and go.sum files into the image
COPY go.mod ./
COPY go.sum ./

# Clean the modcache
# This is not required all the time. You should run this only
# when your modcache contains older versions that you cannot upgrade for some reason.
# RUN go clean -cache -modcache -i -r

# This should always be set. This way, instead of getting packages from Google servers,
# you get them directly from github. This way, there is no sync lag between
# your changes on github and your packages on Google servers. Sometimes, it takes
# more than a day for Google servers to catch up to your changes.
# RUN go env -w GOPROXY=direct 

# Download required packages in the image
# RUN go mod download

# Copy source code into the image
COPY pkg/ /app/pkg/

# This file contains the DNS server information. It is used by the upload
# service to:
#   - Complete the Fully Qualified Domain Name of the request
#   - Locate the IP address of the DNS server
COPY manifests/minikube/resolv.conf /etc/resolv.conf

RUN go build -o upload pkg/main/main.go

RUN mkdir -p /var/lib/postgres/data && \
  chmod 0700 /var/lib/postgres/data

RUN apk add --update util-linux
# RUN whereis initdb

# RUN addgroup -S postgres && adduser -S postgres -G postgres
# USER postgres

# RUN initdb -D /var/lib/postgres/data && \
#    RUN pg_ctl start -D /var/lib/postgres/data

# To see the output of the commands in this file, do:
# docker build --progress=plain --no-cache -t upload -f ./Dockerfile .
# RUN pwd >&2
# RUN ls -l >&2

RUN apk update && \
    apk add tree && \
    apk add bash && \
    apk add curl wget

# By default the ENTRYPOINT is /bin/sh -c.
# We specify CMD to pass to the ENTRYPOINT as an argument,
# so the following command results in /bin/sh -c ./upload
CMD ["./upload"]
