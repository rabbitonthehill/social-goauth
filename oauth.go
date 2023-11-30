package oauth

import (
	"errors"
	"strings"
)

// AuthType is used to represent different types of third-party login methods.
type AuthType string

// Different third-party login methods.
const (
	AuthGoogle   AuthType = "Google"
	AuthApple    AuthType = "Apple"
	AuthFacebook AuthType = "Facebook"
	AuthLine     AuthType = "Line"
)

var (
	ErrInvalidSignature    = errors.New("invalid id signature")
	ErrInvalidIdToken      = errors.New("invalid id token")
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrFetchKeysFail       = errors.New("fetch keys fail")
	ErrInvalidHashType     = errors.New("invalid hash type")
	ErrInvalidIdCode       = errors.New("invalid id code")
	ErrInvalidClientID     = errors.New("invalid client id")
	ErrInvalidClientSecret = errors.New("invalid client secret")
	ErrInvalidRedirectURL  = errors.New("invalid redirect url")
)

// Service represents the basic configuration for OAuth.
type Service struct {
	// ClientID Identifier assigned by the third-party login provider to identify your application.
	ClientID string

	// ClientSecret Secret key used for secure communication with the third-party login provider.
	ClientSecret string

	// RedirectURL URL that the third-party login provider redirects the user to after successful login.
	RedirectURL string

	// ProxyURL Optional proxy URL used during the login process.
	ProxyURL string

	// AuthType Type of third-party login provider being used.
	AuthType AuthType

	// Endpoint where the OAuth server handles the authentication request.
	Endpoint string
}

type Option func(*Service)

// WithRedirectURL sets the RedirectURL option for the Service.
func WithRedirectURL(url string) Option {
	return func(service *Service) {
		service.RedirectURL = url
	}
}

// WithProxyURL sets the ProxyURL option for the Service.
func WithProxyURL(url string) Option {
	return func(service *Service) {
		service.ProxyURL = url
	}
}

// Endpoint returns a URL endpoint given an input string and an endpoint base.
// If the input string begins with "http://" or "https://", it is returned as-is.
// If the input string begins with "/", it is appended to the endpoint base.
// Otherwise, the input string is appended to the endpoint base with a "/" separator.
func Endpoint(endpoint, input string) string {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return input
	}

	if strings.HasPrefix(input, "/") {
		return endpoint + input
	}

	return endpoint + "/" + input
}

// NewService creates a new OAuth service with the provided client ID, client secret, and authentication type.
// It also allows additional options to be applied using the Option functional parameter.
func NewService(clientID, clientSecret string, authType AuthType, options ...Option) (*Service, error) {
	if clientID == "" {
		return nil, ErrInvalidClientID
	}
	if clientSecret == "" {
		return nil, ErrInvalidClientSecret
	}
	service := &Service{ClientID: clientID, ClientSecret: clientSecret, AuthType: authType}
	for _, opt := range options {
		opt(service)
	}

	return service, nil
}
