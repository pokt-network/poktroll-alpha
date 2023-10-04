.SILENT:

POKTROLLD_HOME := ./localnet/poktrolld
POKTROLLD_NODE := tcp://127.0.0.1:36657
SESSION_HEIGHT ?= 1 # Default height when retrieving session data

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
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


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
	poktrolld --home=$(POKTROLLD_HOME) tx bank send $$KEY1 $$KEY2 42069stake

.PHONY: poktroll_balance
poktroll_balance: ## Check the balances of both keys
	KEY1=$$(make -s poktroll_list_keys | awk -F' ' '/address: pokt1/{print $$3}' | head -1); \
	KEY2=$$(make -s poktroll_list_keys | awk -F' ' '/address: pokt1/{print $$3}' | tail -1); \
	poktrolld --home=$(POKTROLLD_HOME) query bank balances $$KEY1 --node $(POKTROLLD_NODE); \
	poktrolld --home=$(POKTROLLD_HOME) query bank balances $$KEY2 --node $(POKTROLLD_NODE);

# Ref: https://rollkit.dev/tutorials/gm-world-frontend
.PHONY: poktroll_frontend_cosmology
poktroll_frontend_cosmology: ## Start the poktroll cosmology frontend
	echo "Visit http://localhost:3000/"
	yarn --cwd ./frontend dev

# Tutorial: https://blog.logrocket.com/create-react-native-app-using-ignite-boilerplate/
.PHONY: poktroll_frontend_react
poktroll_frontend_react: ## Start the poktroll react frontend
	cd ./react && npm install && npm run dev

.PHONY: servicers_get
servicers_get: ## Retrieves all servicers from the poktroll state
	poktrolld --home=$(POKTROLLD_HOME) q servicer list-servicers --node $(POKTROLLD_NODE)

.PHONY: servicer_stake
servicer_stake: ## Stake tokens for the servicer specified (must specify the SERVICER env var)
	poktrolld --home=$(POKTROLLD_HOME) tx servicer stake-servicer ./testutil/json/$(SERVICER).json --keyring-backend test --from $(SERVICER) --node $(POKTROLLD_NODE)

.PHONY: servicer1_stake
servicer1_stake: ## Stake for servicer1
	SERVICER=servicer1 make servicer_stake

.PHONY: servicer2_stake
servicer2_stake: ## Stake for servicer2
	SERVICER=servicer2 make servicer_stake

.PHONY: servicer3_stake
servicer3_stake: ## Stake for servicer3
	SERVICER=servicer3 make servicer_stake

.PHONY: servicer_unstake
servicer_unstake: ## Unstake tokens for the servicer specified
	poktrolld --home=$(POKTROLLD_HOME) tx servicer unstake-servicer --keyring-backend test --from $(SERVICER) --node $(POKTROLLD_NODE)

.PHONY: servicer1_unstake
servicer1_unstake: ## Unstake for servicer1
	SERVICER=servicer1 make servicer_unstake

.PHONY: servicer2_unstake
servicer2_unstake: ## Unstake for servicer2
	SERVICER=servicer2 make servicer_unstake

.PHONY: servicer3_unstake
servicer3_unstake: ## Unstake for servicer3
	SERVICER=servicer3 make servicer_unstake

.PHONY: apps_get
apps_get: ## Retrieves all applications from the poktroll state
	poktrolld --home=$(POKTROLLD_HOME) q application list-application --node $(POKTROLLD_NODE)

.PHONY: app_stake
app_stake: ## Stake tokens for the application specified (must specify the APP and SERVICES env vars)
	poktrolld --home=$(POKTROLLD_HOME) tx application stake-application 1000stake $(SERVICES) --keyring-backend test --from $(APP) --node $(POKTROLLD_NODE)

.PHONY: app1_stake
app1_stake: ## Stake for app1
	SERVICES=svc1,svc2 APP=app1 make app_stake

.PHONY: app2_stake
app2_stake: ## Stake for app2
	SERVICES=svc2,svc3 APP=app2 make app_stake

