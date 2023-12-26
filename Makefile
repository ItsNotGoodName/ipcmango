export DIR=ipcmanview_data
export VITE_HOST=127.0.0.1

-include .env

_:
	mkdir web/dist -p && touch web/dist/index.html

migrate:
	goose -dir internal/migrations/sql sqlite3 "$(DIR)/sqlite.db" up

clean:
	rm -rf $(DIR)

generate:
	go generate ./...

run:
	go run ./cmd/ipcmanview serve

preview: generate run

nightly:
	task nightly

migration:
	atlas migrate diff $(name) --env local

hash:
	atlas migrate hash --env local

# Dev

dev:
	air

dev-web:
	cd internal/web && pnpm install && pnpm run dev

dev-webnext:
	cd internal/webnext && pnpm install && pnpm run dev

# Gen

gen: gen-sqlc gen-pubsub gen-bus gen-proto

gen-sqlc:
	sqlc generate

gen-pubsub:
	sh ./scripts/generate-pubsub-events.sh ./internal/models/event.go

gen-bus:
	go run ./scripts/generate-bus.go -input ./internal/models/event.go -output ./internal/core/bus.gen.go

gen-proto:
	cd rpc && protoc --go_out=. --twirp_out=. rpc.proto
	cd internal/webnext && pnpm exec protoc --ts_out=./src/twirp --ts_opt=generate_dependencies --proto_path=../../rpc rpc.proto

# Tooling

tooling: tooling-air tooling-task tooling-goose tooling-atlas tooling-sqlc tooling-protoc-gen-go

tooling-air:
	go install github.com/cosmtrek/air@latest

tooling-task:
	go install github.com/go-task/task/v3/cmd/task@latest

tooling-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

tooling-atlas:
	go install ariga.io/atlas/cmd/atlas@latest

tooling-sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

tooling-twirp:
	go install github.com/twitchtv/twirp/protoc-gen-twirp@latest

tooling-protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install

install-protoc:
	PROTOC_ZIP=protoc-25.1-linux-x86_64.zip
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v25.1/$PROTOC_ZIP
	sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
	sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
