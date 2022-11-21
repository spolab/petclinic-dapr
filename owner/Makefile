.PHONY: schema build image clean test

VERSION = v0.0.2
DIR_GEN = pkg/api
DIR_SCHEMA = schema
DIR_BIN = bin
PROTOC = protoc --go_out=. --go_opt=module=github.com/spolab/petclinic/owner --go_opt=paths=import
FILE_OUT = $(DIR_BIN)/app

image: build
	nerdctl build --namespace k8s.io -f docker/Dockerfile --tag spolab/petclinic-owner:$(VERSION) .

build: schema
	go vet ./...
	go test ./...
	cd cmd && go build -o ../$(FILE_OUT) .

schema:
	$(PROTOC) $(DIR_SCHEMA)/*.proto

clean:
	rm -Rf bin/ $(DIR_GEN)/command $(DIR_GEN)/model $(DIR_GEN)/event $(DIR_BIN)

test:
	dapr run --app-id owner --components-path dapr/test --log-level debug go test ./...