package globalvar

var ApiKey string
var ApiSecret string

func SetApiKey(apiKey string) {
	ApiKey = apiKey
}

func GetApiKey() string {
	return ApiKey
}

func SetApiSecret(apiSecret string) {
	ApiSecret = apiSecret
}

func GetApiSecret() string {
	return ApiSecret
}
