gopath:=$(shell go env gopath)

.phony: all
all:
	@echo "hello world!"

.phony: vet
vet:
	@go vet cmd/jjsdapi/jjsdapi.go

#启动环境
.phony: run
run:
	@go run cmd/jjsdapi/jjsdapi.go

.phony: drun
drun:
	docker build -f Dockerfile -t jjsdapi/v1 .  && docker run -P --name jjsdapi -d jjsdapi/v1

.phony: admin
admin:
	grass -fb=admin