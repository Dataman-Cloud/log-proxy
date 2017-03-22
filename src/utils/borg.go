package utils

import (
	"encoding/json"

	"github.com/Dataman-Cloud/log-proxy/src/utils/httpclient"

	marathon "github.com/gambol99/go-marathon"
	"golang.org/x/net/context"
)

type BorgAccount struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type BorgResponse struct {
	Code int64           `json:"code"`
	Data json.RawMessage `json:"data"`
}

func LoginBorg(user, pwd, loginUrl string) (string, error) {
	borgAccount := BorgAccount{user, pwd}
	result, err := httpclient.NewDefaultClient().POST(context.TODO(), loginUrl, nil, borgAccount, nil)
	if err != nil {
		return "", err
	}

	var borgResponse BorgResponse
	if err := json.Unmarshal(result, &borgResponse); err != nil {
		return "", err
	}

	var token string
	if err := json.Unmarshal(borgResponse.Data, &token); err != nil {
		return "", err
	}

	return token, nil
}

func ListAppsFromBorg(token, appsUrl string) (*marathon.Applications, error) {
	header := make(map[string][]string)
	header["Authorization"] = []string{token}

	result, err := httpclient.NewDefaultClient().GET(context.TODO(), appsUrl, nil, header)
	if err != nil {
		return nil, err
	}

	var borgResponse BorgResponse
	if err := json.Unmarshal(result, &borgResponse); err != nil {
		return nil, err
	}

	var apps marathon.Applications
	if err := json.Unmarshal(borgResponse.Data, &apps); err != nil {
		return nil, err
	}

	return &apps, nil
}

func ListAppTaskFromBorg(token, appUrl string) ([]*marathon.Task, error) {
	header := make(map[string][]string)
	header["Authorization"] = []string{token}

	result, err := httpclient.NewDefaultClient().GET(context.TODO(), appUrl, nil, header)
	if err != nil {
		return nil, err
	}

	var borgResponse BorgResponse
	if err := json.Unmarshal(result, &borgResponse); err != nil {
		return nil, err
	}

	var app marathon.Application
	if err := json.Unmarshal(borgResponse.Data, &app); err != nil {
		return nil, err
	}

	return app.Tasks, nil
}
