.PHONY: schema build image clean test

VERSION = v0.0.2
DIR_GEN = pkg/api
DIR_SCHEMA = schema
DIR_BIN = bin
PROTOC = protoc --go_out=. --go_opt=module=github.com/spolab/petclinic/owner --go_opt=paths=import
FILE_OUT = $(DIR_BIN)/app

image: build
	docker build -f docker/Dockerfile --tag spolab/petclinic-owner:$(VERSION) .

build: schema
	go vet ./...
	cd cmd && go build -o ../$(FILE_OUT) .

schema:
	$(PROTOC) $(DIR_SCHEMA)/*.proto

clean:
	rm -Rf bin/ $(DIR_GEN)/command $(DIR_GEN)/model $(DIR_GEN)/event $(DIR_BIN)

run:
	dapr run --app-id owner --app-port 3000 --dapr-grpc-port 5000 --components-path dapr/test go run cmd/owner.go