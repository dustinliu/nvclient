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
archive = $(dist_dir)/$(app)-$(os)-$(arch).tar.gz

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	md5 := md5
else
	md5 := md5sum
endif

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

build: generate $(PLATFORMS)

generate:
	@echo 'generate mock codes'
	@go generate ./...

$(PLATFORMS):
	@echo building $(os)/$(arch)...
	@mkdir -p $(target_dir)
	@mkdir -p $(dist_dir)
	@GOOS=$(os) GOARCH=$(arch) go build -ldflags "-X main.version=`cat version`" -o $(target_dir)/$(executable)
	@tar zcf $(dist_dir)/$(app)-$(os)-$(arch).tar.gz -C $(target_dir) $(executable)
	@cd $(dist_dir); $(md5) $(app)-$(os)-$(arch).tar.gz >> checksums.txt

vet: generate
	@echo running go vet...
	@go vet ./...
	@echo

test: generate
	@echo testing...
	@go test -timeout 10s ./...
	@echo

clean:
	@go clean
	@go clean -testcache
	@rm -rf build
	@rm -rf dist

release:
	@echo "git tag `cat version`"
	git push

all: build

.PHONY: build clean test vet generate $(PLATFORMS)
