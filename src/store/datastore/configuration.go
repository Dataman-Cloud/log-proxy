package datastore

import "github.com/Dataman-Cloud/log-proxy/src/models"

func (db *datastore) UpdateConf(conf *models.Configuration) error {
	return db.Model(conf).Where("id = 1").Update(conf).Error
}

func (db *datastore) GetConf() (models.Configuration, error) {
	var (
		conf models.Configuration
		err  error
	)
	err = db.Table("configurations").First(&conf).Error
	return conf, err
}

func (db *datastore) SetDefaultConf() error {
	conf := models.NewConfiguration()
	notfound := db.
		Where("id = 1").
		First(&conf).
		RecordNotFound()
	if notfound {
		return db.Save(conf).Error
	}
	return nil
}
