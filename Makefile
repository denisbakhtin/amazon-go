#Basic makefile

default: build

build: clean vet
	@go build -o amazon-go

doc:
	@godoc -http=:6060 -index

lint:
	@golint ./...

debug_server: 
	@watcher
debug_assets:
	@npm run watch

#run 'make -j2 debug' to launch both servers in parallel
debug: clean debug_server debug_assets 

run: build
	./amazon-go

test:
	@go test ./...

vet:
	@go vet ./...

clean:
	@rm -f ./amazon-go
