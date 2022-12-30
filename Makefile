.PHONY: clean all actor

DIR_BIN = bin
COVERAGE_OUT = .coverage.out
VET_API_MAIN = cmd/vet/api/main.go
VET_API_OUT = $(DIR_BIN)/vet/api
VET_API_DEPS = $(VET_API_MAIN) pkg/vet/api.go pkg/vet/command.go pkg/vet/events.go
VET_ACTOR_MAIN = cmd/vet/actor/main.go
VET_ACTOR_OUT = $(DIR_BIN)/vet/actor
ACTOR_DEPS = cmd/actor/main.go $(wildcard pkg/actor/**/*)

actor: bin/actor 

bin/%: cmd/%/main.go
	go build -o $@ $<

gen/api/%.pb.go: schema/%.proto
	protoc --proto_path=schema --go_out=. --go_opt=M$(notdir $<)=gen/api --go_opt=Mcommon.proto=gen/api $<

test:
	go test cmd/actor pkg/actor pkg/common

clean:
	rm -Rf gen $(DIR_BIN) $(COVERAGE_OUT)

# dot not forget - you can also combine notdir and basename: protoc --proto_path=schema --go_out=. --go_opt=M$(notdir $<)=gen/api/$(notdir $(basename $<)) $<
# note on the use of $< - it always returns the first dependency for which there is matching rule!