all: gen

clean:
	find . -name '*.gen.go' -delete

gen: clean generate
generate:
	GODEBUG=gotypesalias=0 go generate ./...