.PHONY: app3_stake
app3_stake: ## Stake for app3
	SERVICES=svc3,svc4 APP=app3 make app_stake

.PHONY: app_unstake
app_unstake: ## Unstake tokens for the application specified (must specify the APP env var)
	poktrolld --home=$(POKTROLLD_HOME) tx application unstake-application --keyring-backend test --from $(APP) --node $(POKTROLLD_NODE)

.PHONY: app1_unstake
app1_unstake: ## Unstake for app1
	APP=app1 make app_unstake

.PHONY: app2_unstake
app2_unstake: ## Unstake for app2
	APP=app2 make app_unstake

.PHONY: app3_unstake
app3_unstake: ## Unstake for app3
	APP=app3 make app_unstake

.PHONY: delegate
delegate: ## Delegate the application to the specified portal (must specify the APP and PORTAL env vars)
	poktrolld --home=$(POKTROLLD_HOME) tx application delegate-to-portal '$(PORTAL)' --keyring-backend test --from $(APP) --node $(POKTROLLD_NODE)

.PHONY: delegate_app1_portal1
delegate_app1_portal1: ## Delegate app1 to portal1
	APP=app1 PORTAL=pokt15vzxjqklzjtlz7lahe8z2dfe9nm5vxwwmscne4 make delegate

.PHONY: delegate_app2_portal2
delegate_app2_portal2: ## Delegate app2 to portal2
	APP=app2 PORTAL=pokt15w3fhfyc0lttv7r585e2ncpf6t2kl9uh8rsnyz make delegate

.PHONY: delegate_app3_portal3
delegate_app3_portal3: ## Delegate app3 to portal3
	APP=app3 PORTAL=pokt1zhmkkd0rh788mc9prfq0m2h88t9ge0j83gnxya make delegate

.PHONY: undelegate
undelegate: ## Undelegate the application to the specified portal (must specify the APP and PORTAL env vars)
	poktrolld --home=$(POKTROLLD_HOME) tx application undelegate-from-portal '$(PORTAL)' --keyring-backend test --from $(APP) --node $(POKTROLLD_NODE)

.PHONY: undelegate_app1_portal1
undelegate_app1_portal1: ## Undelegate app1 to portal1
	APP=app1 PORTAL=pokt15vzxjqklzjtlz7lahe8z2dfe9nm5vxwwmscne4 make undelegate

.PHONY: undelegate_app2_portal2
undelegate_app2_portal2: ## Delegate app2 to portal2
	APP=app2 PORTAL=pokt15w3fhfyc0lttv7r585e2ncpf6t2kl9uh8rsnyz make undelegate

.PHONY: undelegate_app3_portal3
undelegate_app3_portal3: ## Delegate app3 to portal3
	APP=app3 PORTAL=pokt1zhmkkd0rh788mc9prfq0m2h88t9ge0j83gnxya make undelegate

.PHONY: portals_get
portals_get: ## Retrieves all portals from the poktroll state
	poktrolld --home=$(POKTROLLD_HOME) q portal list-portals --node $(POKTROLLD_NODE)

.PHONY: portal_stake
portal_stake: ## Stake tokens for the portal specified (must specify the PORTAL and SERVICES env vars)
	poktrolld --home=$(POKTROLLD_HOME) tx portal stake-portal 1000stake $(SERVICES) --keyring-backend test --from $(PORTAL) --node $(POKTROLLD_NODE)

.PHONY: portal1_stake
portal1_stake: ## Stake for portal1
	SERVICES=svc1,svc2 PORTAL=portal1 make portal_stake

.PHONY: portal2_stake
portal2_stake: ## Stake for portal2
	SERVICES=svc2,svc3 PORTAL=portal2 make portal_stake

.PHONY: portal3_stake
portal3_stake: ## Stake for portal3
	SERVICES=svc3,svc4 PORTAL=portal3 make portal_stake

