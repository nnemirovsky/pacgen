ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY:
.SILENT:


build:
	go build -v -o bin/server cmd/server/main.go
	go build -v -o bin/migrator cmd/migrator/main.go
	go build -v -o bin/generator cmd/generator/main.go

test:
	go test -v ./...

gofmt:
	go fmt ./...

govet:
	go vet ./...

gofix:
	go fix ./...

# --- development only ---

build_img:
	#docker buildx build --no-cache --platform linux/amd64,linux/arm64 --load -t pacgen:latest .
	docker build --no-cache -t pacgen:latest .

serve:
	go run cmd/server/main.go

generate_pac:
	go run cmd/generator/main.go

create_migration:
	docker run \
		-v $(ROOT_DIR)/migrations:/migrations \
  		migrate/migrate create -dir migrations -ext sql -seq $(name)

migrate_up:
	go run cmd/migrator/main.go up

migrate_down:
	go run cmd/migrator/main.go down
