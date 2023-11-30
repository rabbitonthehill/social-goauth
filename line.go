package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	LineBaseEndpoint          = "https://api.line.me"
	LineURLAccessToken        = LineBaseEndpoint + "/oauth2/v2.1/token"
	LineURLVerifyAccessToken  = LineBaseEndpoint + "/oauth2/v2.1/verify"
	LineURLRefreshAccessToken = LineURLAccessToken
	LineURLRevokeAccessToken  = LineBaseEndpoint + "/oauth2/v2.1/revoke"
	LineURLVerifyIDToken      = LineURLVerifyAccessToken
	LineURLUserInformation    = LineBaseEndpoint + "/oauth2/v2.1/userinfo"
	LineURLProfile            = LineBaseEndpoint + "/v2/profile"
	LineURLFriendshipStatus   = LineBaseEndpoint + "/friendship/v1/status"
)

type LineAccessToken struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LineAccessTokenVerification struct {
	Scope     string `json:"profile"`
	ClientId  string `json:"client_id"`
	ExpiresIn int64  `json:"expires_in"`
}

type LineIDToken struct {
	Iss     string   `json:"iss"`
	Sub     string   `json:"sub"`
	Aud     string   `json:"aud"`
	Exp     int64    `json:"exp"`
	Iat     int64    `json:"iat"`
	Nonce   string   `json:"nonce"`
	Amr     []string `json:"amr"`
	Name    string   `json:"name"`
	Picture string   `json:"picture"`
	Email   string   `json:"email"`
}

