.PHONY: clean all

DIR_BIN = bin
COVERAGE_OUT = .coverage.out
VET_API_MAIN = cmd/vet/api/main.go
VET_API_OUT = $(DIR_BIN)/vet/api
VET_API_DEPS = $(VET_API_MAIN) pkg/vet/api.go pkg/vet/command.go pkg/vet/events.go
VET_ACTOR_MAIN = cmd/vet/actor/main.go
VET_ACTOR_OUT = $(DIR_BIN)/vet/actor
VET_ACTOR_DEPS = $(VET_ACTOR_MAIN) pkg/vet/actor.go pkg/vet/command.go pkg/vet/events.go


all: $(COVERAGE_OUT) $(VET_API_OUT) $(VET_ACTOR_OUT) 

$(VET_ACTOR_OUT): $(VET_ACTOR_DEPS)
	go build -o $@ $(VET_ACTOR_MAIN)

$(VET_API_OUT): $(VET_API_DEPS)
	go build -o $@ $(VET_API_MAIN)

gen/api/%.pb.go: schema/%.proto
	protoc --proto_path=schema --go_out=. --go_opt=M$(notdir $<)=gen/api --go_opt=Mcommon.proto=gen/api $<

$(COVERAGE_OUT): $(wildcard gen/**/*) $(wildcard cmd/**/*) $(wildcard pkg/**/*)
	go test -coverprofile $@ ./...

clean:
	rm -Rf gen $(DIR_BIN) $(COVERAGE_OUT)

# dot not forget - you can also combine notdir and basename: protoc --proto_path=schema --go_out=. --go_opt=M$(notdir $<)=gen/api/$(notdir $(basename $<)) $<
# note on the use of $< - it always returns the first dependency for which there is matching rule!