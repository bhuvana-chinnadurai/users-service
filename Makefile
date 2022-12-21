SHELL = /bin/bash

BIN_DIR = bin
PROTO_DIR = api/proto
SERVER_DIR = api/server
CLIENT_DIR = api/client
RM_F_CMD = rm -f

.PHONY: build

all: build

MODULE = github.com/bhuvana-chinnadurai/users-service

build: codegen dep
	go build -o ${BIN_DIR}/${SERVER_DIR} .

codegen: 
	protoc -I${PROTO_DIR}  --go_out=${PROTO_DIR} --go_opt=module=${MODULE} --go-grpc_out=${PROTO_DIR} --go-grpc_opt=module=${MODULE} ${PROTO_DIR}/*.proto

clean: ## Clean generated files.
	rm -f ${PROTO_DIR}/*.pb.go

dep:
	go mod download
