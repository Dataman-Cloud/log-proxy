package service

import "github.com/Dataman-Cloud/log-proxy/src/models"

type LogSearchService interface {
	Applications(page models.Page) (map[string]int64, error)
	Slots(app string, page models.Page) (map[string]int64, error)
	Tasks(opts map[string]interface{}, page models.Page) (map[string]int64, error)
	Sources(opts map[string]interface{}, page models.Page) (map[string]int64, error)
	Search(keyword string, opts map[string]interface{}, page models.Page) (map[string]interface{}, error)
	Context(opts map[string]interface{}, page models.Page) ([]map[string]interface{}, error)
	Everything(key string, opt map[string]interface{}, page models.Page) (map[string]int64, error)
}
