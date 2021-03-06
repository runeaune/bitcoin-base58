package base58_test

import (
	"bytes"
	"encoding/hex"
	"github.com/runeaune/bitcoin-base58"
	"testing"
)

type Test struct {
	prefix  string
	keyHex  string
	encoded string
}

var vectors = []Test{
	Test{ // 0
		base58.BitcoinPublicKeyHashPrefix,
		"7b2f2061d66d57ffb9502a091ce236ed4c1ede2d",
		"1CELa15H4DMzHtHnuz7LCpSFgFWf61Ra6A",
	},

	// The following public key hash test vectors were generated by:
	// http://gobittest.appspot.com/Address
	Test{ // 1
		base58.BitcoinPublicKeyHashPrefix,
		"010966776006953D5567439E5E39F86A0D273BEE",
		"16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM",
	},
	Test{ // 2
		base58.BitcoinPublicKeyHashPrefix,
		"89C907892A9D4F37B78D5F83F2FD6E008C4F795D",
		"1DZYQ3xEy8mkc7wToQZvKqeLrSLUMVVK41",
	},
	Test{ // 3
		base58.BitcoinPublicKeyHashPrefix,
		"0000000000000000000000000000000000000000",
		"1111111111111111111114oLvT2",
	},
	Test{ // 4
		base58.BitcoinPublicKeyHashPrefix,
		"0000000000000000000000000000000000000001",
		"11111111111111111111BZbvjr",
	},

	// The following private key test vectors were generated by:
	// http://gobittest.appspot.com/PrivateKey
	Test{ // 5
		base58.BitcoinPrivateKeyPrefix,
		"0C28FCA386C7A227600B2FE50B7CAE11EC86D3BF1FBE471BE89827E19D72AA1D",
		"5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ",
	},
	Test{ // 6
		base58.BitcoinPrivateKeyPrefix,
		"00000037FC2B523A9101D653ECB504EBB88FCCE6F7E77548A7B31FA734A00000",
		"5HpHagjigF1P3i1WyFp1uLPEo8gK32CFBRc2ekJU3nytmXnVbYv",
	},
	Test{ // 7
		base58.BitcoinPrivateKeyPrefix,
		"3C80FA4C012E37402C6D43140EC7B111B931C33799C2A07E8172827B12EEA59F",
		"5JGw52y3UuZSwpZKYkuhmat8TNy1nZ6F6mbrHJaPNMt2pkETUVE",
	},
	Test{ // 8
		base58.BitcoinPrivateKeyPrefix,
		"0000000000000000000000000000000000000000000000000000000000000001",
		"5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAnchuDf",
	},

	// The extended keys are taken from the test vectors of BIP32:
	// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
	Test{ // 9
		base58.BitcoinExtendedPublicKeyPrefix,
		"000000000000000000873dff81c02f525623fd1fe5167eac3a55a049de3d314b" +
			"b42ee227ffed37d5080339a36013301597daef41fbe593a02cc513d0" +
			"b55527ec2df1050e2e8ff49c85c2",
		"xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29E" +
			"SFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8",
	},
	Test{ // 10
		base58.BitcoinExtendedPrivateKeyPrefix,
		"0478412e3afffffffe637807030d55d01f9a0cb3a7839515d796bd07706386a6" +
			"eddf06cc29a65a0e2900f1c7c871a54a804afe328b4c83a1c33b8e5f" +
			"f48f5087273f04efa83b247d6a2d",
		"xprvA1RpRA33e1JQ7ifknakTFpgNXPmW2YvmhqLQYMmrj4xJXXWYpDPS3xz7iAxn" +
			"8L39njGVyuoseXzU6rcxFLJ8HFsTjSyQbLYnMpCqE2VbFWc",
	},
}

func TestShortInput(t *testing.T) {
	_, err := base58.CheckDecodeString("16UwLL9")
	if err == nil {
		t.Errorf("Should have failed with checksum mismatch.")
	}
	_, err = base58.CheckDecodeString("1")
	if err == nil {
		t.Errorf("Should have failed with input too short.")
	}
	_, err = base58.CheckDecodeString("")
	if err == nil {
		t.Errorf("Should have failed with input too short.")
	}

	// Just make sure these don't crash.
	key, _ := hex.DecodeString("")
	_ = base58.CheckEncodeToString(key)
	key2, _ := hex.DecodeString("00")
	_ = base58.CheckEncodeToString(key2)
}

func TestBitcoinKeyEncodeDecode(t *testing.T) {
	for i, test := range vectors {
		key, err := hex.DecodeString(test.keyHex)
		if err != nil {
			t.Fatalf("Test %d: Failed to hex-decode test-key: %v", i, err)
		}

		encoded, err := base58.BitcoinCheckEncode(test.prefix, key)
		if err != nil {
			t.Fatalf("Test %d: Failed to encode key %x: %v", i, key, err)
		}
		if encoded != test.encoded {
			t.Errorf("Test %d: Encoded string mismatch; expected %q, got %q.",
				i, test.encoded, encoded)
		}
		decoded, prefix, err := base58.BitcoinCheckDecode(encoded)
		if err != nil {
			t.Fatalf("Test %d: Failed to decode key %s: %v", i, encoded, err)
		}
		if prefix != test.prefix {
			t.Errorf("Test %d: Prefix mismatch; expected %q, got %q.",
				i, test.prefix, prefix)
		}
		if !bytes.Equal(key, decoded) {
			t.Errorf("Key mismatch; expected %x, got %s.", key, decoded)
		}
	}
}
