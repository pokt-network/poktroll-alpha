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

# Intended to be called like so: `export CELESTIA_NODE_AUTH_TOKEN=$(make celestia_localnet_auth_token)`
.PHONY: celestia_localnet_auth_token
celestia_localnet_auth_token: docker_check  ## Get the auth token for the celestia localnet
	CELESTIA_CONTAINER_ID=$$(docker ps -qf "name=celestia"); \
	CELESTIA_AUTH=$$(docker exec $$CELESTIA_CONTAINER_ID celestia bridge --node.store /bridge auth admin); \
	export AUTH_TOKEN=$$CELESTIA_AUTH
	echo $$CELESTIA_AUTH

.PHONY: celestia_localnet_balance_check
celestia_localnet_balance_check: docker_check  ## Check the balance of an account in the celestia localnet
	CELESTIA_NODE_AUTH_TOKEN=$$(make celestia_localnet_auth_token); \
	CELESTIA_CONTAINER_ID=$$(docker ps -qf "name=celestia"); \
	docker exec $$CELESTIA_CONTAINER_ID /bin/sh -c "celestia rpc state Balance --auth $$CELESTIA_NODE_AUTH_TOKEN"

# Useful if you want to run `apk update &&  apk add busybox-extras`
.PHONY: celestia_localnet_exec_root
celestia_localnet_exec_root: docker_check  ## Execu into the container as root user in the celestia localnet
	docker exec -it --user=root celestia /bin/sh

.PHONY: poktroll_start
poktroll_start: docker_check celestia_localnet_auth_token go_version_check  ## Start the poktroll node
	./build/init-local.sh

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
