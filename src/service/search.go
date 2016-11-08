package service

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v3"
)

type SearchService struct {
	ESClient *elastic.Client
}

func GetSearchService() *SearchService {
	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		log.Errorf("new elastic client error: %v", err)
		return nil
	}

	return &SearchService{
		ESClient: client,
	}
}

func (s *SearchService) Applications() {
	info, err := s.ESClient.Count().Do()
	log.Infof("-------------: %d %v", info, err)
}
