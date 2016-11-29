package backends

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Promql struct {
	HttpClient *http.Client
	Server     string
	Path       string
	Query      string
	Time       string
	Start      string
	End        string
	Step       string
}

const (
	QUERYRANGEPATH = "/api/v1/query_range"
	QUERYPATH      = "/api/v1/query"
	RULESPATH      = "/api/v1/rules"
)

func (pm Promql) GetPromqlQuery() (map[string]interface{}, error) {
	response, request, err := pm.getResponse()
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("Failed to parse the response from %s", request)
		return nil, err
	}
	return result, nil
}

func (pm Promql) getResponse() ([]byte, string, error) {

	u, err := url.Parse(pm.Server)
	if err != nil {
		return nil, "", err
	}
	u.Path = strings.TrimRight(u.Path, "/") + pm.Path
	if pm.Path == RULESPATH {
	} else {
		q := u.Query()
		q.Set("query", pm.Query)
		if pm.Path == QUERYPATH {
			q.Set("time", pm.Time)
		} else if pm.Path == QUERYRANGEPATH {
			q.Set("start", pm.Start)
			q.Set("end", pm.End)
			q.Set("step", pm.Step)
		}
		u.RawQuery = q.Encode()
	}

	resp, err := pm.HttpClient.Get(u.String())
	if err != nil {
		err = fmt.Errorf("Failed to get response from %s", u.String())
		return nil, u.String(), err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, u.String(), err
	}

	return body, u.String(), nil
}
