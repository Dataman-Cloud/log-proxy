package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Byte2str array byte parse to string
func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ParseDate judge date is same day
func ParseDate(from, to interface{}) string {
	if from == nil || to == nil {
		return "*"
	}

	f, err := strconv.ParseInt(fmt.Sprint(from), 10, 64)
	if err != nil {
		return time.Now().Format("2006-01-02")
	}

	t, err := strconv.ParseInt(fmt.Sprint(to), 10, 64)
	if err != nil {
		return time.Now().Format("2006-01-02")
	}

	if time.Unix(f/1000, 0).Format("2006-01-02") != time.Unix(t/1000, 0).Format("2006-01-02") {
		return "*"
	}

	return time.Now().Format("2006-01-02")
}

// ReadRequestBody read request body
func ReadRequestBody(request *http.Request) ([]byte, error) {
	defer request.Body.Close()
	return ioutil.ReadAll(request.Body)
}

// ReadResponseBody read response body
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// AlertNotification alert notification interface
func AlertNotification(url string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	return err
}

func ParseTask(taskid string) []string {
	if !strings.Contains(taskid, "-") || len(strings.Split(taskid, "-")) != 2 {
		return strings.Split(taskid, ",")
	}

	taskRange := strings.Split(taskid, "-")
	lower, err := strconv.Atoi(taskRange[0])
	if err != nil {
		return strings.Split(taskid, ",")
	}
	upper, err := strconv.Atoi(taskRange[1])
	if err != nil {
		return strings.Split(taskid, ",")
	}

	var tasks []string
	for lower <= upper {
		tasks = append(tasks, strconv.Itoa(lower))
		lower++
	}
	return tasks
}
