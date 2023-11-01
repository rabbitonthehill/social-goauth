package oauth

import (
	"errors"
	"net/http"
	"net/url"
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
	ErrInvalidRedirectURI  = errors.New("invalid redirect uri")
)

// Provider interface defines the methods that a third-party login provider should implement.
type Provider interface {
	Do() error
	IDToken(token string) error
	IdentityCode(code string) error
}

// Service represents the basic configuration for OAuth.
type Service struct {
	// ClientID is the identifier that is used to identify your application.
	ClientID string

	// ClientSecret is the secret key that is used to communicate securely with the third-party login provider.
	ClientSecret string

	// RedirectURL is the URL that the third-party login provider will redirect the user to after a successful login.
	RedirectURL string

	// ProxyURL is an optional field that specifies a proxy URL to be used during the login process.
	ProxyURL string

	// AuthType is the type of third-party login provider being used.
	AuthType AuthType
}

// UserInfo represents user information.
type UserInfo struct {
	// ID is the unique identifier of the user.
	ID string

	// Name is the user's full name, where the first name and last name are combined.
	Name string

	// Avatar is the file path or URL of the user's profile picture.
	Avatar string

	// Email is the email address associated with the user's account.
	Email string

	// Gender indicates the user's gender. Possible values are:
	// 0: Unknown or not specified
	// 1: Male
	// 2: Female
	Gender int8
}

type Option func(*Service)

func WithRedirectURL(url string) Option {
	return func(service *Service) {
		service.RedirectURL = url
	}
}

func WithProxyURL(url string) Option {
	return func(service *Service) {
		service.ProxyURL = url
	}
}

func New(clientId, clientSecret string, authType AuthType, options ...Option) *Service {
	service := &Service{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		AuthType:     authType,
	}
	for _, opt := range options {
		opt(service)
	}

	return service
}

func (o Service) setProxy(client *http.Client) {
	if "" != o.ProxyURL {
		proxy, _ := url.Parse(o.ProxyURL)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}
}
