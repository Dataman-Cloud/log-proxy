package backends

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
