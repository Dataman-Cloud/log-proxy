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
	SILENCES_API = "/api/v1/silences"
	GET_SILENCE  = "/api/v1/silence/"
)

type AlertManager struct {
	HttpClient *http.Client
	Server     string
	Path       string
}

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

func (am AlertManager) getResponse() ([]byte, error) {
	u, err := url.Parse(am.Server)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimRight(u.Path, "/") + am.Path
	resp, err := am.HttpClient.Get(u.String())
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

func (am *AlertManager) GetSilences() ([]interface{}, error) {
	u, err := url.Parse(am.Server)
	if err != nil {
		return nil, err
	}
	u.Path = SILENCES_API
	resp, err := am.HttpClient.Get(u.String())
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

	fmt.Println(m["data"])

	return m["data"].(map[string]interface{})["silences"].([]interface{}), nil
}

func (am *AlertManager) CreateSilence(silence map[string]interface{}) error {
	u, err := url.Parse(am.Server)
	if err != nil {
		return err
	}
	u.Path = SILENCES_API

	data, _ := json.Marshal(silence)
	_, err = am.HttpClient.Post(u.String(), "application/json", bytes.NewBuffer(data))

	return err
}

func (am *AlertManager) GetSilence(id string) (map[string]interface{}, error) {
	u, err := url.Parse(am.Server)
	if err != nil {
		return nil, err
	}
	u.Path = GET_SILENCE + id
	resp, err := am.HttpClient.Get(u.String())
	body, _ := utils.ReadResponseBody(resp)
	var m map[string]interface{}

	json.Unmarshal(body, &m)

	return m["data"].(map[string]interface{}), nil
}

func (am *AlertManager) DeleteSilence(id string) error {
	u, err := url.Parse(am.Server)
	if err != nil {
		return err
	}
	u.Path = GET_SILENCE + id
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	_, err = am.HttpClient.Do(req)
	return err
}
