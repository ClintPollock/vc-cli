package initialize

import (
	"fmt"
	"os"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/docker/distribution/uuid"
	"github.com/spf13/viper"

	"github.com/veracode/veracode-cli/internal/api/identity/logging"
	identityUser "github.com/veracode/veracode-cli/internal/api/identity/user"
	"github.com/veracode/veracode-cli/internal/cache"
	"github.com/veracode/veracode-cli/internal/globalvar"
	"github.com/veracode/veracode-cli/internal/hmac"
	"github.com/veracode/veracode-cli/internal/verascanner"
)

func ValidateUserAndAddLogs(command string, subCommand []string,
	logUUID uuid.UUID, userData *identityUser.UserData) *identityUser.UserData {

	creds := hmac.HmacCredentials{globalvar.GetApiKey(), globalvar.GetApiSecret()}

	if creds.Id == "" || creds.Secret == "" {
		// return empty userData
		return nil
	}

	// Format the API URL.
	apiURL := url.URL{
		Scheme: "https",
		Host:   "api.veracode.com",
		Path:   "/api/authn/v2/users/self",
	}
	httpMethod := "GET"

	authHeader, err := hmac.CalculateAuthorizationHeader(&apiURL, httpMethod, &creds)
	if err != nil {
		return nil
	}

	client := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(httpMethod, apiURL.String(), nil)
	if err != nil {
		return nil
	}

	userAgentHeader := logging.BuildLogJson(command, subCommand, logUUID, userData)
	if userData == nil {
		req.Header = http.Header{
			"Authorization": {authHeader},
			"User-Agent":    {userAgentHeader},
		}
	} else {
		req.Header = http.Header{
			"Authorization": {authHeader},
			"User-Agent":    {userAgentHeader},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		if err := json.Unmarshal(body, &userData); err != nil {
			panic(err)
		}
	}
	// return if userID and username is empty
	if &userData.UserID == nil && &userData.UserName == nil {
		return nil
	}
	return userData
}

func ValidateWithCreds(command string, subCommand []string,
	logUUID uuid.UUID, userData *identityUser.UserData) *identityUser.UserData {
	// Verify user via Identity API.
	userData = ValidateUserAndAddLogs(command, subCommand, logUUID, userData)
	if userData == nil || userData.UserID == "" || userData.UserName == "" {
		fmt.Println("HMAC credentials not associated with a valid Veracode account.")
		os.Exit(0)
	}
	return userData
}

func Validate(command string, subCommand []string, logUUID uuid.UUID, userData *identityUser.UserData) *identityUser.UserData {
	if err := viper.ReadInConfig(); err == nil {
		// Should be able to marshal / unmarshal config values using mapstructure:
		// https://dev.to/techschoolguru/load-config-from-file-environment-variables-in-golang-with-viper-2j2d
		globalvar.SetApiKey(viper.GetString("default.veracode_api_key_id"))
		globalvar.SetApiSecret(viper.GetString("default.veracode_api_key_secret"))
		// Verify user via Identity API.
		userData = ValidateWithCreds(command, subCommand, logUUID, userData)
	} else {
		fmt.Println("No config file found. Use `veracode init` to set up Veracode HMAC credentials.")
		os.Exit(0)
	}
	viper.AutomaticEnv()
	return userData
}

func InitializeDocker() {
	// Pull veracode scanner image
	err := verascanner.ImagePull()
	if err != nil {
		panic(err)
	}
}

func InitializeCache() {
	_ = cache.Path()
}
