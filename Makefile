all: build

build:
	go build -o out/hydrun main.go

release:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o out/release/hydrun.linux-$$(uname -m) main.go

clean:
	rm -rf out
