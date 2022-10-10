
PYPKGNAME=pytdproxy
GOROOT=$(which go)
GOPATH=$(HOME)/go
GOBIN=$(GOPATH)/bin
PATH:=$(PATH):$(GOROOT):$(GOPATH):$(GOBIN)
MAKEFILE_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
APP_SRC_ROOT=$(MAKEFILE_DIR)
# Generate the python client too
PY_OUT_DIR=$(APP_SRC_ROOT)/protos/$(PYPKGNAME)


all: protos

protos: printenv goprotos pyprotos

test:
	cd $(APP_SRC_ROOT)/ && go test ./... -cover

goprotos:
	echo "Generating GO bindings"
	protoc --go_out=$(APP_SRC_ROOT) --go_opt=paths=source_relative          \
       --go-grpc_out=$(APP_SRC_ROOT) --go-grpc_opt=paths=source_relative	\
       --proto_path=$(APP_SRC_ROOT)                     			\
      $(APP_SRC_ROOT)/protos/auth.proto 				\
      $(APP_SRC_ROOT)/protos/chains.proto 			\
      $(APP_SRC_ROOT)/protos/trades.proto 			\
      $(APP_SRC_ROOT)/protos/streamer.proto 		\
      $(APP_SRC_ROOT)/protos/tickers.proto

pyprotos:
	echo "Generating Python bindings"
	# git submodule update --init
	mkdir -p $(PY_OUT_DIR) $(APP_SRC_ROOT)/$(PYPKGNAME)
	python3 -m grpc_tools.protoc -I./protos   \
      --python_out="$(PY_OUT_DIR)"          \
      --grpc_python_out="$(PY_OUT_DIR)"     \
      --proto_path=$(APP_SRC_ROOT)           \
      $(APP_SRC_ROOT)/protos/auth.proto \
      $(APP_SRC_ROOT)/protos/chains.proto \
      $(APP_SRC_ROOT)/protos/trades.proto \
      $(APP_SRC_ROOT)/protos/streamer.proto \
      $(APP_SRC_ROOT)/protos/tickers.proto
	@mv $(PY_OUT_DIR)/protos/*.py $(APP_SRC_ROOT)/$(PYPKGNAME)
	@echo "Cleaning up files..."
	# rm -Rf $(PY_OUT_DIR)

printenv:
	@echo MAKEFILE_DIR=$(MAKEFILE_DIR)
	@echo APP_SRC_ROOT=$(APP_SRC_ROOT)
	@echo MAKEFILE_LIST=$(MAKEFILE_LIST)
	@echo APP_SRC_ROOT=$(APP_SRC_ROOT)
	@echo GOROOT=$(GOROOT)
	@echo GOPATH=$(GOPATH)
	@echo GOBIN=$(GOBIN)
	@echo PYPKGNAME=$(PYPKGNAME)
	@echo PY_OUT_DIR=$(PY_OUT_DIR)

## Setting up the dev db
pgdb:
	docker build -t tdproxypgdb -f Dockerfile.pgdb .

runtestdb:
	mkdir -p $(MAKEFILE_DIR)/pgdata_test
	docker run --rm --name tdproxy-pgdb-container -v ${MAKEFILE_DIR}/pgdata_test:/var/lib/postgresql/data -e POSTGRES_PASSWORD=password -p 5432:5432 tdproxypgdb

rundb:
	mkdir -p $(MAKEFILE_DIR)/pgdata
	docker run --rm --name tdproxy-pgdb-container -p 5432:5432 -v $(MAKEFILE_DIR)/pgdata:/var/lib/postgresql/data tdproxypgdb

