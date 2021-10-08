.PHONY: build

GO_BUILD_CONTAINER := "golang:1.17.2-buster"

build:
ifndef GOOS
	@echo "[error] \$$GOOS must be specified"
	@exit 1
endif
ifndef GOARCH
	@echo "[error] \$$GOARCH must be specified"
	@exit 1
endif
	$(eval TAG := $(shell git describe --tags --abbrev=0))
	docker run -it -v $(PWD):/conk -w /conk \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		$(GO_BUILD_CONTAINER) \
		go build -o ./bin/conk_$(GOOS)_$(GOARCH)_$(TAG) \
		-ldflags '-X "github.com/moznion/conk/internal.Revision=$(shell git rev-parse HEAD)" -X "github.com/moznion/conk/internal.Version=$(TAG)"' \
		./cmd/conk/main.go

