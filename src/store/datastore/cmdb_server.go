package datastore

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func (db *datastore) CreateCmdbServer(cmdb *models.CmdbServer) error {
	var result models.CmdbServer
	if db.Where("cmdb_servers.app_id = ?", cmdb.AppID).First(&result).RecordNotFound() {
		return db.Create(cmdb).Error
	}

	return db.Model(cmdb).Updates(cmdb).Error
}

func (db *datastore) GetCmdbServer(appID string) (*models.CmdbServer, error) {
	var cmdb models.CmdbServer
	var camaConfig = db.GetConfigFromDB()
	if db.Where("cmdb_servers.app_id = ?", appID).First(&cmdb).RecordNotFound() {
		return &models.CmdbServer{appID, camaConfig.CamaCmdbDefaultAppID}, nil
	}

	if err := db.Where("cmdb_servers.app_id = ?", appID).First(&cmdb).Error; err != nil {
		return nil, err
	}

	return &cmdb, nil
}
