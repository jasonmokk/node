package chains

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestChain_Validate(t *testing.T) {
	tests := []struct {
		name   string
		chain  Chain
		errStr string
	}{
		{
			name: "should pass if chain is valid",
			chain: Chain{
				ChainId:     42,
				ChainName:   ChainName_empty,
				Network:     Network_optimism,
				NetworkType: NetworkType_testnet,
				Vm:          Vm_evm,
				Consensus:   Consensus_op_stack,
				IsExternal:  true,
			},
		},
		{
			name: "should error if chain ID is zero",
			chain: Chain{
				ChainId:     0,
				ChainName:   ChainName_empty,
				Network:     Network_optimism,
				NetworkType: NetworkType_testnet,
				Vm:          Vm_evm,
				Consensus:   Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "chain ID must be positive",
		},
		{
			name: "should error if chain ID is negative",
			chain: Chain{
				ChainId:     0,
				ChainName:   ChainName_empty,
				Network:     Network_optimism,
				NetworkType: NetworkType_testnet,
				Vm:          Vm_evm,
				Consensus:   Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "chain ID must be positive",
		},
		{
			name: "should error if chain name invalid",
			chain: Chain{
				ChainId:     42,
				ChainName:   ChainName_base_sepolia + 1,
				Network:     Network_optimism,
				NetworkType: NetworkType_testnet,
				Vm:          Vm_evm,
				Consensus:   Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "invalid chain name",
		},
		{
			name: "should error if network invalid",
			chain: Chain{
				ChainId:     42,
				ChainName:   ChainName_empty,
				Network:     Network_base + 1,
				NetworkType: NetworkType_testnet,
				Vm:          Vm_evm,
				Consensus:   Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "invalid network",
		},
		{
			name: "should error if network type invalid",
			chain: Chain{
				ChainId:     42,
				ChainName:   ChainName_empty,
				Network:     Network_base,
				NetworkType: NetworkType_devnet + 1,
				Vm:          Vm_evm,
				Consensus:   Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "invalid network type",
		},
		{
			name: "should error if vm invalid",
			chain: Chain{
				ChainId:     42,
				ChainName:   ChainName_empty,
				Network:     Network_base,
				NetworkType: NetworkType_devnet,
				Vm:          Vm_evm + 1,
				Consensus:   Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "invalid vm",
		},
		{
			name: "should error if consensus invalid",
			chain: Chain{
				ChainId:     42,
				ChainName:   ChainName_empty,
				Network:     Network_base,
				NetworkType: NetworkType_devnet,
				Vm:          Vm_evm,
				Consensus:   Consensus_op_stack + 1,
				IsExternal:  true,
			},
			errStr: "invalid consensus",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.errStr != "" {
				require.ErrorContains(t, tt.chain.Validate(), tt.errStr)
			} else {
				require.NoError(t, tt.chain.Validate())
			}
		})
	}

	t.Run("all default chains are valid", func(t *testing.T) {
		for _, chain := range DefaultChainsList() {
			require.NoError(t, chain.Validate())
		}
	})
}

func TestChain_EncodeAddress(t *testing.T) {
	tests := []struct {
		name    string
		chain   Chain
		b       []byte
		want    string
		wantErr bool
	}{
		{
			name: "should error if b is not a valid address on the bitcoin network",
			chain: Chain{
				ChainName: ChainName_btc_testnet,
				ChainId:   18332,
			},
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "",
			wantErr: true,
		},
		{
			name: "should pass if b is a valid address on the network",
			chain: Chain{
				ChainName: ChainName_btc_mainnet,
				ChainId:   8332,
			},
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c",
			wantErr: false,
		},
		{
			name: "should error if b is not a valid address on the evm network",
			chain: Chain{
				ChainName: ChainName_goerli_testnet,
				ChainId:   5,
			},
			b:       ethcommon.Hex2Bytes("0x321"),
			want:    "",
			wantErr: true,
		},
		{
			name: "should pass if b is a valid address on the evm network",
			chain: Chain{
				ChainName: ChainName_goerli_testnet,
				ChainId:   5,
			},
			b:       []byte("0x321"),
			want:    "0x0000000000000000000000000000003078333231",
			wantErr: false,
		},
		{
			name: "should error if chain not supported",
			chain: Chain{
				ChainName: 999,
				ChainId:   999,
			},
			b:       ethcommon.Hex2Bytes("0x321"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s, err := tc.chain.EncodeAddress(tc.b)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, tc.want, s)
		})
	}
}

func TestChain_DecodeAddress(t *testing.T) {
	tests := []struct {
		name    string
		chain   Chain
		b       string
		want    []byte
		wantErr bool
	}{
		{
			name: "should decode on btc chain",
			chain: Chain{
				ChainName: ChainName_btc_testnet,
				ChainId:   18332,
			},
			want:    []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			b:       "bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c",
			wantErr: false,
		},
		{
			name: "should decode on evm chain",
			chain: Chain{
				ChainName: ChainName_goerli_testnet,
				ChainId:   5,
			},
			want:    ethcommon.HexToAddress("0x321").Bytes(),
			b:       "0x321",
			wantErr: false,
		},
		{
			name: "should error if chain not supported",
			chain: Chain{
				ChainName: 999,
				ChainId:   999,
			},
			want:    ethcommon.Hex2Bytes("0x321"),
			b:       "",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s, err := tc.chain.DecodeAddress(tc.b)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, tc.want, s)
		})
	}
}

func TestChain_InChainList(t *testing.T) {
	require.True(t, ZetaChainMainnet.InChainList(ChainListByNetwork(Network_zeta)))
	require.True(t, ZetaChainDevnet.InChainList(ChainListByNetwork(Network_zeta)))
	require.True(t, ZetaChainPrivnet.InChainList(ChainListByNetwork(Network_zeta)))
	require.True(t, ZetaChainTestnet.InChainList(ChainListByNetwork(Network_zeta)))
	require.False(t, Ethereum.InChainList(ChainListByNetwork(Network_zeta)))
}

func TestIsZetaChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Zeta Mainnet", ZetaChainMainnet.ChainId, true},
		{"Zeta Testnet", ZetaChainTestnet.ChainId, true},
		{"Zeta Mocknet", ZetaChainDevnet.ChainId, true},
		{"Zeta Privnet", ZetaChainPrivnet.ChainId, true},
		{"Non-Zeta", Ethereum.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsZetaChain(tt.chainID))
		})
	}
}

func TestIsEVMChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Ethereum Mainnet", Ethereum.ChainId, true},
		{"Goerli Testnet", Goerli.ChainId, true},
		{"Sepolia Testnet", Sepolia.ChainId, true},
		{"Non-EVM", BitcoinMainnet.ChainId, false},
		{"Zeta Mainnet", ZetaChainMainnet.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsEVMChain(tt.chainID))
		})
	}
}

func TestIsHeaderSupportedChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Ethereum Mainnet", Ethereum.ChainId, true},
		{"Goerli Testnet", Goerli.ChainId, true},
		{"Goerli Localnet", GoerliLocalnet.ChainId, true},
		{"Sepolia Testnet", Sepolia.ChainId, true},
		{"BSC Testnet", BscTestnet.ChainId, true},
		{"BSC Mainnet", BscMainnet.ChainId, true},
		{"BTC", BitcoinMainnet.ChainId, true},
		{"Zeta Mainnet", ZetaChainMainnet.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsHeaderSupportedChain(tt.chainID))
		})
	}
}

func TestSupportMerkleProof(t *testing.T) {
	tests := []struct {
		name  string
		chain Chain
		want  bool
	}{
		{"Ethereum Mainnet", Ethereum, true},
		{"BSC Testnet", BscTestnet, true},
		{"BSC Mainnet", BscMainnet, true},
		{"Non-EVM", BitcoinMainnet, true},
		{"Zeta Mainnet", ZetaChainMainnet, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.chain.SupportMerkleProof())
		})
	}
}

func TestIsBitcoinChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Bitcoin Mainnet", BitcoinMainnet.ChainId, true},
		{"Bitcoin Testnet", BitcoinTestnet.ChainId, true},
		{"Bitcoin Regtest", BitcoinRegtest.ChainId, true},
		{"Non-Bitcoin", Ethereum.ChainId, false},
		{"Zeta Mainnet", ZetaChainMainnet.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsBitcoinChain(tt.chainID))
		})
	}
}

func TestIsEthereumChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Ethereum Mainnet", Ethereum.ChainId, true},
		{"Goerli Testnet", Goerli.ChainId, true},
		{"Sepolia Testnet", Sepolia.ChainId, true},
		{"Non-Ethereum", BitcoinMainnet.ChainId, false},
		{"Zeta Mainnet", ZetaChainMainnet.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsEthereumChain(tt.chainID))
		})
	}
}

func TestChain_IsExternalChain(t *testing.T) {
	require.False(t, ZetaChainMainnet.IsExternalChain())
	require.True(t, Ethereum.IsExternalChain())
}

func TestChain_IsZetaChain(t *testing.T) {
	require.True(t, ZetaChainMainnet.IsZetaChain())
	require.False(t, Ethereum.IsZetaChain())
}

func TestChain_IsEmpty(t *testing.T) {
	require.True(t, Chain{}.IsEmpty())
	require.False(t, ZetaChainMainnet.IsEmpty())
}

func TestChain_WitnessProgram(t *testing.T) {
	// Ordinarily the private key would come from whatever storage mechanism
	// is being used, but for this example just hard code it.
	privKeyBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2" +
		"d4f8720ee63e502ee2869afab7de234b80c")
	require.NoError(t, err)

	t.Run("should return btc address", func(t *testing.T) {
		_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
		pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
		addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
		require.NoError(t, err)

		chain := BitcoinTestnet
		_, err = chain.BTCAddressFromWitnessProgram(addr.WitnessProgram())
		require.NoError(t, err)
	})

	t.Run("should fail for wrong chain id", func(t *testing.T) {
		_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
		pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
		addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
		require.NoError(t, err)

		chain := Goerli
		_, err = chain.BTCAddressFromWitnessProgram(addr.WitnessProgram())
		require.Error(t, err)
	})

	t.Run("should fail for wrong witness program", func(t *testing.T) {
		_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
		pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
		addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
		require.NoError(t, err)

		chain := BitcoinTestnet
		_, err = chain.BTCAddressFromWitnessProgram(addr.WitnessProgram()[0:19])
		require.Error(t, err)
	})
}

func TestChains_Has(t *testing.T) {
	chains := Chains{ZetaChainMainnet, ZetaChainTestnet}
	require.True(t, chains.Has(ZetaChainMainnet))
	require.False(t, chains.Has(Ethereum))
}

func TestChains_Distinct(t *testing.T) {
	chains := Chains{ZetaChainMainnet, ZetaChainMainnet, ZetaChainTestnet}
	distinctChains := chains.Distinct()
	require.Len(t, distinctChains, 2)
}

func TestChains_Strings(t *testing.T) {
	chains := Chains{ZetaChainMainnet, ZetaChainTestnet}
	strings := chains.Strings()
	expected := []string{chains[0].String(), chains[1].String()}
	require.Equal(t, expected, strings)
}

func TestGetChainFromChainID(t *testing.T) {
	chain := GetChainFromChainID(ZetaChainMainnet.ChainId)
	require.Equal(t, ZetaChainMainnet, *chain)
	require.Nil(t, GetChainFromChainID(9999))
}

func TestGetBTCChainParams(t *testing.T) {
	params, err := GetBTCChainParams(BitcoinMainnet.ChainId)
	require.NoError(t, err)
	require.Equal(t, &chaincfg.MainNetParams, params)

	_, err = GetBTCChainParams(9999)
	require.Error(t, err)
}

func TestGetBTCChainIDFromChainParams(t *testing.T) {
	chainID, err := GetBTCChainIDFromChainParams(&chaincfg.MainNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(8332), chainID)

	chainID, err = GetBTCChainIDFromChainParams(&chaincfg.RegressionNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(18444), chainID)

	chainID, err = GetBTCChainIDFromChainParams(&chaincfg.TestNet3Params)
	require.NoError(t, err)
	require.Equal(t, int64(18332), chainID)

	_, err = GetBTCChainIDFromChainParams(&chaincfg.Params{Name: "unknown"})
	require.Error(t, err)
}

func TestChainIDInChainList(t *testing.T) {
	require.True(t, ChainIDInChainList(ZetaChainMainnet.ChainId, ChainListByNetwork(Network_zeta)))
	require.False(t, ChainIDInChainList(Ethereum.ChainId, ChainListByNetwork(Network_zeta)))
}
