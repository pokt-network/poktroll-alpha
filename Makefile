.SILENT:

POKTROLLD_HOME := ./localnet/.poktrolld
POKTROLLD_NODE := tcp://127.0.0.1:36657

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
	if [ "$$MAJOR_VERSION" -ne 1 ] || [ "$$MINOR_VERSION" -ge 21 ] ||  [ "$$MINOR_VERSION" -le 18 ] ; then \
		echo "Invalid Go version. Expected 1.19.x or 1.20.x but found $$GO_VERSION"; \
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

.PHONY: celestia_light_client_auth_token
celestia_light_client_auth_token: ## Get the auth token for the celestia light client on arabica testnet
	AUTH_TOKEN=$$(celestia light auth admin --p2p.network arabica); \
	echo $$AUTH_TOKEN

# Reference: https://github.com/rollkit/local-celestia-devnet
.PHONY: celestia_localnet
celestia_localnet: docker_check  ## Run a celestia localnet
	docker run --name celestia --platform linux/amd64 -p 26657:26657 -p 26658:26658 -p 26659:26659 ghcr.io/rollkit/local-celestia-devnet:v0.11.0-rc8

.PHONY: celestia_localnet_stop
celestia_localnet_stop: docker_check  ## Stop the celestia localnet
	docker stop celestia

# Intended to be called like so: `export AUTH_TOKEN=$(make celestia_localnet_auth_token)`
.PHONY: celestia_localnet_auth_token
celestia_localnet_auth_token: docker_check  ## Get the auth token for the celestia localnet
	CONTAINER_ID=$$(docker ps -qf "name=celestia"); \
	AUTH_TOKEN=$$(docker exec $$CONTAINER_ID celestia bridge --node.store /bridge auth admin); \
	echo $$AUTH_TOKEN

.PHONY: celestia_localnet_balance_check
celestia_localnet_balance_check: docker_check  ## Check the balance of an account in the celestia localnet
	AUTH_TOKEN=$$(make -s celestia_localnet_auth_token); \
	CONTAINER_ID=$$(docker ps -qf "name=celestia"); \
	docker exec $$CONTAINER_ID /bin/sh -c "celestia rpc state Balance --auth $$AUTH_TOKEN"

.PHONY: celestia_light_client_start
celestia_light_client_start: docker_check  ## Start the celestia light client
	echo "See the following link if there's an error https://docs.celestia.org/nodes/light-node/#install-celestia-node"
	celestia light start --core.ip consensus-validator-arabica-9.celestia-arabica.com --p2p.network arabica

.PHONY: celestia_testnet_balance_check
celestia_testnet_balance_check: ## Check the balance of the light client account on celestia arabica testnet
	AUTH_TOKEN=$$(make -s celestia_light_client_auth_token); \
	celestia rpc state Balance --auth $$AUTH_TOKEN

# Useful if you want to run `apk update &&  apk add busybox-extras`
.PHONY: celestia_localnet_exec_root
celestia_localnet_exec_root: docker_check  ## Execu into the container as root user in the celestia localnet
	docker exec -it --user=root celestia /bin/sh

.PHONY: poktroll_local_start
poktroll_local_start: docker_check go_version_check ## Start the localnet poktroll node
	@AUTH_TOKEN=$$(make -s celestia_localnet_auth_token) ./build/init-local.sh

.PHONY: poktroll_testnet_start
poktroll_testnet_start: docker_check go_version_check ## Start the testnet poktroll node
	@AUTH_TOKEN=$$(make -s celestia_localnet_auth_token) ./build/init-testnet.sh

.PHONY: poktroll_clear
poktroll_clear: ## Clear the poktroll state
	rm -rf ${HOME}/.poktroll
	rm ${HOME}/go/bin/poktrolld

.PHONY: poktroll_list_keys
poktroll_list_keys: ## List the poktroll keys
	poktrolld --home=$(POKTROLLD_HOME) keys list --keyring-backend test

