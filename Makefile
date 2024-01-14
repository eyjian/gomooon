# Package gomooon
# Wrote by yijian on 2024/01/04

all: test hmac_sha256_sign

test: main.go
	go build -o $@ main.go

hmac_sha256_sign: hmac_sha256_sign.go
	go build -o $@ $<

.PHONY: clean tidy

clean:
	rm -f test hmac_sha256_sign

tidy:
	go mod tidy
