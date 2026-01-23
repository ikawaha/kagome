prebuild:
	@mkdir -p ./bin
	@rm -rf ./bin/lib*
	@if [ "$$(go env GOOS)" = "darwin" ]; then \
		GOOS=darwin GOARCH=$$(go env GOARCH) go build -buildmode=c-archive -o ./bin/libkagome.a ./go/main.go; \
	elif [ "$$(go env GOOS)" = "windows" ]; then \
		GOOS=windows GOARCH=$$(go env GOARCH) go build -buildmode=c-archive -o ./bin/libkagome.a ./go/main.go; \
	elif [ "$$(go env GOOS)" = "linux" ]; then \
		GOOS=linux GOARCH=$$(go env GOARCH) go build -buildmode=c-archive -o ./bin/libkagome.a ./go/main.go; \
	else \
		echo "Unsupported GOOS: $$(go env GOOS)" ; exit 1 ; \
	fi
	@sleep 1; # cooldown for file system

# Use -fvisibility=hidden to only export symbols explicitly marked in the wrapper (for a clean public API)
build: prebuild
	@if [ "$$(go env GOOS)" = "darwin" ]; then \
		cc -fPIC -fvisibility=hidden -shared ./c_wrapper/kagome_wrapper.c ./bin/libkagome.a -o ./bin/libkagome.dylib; \
	elif [ "$$(go env GOOS)" = "windows" ]; then \
		cc -fPIC -fvisibility=hidden -shared ./c_wrapper/kagome_wrapper.c ./bin/libkagome.a -o ./bin/libkagome.dll; \
	elif [ "$$(go env GOOS)" = "linux" ]; then \
		cc -fPIC -fvisibility=hidden -shared ./c_wrapper/kagome_wrapper.c ./bin/libkagome.a -o ./bin/libkagome.so; \
	else \
		echo "Unsupported GOOS: $$(go env GOOS)" ; exit 1 ; \
	fi
	@sleep 1; # cooldown for file system

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

test: python-test php-test

clean:
	@rm -rf ./bin

docker-clean: clean
	@docker compose down --rmi all --volumes --remove-orphans
	@docker container prune -f
	@docker image prune -f
	@docker volume prune -f
	@docker network prune -f
	@docker volume rm -f clib_kagome-dist
	@docker image rm -f clib-kagome-build:latest clib-php-test:latest

docker-build:
	# Pull latest base images for stability
	@docker pull golang:latest
	@docker pull python:latest
	@docker pull php:latest
	# Build all images without cache
	@docker compose build --no-cache

docker-test:
	@docker compose up --exit-code-from python-test --exit-code-from php-test
