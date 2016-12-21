package backends

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/utils"
)

const (
	// SilencesAPI api path
	SilencesAPI = "/api/v1/silences"
	// GetSilence get silence api path
	GetSilence = "/api/v1/silence/"
	// AlertsPath alerts api path
	AlertsPath = "/api/v1/alerts"
	// AlertsGroupsPath alert group api path
	AlertsGroupsPath = "/api/v1/alerts/groups"
	// AlertsStatusPath alert status api path
	AlertsStatusPath = "/api/v1/status"
)

// AlertManager define the query info of AlertManager
type AlertManager struct {
	HTTPClient *http.Client
	Server     string
	Path       string
}

// GetAlertManagerResponse return the result from func getResponse
func (am AlertManager) GetAlertManagerResponse() (map[string]interface{}, error) {
	response, err := am.getResponse()
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("Failed to parse the response from %s%s", am.Server, am.Path)
		return nil, err
	}
	return result, nil
}

// getResponse return the response from AlertManager
func (am AlertManager) getResponse() ([]byte, error) {
	u, err := url.Parse(am.Server)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimRight(u.Path, "/") + am.Path
	resp, err := am.HTTPClient.Get(u.String())
	if err != nil {
		err = fmt.Errorf("Failed to get response from %s", u.String())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetSilences return the list of Silences from AlertManager
func (am *AlertManager) GetSilences() ([]interface{}, error) {
	u, err := url.Parse(am.Server)
	if err != nil {
		return nil, err
	}
	u.Path = SilencesAPI
	resp, err := am.HTTPClient.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := utils.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return m["data"].(map[string]interface{})["silences"].([]interface{}), nil
}

// CreateSilence post the json silence to AlertManager
func (am *AlertManager) CreateSilence(silence map[string]interface{}) error {
	u, err := url.Parse(am.Server)
	if err != nil {
		return err
	}
	u.Path = SilencesAPI

	data, _ := json.Marshal(silence)
	_, err = am.HTTPClient.Post(u.String(), "application/json", bytes.NewBuffer(data))

	return err
}

// GetSilence return the silence from AlertManager
func (am *AlertManager) GetSilence(id string) (map[string]interface{}, error) {
	u, err := url.Parse(am.Server)
	if err != nil {
		return nil, err
	}
	u.Path = GetSilence + id
	resp, err := am.HTTPClient.Get(u.String())
	if err != nil {
		return nil, err
	}

	body, _ := utils.ReadResponseBody(resp)
	var m map[string]interface{}

	json.Unmarshal(body, &m)
	if m["data"] == nil {
		return nil, nil
	}

	return m["data"].(map[string]interface{}), nil
}

// DeleteSilence send the DELETE request to AlertManager
func (am *AlertManager) DeleteSilence(id string) error {
	u, err := url.Parse(am.Server)
	if err != nil {
		return err
	}
	u.Path = GetSilence + id
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	_, err = am.HTTPClient.Do(req)
	return err
}