.PHONY: poktroll_send
poktroll_send: ## Send tokens from one key to another
	KEY1=$$(make -s poktroll_list_keys | awk -F' ' '/address: pokt1/{print $$3}' | head -1); \
	KEY2=$$(make -s poktroll_list_keys | awk -F' ' '/address: pokt1/{print $$3}' | tail -1); \
	poktrolld --home=$(POKTROLLD_HOME) tx bank send $$KEY1 $$KEY2 42069stake --keyring-backend test --node $(POKTROLLD_NODE)

.PHONY: poktroll_balance
poktroll_balance: ## Check the balances of both keys
	KEY1=$$(make -s poktroll_list_keys | awk -F' ' '/address: pokt1/{print $$3}' | head -1); \
	KEY2=$$(make -s poktroll_list_keys | awk -F' ' '/address: pokt1/{print $$3}' | tail -1); \
	poktrolld --home=$(POKTROLLD_HOME) query bank balances $$KEY1 --node $(POKTROLLD_NODE); \
	poktrolld --home=$(POKTROLLD_HOME) query bank balances $$KEY2 --node $(POKTROLLD_NODE);

.PHONY: poktroll_get_session
poktroll_get_session: ## Queries the poktroll node for session data
	poktrolld --home=$(POKTROLLD_HOME) query poktroll get-session --node $(POKTROLLD_NODE)

# Ref: https://rollkit.dev/tutorials/gm-world-frontend
.PHONY: poktroll_cosmology_frontend
poktroll_cosmology_frontend: ## Start the poktroll cosmology frontend
	echo "Visit http://localhost:3000/"
	yarn --cwd ./frontend dev

.PHONY: poktroll_servicer_stake
poktroll_servicer_stake: ## Stake tokens for the servicer specified
	poktrolld --home=$(POKTROLLD_HOME) tx poktroll stake 1000stake servicer --keyring-backend test --from poktroll-key --node $(POKTROLLD_NODE)

.PHONY: poktroll_servicers_get
poktroll_servicers_get: ## Retrieves all servicers from the poktroll state
	poktrolld --home=$(POKTROLLD_HOME) q poktroll servicers --node $(POKTROLLD_NODE)

.PHONY: poktroll_servicer_unstake
poktroll_servicer_unstake: ## Unstake tokens for the servicer specified
	poktrolld --home=$(POKTROLLD_HOME) tx poktroll unstake 1000stake servicer --keyring-backend test --from poktroll-key --node $(POKTROLLD_NODE)

.PHONY: poktroll_app_stake
poktroll_app_stake: ## Stake tokens for the application specified
	poktrolld --home=$(POKTROLLD_HOME) tx poktroll stake 1000stake application --keyring-backend test --from poktroll-key --node $(POKTROLLD_NODE)

.PHONY: poktroll_apps_get
poktroll_apps_get: ## Retrieves all applications from the poktroll state
	poktrolld --home=$(POKTROLLD_HOME) q poktroll application --node $(POKTROLLD_NODE)

.PHONY: poktroll_app_unstake
poktroll_app_unstake: ## Unstake tokens for the application specified
	poktrolld --home=$(POKTROLLD_HOME) tx poktroll unstake 1000stake application --keyring-backend test --from poktroll-key --node $(POKTROLLD_NODE)

.PHONY: test_unit_all
test_unit_all: ## Run all unit tests
	go test -v ./...

.PHONY: localnet_up
localnet_up: ## Starts localnet
	tilt up

.PHONY: localnet_down
localnet_down: ## Delete resources created by localnet
	tilt down
	kubectl delete secret celestia-secret || exit 1

.PHONY: localnet_poktrolld_dlv_attach
localnet_poktrolld_dlv_attach: ## Attaches dlv to the process and listens on `40004` for you to connect with debug tool of a choice (gdlv/visual studio)
	kubectl exec deploy/poktrolld -- sh -c "dlv attach \$(pgrep poktrolld) --listen :40004 --headless --api-version=2 --accept-multiclient"
