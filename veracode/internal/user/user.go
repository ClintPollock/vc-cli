package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	identityUser "github.com/veracode/veracode-cli/internal/api/identity/user"
	"github.com/veracode/veracode-cli/internal/hmac"
)

type User struct {
	Data *identityUser.UserData
}

func (u *User) Login(creds *hmac.HmacCredentials) error {
	err := errors.New("")
	u.Data, err = retrieveVeracodeUser(creds)
	if err != nil {
		panic(err)
	}
	return err
}

func (u *User) Validate() error {
	if u.Data == nil || u.Data.UserID == "" || u.Data.UserName == "" {
		return errors.New("HMAC credentials not associated with a valid Veracode account.")
	}
	return nil
}

func (u *User) Logout(apiId string, apiSecretKey string) {

}

func retrieveVeracodeUser(creds *hmac.HmacCredentials) (*identityUser.UserData, error) {
	var userData identityUser.UserData

	if creds.Id == "" || creds.Secret == "" {
		// return empty userData
		return nil, nil
	}

	// Format the API URL.
	apiUrl := url.URL{
		Scheme: "https",
		Host:   "api.veracode.com",
		Path:   "/api/authn/v2/users/self",
	}
	httpMethod := "GET"

	authHeader, err := hmac.CalculateAuthorizationHeader(&apiUrl, httpMethod, creds)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(httpMethod, apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization": {authHeader},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(body, &userData)
	} else {
		err = errors.New("Non-200 status code when authenticating")
		return nil, err
	}

	// return if userID and username is empty
	if &userData.UserID == nil && &userData.UserName == nil {
		err = errors.New("Empty user ID and username found when authenticating.")
		return nil, err
	}

	return &userData, nil
}
