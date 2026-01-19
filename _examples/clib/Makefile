prebuild:
	mkdir -p ./bin
	if [ "$$(go env GOOS)" = "darwin" ]; then \
		GOOS=darwin GOARCH=$$(go env GOARCH) go build -buildmode=c-archive -o ./bin/libkagome.a ./go/main.go; \
	elif [ "$$(go env GOOS)" = "windows" ]; then \
		GOOS=windows GOARCH=$$(go env GOARCH) go build -buildmode=c-archive -o ./bin/libkagome.a ./go/main.go; \
	elif [ "$$(go env GOOS)" = "linux" ]; then \
		GOOS=linux GOARCH=$$(go env GOARCH) go build -buildmode=c-archive -o ./bin/libkagome.a ./go/main.go; \
	else \
		echo "Unsupported GOOS: $$(go env GOOS)" ; exit 1 ; \
	fi

build: prebuild
	if [ "$$(go env GOOS)" = "darwin" ]; then \
		cc -fPIC -shared ./c/kagome_wrapper.c ./bin/libkagome.a -o ./bin/libkagome.dylib; \
	elif [ "$$(go env GOOS)" = "windows" ]; then \
		cc -fPIC -shared ./c/kagome_wrapper.c ./bin/libkagome.a -o ./bin/libkagome.dll; \
	elif [ "$$(go env GOOS)" = "linux" ]; then \
		cc -fPIC -shared ./c/kagome_wrapper.c ./bin/libkagome.a -o ./bin/libkagome.so; \
	else \
		echo "Unsupported GOOS: $$(go env GOOS)" ; exit 1 ; \
	fi

test:
	@if [ -f ./bin/libkagome.so ] || [ -f ./bin/libkagome.dll ] || [ -f ./bin/libkagome.dylib ]; then \
		python3 ./python3/main.py ; \
	else \
		echo "No libkagome.(so|dll|dylib) found. Please run make build first." ; exit 1 ; \
	fi
