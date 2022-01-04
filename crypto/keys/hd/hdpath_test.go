package hd

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cosmos/go-bip39"
)

var defaultBIP39Passphrase = ""

// return bip39 seed with empty passphrase
func mnemonicToSeed(mnemonic string) []byte {
	return bip39.NewSeed(mnemonic, defaultBIP39Passphrase)
}

//nolint
func ExampleStringifyPathParams() {
	path := NewParams(44, 0, 0, false, 0)
	fmt.Println(path.String())
	// Output: 44'/0'/0'/0/0
}

func TestParamsFromPath(t *testing.T) {
	goodCases := []struct {
		params *BIP44Params
		path   string
	}{
		{&BIP44Params{44, 0, 0, false, 0}, "44'/0'/0'/0/0"},
		{&BIP44Params{44, 1, 0, false, 0}, "44'/1'/0'/0/0"},
		{&BIP44Params{44, 0, 1, false, 0}, "44'/0'/1'/0/0"},
		{&BIP44Params{44, 0, 0, true, 0}, "44'/0'/0'/1/0"},
		{&BIP44Params{44, 0, 0, false, 1}, "44'/0'/0'/0/1"},
		{&BIP44Params{44, 1, 1, true, 1}, "44'/1'/1'/1/1"},
		{&BIP44Params{44, 714, 52, true, 41}, "44'/714'/52'/1/41"},
	}

	for i, c := range goodCases {
		params, err := NewParamsFromPath(c.path)
		errStr := fmt.Sprintf("%d %v", i, c)
		assert.NoError(t, err, errStr)
		assert.EqualValues(t, c.params, params, errStr)
		assert.Equal(t, c.path, c.params.String())
	}

	badCases := []struct {
		path string
	}{
		{"43'/0'/0'/0/0"},   // doesnt start with 44
		{"44'/1'/0'/0/0/5"}, // too many fields
		{"44'/0'/1'/0"},     // too few fields
		{"44'/0'/0'/2/0"},   // change field can only be 0/1
		{"44/0'/0'/0/0"},    // first field needs '
		{"44'/0/0'/0/0"},    // second field needs '
		{"44'/0'/0/0/0"},    // third field needs '
		{"44'/0'/0'/0'/0"},  // fourth field must not have '
		{"44'/0'/0'/0/0'"},  // fifth field must not have '
		{"44'/-1'/0'/0/0"},  // no negatives
		{"44'/0'/0'/-1/0"},  // no negatives
	}

	for i, c := range badCases {
		params, err := NewParamsFromPath(c.path)
		errStr := fmt.Sprintf("%d %v", i, c)
		assert.Nil(t, params, errStr)
		assert.Error(t, err, errStr)
	}

}

//nolint
func ExampleSomeBIP32TestVecs() {

	seed := mnemonicToSeed("barrel original fuel morning among eternal " +
		"filter ball stove pluck matrix mechanic")
	master, ch := ComputeMastersFromSeed(seed)
	fmt.Println("keys from fundraiser test-vector (cosmos, bitcoin, ether)")
	fmt.Println()
	// cosmos
	priv, _ := DerivePrivateKeyForPath(master, ch, FullFundraiserPath)
	fmt.Println(hex.EncodeToString(priv[:]))
	// bitcoin
	priv, _ = DerivePrivateKeyForPath(master, ch, "44'/0'/0'/0/0")
	fmt.Println(hex.EncodeToString(priv[:]))
	// ether
	priv, _ = DerivePrivateKeyForPath(master, ch, "44'/60'/0'/0/0")
	fmt.Println(hex.EncodeToString(priv[:]))

	fmt.Println()
	fmt.Println("keys generated via https://coinomi.com/recovery-phrase-tool.html")
	fmt.Println()

	seed = mnemonicToSeed(
		"advice process birth april short trust crater change bacon monkey medal garment " +
			"gorilla ranch hour rival razor call lunar mention taste vacant woman sister")
	master, ch = ComputeMastersFromSeed(seed)
	priv, _ = DerivePrivateKeyForPath(master, ch, "44'/1'/1'/0/4")
	fmt.Println(hex.EncodeToString(priv[:]))

	seed = mnemonicToSeed("idea naive region square margin day captain habit " +
		"gun second farm pact pulse someone armed")
	master, ch = ComputeMastersFromSeed(seed)
	priv, _ = DerivePrivateKeyForPath(master, ch, "44'/0'/0'/0/420")
	fmt.Println(hex.EncodeToString(priv[:]))

	fmt.Println()
	fmt.Println("BIP 32 example")
	fmt.Println()

	// bip32 path: m/0/7
	seed = mnemonicToSeed("monitor flock loyal sick object grunt duty ride develop assault harsh history")
	master, ch = ComputeMastersFromSeed(seed)
	priv, _ = DerivePrivateKeyForPath(master, ch, "0/7")
	fmt.Println(hex.EncodeToString(priv[:]))

	// Output: keys from fundraiser test-vector (cosmos, bitcoin, ether)
	//
	// 01dcb36acfd5de52ac1f00daf231e64637388202f1fce7bdc64f6bb3199d270d
	// e77c3de76965ad89997451de97b95bb65ede23a6bf185a55d80363d92ee37c3d
	// 7fc4d8a8146dea344ba04c593517d3f377fa6cded36cd55aee0a0bb968e651bc
	//
	// keys generated via https://coinomi.com/recovery-phrase-tool.html
	//
	// a61f10c5fecf40c084c94fa54273b6f5d7989386be4a37669e6d6f7b0169c163
	// 32c4599843de3ef161a629a461d12c60b009b676c35050be5f7ded3a3b23501f
	//
	// BIP 32 example
	//
	// c4c11d8c03625515905d7e89d25dfc66126fbc629ecca6db489a1a72fc4bda78
}

// Ensuring that we don't crash if values have trailing slashes
func TestDerivePrivateKeyForPathDoNotCrash(t *testing.T) {
	paths := []string{
		"m/5/",
		"m/5",
		"/44",
		"m//5",
		"m/0/7",
		"/",
		" m       /0/7",
		"              /       ",
		"m///7//////",
	}

	for _, path := range paths {
		path := path
		t.Run(path, func(t *testing.T) {
			DerivePrivateKeyForPath([32]byte{}, [32]byte{}, path)
		})
	}
}
