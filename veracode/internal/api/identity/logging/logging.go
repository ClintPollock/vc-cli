package logging

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/veracode/veracode-cli/cmd/version"
	identityUser "github.com/veracode/veracode-cli/internal/api/identity/user"
)

func BuildLogJson(command string, subCommand []string, logUUID uuid.UUID, userData *identityUser.UserData) string {
	jsonMap := map[string]string{
		"appName":      "MAERSK",
		"version":      version.Version,
		"os":           runtime.GOOS,
		"architecture": runtime.GOARCH,
		"command":      command,
		"subCommand":   strings.Join(subCommand, " "),
		"UUID":         logUUID.String(),
		"timestamp":    time.Now().Format("2006-01-02 15:04:05 MST"),
	}

	if userData != nil {
		jsonMap["status"] = fmt.Sprintf("%s success", command)
		jsonMap["userId"] = userData.UserID
		jsonMap["organizationId"] = userData.Organization.OrgID
	}
	jsonStr, _ := json.Marshal(jsonMap)

	return string(jsonStr)
}
