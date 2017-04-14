package datastore

import "github.com/Dataman-Cloud/log-proxy/src/models"

var camaConfig *models.Configuration

func (db *datastore) GetConfigFromDB() *models.Configuration {
	return camaConfig
}

func (db *datastore) UpdateConf(conf *models.Configuration) error {
	err := db.Model(conf).Where("id = 1").Update(conf).Error
	if err == nil {
		camaConfig = conf
	}
	return err
}

func (db *datastore) SetDefaultConf() error {
	var err error
	conf := models.NewConfiguration()
	notfound := db.
		Where("id = 1").
		First(&conf).
		RecordNotFound()
	if notfound {
		err = db.Save(conf).Error
	}
	camaConfig = conf
	return err
}
