# Package gomooon
# Wrote by yijian on 2024/01/04

all: test

test: main.go
	go build -o $@ main.go

.PHONY: clean tidy

clean:
	rm -f test

tidy:
	go mod tidy
