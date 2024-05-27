export SERVICE_DEBUG=true
export SERVICE_DIR=ipcmanview_data

-include .env

TOOL_AIR=github.com/cosmtrek/air@v1.51.0
TOOL_GOOSE=github.com/pressly/goose/v3/cmd/goose@v3.20.0
TOOL_RESTISH=github.com/danielgtaylor/restish@v0.20.0

dev:
	air

clean:
	rm -rf $(SERVICE_DIR)

restish-sync:
	restish api sync huma

migration:
	atlas migrate diff $(name) --env local

hash:
	atlas migrate hash --env local

nuke: clean
	rm -rf ./internal/sqlite/migrations
	atlas migrate diff initial --env local

# ---------- Tooling is only required for development

tooling: tooling-air tooling-goose tooling-restish tooling-atlas

tooling-air:
	go install $(TOOL_AIR)

tooling-goose:
	go install $(TOOL_GOOSE)

tooling-restish:
	go install $(TOOL_GOOSE)

tooling-atlas:
	# TODO: pin atlas version
	curl -sSf https://atlasgo.sh | sh
