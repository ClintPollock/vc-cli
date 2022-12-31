package hmac

/*
 *
 * Pulled from work done by Zachary Estrella:
 *
 *   https://gitlab.laputa.veracode.io/roboyle/veracode-cli-hack
 *
 */
import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Included in the signature to inform Veracode of the signature version.
const veracodeRequestVersionString = "vcode_request_version_1"

// Expected format for the unencrypted data string.
const dataFormat = "id=%s&host=%s&url=%s&method=%s"

// Expected format for the Authorization header.
const headerFormat = "%s id=%s,ts=%s,nonce=%X,sig=%X"

// Expect prefix to the Authorization header.
const veracodeHMACSHA256 = "VERACODE-HMAC-SHA-256"

type HmacCredentials struct {
	Id     string
	Secret string
}

func CalculateAuthorizationHeader(url *url.URL, httpMethod string, creds *HmacCredentials) (string, error) {

	nonce, err := createNonce(16)

	if err != nil {
		return "", err
	}

	secret, err := fromHexString(creds.Secret)

	if err != nil {
		return "", err
	}

	timestampMilliseconds := strconv.FormatInt(time.Now().UnixNano()/int64(1000000), 10)
	data := fmt.Sprintf(dataFormat, creds.Id, url.Hostname(), url.RequestURI(), httpMethod)
	dataSignature := calculateSignature(secret, nonce, []byte(timestampMilliseconds), []byte(data))
	return fmt.Sprintf(headerFormat, veracodeHMACSHA256, creds.Id, timestampMilliseconds, nonce, dataSignature), nil
}

func createNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)

	_, err := rand.Read(nonce)

	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func fromHexString(input string) ([]byte, error) {
	decoded, err := hex.DecodeString(input)

	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func calculateSignature(key, nonce, timestamp, data []byte) []byte {
	encryptedNonce := hmac256(nonce, key)
	encryptedTimestampMilliseconds := hmac256(timestamp, encryptedNonce)
	signingKey := hmac256([]byte(veracodeRequestVersionString), encryptedTimestampMilliseconds)
	return hmac256(data, signingKey)
}

func hmac256(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}
