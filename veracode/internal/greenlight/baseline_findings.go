package greenlight

/*
 *
 * Static Pipeline basic flow implementation
 *
 */

import (
	"encoding/json"
	"io/ioutil"

	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
)

func LoadBaselineFindings(app greenlight_api.AppContext) (*greenlight_api.ScanFindings, error) {

	filename := app.BaselineFile

	if filename == "" {
		return nil, nil
	}

	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var findings greenlight_api.ScanFindings

	json.Unmarshal(content, &findings)

	return &findings, err

}