.PHONY: portal_unstake
portal_unstake: ## Unstake tokens for the portal specified (must specify the PORTAL env var)
	poktrolld --home=$(POKTROLLD_HOME) tx portal unstake-portal --keyring-backend test --from $(PORTAL) --node $(POKTROLLD_NODE)

.PHONY: portal1_unstake
portal1_unstake: ## Unstake for portal1
	PORTAL=portal1 make portal_unstake

.PHONY: portal2_unstake
portal2_unstake: ## Unstake for portal2
	PORTAL=portal2 make portal_unstake

.PHONY: portal3_unstake
portal3_unstake: ## Unstake for portal3
	PORTAL=portal3 make portal_unstake

.PHONY: portal_allowlist
portal_allowlist: ## allowlist the application for the portal specified (must specify the PORTAL and APP env vars)
	poktrolld --home=$(POKTROLLD_HOME) tx portal allowlist-application '$(APP)' --keyring-backend test --from $(PORTAL) --node $(POKTROLLD_NODE)

.PHONY: portal1_allowlist_app1
portal1_allowlist_app1: ## Allowlist app1 for portal1
	PORTAL=portal1 APP=pokt1mrqt5f7qh8uxs27cjm9t7v9e74a9vvdnq5jva4 make portal_allowlist

.PHONY: portal2_allowlist_app2
portal2_allowlist_app2: ## Allowlist app2 for portal2
	PORTAL=portal2 APP=pokt184zvylazwu4queyzpl0gyz9yf5yxm2kdhh9hpm make portal_allowlist

.PHONY: portal3_allowlist_app3
portal3_allowlist_app3: ## Allowlist app3 for portal3
	PORTAL=portal3 APP=pokt1lqyu4v88vp8tzc86eaqr4lq8rwhssyn6rfwzex make portal_allowlist

.PHONY: portal_unallowlist
portal_unallowlist: ## Unallowlist the application for the portal specified (must specify the PORTAL and APP env vars)
	poktrolld --home=$(POKTROLLD_HOME) tx portal unallowlist-application '$(APP)' --keyring-backend test --from $(PORTAL) --node $(POKTROLLD_NODE)

.PHONY: portal1_unallowlist_app1
portal1_unallowlist_app1: ## Unallowlist app1 for portal1
	PORTAL=portal1 APP=pokt1mrqt5f7qh8uxs27cjm9t7v9e74a9vvdnq5jva4 make portal_unallowlist

.PHONY: portal2_unallowlist_app2
portal2_unallowlist_app2: ## Unallowlist app2 for portal2
	PORTAL=portal2 APP=pokt184zvylazwu4queyzpl0gyz9yf5yxm2kdhh9hpm make portal_unallowlist

.PHONY: portal3_unallowlist_app3
portal3_unallowlist_app3: ## Unallowlist app3 for portal3
	PORTAL=portal3 APP=pokt1lqyu4v88vp8tzc86eaqr4lq8rwhssyn6rfwzex make portal_unallowlist:pokt1lqyu4v88vp8tzc86eaqr4lq8rwhssyn6rfwzex

.PHONY: app_delegatees
app_delegatees: ## Retrieves all delegatees for the application specified (must specify the APP env var)
	poktrolld --home=$(POKTROLLD_HOME) query portal get-delegated-portals '$(APP)' --node $(POKTROLLD_NODE)

.PHONY: app1_delegatees
app1_delegatees: ## Retrieves all delegatees for app1
	APP=pokt1mrqt5f7qh8uxs27cjm9t7v9e74a9vvdnq5jva4 make app_delegatees

.PHONY: app2_delegatees
app2_delegatees: ## Retrieves all delegatees for app2
	APP=pokt184zvylazwu4queyzpl0gyz9yf5yxm2kdhh9hpm make app_delegatees

.PHONY: app3_delegatees
app3_delegatees: ## Retrieves all delegatees for app3
	APP=pokt1lqyu4v88vp8tzc86eaqr4lq8rwhssyn6rfwzex make app_delegatees

