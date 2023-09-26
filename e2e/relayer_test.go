package e2e

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testing"

	cometClient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/stretchr/testify/require"

	servicerTypes "poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

const (
	castGetBlockPayload = `{
		"id":1,
		"jsonrpc":"2.0",
		"method":"eth_getBlockByNumber",
		"params":["latest",false]
	}`
)

func TestRelayer_Relay(t *testing.T) {
	ctx := context.Background()
	anvilPort := 8547
	// TECHDEBT: this should be a config of some sort.
	//relayerProxyPort := 8545
	relayerProxyPort := anvilPort
	grpcPort := 36657
	grpcHost := fmt.Sprintf("localhost:%d", grpcPort)
	grpcAddress := fmt.Sprintf("tcp://%s", grpcHost)
	// TODO_THIS_COMMIT: load correct app address
	appAddress := "pokt1mrqt5f7qh8uxs27cjm9t7v9e74a9vvdnq5jva4"
	// TODO_THIS_COMMIT: load this from servicer1.json
	serviceId := "svc1"

	behavesLikeStakedRelayer(t, 1, 1, anvilPort)

	abciClient, err := cometClient.New(grpcAddress, "/websocket")
	require.NoError(t, err)

	// NB: nil height means latest latestBlockResponse.
	height := int64(100)
	latestBlockResponse, err := abciClient.Block(ctx, &height)
	require.NoError(t, err)

	//conn, err := grpc.Dial("localhost:36657", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	sessionClient := sessionTypes.NewQueryClient(conn)
	sessionResponse, err := sessionClient.GetSession(ctx, &sessionTypes.QueryGetSessionRequest{
		BlockHeight: uint64(latestBlockResponse.Block.Height),
		AppAddress:  appAddress,
		ServiceId:   serviceId,
	})
	require.NoError(t, err)

	relay := &servicerTypes.Relay{
		Req: &servicerTypes.RelayRequest{
			Headers: map[string]string{
				"content-type": "application/json",
				"accept":       "*/*",
				"Host":         fmt.Sprintf("localhost:%d", relayerProxyPort),
			},
			Method:             "POST",
			Url:                fmt.Sprintf("http://localhost:%d/", relayerProxyPort),
			Payload:            []byte(castGetBlockPayload),
			SessionId:          sessionResponse.Session.SessionId,
			ApplicationAddress: appAddress,
			//ApplicationSignature: nil,
		},
	}
	relay.Req.Headers["content-length"] = fmt.Sprintf("%d", len(relay.Req.Payload))

	requestRelay(t, relay)
	//cast(t, "block")
}

func requestRelay(t *testing.T, relay *servicerTypes.Relay) {
	httpRelayReq, err := http.NewRequest(
		relay.Req.Method,
		relay.Req.Url,
		io.NopCloser(bytes.NewBuffer(relay.Req.Payload)),
	)
	require.NoError(t, err)

	httpRelayRes, err := http.DefaultClient.Do(httpRelayReq)
	require.NoError(t, err)

	require.Equalf(t, 200, httpRelayRes.StatusCode, "HTTP status code not OK")

	httpRelayResBody, err := io.ReadAll(httpRelayRes.Body)
	require.NoError(t, err)

	t.Log(string(httpRelayResBody))
}

func behavesLikeStakedRelayer(t *testing.T, appIndex int, svcIndex int, anvilPort int) {
	// Stake application1
	// make app1_stake
	stakeApplication(t, appIndex)

	// Stake servicer1
	// make svc1_stake
	stakeServicer(t, svcIndex)

	// Start anvil
	// anvil -p 8546
	startAnvil(t, anvilPort)

	// Start relayer
	// TODO: configure upstream (relay-chain) port
	// poktrolld relayer start --keyring-backend test --node tcp://localhost:36657
	startRelayer(t)
}

