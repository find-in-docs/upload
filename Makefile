# This makes sure the commands are run within a BASH shell.
SHELL := /bin/bash
EXEDIR := ./bin
EXENAME := upload
BIN_NAME=./${EXEDIR}/${EXENAME}

# The .PHONY target will ignore any file that exists with the same name as the target
# in your makefile, and built it regardless.
.PHONY: all init build run clean upload

# The all target is the default target when make is called without any arguments.
all: clean | run

init:
	echo "Setting up local ..."
	go env -w GOPROXY=direct 
	echo "----------------------------------------------------"
	echo "To get protoc, look here:"
	echo "  example: https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip"
	echo "To install protoc-gen-go-grpc, do this:"
	echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	echo "To install protoc-gen-go, do this:"
	echo "  go install google.golang.org/grpc/cmd/protoc-gen-go"
	echo "----------------------------------------------------"
	- rm go.mod
	- rm go.sum
	go mod init github.com/find-in-docs/upload
	go mod tidy

${EXEDIR}:
	echo "Building exe directory ..."
	mkdir ${EXEDIR}

build: | ${EXEDIR}
	echo ">>>>>>>>>>>>>>>>>"
	echo "  Get latest tagged version for your code (ex. sidecar) that you depend on."
	echo "  Use: go get github.com/find-in-docs/sidecar@v0.0.0-beta.10-lw (for example)."
	echo "  This will ensure your packages from github are synced with Google's servers"
	echo "  (https://proxy.golang.org). You can change this value using:"
	echo "  go env -w GOPROXY=direct"
	echo "  to get it from the github repo directly."
	echo "  Google's server might take 30 minutes to sync up with github after you request"
	echo "  your package from them. So the first time you request it after the package change occurs"
	echo "  in github, it will get it from github directly, then add your repo to their syncing process."
	echo "<<<<<<<<<<<<<<<<<"
	
	echo "Building locally ..."
	go build -o ${BIN_NAME} pkg/main/main.go

run: build
	echo "Running locally ..."
	./${BIN_NAME}

clean:
	echo "Cleaning locally ..."
	go clean
	- rm ${BIN_NAME}
	go clean -cache -modcache -i -r
	go mod tidy

upload: build
	echo "Start building on minikube ..."
	# echo "Get each of these packages in the Dockerfile"
	# rg --iglob "*.go" -o -I -N "[\"]github([^\"]+)[\"]" | sed '/^$/d' | sed 's/\"//g' | awk '{print "RUN go get " $0}'
	# docker build --progress=plain --no-cache -t upload -f ./Dockerfile .
	docker pull nats:latest
	docker pull postgres:latest
	# docker build -t postgres -f Dockerfile_postgres .

	# docker build -t upload -f ./Dockerfile .
	# If you want to see the output of your RUN commands present in the Dockerfile, do:
	# docker build --progress=plain --no-cache -t upload -f ./Dockerfile .
	docker build -t upload -f ./Dockerfile .
	
	# We specify image-pull-policy=Never because we're actually building the image on minikube.
	# kubectl run upload --image=persistlogs:latest --image-pull-policy=Never --restart=Never

	# kubectl apply -f manifests/minikube
