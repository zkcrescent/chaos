build:
	go build ./...
.PHONY: build
build:
	go build -o ./bin/ ./gorp-util/gen/...
.PHONY: gen
gen:
	cd ./gorp-util/gen/sample/ && ../../../bin/gen 