func stakeApplication(t *testing.T, appIndex int) {
	stakeCmd := exec.Command(
		"make",
		fmt.Sprintf("app%d_stake", appIndex),
	)

	cmdWithoutError(t, stakeCmd)
	cmdWithRepoRootWorkingDir(t, stakeCmd)
	cmdWithConfirmation(t, stakeCmd, "[y/N]: ", "y")

	err := stakeCmd.Start()
	require.NoError(t, err)
}

func stakeServicer(t *testing.T, svcIndex int) {
	stakeCmd := exec.Command(
		"make",
		fmt.Sprintf("servicer%d_stake", svcIndex),
	)

	cmdWithoutError(t, stakeCmd)
	cmdWithRepoRootWorkingDir(t, stakeCmd)
	cmdWithConfirmation(t, stakeCmd, "[y/N]: ", "y")

	err := stakeCmd.Start()
	require.NoError(t, err)
}

func startAnvil(t *testing.T, port int) {
	anvilCmd := exec.Command(
		"anvil",
		strings.Split(
			fmt.Sprintf("-p %d", port), " ",
		)...,
	)

	cmdWithoutError(t, anvilCmd)

	// Kill process if still running when tests complete
	t.Cleanup(func() {
		if anvilCmd.ProcessState == nil {
			err := anvilCmd.Process.Kill()
			require.NoError(t, err)
		}
	})

	err := anvilCmd.Start()
	require.NoError(t, err)
}

// TODO_CONSIDERATION: it seems like we would mainly want to use cast to test
// the portal use case.
func cast(t *testing.T, args ...string) {
	castCmd := exec.Command(
		"cast",
		args...,
	)

	cmdWithoutError(t, castCmd)
	logStdoutStream(t, castCmd)

	err := castCmd.Run()
	require.NoError(t, err)
}

func startRelayer(t *testing.T) {
	relayerCmd := exec.Command(
		"poktrolld",
		"relayer", "start",
		"--signing-key", "servicer1",
		"--keyring-backend", "test",
		"--node", "tcp://localhost:36657",
	)

	cmdWithoutError(t, relayerCmd)
	cmdWithRepoRootWorkingDir(t, relayerCmd)

	// Kill process if still running when tests complete
	t.Cleanup(func() {
		if relayerCmd.ProcessState == nil {
			err := relayerCmd.Process.Kill()
			require.NoError(t, err)
		}
	})

	err := relayerCmd.Start()
	require.NoError(t, err)
}
func cmdWithConfirmation(t *testing.T, cmd *exec.Cmd, prompt string, response string) {
	cmdStdin, err := cmd.StdinPipe()
	require.NoError(t, err)

	cmdStdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	// Wait for confirmation prompt & confirm
	cmdStdoutScanner := bufio.NewScanner(cmdStdout)
	go func() {
		for cmdStdoutScanner.Scan() {
			line := cmdStdoutScanner.Text()
			if strings.HasSuffix(line, prompt) {
				_, err := cmdStdin.Write([]byte(response + "\n"))
				require.NoError(t, err)
			}
		}
	}()
}

func cmdWithoutError(t *testing.T, cmd *exec.Cmd) {
	cmdStderr, err := cmd.StderrPipe()
	require.NoError(t, err)

	// Read stderr and assert empty
	go func() {
		cmdStderrScanner := bufio.NewScanner(cmdStderr)
		for cmdStderrScanner.Scan() {
			line := cmdStderrScanner.Text()
			require.Empty(t, line)
		}
	}()
}

func cmdWithRepoRootWorkingDir(t *testing.T, cmd *exec.Cmd) {
	cmd.Dir = ".."
}

func logStdoutStream(t *testing.T, cmd *exec.Cmd) {
	cmdStdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	cmdStdoutScanner := bufio.NewScanner(cmdStdout)
	go func() {
		for cmdStdoutScanner.Scan() {
			line := cmdStdoutScanner.Text()
			t.Log(line)
		}
	}()
}