.PHONY: test_unit_all
test_unit_all: ## Run all unit tests
	go test -v ./...

.PHONY: session_get
session_get: ## Queries the poktroll node for session data
	poktrolld --home=$(POKTROLLD_HOME) query session get-session $(APP) $(SVC) $(HEIGHT) --node $(POKTROLLD_NODE)

.PHONY: session_get_app1_svc1
session_get_app1_svc1: ## Getting the session for app1 and svc1 and height1
	APP=pokt1mrqt5f7qh8uxs27cjm9t7v9e74a9vvdnq5jva4 SVC=svc1 HEIGHT=$(SESSION_HEIGHT) make session_get

.PHONY: session_get_app2_svc2
session_get_app2_svc2: ## Getting the session for app2 and svc2 and height1
	APP=pokt184zvylazwu4queyzpl0gyz9yf5yxm2kdhh9hpm SVC=svc2 HEIGHT=$(SESSION_HEIGHT) make session_get

.PHONY: session_get_app3_svc3
session_get_app3_svc3: ## Getting the session for app3 and svc3 and height1
	APP=pokt1lqyu4v88vp8tzc86eaqr4lq8rwhssyn6rfwzex SVC=svc3 HEIGHT=$(SESSION_HEIGHT) make session_get

.PHONY: relayer_start
relayer_start: ## Start the relayer
	poktrolld relayer \
	--node $(POKTROLLD_NODE) \
	--signing-key servicer1 \
	--keyring-backend test

.PHONY: claims_query
claims_query: ## Query the poktroll node for claims data
	SERVICER_ADDR=$(shell poktrolld keys show servicer1 -a --keyring-backend test); \
	poktrolld query servicer claims $$SERVICER_ADDR

.PHONY: anvil_start
anvil_start: ## Start the anvil
	anvil -p 8547 -b 5

.PHONY: cast_relay
cast_relay: ## Cast a relay
	cast block

.PHONY: ws_subscribe
ws_subscribe: ## Subscribe to the websocket for new blocks
	echo "Copy paste the following: {"id":1,"jsonrpc":"2.0","method":"eth_subscribe","params":["newHeads"]}"
	wscat --connect ws://localhost:8546

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


### Inspired by @goldinguy_ in this post: https://goldin.io/blog/stop-using-todo ###
# TODO          - General Purpose catch-all.
# DECIDE        - A TODO indicating we need to make a decision and document it using an ADR in the future; https://github.com/pokt-network/pocket-network-protocol/tree/main/ADRs
# TECHDEBT      - Not a great implementation, but we need to fix it later.
# IMPROVE       - A nice to have, but not a priority. It's okay if we never get to this.
# OPTIMIZE      - An opportunity for performance improvement if/when it's necessary
# DISCUSS       - Probably requires a lengthy offline discussion to understand next steps.
# INCOMPLETE    - A change which was out of scope of a specific PR but needed to be documented.
# INVESTIGATE   - TBD what was going on, but needed to continue moving and not get distracted.
# CLEANUP       - Like TECHDEBT, but not as bad.  It's okay if we never get to this.
# HACK          - Like TECHDEBT, but much worse. This needs to be prioritized
# REFACTOR      - Similar to TECHDEBT, but will require a substantial rewrite and change across the codebase
# CONSIDERATION - A comment that involves extra work but was thoughts / considered as part of some implementation
# CONSOLIDATE   - We likely have similar implementations/types of the same thing, and we should consolidate them.
# ADDTEST       - Add more tests for a specific code section
# DEPRECATE     - Code that should be removed in the future
# RESEARCH      - A non-trivial action item that requires deep research and investigation being next steps can be taken
# DOCUMENT		- A comment that involves the creation of a README or other documentation
# BUG           - There is a known existing bug in this code
# NB            - An important note to reference later
# DISCUSS_IN_THIS_COMMIT - SHOULD NEVER BE COMMITTED TO MASTER. It is a way for the reviewer of a PR to start / reply to a discussion.
# TODO_IN_THIS_COMMIT    - SHOULD NEVER BE COMMITTED TO MASTER. It is a way to start the review process while non-critical changes are still in progress
TODO_KEYWORDS = -e "TODO" -e "DECIDE" -e "TECHDEBT" -e "IMPROVE" -e "OPTIMIZE" -e "DISCUSS" -e "INCOMPLETE" -e "INVESTIGATE" -e "CLEANUP" -e "HACK" -e "REFACTOR" -e "CONSIDERATION" -e "TODO_IN_THIS_COMMIT" -e "DISCUSS_IN_THIS_COMMIT" -e "CONSOLIDATE" -e "DEPRECATE" -e "ADDTEST" -e "RESEARCH" -e "BUG"

