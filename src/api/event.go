package api

import (
	"bufio"
	"net/http"
	"net/url"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	EVENT_URL = "/v2/events"
)

func (s *search) GetEvents(ctx *gin.Context) {

	result, err := s.Service.GetEvents(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_EVENTS_ERROR, err))
		return
	}

	utils.Ok(ctx, result)
}

func (s *search) ReceiverMarathonEvent() {
	if config.GetConfig().MARATHON_URL == "" {
		log.Info("not setting marathon recive url")
		return
	}
	u, err := url.Parse(config.GetConfig().MARATHON_URL)
	if err != nil {
		log.Errorf("invalid marathon url: %s", config.GetConfig().MARATHON_URL)
		return
	}
	u.Path = EVENT_URL

Again:
	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Accept", `text/event-stream`)
	resp, err := client.Do(req)
	if err != nil {
		time.Sleep(time.Second * 3)
		goto Again
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		s.Service.SaveFaildEvent(scanner.Bytes())
	}

	goto Again
}
