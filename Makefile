.SILENT:

.PHONY: prompt_user
# Internal helper target - prompt the user before continuing
prompt_user:
	@echo "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: list ## List all make targets
list:
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help ## Prints all the targets in all the Makefiles
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: go_version_check
# Internal helper target - check go version
go_version_check:
	@# Extract the version number from the `go version` command.
	@GO_VERSION=$$(go version | cut -d " " -f 3 | cut -c 3-) && \
	MAJOR_VERSION=$$(echo $$GO_VERSION | cut -d "." -f 1) && \
	MINOR_VERSION=$$(echo $$GO_VERSION | cut -d "." -f 2) && \
	\
	if [ "$$MAJOR_VERSION" -gt 1 ] || ( [ "$$MAJOR_VERSION" -eq 1 ] && [ "$$MINOR_VERSION" -ge 20 ] ); then \
		echo "Invalid Go version. Expected 1.19.x but found $$GO_VERSION"; \
		exit 1; \
	fi

.PHONY: docker_check
# Internal helper target - check if docker is installed
docker_check:
	{ \
	if ( ! ( command -v docker >/dev/null && (docker compose version >/dev/null || command -v docker-compose >/dev/null) )); then \
		echo "Seems like you don't have Docker or docker-compose installed. Make sure you review build/localnet/README.md and docs/development/README.md  before continuing"; \
		exit 1; \
	fi; \
	}

.PHONY: protoc_check
protoc_check: ## Checks if protoc is installed
	{ \
	if ! command -v protoc >/dev/null; then \
		echo "Follow instructions to install 'protoc': https://grpc.io/docs/protoc-installation/"; \
	fi; \
	}


.PHONY: warn_destructive
warn_destructive: ## Print WARNING to the user
	@echo "This is a destructive action that will affect docker resources outside the scope of this repo!"


.PHONY: docker_wipe
docker_wipe: docker_check warn_destructive prompt_user ## [WARNING] Remove all the docker containers, images and volumes.
	docker ps -a -q | xargs -r -I {} docker stop {}
	docker ps -a -q | xargs -r -I {} docker rm {}
	docker images -q | xargs -r -I {} docker rmi {}
	docker volume ls -q | xargs -r -I {} docker volume rm {}

# Reference: https://github.com/rollkit/local-celestia-devnet
.PHONY: celestia_localnet
celestia_localnet: docker_check  ## Run a celestia localnet
	docker run --name celestia --platform linux/amd64 -p 26657:26657 -p 26658:26658 -p 26659:26659 ghcr.io/rollkit/local-celestia-devnet:v0.11.0-rc8

.PHONY: celestia_localnet_stop
celestia_localnet_stop: docker_check  ## Stop the celestia localnet
	docker stop celestia

.PHONY: celestia_light_client_start
celestia_light_client_start: docker_check  ## Start the celestia light client
	echo "See the following link if there's an error https://docs.celestia.org/nodes/light-node/#install-celestia-node"
	celestia light start --core.ip consensus-validator-arabica-9.celestia-arabica.com --p2p.network arabica

# Intended to be called like so: `export AUTH_TOKEN=$(make celestia_localnet_auth_token)`
.PHONY: celestia_localnet_auth_token
celestia_localnet_auth_token: docker_check  ## Get the auth token for the celestia localnet
	CONTAINER_ID=$$(docker ps -qf "name=celestia"); \
	AUTH_TOKEN=$$(docker exec $$CONTAINER_ID celestia bridge --node.store /bridge auth admin); \
	echo $$AUTH_TOKEN

.PHONY: celestia_localnet_balance_check
celestia_localnet_balance_check: docker_check  ## Check the balance of an account in the celestia localnet
	AUTH_TOKEN=$$(make celestia_localnet_auth_token); \
	CONTAINER_ID=$$(docker ps -qf "name=celestia"); \
	docker exec $$CONTAINER_ID /bin/sh -c "celestia rpc state Balance --auth $$AUTH_TOKEN"

# Useful if you want to run `apk update &&  apk add busybox-extras`
.PHONY: celestia_localnet_exec_root
celestia_localnet_exec_root: docker_check  ## Execu into the container as root user in the celestia localnet
	docker exec -it --user=root celestia /bin/sh

.PHONY: poktroll_local_start
poktroll_local_start: docker_check go_version_check ## Start the localnet poktroll node
	@AUTH_TOKEN=$$(make celestia_localnet_auth_token) ./build/init-local.sh

.PHONY: poktroll_testnet_start
poktroll_testnet_start: docker_check go_version_check ## Start the testnet poktroll node
	@AUTH_TOKEN=$$(make celestia_localnet_auth_token) ./build/init-testnet.sh

.PHONY: poktroll_clear
poktroll_clear: ## Clear the poktroll state
	rm -rf ${HOME}/.poktroll
	rm ${HOME}/go/bin/poktrolld

.PHONY: poktroll_list_keys
poktroll_list_keys: ## List the poktroll keys
	poktrolld keys list --keyring-backend test

.PHONY: poktroll_send
poktroll_send: ## Send tokens from one key to another
	KEY1=$$(make poktroll_list_keys | awk -F' ' '/address: pokt1/{print $$3}' | head -1); \
	KEY2=$$(make poktroll_list_keys | awk -F' ' '/address: pokt1/{print $$3}' | tail -1); \
	poktrolld tx bank send $$KEY1 $$KEY2 42069stake --keyring-backend test --node tcp://127.0.0.1:36657

# Protobuf convenience targets

.PHONY: go_protoc-go-inject-tag
go_protoc-go-inject-tag: ## Checks if protoc-go-inject-tag is installed
	{ \
	if ! command -v protoc-go-inject-tag >/dev/null; then \
		echo "Install with 'go install github.com/favadi/protoc-go-inject-tag@latest'"; \
	fi; \
	}

.PHONY: protogen_show
protogen_show: ## A simple `find` command that shows you the generated protobufs.
	find . -name "*.pb.go" | grep -v -e "prototype" -e "vendor"

.PHONY: protogen_clean
protogen_clean: ## Remove all the generated protobufs.
	find . -name "*.pb.go" | grep -v -e "prototype" -e "vendor" | xargs -r rm

# IMPROVE: Look into using buf in the future; https://github.com/bufbuild/buf.
PROTOC = protoc --experimental_allow_proto3_optional --go_opt=paths=source_relative
PROTOC_SHARED = $(PROTOC) -I=./types/proto

.PHONY: protogen_local
protogen_local: go_protoc-go-inject-tag ## Generate go structures for all of the protobufs
# $(PROTOC) -I=./codec/proto           --go_out=./codec           ./codec/proto/*.proto
	$(PROTOC_SHARED)                     --go_out=./types           ./types/proto/*.proto

	$(PROTOC) -I=./runtime/configs/types/proto				              --go_out=./runtime/configs/types ./runtime/configs/types/proto/*.proto
	$(PROTOC) -I=./runtime/configs/proto -I=./runtime/configs/types/proto --go_out=./runtime/configs       ./runtime/configs/proto/*.proto


	# echo "View generated proto files by running: make protogen_show"
