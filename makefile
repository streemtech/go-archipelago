all: gen build run

clean:
	find . -name '*.gen.go' -delete

gen: clean generate
generate:
	GODEBUG=gotypesalias=0 go generate ./...



run: 
	./go-archipelago

build: 
	go build

