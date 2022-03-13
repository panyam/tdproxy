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
PYPKGNAME=pytdproxy
# Generate the python client too
git submodule update --init
PY_OUT_DIR="./protos/$PYPKGNAME"
mkdir -p $PY_OUT_DIR/src/$PYPKGNAME
# This is silly - needed to preserve python's directory structure needs
cp protos/*.proto protos/$PYPKGNAME
python3 -m grpc_tools.protoc -I./protos   \
      --python_out="$PY_OUT_DIR"          \
      --grpc_python_out="$PY_OUT_DIR"     \
      protos/$PYPKGNAME/auth.proto \
      protos/$PYPKGNAME/chains.proto \
      protos/$PYPKGNAME/streamer.proto \
      protos/$PYPKGNAME/tickers.proto
rm protos/$PYPKGNAME/*.proto
mv protos/$PYPKGNAME/$PYPKGNAME/* $PY_OUT_DIR/src/$PYPKGNAME
touch $PY_OUT_DIR/src/$PYPKGNAME/__init__.py
