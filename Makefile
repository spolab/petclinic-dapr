.PHONY: clean all

all: bin/owner/actor 

bin/owner/actor: gen/api/owner.pb.go coverage.out cmd/owner/actor/main.go pkg/owner/server.go
	go build -o $@ $<

gen/api/%.pb.go: schema/%.proto
	protoc --proto_path=schema --go_out=. --go_opt=M$(notdir $<)=gen/api $<

coverage.out: $(wildcard gen/**/*) $(wildcard cmd/**/*) $(wildcard pkg/**/*)
	go test -coverprofile $@ ./...

clean:
	rm -Rf gen coverage.out

# dot not forget - you can also combine notdir and basename: protoc --proto_path=schema --go_out=. --go_opt=M$(notdir $<)=gen/api/$(notdir $(basename $<)) $<
# note on the use of $< - it always returns the first dependency for which there is matching rule!