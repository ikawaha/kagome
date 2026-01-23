# Detect platform once
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Platform-specific library extension
ifeq ($(GOOS),darwin)
    LIB_EXT := dylib
else ifeq ($(GOOS),windows)
    LIB_EXT := dll
else ifeq ($(GOOS),linux)
    LIB_EXT := so
else
    $(error Unsupported GOOS: $(GOOS))
endif

prebuild:
	@mkdir -p ./bin
	@rm -rf ./bin/lib*
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -buildmode=c-archive -o ./bin/libkagome.a ./go/main.go

# Use -fvisibility=hidden to only export symbols explicitly marked in the wrapper (for a clean public API)
build: prebuild
	cc -fPIC -fvisibility=hidden -shared ./c_wrapper/kagome_wrapper.c ./bin/libkagome.a -o ./bin/libkagome.$(LIB_EXT)

go-test:
	@if ! command -v go >/dev/null 2>&1; then \
		echo "[SKIP] Go is not installed."; \
		exit 0; \
	fi; \
	echo "**NOTE**: On macOS, 'ld: warning: ...' can be ignored if tests pass."; \
	go clean -testcache && go test -race ./...

python-test:
	@if ! command -v python3 >/dev/null 2>&1; then \
		echo "[SKIP] Python3 is not installed."; \
		exit 0; \
	fi; \
	if [ -f ./bin/libkagome.so ] || [ -f ./bin/libkagome.dll ] || [ -f ./bin/libkagome.dylib ]; then \
		python3 ./python3/main.py ; \
	else \
		echo "No libkagome.(so|dll|dylib) found. Please run make build first." ; exit 1 ; \
	fi

php-test:
	@if ! command -v php >/dev/null 2>&1; then \
		echo "[SKIP] PHP is not installed."; \
		exit 0; \
	fi; \
	if [ -f ./bin/libkagome.so ] || [ -f ./bin/libkagome.dll ] || [ -f ./bin/libkagome.dylib ]; then \
		php ./php8/main.php ; \
	else \
		echo "No libkagome.(so|dll|dylib) found. Please run make build first." ; exit 1 ; \
	fi

test: go-test python-test php-test

clean:
	@rm -rf ./bin

docker-clean: clean
	@docker compose down --rmi all --volumes --remove-orphans
	@docker container prune -f
	@docker image prune -f
	@docker volume prune -f
	@docker network prune -f

docker-build:
	# Pull latest base images for stability
	@docker pull golang:latest
	@docker pull python:latest
	@docker pull php:latest
	# Build all images without cache
	@docker compose build --no-cache

docker-test:
	@docker compose up --exit-code-from python-test --exit-code-from php-test
