export SERVICE_DEBUG=true
export SERVICE_DIR=ipcmanview_data

-include .env

TOOL_AIR=github.com/cosmtrek/air@v1.51.0
TOOL_GOOSE=github.com/pressly/goose/v3/cmd/goose@v3.23.0
TOOL_RESTISH=github.com/danielgtaylor/restish@v0.20.0

migration:
	atlas migrate diff $(name) --env local

hash:
	atlas migrate hash --env local

nuke:
	rm -rf $(SERVICE_DIR)
	rm -rf ./internal/sqlite/migrations
	atlas migrate diff initial --env local

web:
	cd internal/web && pnpm install

# ---------- Dev

dev:
	air

dev-proxy:
	go run ./cmd/dev-proxy

dev-web:
	cd internal/web && pnpm run dev

# ---------- OpenAPI

openapi: openapi-json openapi-web

openapi-json:
	go run ./cmd/ipcmanview openapi > ./docs/openapi.json

openapi-web:
	cd ./internal/web && pnpm run generate-openapi-ts

# ---------- Tooling is only required for development

tooling: tooling-air tooling-goose tooling-restish tooling-atlas

tooling-air:
	go install $(TOOL_AIR)

tooling-goose:
	go install $(TOOL_GOOSE)

tooling-restish:
	go install $(TOOL_RESTISH)

tooling-atlas:
	# TODO: pin atlas version
	curl -sSf https://atlasgo.sh | sh
