package esclient

import (
	"fmt"

	"github.com/Dataman-Cloud/log-proxy/src/config"

	log "github.com/Sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v3"
)

var client *elastic.Client

func New(urls []string) (*elastic.Client, error) {
	var err error
	var ofs []elastic.ClientOptionFunc
	ofs = append(ofs, elastic.SetURL(urls...))
	ofs = append(ofs, elastic.SetMaxRetries(3))
	if config.GetConfig().SearchDebug {
		ofs = append(ofs, elastic.SetTraceLog(log.StandardLogger()))
	}
	client, err = elastic.NewClient(ofs...)
	if err != nil {
		return nil, fmt.Errorf("new elastic client error: %v", err)
	}
	log.Infof("successfully connected elasticsearch uri: %v", urls)

	return client, err
}

// GetClient returns the struct gorm.DB
func GetClient(urls []string) (*elastic.Client, error) {
	if client == nil {
		return New(urls)
	}
	return client, nil
}
