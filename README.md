# Kratos Project Template

## Install Kratos
```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```bash
# Create a template project
kratos new server -r https://github.com/ofavor/kratos-layout.git
# OR 
kratos new server -r https://gitee.com/ofavor/kratos-layout.git

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```
## Generate other auxiliary files by Makefile
```bash
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```
## Automated Initialization (wire)
```bash
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

OR

```bash
# $(SERVER) is the name of server (under cmd/) you want to generate
make wire SERVER=server
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

## OpenAPI
http://localhost:8000/q/swagger-ui/

Select a definition on the top right.

## VSCode debug

launch.json

```json
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Server",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "cwd": "${workspaceFolder}",
      "program": "${workspaceFolder}/cmd/server",
      "args": ["-conf", "${workspaceFolder}/configs"]
    }
  ]
}
```