type LineUserInformation struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type LineUserProfile struct {
	UserId        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureUrl    int64  `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

type Line struct {
	service *Service
}

func NewLine(service *Service) *Line {
	service.Endpoint = LineBaseEndpoint
	return &Line{service: service}
}

// AccessToken Verifies if an access token is valid.
//
// For general recommendations on how to securely handle user registration and login with access tokens,
// see Creating a secure login process between your app and server(https://developers.line.biz/en/docs/line-login/secure-login-process/) in the LINE Login documentation.
//
// Note:
// This is the reference for the LINE Login v2.1 endpoint. For information on the v2.0 endpoint,
// see Verify access token validity(https://developers.line.biz/en/reference/line-login/#verify-access-token) in the LINE Login v2.0 API reference.
//
// documentation https://developers.line.biz/en/reference/line-login/#verify-access-token
func (p *Line) AccessToken(accessToken string) (*LineAccessTokenVerification, error) {
	if "" == accessToken {
		return nil, ErrInvalidAccessToken
	}
	u := fmt.Sprintf("%s?access_token=%s", LineURLVerifyAccessToken, accessToken)
	resp, err := New(u, http.MethodGet, p.service.ProxyURL).Get()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// If the access token has expired, a 400 Bad Request HTTP status code and a JSON response are returned
	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("the status code is : %d", resp.StatusCode)
	}
	fmt.Println(string(value))
	data := &LineAccessTokenVerification{}
	err = json.Unmarshal(value, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RefreshAccessToken Gets a new access token using a refresh token.
//
// A refresh token is returned along with an access token once user authentication is complete.
// Note:
// This is the reference for the LINE Login v2.1 endpoint. For information on the v2.0 endpoint,
// see Refresh access token in the LINE Login v2.0 API reference.
// You can't use this to refresh a channel access token for the Messaging API.
//
// documentation https://developers.line.biz/en/reference/line-login/#refresh-access-token
func (p *Line) RefreshAccessToken(RefreshToken string) (*LineAccessToken, error) {
	if "" == RefreshToken {
		return nil, ErrInvalidRefreshToken
	}
	params := url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{RefreshToken},
		"client_id":     []string{p.service.ClientID},
		"client_secret": []string{p.service.ClientSecret},
	}
	header := http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
	}
	resp, err := New(LineURLRefreshAccessToken, http.MethodPost, p.service.ProxyURL,
		WithData(params),
		WithHeader(header),
		WithTimeout(30*time.Second),
	).Post()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("the status code is : %d", resp.StatusCode)
	}
	fmt.Println(string(value))
	data := &LineAccessToken{}
	err = json.Unmarshal(value, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RevokeAccessToken Invalidates a user's access token.
//
// Note:
// This is the reference for the LINE Login v2.1 endpoint. For information on the v2.0 endpoint,
// see Revoke access token(https://developers.line.biz/en/reference/line-login-v2/#revoke-access-token) in the LINE Login v2.0 API reference.
// You can't use this to invalidate a channel access token for the Messaging API.
//
// documentation https://developers.line.biz/en/reference/line-login/#revoke-access-token
func (p *Line) RevokeAccessToken(accessToken string) (bool, error) {
	if "" == accessToken {
		return false, ErrInvalidAccessToken
	}
	header := http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
	}
	params := url.Values{
		"access_token": []string{accessToken},
		"client_id":    []string{p.service.ClientID},
		// "client_secret": []string{o.ClientSecret},
	}
	resp, err := New(LineURLRevokeAccessToken, http.MethodPost, p.service.ProxyURL,
		WithData(params),
		WithHeader(header),
		WithTimeout(30*time.Second),
	).Post()
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if http.StatusOK != resp.StatusCode {
		return false, fmt.Errorf("the status code is : %d", resp.StatusCode)
	}
	return true, nil
}

// IDToken ID tokens are JSON web tokens (JWT) with information about the user.
// It's possible for an attacker to spoof an ID token(https://developers.line.biz/en/docs/line-login/verify-id-token/#id-tokens).
// Use this call to verify that a received ID token is authentic,
// meaning you can use it to obtain the user's profile information and email.
//
// documentation https://developers.line.biz/en/reference/line-login/#verify-id-token
func (p *Line) IDToken(idToken string) (*LineIDToken, error) {
	if "" == idToken {
		return nil, ErrInvalidIdToken
	}
	header := http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded"},
	}
	params := url.Values{
		"id_token":  []string{idToken},
		"client_id": []string{p.service.ClientID},
	}
	resp, err := New(LineURLVerifyIDToken, http.MethodPost, p.service.ProxyURL,
		WithData(params),
		WithHeader(header),
		WithTimeout(30*time.Second),
	).Post()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(value))
	data := &LineIDToken{}
	err = json.Unmarshal(value, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UserInformation Gets a user's ID, display name, and profile image.
// The scope required for the access token is different for the Get user profile(https://developers.line.biz/en/reference/line-login/#get-user-profile) endpoint.
//
// Note:
// Requires an access token with the openid scope.
// For more information,see Authenticating users and making authorization requests(https://developers.line.biz/en/docs/line-login/integrate-line-login/#making-an-authorization-request)
// and Scopes(https://developers.line.biz/en/docs/line-login/integrate-line-login/#scopes) in the LINE Login documentation.
//
// documentation https://developers.line.biz/en/reference/line-login/#userinfo
func (p *Line) UserInformation(accessToken string) (*LineUserInformation, error) {
	if "" == accessToken {
		return nil, ErrInvalidAccessToken
	}
	//_, err := p.VerifyAccessToken(accessToken)
	//if nil != err {
	//	return nil, err
	//}
	header := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", accessToken)},
	}
	resp, err := New(LineURLProfile, http.MethodGet, p.service.ProxyURL, WithTimeout(30*time.Second), WithHeader(header)).Get()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(value))
	data := &LineUserInformation{}
	err = json.Unmarshal(value, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UserProfile Gets a user's ID, display name, profile image, and status message.
// The scope required for the access token is different for the Get user information(https://developers.line.biz/en/reference/line-login/#userinfo) endpoint.
//
// Note: Requires an access token with the profile scope.
// For more information, see Authenticating users and making authorization requests(https://developers.line.biz/en/docs/line-login/integrate-line-login/#making-an-authorization-request)
// and Scopes(https://developers.line.biz/en/docs/line-login/integrate-line-login/#scopes) in the LINE Login documentation.
//
// documentation https://developers.line.biz/en/reference/line-login/#get-user-profile
func (p *Line) UserProfile(accessToken string) (*LineUserProfile, error) {
	if "" == accessToken {
		return nil, ErrInvalidAccessToken
	}
	header := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", accessToken)},
	}
	resp, err := New(LineURLProfile, http.MethodGet, p.service.ProxyURL, WithTimeout(30*time.Second), WithHeader(header)).Get()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(value))
	data := &LineUserProfile{}
	err = json.Unmarshal(value, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// FriendshipStatus Gets the friendship status between a user and the LINE Official Account linked to your LINE Login channel.
//
// For more information on how to link a LINE Official Account to a LINE Login channel,
// see Add a LINE Official Account as a friend when logged in (bot link)(https://developers.line.biz/en/docs/line-login/link-a-bot/) in the LINE Login documentation.
//
// Note:
// Requires an access token with the profile scope.
// For more information, see Authenticating users and making authorization requests(https://developers.line.biz/en/docs/line-login/integrate-line-login/#making-an-authorization-request)
// and Scopes(https://developers.line.biz/en/docs/line-login/integrate-line-login/#scopes) in the LINE Login documentation
//
// https://developers.line.biz/en/reference/line-login/#get-friendship-status
func (p *Line) FriendshipStatus(accessToken string) (bool, error) {
	if "" == accessToken {
		return false, ErrInvalidAccessToken
	}
	header := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", accessToken)},
	}
	resp, err := New(LineURLFriendshipStatus, http.MethodGet, p.service.ProxyURL,
		WithTimeout(30*time.Second),
		WithHeader(header),
	).Get()
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	value, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	fmt.Println(string(value))
	var data map[string]bool
	err = json.Unmarshal(value, &data)
	if err != nil {
		return false, err
	}

	return data["friendFlag"], nil
}
