.PHONY: image clean
.DEFAULT_GOAL: build

VERSION = v0.0.2
DIR_GEN = pkg/api
DIR_PROTO = schema
DIR_BIN = bin
MODEL_PROTO = $(DIR_PROTO)/model.proto
MODEL_GO = $(DIR_GEN)/model/model.pb.go
COMMAND_PROTO = $(DIR_PROTO)/command.proto
COMMAND_GO = $(DIR_GEN)/command/command.pb.go
EVENT_PROTO = $(DIR_PROTO)/event.proto
EVENT_GO = $(DIR_GEN)/event/event.pb.go
PROTOC = protoc --go_out=. --go_opt=module=github.com/spolab/petclinic/owner --go_opt=paths=import

$(DIR_BIN)/owner: $(COMMAND_GO) $(EVENT_GO)
	go vet ./...
	go test ./...
	go build -o bin/owner ./...

$(MODEL_GO): $(MODEL_PROTO)
	$(PROTOC) $<

$(COMMAND_GO): $(MODEL_GO) $(COMMAND_PROTO)
	$(PROTOC) $(COMMAND_PROTO)

$(EVENT_GO): $(MODEL_GO) $(EVENT_PROTO)
	$(PROTOC) $(EVENT_PROTO)

image:
	nerdctl build --namespace k8s.io -f docker/Dockerfile --tag spolab/petclinic-owner:$(VERSION) .

clean:
	rm -Rf bin/ $(DIR_GEN)/command $(DIR_GEN)/model $(DIR_GEN)/event