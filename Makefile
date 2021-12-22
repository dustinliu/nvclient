app := nvclient
build_dir := build
dist_dir := dist

windows = $(app).exe
linux = $(app)
darwin = $(app)

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
target_dir = '$(build_dir)/$(os)-$(arch)'
executable = $($(os))

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

build: $(PLATFORMS)

$(PLATFORMS):
	@echo building $(os)/$(arch)...
	@mkdir -p $(target_dir)
	@mkdir -p $(dist_dir)
	@GOOS=$(os) GOARCH=$(arch) go build -o $(target_dir)/$(executable)
	@tar zcf $(dist_dir)/$(app)-$(os)-$(arch).tar.gz -C $(target_dir) $(executable)

vet:
	@echo running go vet...
	@go vet ./...
	@echo

test:
	@echo testing...
	@go test -timeout 10s ./...
	@echo

clean:
	@go clean
	@go clean -testcache
	@rm -rf build
	@rm -rf dist

all: build

.PHONY: build clean test vet $(PLATFORMS)
