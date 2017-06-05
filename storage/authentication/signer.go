package oldauthentication

const authorizationHeaderFormat = `Signature keyId="%s",algorithm="%s",headers="%s",signature="%s"`

type Signer interface {
	Sign(dateHeader string) (string, error)
	SignRaw(toSign string) (string, string, error)
	KeyFingerprint() string
	DefaultAlgorithm() string
}
