package service

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	var head string
	msgLength := len(message)
	switch {
	case msgLength < 10:
		head = fmt.Sprintf("000%d", msgLength)
	case msgLength >= 10 && msgLength < 100:
		head = fmt.Sprintf("00%d", msgLength)
	case msgLength >= 100 && msgLength < 1000:
		head = fmt.Sprintf("0%d", msgLength)
	case msgLength >= 1000 && msgLength < 10000:
		head = fmt.Sprintf("%d", msgLength)
	default:
		head = "9999"
	}

	byteArray := [][]byte{
		[]byte(head),
		message,
	}
	content := bytes.Join(byteArray, nil)

	if _, err := conn.Write(content); err != nil {
		log.Errorf("write message to conn failed. Error %s: ", err)
		return
	}
	log.Infof("sent alert message to %s with content %s", camaAddr, string(content))

	return
}
