package oauth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Constants for Apple URLs
const (
	AppleBaseEndpoint  = "https://appleid.apple.com"
	AppleURLAuthKeys   = AppleBaseEndpoint + "/auth/keys"
	AppleURLAuthToken  = AppleBaseEndpoint + "/auth/token"
	AppleURLAuthRevoke = AppleBaseEndpoint + "/auth/revoke"
)

// Apple struct represents the Apple OAuth provider.
type Apple struct {
	service *Service
}

// AppleClaims struct represents the claims in Apple Identity Token.
type AppleClaims struct {
	// Expiration time of the token
	Exp int64 `json:"exp"`

	// Issued at time of the token
	Iat int64 `json:"iat"`

	// Time when the user authenticated
	AuthTime int64 `json:"auth_time"`

	// Issuer of the token
	Iss string `json:"iss"`

	// Audience of the token
	Aud string `json:"aud"`

	// Subject of the token
	Sub string `json:"sub"`

	// Code hash
	CHash string `json:"c_hash"`
	// Email address of the user

	Email string `json:"email"`
	// Indicates if the email is verified

	EmailVerified string `json:"email_verified"`

	// Indicates if nonce is supported
	NonceSupported bool `json:"nonce_supported"`
}

// ApplePublicKey struct represents the public key used for signature verification.
type ApplePublicKey struct {
	// Key type
	Kty string `json:"kty"`

	// Key ID
	Kid string `json:"kid"`

	// Key usage
	Use string `json:"use"`

	// Key algorithm
	Alg string `json:"alg"`

	// Modulus
	N string `json:"n"`

	// Exponent
	E string `json:"e"`
}

// ApplePublicKeyResponse struct represents the response containing the Apple public keys.
type ApplePublicKeyResponse struct {
	Keys []*ApplePublicKey `json:"keys"` // List of Apple public keys
}

// NewApple creates a new instance of the Apple OAuth provider.
func NewApple(service *Service) *Apple {
	service.Endpoint = AppleBaseEndpoint
	return &Apple{service: service}
}

// getPublicKey retrieves the Apple public keys.
func (p *Apple) getPublicKey() (ApplePublicKeyResponse, error) {
	resp, err := New(AppleURLAuthKeys, http.MethodGet, p.service.ProxyURL, WithTimeout(30*time.Second)).Do()
	if nil != err {
		return ApplePublicKeyResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ApplePublicKeyResponse{}, fmt.Errorf("the status code is: %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApplePublicKeyResponse{}, err
	}
	var value ApplePublicKeyResponse
	if err = json.Unmarshal(data, &value); err != nil {
		return ApplePublicKeyResponse{}, err
	}

	return value, nil
}

// decodePayload decodes the payload of the Identity Token.
func (p *Apple) decodePayload(str string) (*AppleClaims, error) {
	payload, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("failed to base64url decode ID Token: %s", err.Error())
	}
	var claims *AppleClaims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ID Token claims: %s", err.Error())
	}
	return claims, nil
}

// VerifySignature verifies the signature of the Identity Token.
func (p *Apple) VerifySignature(val []string) error {
	// Step 1: Get the public key
	keys, err := p.getPublicKey()
	if err != nil {
		return err
	}

	// Step 2: Extract the encryption algorithm from the header
	headerBytes, err := base64.RawURLEncoding.DecodeString(val[0])
	if err != nil {
		return err
	}
	var header struct {
		Alg string `json:"alg"` // Encryption algorithm
		Kid string `json:"kid"` // Key ID
	}
	if err = json.Unmarshal(headerBytes, &header); err != nil {
		return err
	}

	// Step 3: Find the matching public key in the collection
	publicKey := &ApplePublicKey{}
	for _, key := range keys.Keys {
		if key.Alg == header.Alg && key.Kid == header.Kid {
			publicKey = key
			break
		}
	}
	// No matching public key found
	if publicKey.Kid == "" {
		return ErrInvalidSignature
	}

	// Step 4: Verify the signature using the public key
	data := val[0] + "." + val[1]
	signature, err := base64.RawURLEncoding.DecodeString(val[2])
	if err != nil {
		return err
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(publicKey.N)
	if err != nil {
		return err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(publicKey.E)
	if err != nil {
		return err
	}

	pubKey := &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(nBytes),
		E: int(big.NewInt(0).SetBytes(eBytes).Int64()),
	}

	hashed := sha256.Sum256([]byte(data))
	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature); nil != err {
		return err
	}

	return nil
}

// IdToken verifies the Apple Identity Token.
func (p *Apple) IdToken(token string) (*AppleClaims, error) {
	if token == "" {
		return nil, ErrInvalidIdToken
	}
	// Split the token into header, payload, and signature (arr[0], arr[1], arr[2])
	arr := strings.Split(token, ".")
	if err := p.VerifySignature(arr); nil != err {
		return nil, err
	}
	claims, err := p.decodePayload(arr[1])
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// IdentityCode verifies the Apple Identity Code.
func (p *Apple) IdentityCode(code string) (int, error) {
	if code == "" {
		return -1, ErrInvalidIdCode
	}
	// The redirect_uri parameter must be provided when verifying the Identity Code,
	// and it must use the HTTPS protocol.
	// if uri := strings.ToLower(o.RedirectUri); strings.HasPrefix(uri, "https://") {
	// 	return nil, ErrInvalidRedirectURI
	//}
	params := url.Values{
		"client_id":     []string{p.service.ClientID},
		"client_secret": []string{p.service.ClientSecret},
		"code":          []string{code},
		"grant_type":    []string{"authorization_code"},
		"redirect_uri":  []string{p.service.RedirectURL},
	}
	header := http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
	}
	resp, err := New(AppleURLAuthToken, http.MethodPost, p.service.ProxyURL,
		WithTimeout(30*time.Second),
		WithHeader(header),
		WithData(params),
	).Post()
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("the status code is: %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var result map[string]interface{}
	if err = json.Unmarshal(data, &result); err != nil {
		return -1, err
	}
	return result["code"].(int), nil
}
