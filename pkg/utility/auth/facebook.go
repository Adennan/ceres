package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/// Facebook REST client
/// implemnetes the OauthClient interface
/// https://developers.facebook.com/docs/facebook-login/advanced
type Facebook struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	client       *http.Client
}

type facebookAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint64 `json:"expires_in"`
}

/// GetAccessToken
/// use the requestToken to get the access token which will be used to get the github user information
func (facebook *Facebook) GetAccessToken(requestToken string) (accessToken string, err error) {
	u := fmt.Sprintf(
		"https://graph.facebook.com/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		facebook.ClientID,
		facebook.RedirectURI,
		facebook.ClientSecret,
		requestToken,
	)
	request, _ := http.NewRequest(
		"GET",
		u,
		nil,
	)
	response, err := facebook.client.Do(request)
	if err != nil {
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	r := facebookAccessTokenResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}
	accessToken = r.AccessToken
	return

}

type facebookInspectBody struct {
	UserID string `json:"user_id"`
}

type facebookInspectResposne struct {
	Data facebookInspectBody `json:"data"`
}

func (response *facebookInspectResposne) GetUserID() string {
	return response.Data.UserID
}

/// FacebookOauthAccount
/// Facebook Oauth account profile
type FacebookOauthAccount struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"` // FIXME: should check if this field is String type or object 
}

/// implement the OauthAccount interface
func (account *FacebookOauthAccount) GetUserID() string {
	return account.ID
}

/// implement the OauthAccount interface
func (account *FacebookOauthAccount) GetUserAvatar() string {
	return account.Picture
}

/// implement the OauthAccount interface
func (account *FacebookOauthAccount) GetUserNick() string {
	return account.Name
}

/// GetUserProfile
/// Facebook Oauth get user profile logic
func (facebook *Facebook) GetUserProfile(accessToken string, _userId string) (account OauthAccount, err error) {
	// firstly to intspect the access token to get the facebook user ID
	request, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://graph.facebook.com/debug_token?input_token={token-to-inspect}&access_token=%s", accessToken),
		nil,
	)
	response, err := facebook.client.Do(request)
	if err != nil {
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	r := facebookInspectResposne{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}

	// get user avatar result have to try 3 times to confirm the avatar get it
	request, _ = http.NewRequest(
		"GET",
		fmt.Sprintf("https://graph.facebook.com/v10.0/%s?fields=id,name,picture&access_token=%s", r.GetUserID(), accessToken),
		nil,
	)
	response, err = facebook.client.Do(request)
	if err != nil {
		return
	}
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	account = &FacebookOauthAccount{}
	err = json.Unmarshal(body, account)
	if err != nil {
		return
	}
	return
}