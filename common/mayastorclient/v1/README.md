To generate the client code using nix shell
1. Edit the proto file and add the line
    * `option go_package = "github.com/openebs/mayastor-api/protobuf/v1";`

2. run below command to update dependency
	`nix-shell -p protoc-gen-go-grpc protobuf protoc-gen-go`

3. generate code using the following command, run from the `protobuf` directory
    `protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative <protofile_name>`

