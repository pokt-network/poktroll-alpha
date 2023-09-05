package slip

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	// Test Vectors
	testVector1SeedHex = "000102030405060708090a0b0c0d0e0f"
	testVector2SeedHex = "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542"

	// SLIPS-0010
	testSeedHex = "045e8380086abc6f6e941d6fe47ca93b86723bc246ec8c4beee411b410028675"
)

func TestSlip_DeriveChild_TestVectors(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		seed        string
		wantPrivHex string
		wantPubHex  string

		wantErr bool
	}{
		// https://github.com/satoshilabs/slips/blob/master/slip-0010.md#test-vector-1-for-ed25519
		// Note that ed25519 public keys normally don't have a "00" prefix, but we are reflecting the
		// test vectors from the spec which do
		{
			name:        "TestVector1 Key derivation is deterministic for path `m` (master key)",
			path:        "m",
			seed:        testVector1SeedHex,
			wantPrivHex: "9662541b93982b55b7ef6d8b43dfbe84e3b88c96bb76be08f90b0c099aa9ceb9",
			wantPubHex:  "00023c3abe95d897e14fe6bb9cb6e16829ed618061a8655914a4a6eb95d70b654fb2",
			wantErr:     false,
		},
		{
			name:        "TestVector1 Key derivation is deterministic for path `m/0'`",
			path:        "m/0'",
			seed:        testVector1SeedHex,
			wantPrivHex: "7d77c00fdc06ded7e7b6037c143a946bc4b9cd9b37ca6d62c19d14da752d5d88",
			wantPubHex:  "0003ddbf673db5cdb766f9a319aa92d12a1b7f01cee38bec06e21ada7cb482fa10e8",
			wantErr:     false,
		},
		{
			name:        "TestVector1 Key derivation is deterministic for path `m/0'/1'`",
			path:        "m/0'/1'",
			seed:        testVector1SeedHex,
			wantPrivHex: "1fe7aede8150fa459f610fe0629ed6f2869cfdd210b44a0cf2c6785a57be9c85",
			wantPubHex:  "0003d2e9a7d5f60c7ac4c50dbd6672bd9643e7704fdb80a71f978fb3d4e5f0ae7f7a",
			wantErr:     false,
		},
		{
			name:        "TestVector1 Key derivation is deterministic for path `m/0'/1'/2'`",
			path:        "m/0'/1'/2'",
			seed:        testVector1SeedHex,
			wantPrivHex: "fabe59803d526eb9269951b48fa0bdcd453db0afe8f9ee6f97a15836855b51ea",
			wantPubHex:  "000349ed160aebd46d394352c1665cbaff4300f0e7bb8ad508965a8aab3e74f76fac",
			wantErr:     false,
		},
		{
			name:        "TestVector1 Key derivation is deterministic for path `m/0'/1'/2'/2'`",
			path:        "m/0'/1'/2'/2'",
			seed:        testVector1SeedHex,
			wantPrivHex: "3a7a030ff32ea2ef472d80a463b7cfe5f252a7a553cb91c10f9ba3fd72c8beea",
			wantPubHex:  "00031a71863492e0df1324eb3bb3dc91168ba0712866a0c3ec8768a4df4bc5d805da",
			wantErr:     false,
		},
		{
			name:        "TestVector1 Key derivation is deterministic for path `m/0'/1'/2'/2'/1000000000'`",
			path:        "m/0'/1'/2'/2'/1000000000'",
			seed:        testVector1SeedHex,
			wantPrivHex: "9c97be512c86210fa544185b5adb0bff83c01b7b044c1ab6bf3f6da86ef17276",
			wantPubHex:  "000249c416e004d6246c44284df8d21dd7bbc714086167254bbd96204cf67b536a29",
			wantErr:     false,
		},
		{
			name:        "TestVector1 Key derivation fails with invalid path `m/0`",
			path:        "m/0",
			seed:        testVector1SeedHex,
			wantPrivHex: "",
			wantErr:     true,
		},
		// https://github.com/satoshilabs/slips/blob/master/slip-0010.md#test-vector-2-for-ed25519
		{
			name:        "TestVector2 Key derivation is deterministic for path `m` (master key)",
			path:        "m",
			seed:        testVector2SeedHex,
			wantPrivHex: "f955bff277960699b72a6202ba6905b2b7c4ca2cb39640e848faafa3d20be45c",
			wantPubHex:  "00035ac08e72720c4f2c9809412b042fffdac7827dd2d7e8524b26d18f1507da21c1",
			wantErr:     false,
		},
		{
			name:        "TestVector2 Key derivation is deterministic for path `m/0'`",
			path:        "m/0'",
			seed:        testVector2SeedHex,
			wantPrivHex: "43f588e813b1aa36ff599f957d2ae8facc7b55c4f97a282c421020e4f3ce936d",
			wantPubHex:  "0002be84c291e76c23a5632e1cbafe5d646d886236f6ae8d0d4b76a84e153db8b3c5",
			wantErr:     false,
		},
		{
			name:        "TestVector2 Key derivation is deterministic for path `m/0'/2147483647'`",
			path:        "m/0'/2147483647'",
			seed:        testVector2SeedHex,
			wantPrivHex: "d34e0ff88f662b0811d049a99942fa7f50800c32d3535b318452ae981db0019b",
			wantPubHex:  "00039f4ccadf338414b1077e9772eb72c74b3c4d066010d5d2224051537d1b78f3dd",
			wantErr:     false,
		},
		{
			name:        "TestVector2 Key derivation is deterministic for path `m/0'/2147483647'/1'`",
			path:        "m/0'/2147483647'/1'",
			seed:        testVector2SeedHex,
			wantPrivHex: "bc7efd11b95e9102376fc394823c86f2c0e71cf887c1774f117e5d7c59e900b1",
			wantPubHex:  "00024212408b9255935382d26399c057923721a3434f9dfc96f1943c3f58a9449222",
			wantErr:     false,
		},
		{
			name:        "TestVector2 Key derivation is deterministic for path `m/0'/2147483647'/1'/2147483646'`",
			path:        "m/0'/2147483647'/1'/2147483646'",
			seed:        testVector2SeedHex,
			wantPrivHex: "047e0dac03a680cc0f9033d22029546d648cd83fe910ee70e39f3e33a681ce52",
			wantPubHex:  "0002edad0c57640eab51cfe358b51933ce2d71ee5a09e39aa5fa0cfe05005496e010",
			wantErr:     false,
		},
		{
			name:        "TestVector2 Key derivation is deterministic for path `m/0'/2147483647'/1'/2147483646'/2'`",
			path:        "m/0'/2147483647'/1'/2147483646'/2'",
			seed:        testVector2SeedHex,
			wantPrivHex: "674601c26bf4a4ec475e047b9b5a68b339145c2954509b23c60c69c45ba117de",
			wantPubHex:  "000310cd8b01610c7de0beea48ae7163b90a3a6243f35a191e7ef2d65e36c4547917",
			wantErr:     false,
		},
		// Pocket specific test vectors
		{
			name:        "PoktTestVector Key derivation is deterministic for path `m` (master key)",
			path:        "m",
			seed:        testSeedHex,
			wantPrivHex: "875cfa863856778dd9e86e3d9d9db08774fd7d0c0643f795bf5d3d3d7851d57a",
			wantPubHex:  "000209d3d7fcd6c35b9c6dba8e6ffa9380934dbaa9f87659b174fca4d7b815977103",
			wantErr:     false,
		},
		{
			name:        "PoktTestVector Key derivation is deterministic for path `m/44'`",
			path:        "m/44'",
			seed:        testSeedHex,
			wantPrivHex: "c664c83b87ec130707142b73b3a7732b10c7ef73543e0fa77530ebae31403739",
			wantPubHex:  "0002b0b0c348af7a7800f99251f8eeb3b8bc81bbad2aadfc6122f9ddf588c19a416b",
			wantErr:     false,
		},
		{
			name:        "PoktTestVector Key derivation is deterministic for path `m/44'/635'`",
			path:        "m/44'/635'",
			seed:        testSeedHex,
			wantPrivHex: "0158f66384fbc009e4e3a3b667b002593c0771149b018562925b1518ac0ae796",
			wantPubHex:  "000237ebb9ad901619c9a4f28361fc2b7a2e94880e0ac8070cf7117ad8ac5d13ca7d",
			wantErr:     false,
		},
		{
			name:        "PoktTestVector Child key derivation is deterministic for index `0` (first child)",
			path:        fmt.Sprintf(PoktAccountPathFormat, 0),
			seed:        testSeedHex,
			wantPrivHex: "8d681f1af7e6da1c12c2b0eb70ae109cf38f6a0a3c37eccc16616641aba4a8e4",
			wantPubHex:  "00024ba2ce7db5a16c9583d2c2e8a5cff0bf0cb230f2b7b26eeee4fc9f4ec79e2ea4",
			wantErr:     false,
		},
		{
			name:        "PoktTestVector Child key derivation is deterministic for index `1000000`",
			path:        fmt.Sprintf(PoktAccountPathFormat, 1000000),
			seed:        testSeedHex,
			wantPrivHex: "c1dce7ac6a08191ecdf2a9b56e26963743b75a3fd2f168443509596b8ae4951d",
			wantPubHex:  "00031f926f16d99a3aa011d4e6f3052cf911de5a3d58f04827987262fb1a39989be7",
			wantErr:     false,
		},
		{
			name:        "PoktTestVector Child key derivation is deterministic for index `2147483647` (last child)",
			path:        fmt.Sprintf(PoktAccountPathFormat, 2147483647),
			seed:        testSeedHex,
			wantPrivHex: "653eacdb08d7c35103211bab0613a66511e1da3fa9500650d055564bf9099717",
			wantPubHex:  "000260ee3238a95b9bb5d858892351d67c2e21c1ce45589a1a9ad64a37028e2fdce6",
			wantErr:     false,
		},
		{
			name:        "PoktTestVector Child index is too large to derive ed25519 key for index `2147483648` ",
			path:        fmt.Sprintf(PoktAccountPathFormat, 2147483648),
			seed:        testSeedHex,
			wantPrivHex: "",
			wantPubHex:  "",
			wantErr:     true,
		},
		{
			name:        "PoktTestVector Child index is too large to derive ed25519 key for index `4294967295` ",
			path:        fmt.Sprintf(PoktAccountPathFormat, ^uint32(0)),
			seed:        testSeedHex,
			wantPrivHex: "",
			wantPubHex:  "",
			wantErr:     true,
		},
	}
	for _, tv := range tests {
		t.Run(tv.name, func(t *testing.T) {
			seed, err := hex.DecodeString(tv.seed)
			require.NoError(t, err)
			childKey, err := DeriveChild(tv.path, seed)
			if tv.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if err != nil {
				return
			}

			// Slip-0010 private keys in test vector are only the seed of the full private key
			// This is equivalent to the SecretKey of the HMAC key used to generate the ed25519 key
			privSeed, err := childKey.GetSeed("")
			require.NoError(t, err)
			privHex := hex.EncodeToString(privSeed)
			require.Equal(t, tv.wantPrivHex, privHex)

			// Slip-0010 keys are prefixed with "00" in the test vectors
			pubHex := childKey.GetPublicKey().String()
			require.Equal(t, tv.wantPubHex, "00"+pubHex)
		})
	}
}
