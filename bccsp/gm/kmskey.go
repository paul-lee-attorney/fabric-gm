package gm

import (
	"crypto/sha256"
	"crypto/elliptic"
	"github.com/Hyperledger-TWGC/tjfoc-gm/sm2"

	"github.com/pkg/errors"
	kmssm2 "github.com/tw-bc-group/aliyun-kms/sm2"
	"github.com/tw-bc-group/fabric-gm/bccsp"
)

type kmsSm2PrivateKey struct {
	adapter *kmssm2.KeyAdapter
}

func (pri *kmsSm2PrivateKey) Bytes() ([]byte, error) {
	return nil, errors.Errorf("Unsupported")
}

func (pri *kmsSm2PrivateKey) SKI() []byte {
	publicKey := pri.adapter.Public().(*sm2.PublicKey)
	raw := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

func (pri *kmsSm2PrivateKey) Symmetric() bool {
	return false
}

func (pri *kmsSm2PrivateKey) Private() bool {
	return true
}

func (pri *kmsSm2PrivateKey) PublicKey() (bccsp.Key, error) {
	pubKey := pri.adapter.Public().(*sm2.PublicKey)
	return &gmsm2PublicKey{pubKey: pubKey}, nil
}

func createKmsSm2PrivateKey() (*kmsSm2PrivateKey, error) {
	adapter, err := kmssm2.CreateSm2KeyAdapter("", kmssm2.SignAndVerify)
	if err != nil {
		return nil, err
	}

	return &kmsSm2PrivateKey{
		adapter: adapter,
	}, nil
}

type kmssm2ImportKeyOptsKeyImporter struct{}

func (*kmssm2ImportKeyOptsKeyImporter) KeyImport(raw interface{}, opts bccsp.KeyImportOpts) (k bccsp.Key, err error) {
	adapter, err := kmssm2.CreateSm2KeyAdapter(raw.(string), kmssm2.SignAndVerify)
	if err != nil {
		return nil, err
	}
	return &kmsSm2PrivateKey{
		adapter: adapter,
	}, nil
}

type kmssm2PrivateKeySigner struct{}

func (s *kmssm2PrivateKeySigner) Sign(k bccsp.Key, digest []byte, opts bccsp.SignerOpts) (signature []byte, err error) {
	return k.(*kmsSm2PrivateKey).adapter.AsymmetricSign(digest)
}
