GOROOT=`which go`
GOPATH=$HOME/go
GOBIN=$GOPATH/bin
PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
echo "Generating GO bindings..."
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      protos/auth.proto \
      protos/chains.proto \
      protos/streamer.proto \
      protos/tickers.proto

echo "Generating Python bindings..."
# Generate the python client too
git submodule update --init
PY_OUT_DIR="./protos/pytdproxy/"
# This is silly - needed to preserve python's directory structure needs
cp protos/*.proto protos/pytdproxy
python3 -m grpc_tools.protoc -I./protos  \
      --python_out="$PY_OUT_DIR"  \
      --grpc_python_out="$PY_OUT_DIR"  \
      protos/pytdproxy/auth.proto \
      protos/pytdproxy/chains.proto \
      protos/pytdproxy/streamer.proto \
      protos/pytdproxy/tickers.proto
