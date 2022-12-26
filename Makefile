# .PHONY: clean all toolbox

TOOLBOX_IMAGE = toremo/golang-build:latest
COVERAGE_OUT = .coverage.out

all: bin/owner/actor 

bin/owner/actor: gen/api/owner.pb.go $(COVERAGE_OUT) cmd/owner/actor/main.go pkg/owner/server.go
	go build -o $@ $<

gen/api/%.pb.go: schema/%.proto
	protoc --proto_path=schema --go_out=. --go_opt=M$(notdir $<)=gen/api $<

$(COVERAGE_OUT): $(wildcard gen/**/*) $(wildcard cmd/**/*) $(wildcard pkg/**/*)
	go test -coverprofile $@ ./...

clean:
	rm -Rf gen $(COVERAGE_OUT)

# dot not forget - you can also combine notdir and basename: protoc --proto_path=schema --go_out=. --go_opt=M$(notdir $<)=gen/api/$(notdir $(basename $<)) $<
# note on the use of $< - it always returns the first dependency for which there is matching rule!