# How do I use TODOs?
# 1. <KEYWORD>: <Description of follow up work>;
# 	e.g. HACK: This is a hack, we need to fix it later
# 2. If there's a specific issue, or specific person, add that in paranthesiss
#   e.g. TODO(@Olshansk): Automatically link to the Github user https://github.com/olshansk
#   e.g. INVESTIGATE(#420): Automatically link this to github issue https://github.com/pokt-network/pocket/issues/420
#   e.g. DISCUSS(@Olshansk, #420): Specific individual should tend to the action item in the specific ticket
#   e.g. CLEANUP(core): This is not tied to an issue, or a person, but should only be done by the core team.
#   e.g. CLEANUP: This is not tied to an issue, or a person, and can be done by the core team or external contributors.
# 3. Feel free to add additional keywords to the list above

.PHONY: todo_list
todo_list: ## List all the TODOs in the project (excludes vendor and prototype directories)
	grep --exclude-dir={.git,vendor,prototype} -r ${TODO_KEYWORDS}  .


TODO_SEARCH ?= $(shell pwd)

.PHONY: todo_search
todo_search: ## List all the TODOs in a specific directory specific by `TODO_SEARCH`
	grep --exclude-dir={.git,vendor,prototype} -r ${TODO_KEYWORDS} ${TODO_SEARCH}

.PHONY: todo_count
todo_count: ## Print a count of all the TODOs in the project
	grep --exclude-dir={.git,vendor,prototype} -r ${TODO_KEYWORDS} . | wc -l

.PHONY: todo_this_commit
todo_this_commit: ## List all the TODOs needed to be done in this commit
	grep --exclude-dir={.git,vendor,prototype,.vscode} --exclude=Makefile -r -e "TODO_IN_THIS_COMMIT" -e "DISCUSS_IN_THIS_COMMIT"

.PHONY: mockgen
mockgen: ## Use `mockgen` to generate mocks used for testing purposes of all the modules.
	go generate ./x/application/keeper
	go generate ./x/servicer/keeper
	go generate ./x/session/keeper

.PHONY: ignite_regenerate
ignite_regenerate: ## Regenerate the ignite boilerplate
	ignite generate proto-go --yes && ignite generate openapi --yes

# Create new accounts with:
#  - ignite account create {KEY_NAME} --keyring-dir ./localnet/poktrolld --keyring-backend test
#  - poktrolld --home=${POKTROLLD_HOME} --node ${POKTROLLD_NODE} --keyring-backend test add-genesis-account {KEY_NAME} 1000000000000000pokt
.PHONY: ignite_acc_list
ignite_acc_list: ## List all the accounts in the ignite boilerplate
	ignite account list  --keyring-dir $(POKTROLLD_HOME) --keyring-backend test

.PHONY: localnet_regenesis
localnet_regenesis:
	# NOTE: intentionally not using --home <dir> flag to avoid overwriting the test keyring
	ignite chain init --skip-proto
	rm -rf $(POKTROLLD_HOME)/keyring-test
	cp -r ${HOME}/.poktroll/keyring-test $(POKTROLLD_HOME)
	cp ${HOME}/.poktroll/config/*_key.json $(POKTROLLD_HOME)/config/
	cp ${HOME}/.poktroll/config/genesis.json ./localnet/