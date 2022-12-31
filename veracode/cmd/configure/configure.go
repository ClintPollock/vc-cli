/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package configure

import (
	"bufio"
	"context"

	//"flag"
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/spf13/viper"
	"github.com/veracode/veracode-cli/internal/hmac"
	veracodeUser "github.com/veracode/veracode-cli/internal/user"
)

func Run(ctx context.Context, args []string, flags *flag.FlagSet) error {

	//creds, err := cmd.ReadCredentials(ctx)

	creds := hmac.HmacCredentials{}

	// Prompt for new values.
	creds.Id = promptForString("API ID", viper.GetString("credentials.veracode_api_key_id"))
	creds.Secret = promptForString("API Secret Key", viper.GetString("credentials.veracode_api_key_secret"))

	// Try to login with these credentials
	var u veracodeUser.User
	err := u.Login(&creds)
	if err != nil {
		panic(err)
	}
	err = u.Validate()

	fmt.Printf("Validated credentials for user %s (%s)\n", u.Data.UserName, u.Data.UserID)
	if err != nil {
		panic(err)
	}

	viper.Set("credentials.veracode_api_key_id", creds.Id)
	viper.Set("credentials.veracode_api_key_secret", creds.Secret)

	// Setup defaults
	viper.SetDefault("ui.prettyprint", "true")
	viper.SetDefault("ui.verbose", "false")

	viper.SetDefault("urls.scheme", "https")
	viper.SetDefault("urls.host", "api.veracode.com")
	viper.SetDefault("urls.pipeline_api_path", "/pipeline_scan/v1")
	viper.SetDefault("urls.events_api_path", "/events/v1")
	viper.SetDefault("urls.policy_api_path", "/appsec/v1/policies")

	viper.SetDefault("ignoreauth", "false")

	// Save
	viper.WriteConfig()
	viper.AutomaticEnv()

	return nil
}

func promptForString(valueName string, defaultValue string) string {
	fmt.Print(valueName + " " + "[" + defaultValue + "]" + " ")
	var r = bufio.NewReader(os.Stdin)
	value, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultValue
	}
	return value
}

func Set(ctx context.Context, args []string, flags *flag.FlagSet) {

	key := strings.TrimSpace(args[0])
	value := strings.TrimSpace(args[1])

	// sanitization
	key = strings.TrimRight(key, ":")

	viper.Set(key, value)
	viper.WriteConfig()

}

func Delete(ctx context.Context, args []string, flags *flag.FlagSet) {

	key := strings.TrimSpace(args[0])

	viper.Set(key, nil) //Delete(key)
	viper.WriteConfig()
}
