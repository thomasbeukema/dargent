package address

import (
    "encoding/base64"

    "github.com/Yawning/sphincs256"
    "github.com/mr-tron/base58/base58"
    "github.com/thomasbeukema/fastrand"
)

type SPHINCSKeyPair struct {
	PrivateKey	*[sphincs256.PrivateKeySize]byte
	PublicKey	*[sphincs256.PublicKeySize]byte
    Entropy     *[32]byte
}

func GenerateSPHINCSKeyPair(ent []byte) SPHINCSKeyPair {

    fastrand.New()

    var entropy [32]byte

    if ent != nil {
        for i := range ent {
            entropy[i] = ent[i]
        }

        fastrand.SetEntropy(entropy)
    } else {
        entropy = fastrand.GetEntropy()
    }

	pub, priv, err := sphincs256.GenerateKey(fastrand.Reader)
	if err != nil {
		panic(1)
	}

	return SPHINCSKeyPair{priv, pub, &entropy}
}

func SPHINCSPubKeyToAddress(pubkey []byte) string {
    pubkey = append(HashPubKey(pubkey), sphincsPadding...)

    b32Pubkey := base58.Encode(pubkey)                                 // Generate Base32 of public key
    b32Checksum := base58.Encode(generateChecksum(HashPubKey(pubkey))) // Get checksum for the public key and generate Base32 of it

    return "999" + b32Pubkey + b32Checksum + "666"
}

func (kp *SPHINCSKeyPair) Mnemonic() string {
    return getMnemonic(kp.Entropy[:])
}

func validateSPHINCSAddress(address string) bool {

    if address[:3] != "999" || address[len(address)-3:] != "666" {
        return false
    }

    address = address[3:len(address)-3] // Strip '999' & '666'

    checksum := address[len(address)-8:]
    generatedChecksum := base58.Encode(generateChecksum([]byte(address[:len(address)-8])))

    return checksum == generatedChecksum
}

func (kp *SPHINCSKeyPair) GetAddress() string {
    return SPHINCSPubKeyToAddress((*kp.PublicKey)[:])
}

func (kp *SPHINCSKeyPair) Sign(hash []byte) string {
    signature := sphincs256.Sign(kp.PrivateKey, hash)
    b64sig := base64.StdEncoding.EncodeToString(signature[:])

    return b64sig
}

func ValidateSPHINCSSignature(sig string, hash []byte, pubkey *[1056]byte) bool {
    rawSig, err := base64.StdEncoding.DecodeString(sig)
    if err != nil {
        // TODO: Proper err handling
        panic(err)
    }

    var finalSig [41000]byte
    for i := range rawSig {
        finalSig[i] = rawSig[i]
    }

    return sphincs256.Verify(pubkey, hash, &finalSig)
}
