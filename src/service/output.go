package service

import (
	"encoding/json"
	"net"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"

	log "github.com/Sirupsen/logrus"
)

func SendCamaEvent(event *models.CamaEvent) {
	camaAddr := config.GetConfig().CamaNotifactionADDR
	tcpAddr, err := net.ResolveTCPAddr("tcp4", camaAddr)
	if err != nil {
		log.Errorf("resolve tcp addr failed. Error: %v", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Errorf("DialTCP to %s failed. Error %s: ", tcpAddr.String(), err)
		return
	}

	defer conn.Close()

	message, err := json.Marshal(event)
	if err != nil {
		log.Errorf("Marshal event failed. Error %s: ", err)
		return
	}

	if _, err := conn.Write(message); err != nil {
		log.Errorf("write message to conn failed. Error %s: ", err)
		return
	}
	log.Infof("sent alert message to %s", camaAddr)

	return
}
