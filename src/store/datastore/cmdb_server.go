package datastore

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) CreateCmdbServer(cmdb *models.CmdbServer) error {
	var result models.CmdbServer
	notfound := db.Where("cmdb_servers.app_id = ?", cmdb.AppID).First(&result).RecordNotFound()
	if notfound {
		return db.Create(cmdb).Error
	}

	return db.Model(cmdb).Updates(cmdb).Error
}

func (db *datastore) GetCmdbServer(appID string) (*models.CmdbServer, error) {
	var cmdb models.CmdbServer
	if err := db.Where("cmdb_servers.appid = ?", appID).First(&cmdb).Error; err != nil {
		return nil, err
	}

	return &cmdb, nil
}
