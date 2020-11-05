
pack:
	go build -o	$(GOPATH)/bin/wengine-pack cmd/pack.go

unpack:
	go build -o	$(GOPATH)/bin/wengine-unpack cmd/unpack.go