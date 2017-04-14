package service

import (
	"encoding/json"
	"net"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"

	log "github.com/Sirupsen/logrus"
)

func SendCamaEvent(event *models.CamaEvent) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.GetConfig().CamaNotifactionADDR)
	if err != nil {
		log.Error("resolve tcp addr failed. Error: ", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Error("DialTCP to %s failed. Error %s: ", tcpAddr.String(), err)
		return
	}

	defer conn.Close()

	message, err := json.Marshal(event)
	if err != nil {
		log.Error("Marshal event failed. Error %s: ", err)
		return
	}

	if _, err := conn.Write(message); err != nil {
		log.Error("write message to conn failed. Error %s: ", err)
		return
	}

	return
}
