package e2e

import (
	"fmt"
	"testing"

	"github.com/regen-network/gocuke"
	"github.com/stretchr/testify/require"

	cryptoPocket "poktroll/shared/crypto"
)

func TestUtility(t *testing.T) {
	// a new step definition utilitySuite is constructed for every scenario
	gocuke.NewRunner(t, &utilitySuite{}).Path("features/utility/*.feature").Run()
}

const (
	// 001 servicer is in session 0 for applicatio 000
	//	The list of servicers in the session is decided by the 'servicers' section of the genesis, from 'build/localnet/manifest/configs.yaml' file
	servicerA      = "001"
	appA           = "000"
	serviceA       = "0001"
	timeoutService = "9999"

	relaychainEth = "RelayChainETH" // used to refer to Ethereum chain when retrieving relaychain settings
)

type utilitySuite struct {
	// special arguments like TestingT are injected automatically into exported fields
	gocuke.TestingT

	// validator holds command results between runs and reports errors to the test utilitySuite
	// TECHDEBT: Rename `validator` to something more appropriate
	validator *validatorPod
	// validatorA maps to suffix ID 001 of the kube pod that we use as our control agent

	// servicerKeys is hydrated by the clientset with credentials for all servicers.
	// servicerKeys maps servicer IDs to their private key as a hex string.
	servicerKeys map[string]string

	// appKeys is hydrated by the clientset with credentials for all apps.
	// appKeys maps app IDs to their private key as a hex string.
	appKeys map[string]string

	// relaychains holds settings for all relaychains used in the tests
	//	the map key is a constant selected as the identifier for the relaychain, e.g. "RelayChainETH" represented as "0001" in other parts of the codebase for Ethereum
	relaychains map[string]*relaychainSettings

	// servicer holds the key for the servicer that should received the relay
	servicerKey string
}

// relaychainSettings holds the settings for a specific relaychain
type relaychainSettings struct {
	account string
	height  string
}

func (s *utilitySuite) TheRelayResponseHasValidId() {
	require.Contains(s, s.validator.result.Stdout, `"id":1`)
}

func (s *utilitySuite) TheValidatorShouldHaveExitedWithoutError() {
	require.NoError(s, s.validator.result.Err)
}

// Then the request times out without a response
func (s *utilitySuite) TheRequestTimesOutWithoutAResponse() {
	require.Contains(s, s.validator.result.Stdout, "HTTP status code: 500")
}

// TheApplicationHasAValidEthereumRelaychainAccount fullfils the following condition from feature file:
//
//	"Given the application has a valid ethereum relaychain account"
func (s *utilitySuite) TheApplicationHasAValidEthereumRelaychainAccount() {
	// Account: 0x8315177aB297bA92A06054cE80a67Ed4DBd7ed3a   (Arbitrum Bridge)
	s.relaychains[relaychainEth].account = "0x8315177aB297bA92A06054cE80a67Ed4DBd7ed3a"
}

func (s *utilitySuite) TheRelayResponseIsValidJsonRpc() {
	require.Contains(s, s.validator.result.Stdout, `"jsonrpc":"2.0"`)
}

// An Application requests the account balance of a specific address at a specific height
func (s *utilitySuite) TheApplicationSendsAGetBalanceRelayAtASpecificHeightToAnEthereumServicer() {
	// ADD_IN_THIS_PR: Add a servicer staked for the Ethereum RelayChain
	params := fmt.Sprintf("%q: [%q, %q]", "params", s.relaychains[relaychainEth].account, s.relaychains[relaychainEth].height)
	checkBalanceRelay := fmt.Sprintf("{%s, %s}", `"method": "eth_getBalance", "id": "1", "jsonrpc": "2.0"`, params)

	servicerPrivateKey := s.getServicerPrivateKey(s.servicerKey)
	appPrivateKey := s.getAppPrivateKey(appA)

	s.sendTrustlessRelay(checkBalanceRelay, servicerPrivateKey.Address().String(), appPrivateKey.Address().String(), serviceA, true)
}
func (s *utilitySuite) TheRelayResponseContains(relayResponse string) {
	require.Contains(s, s.validator.result.Stdout, relayResponse)
}

// An Application requests the account balance of a specific address at a specific height on "ServiceWithTimeout", i.e. timing out, service
func (s *utilitySuite) TheApplicationSendsAGetBalanceRelayAtASpecificHeightToTheServicewithtimeoutService() {
	params := fmt.Sprintf("%q: [%q, %q]", "params", s.relaychains[relaychainEth].account, s.relaychains[relaychainEth].height)
	checkBalanceRelay := fmt.Sprintf("{%s, %s}", `"method": "eth_getBalance", "id": "1", "jsonrpc": "2.0"`, params)

	servicerPrivateKey := s.getServicerPrivateKey(s.servicerKey)
	appPrivateKey := s.getAppPrivateKey(appA)

	s.sendTrustlessRelay(checkBalanceRelay, servicerPrivateKey.Address().String(), appPrivateKey.Address().String(), timeoutService, false)
}

// TheApplicationHasAValidEthereumRelaychaindHeight fullfils the following condition from feature file:
//
//	"Given the application has a valid ethereum relaychain height"
func (s *utilitySuite) TheApplicationHasAValidEthereumRelaychainHeight() {
	// Ethereum relaychain BlockNumber: 17605670 = 0x10CA426
	s.relaychains[relaychainEth].height = "0x10CA426"
}

// TheApplicationHasAValidServicer fullfils the following condition from feature file:
//
//	"Given the application has a valid servicer"
func (s *utilitySuite) TheApplicationHasAValidServicer() {
	s.servicerKey = servicerA
}

// getServicerPrivateKey generates a new keypair from the servicer private hex key that we get from the clientset
func (s *utilitySuite) getServicerPrivateKey(
	servicerId string,
) cryptoPocket.PrivateKey {
	privHexString := s.servicerKeys[servicerId]
	privateKey, err := cryptoPocket.NewPrivateKey(privHexString)
	require.NoErrorf(s, err, "failed to extract privkey for servicer with id %s", servicerId)

	return privateKey
}

// getAppPrivateKey generates a new keypair from the application private hex key that we get from the clientset
func (s *utilitySuite) getAppPrivateKey(
	appId string,
) cryptoPocket.PrivateKey {
	privHexString := s.appKeys[appId]
	privateKey, err := cryptoPocket.NewPrivateKey(privHexString)
	require.NoErrorf(s, err, "failed to extract privkey for app with id %s", appId)

	return privateKey
}

func (s *utilitySuite) sendTrustlessRelay(relayPayload string, servicerAddr, appAddr, serviceId string, shouldSucceed bool) {
	//args := []string{
	//	"Servicer",
	//	"Relay",
	//	appAddr,
	//	servicerAddr,
	//	// IMPROVE: add ETH_Goerli as a chain/service to genesis
	//	serviceId,
	//	relayPayload,
	//}

	// TODO:
	// Use ServicerClient to send the relay

	//// TECHDEBT: run the command from a client, i.e. not a validator, pod.
	//res, err := s.validator.RunCommand(args...)

	if shouldSucceed {
		require.NoError(s, err)
	}

	//s.validator.result = res
}
