# Proto Files

The proto files is used for gRPC communication between the client and the server. The proto files are located in the `proto` directory.

## Usage

To generate the proto files, you need to have the `protoc` compiler installed. We use the library `protoc-gen-go` to generate the Go files. You can install the library by running:

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
```

For macOS, you also need to set the `PATH` variable to include the `bin` directory of the `protoc-gen-go` library. You can do this by running:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Here is an example of how to generate the proto files:

```bash
protoc --go_out=. --go-grpc_out=. proto/fingerprint.proto
```
