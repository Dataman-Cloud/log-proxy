# Logging/Monitoring Proxy API Guide

## Monitoring

### Metric(监控数据)


#### Get the metric values (CPU/Memory/Network/Filesystem)

```GET /v1/monitor/query```
- path: /v1/monitor/query
- HTTP Method: GET
- URL Params: Null
- Query Params: metric, clusterid, userid, appid, taskid, start, end, step, expr
  - metric=[cpu/memory/memory_usage/memory_total/network_rx/network_tx/fs_read/fs_write]: the metric string.
  - clusterid=<string>: the name of cluster.
  - userid=<string>: the name of user.
  - appid=<string>: the name of application.
  - taskid=<number>: the id string of the task instance, support format as "1,2,3" and "1-3"
  - start=<2016-12-02T00:00:01.781Z>: the start time of the query range.
  - end=<2016-12-02T00:00:01.781Z>: the end time of the query range.
  - step=<duration>: Query resolution step width.
  - expr=<string>: Prometheus expression query string. it is conflict with metric.

For example:

Get the metrics by Expr, Refer to "https://prometheus.io/docs/querying/api/"
```
http://127.0.0.1:5098/v1/monitor/query?expr=avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID='work-web',id=~'/docker/.*',name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)&start=2016-12-05T00:00:01.781Z&end=2016-12-05T00:01:00.781Z&step=30s
```

Get the metrics by URL
```
http://127.0.0.1:5098/v1/monitor/query?start=1483942403&end=1483942403&step=30s&metric=memory&appid=nginx0051-xcm-datamanmesos&clusterid=datamanmesos&userid=xcm&taskid=1-3
```

#### Get the info of clusters, cluster, application

```GET /v1/monitor/info```

- path: /v1/monitor/applications
- HTTP Method: GET
- URL Params: Null
- Query Params: clusterid, appid
  - clusterid=<string>: the name of cluster.
  - userid=<string>: the name of user.
  - appid=<string>: the name of application.

For example:
Get the info of Clusters
```
http://127.0.0.1:5098/v1/monitor/info
```
Get the info of cluster
```
http://127.0.0.1:5098/v1/monitor/info?clusterid=work
```
Get the info of application
```
http://127.0.0.1:5098/v1/monitor/info?clusterid=datamanmesos&userid=xcm&appid=nginx0051-xcm-datamanmesos
```

#### Get the metric data of nodes

```GET /v1/monitor/nodes```

- path: /v1/monitor/nodes
- HTTP Method: GET
- URL Params: Null
- Query Params: nodeid
  - nodeid=<string>: the IP address of node.

For example:
Get the metric data of all nodes
```
http://127.0.0.1:5098/v1/monitor/nodes
```
Get the metric data of one node
```
http://127.0.0.1:5098/v1/monitor/nodes?nodeid=192.168.1.101
```

### AlertManager API

#### 获取当前活动的报警
`GET /v1/monitor/alerts`

#### 获取当前活动的报警，分组
`GET /v1/monitor/alerts/groups`

#### 获取AlertManager的状态信息
`GET /v1/monitor/alerts/status`

## 日志

### 获取某个时间段内有日志的所有集群

```
 GET /v1/log/clusters
```

For example:

```
http://192.168.59.3:5098/v1/log/clusters?from=1451577600000&to=1490096284000
```

- Query Params:from, to
 - 起始时间: from=1478769333000
 - 结束时间: to=1478769333000

return

```
{
  "code": 0,
  "data": {
    "yaoyun": 84
  }
}
```

### 获取指定集群在指定时间段内有日志的所有应用

```
GET /v1/log/clusters/:cluster/apps
```

For example

```
http://192.168.59.3:5098/v1/log/clusters/yaoyun/apps
```

- URL Params: cluster
 - cluster: 集群ID
- Query Params:from, to
 - 起始时间: from=1478769333000
 - 结束时间: to=1478769333000
 
 return
 
 ```
 {
  "code": 0,
  "data": {
    "yaoyun-nginx": 73,
    "yaoyun-nginx2": 11
  }
}
 ```
 
### 获取指定集群和应用在指定时间段内有日志的所有实例

```
http://192.168.59.3:5098/v1/log/clusters/:cluser/apps/:app/tasks
```

For example:

```
curl -XGET http://192.168.59.3:5098/v1/log/clusters/yaoyun/apps/yaoyun-nginx/tasks
```

- URL Params: cluster, app
 - cluster: 集群ID
 - app: appID
- Query Params:from, to
 - 起始时间: from=1478769333000
 - 结束时间: to=1478769333000


 return

```
{
  "code": 0,
  "data": [
    {
      "id": "yaoyun-nginx.26f179b7-0a8f-11e7-a287-02427ffc9690",
      "status": "running",
      "logCount": 7
    },
    {
      "id": "yaoyun-nginx.383c3939-0a33-11e7-b151-0242933964c8",
      "status": "died",
      "logCount": 9
    },
    {
      "id": "yaoyun-nginx.71fe7d3b-0a34-11e7-b151-0242933964c8",
      "status": "died",
      "logCount": 9
    },
    {
      "id": "yaoyun-nginx.26ff5c69-0a8f-11e7-a287-02427ffc9690",
      "status": "died",
      "logCount": 8
    },
    {
      "id": "yaoyun-nginx.45f8be4a-0a8f-11e7-a287-02427ffc9690",
      "status": "died",
      "logCount": 8
    },
    {
      "id": "yaoyun-nginx.be47fa68-0a32-11e7-b151-0242933964c8",
      "status": "died",
      "logCount": 9
    }
  ]
}
```

* 说明: status 标记日志产生的实例是否正在运行


### 查询指定集群和应用在指定时间段内有日志的来源
```/clusters/:cluster/apps/:app/sources```

- URL Params: cluster,app
 - cluster: 集群ID
 - app: appid
- Query Params:task, from, to
 -task: taskid
 - 起始时间: from=1478769333000
 - 结束时间: to=1478769333000

for example
```
curl -XGET http://192.168.59.3:5098/v1/log/clusters/yaoyun/apps/yaoyun-nginx/source?task=yaoyun-nginx.26f179b7-0a8f-11e7-a287-02427ffc9690
```

return

```
{
  "code": 0,
  "data": {
    "stdout": 21
  }
}
```

### 日志搜索

`GET /v1/log/clusters/:cluster/apps/:app/index`

- URL Params: cluster,app
 - cluster: 集群ID
 - app: appid
- Query Params: task,source,keyword,from,to
 - task=tasktest
 - source=stdout
 - keyword=test
 - from=now-7d
 - to=now

```
curl -XGET http://192.168.59.3:5098/v1/log/clusters/yaoyun/apps/yaoyun-nginx/index?task=yaoyun-nginx.67202d60-1041-11e7-aea5-024292f41627&keyword=Starting
```

return

```
{
  "code": 0,
  "data": {
    "count": 1,
    "results": [
      {
        "@timestamp": "2017-03-24T03:24:40.077Z",
        "@version": "1",
        "appid": "yaoyun-nginx",
        "clusterid": "yaoyun",
        "containerid": "3978677ccad9c22b2d7140a46abb81522cd15272f890d01d387e519e9bdaf573",
        "groupid": "1",
        "host": "192.168.59.104",
        "logtime": "2017-03-24T11:24:40.021974705+08:00",
        "message": "<mark>Starting</mark> task yaoyun-nginx.67202d60-1041-11e7-aea5-024292f41627\n",
        "offset": 1490325880021974800,
        "path": "stdout",
        "port": 48398,
        "taskid": "yaoyun-nginx.67202d60-1041-11e7-aea5-024292f41627",
        "userid": "1"
      }
    ]
  }
}
```

### 查询指定日志上下文

`GET /v1/log/clusters/:cluster/apps/:app/context`

- URL Params: cluster,app
 - cluster: 集群ID
 - app: appid
- Query Params: task,source, offset,from,to
 - task=tasktest
 - source=stdout
 - offset=1489720851875141600(从index返回的接口中取)
 - from=now-7d
 - to=now


for example
```
curl -XGET http://192.168.59.3:5098/v1/log/clusters/yaoyun/apps/yaoyun-nginx2/context?offset=1489720851875141600&page=1&source=stdout&size=100&task=yaoyun-nginx2.719226da-0a34-11e7-b151-0242933964c8
```

return

```
{
  "code": 0,
  "data": [
    {
      "@timestamp": "2017-03-17T15:52:59.657Z",
      "@version": "1",
      "appid": "yaoyun-nginx2",
      "clusterid": "yaoyun",
      "containerid": "461d7a408b55b97aa401ff22f24a3365422928840069f345c7b57d3ef2cf895e",
      "groupid": "1",
      "host": "192.168.59.104",
      "logtime": "2017-03-17T23:52:59.609401568+08:00",
      "message": "192.168.59.3 - - [17/Mar/2017:15:52:55 +0000] \"GET / HTTP/1.1\" 200 612 \"http://192.168.59.103:5013/ui/app/yaoyun-nginx2/instance?from=my\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36\" \"-\"\n",
      "offset": 1489765979609401600,
      "path": "stdout",
      "port": 51698,
      "taskid": "yaoyun-nginx2.719226da-0a34-11e7-b151-0242933964c8",
      "userid": "1"
    },
    {
      "@timestamp": "2017-03-17T15:52:59.664Z",
      "@version": "1",
      "appid": "yaoyun-nginx2",
      "clusterid": "yaoyun",
      "containerid": "461d7a408b55b97aa401ff22f24a3365422928840069f345c7b57d3ef2cf895e",
      "groupid": "1",
      "host": "192.168.59.104",
      "logtime": "2017-03-17T23:52:59.609529309+08:00",
      "message": "192.168.59.3 - - [17/Mar/2017:15:52:56 +0000] \"GET / HTTP/1.1\" 304 0 \"http://192.168.59.103:5013/ui/app/yaoyun-nginx2/instance?from=my\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36\" \"-\"\n",
      "offset": 1489765979609529300,
      "path": "stdout",
      "port": 51698,
      "taskid": "yaoyun-nginx2.719226da-0a34-11e7-b151-0242933964c8",
      "userid": "1"
    }
 ]
}
```

## 日志报警

### 创建一条报警规则

```
POST http://192.168.59.3:5098/v1/log/rules -d 
'{
	"app":"nginx",
	"source":"stderr",
	"keyword":"ABBBAcccd"
}'
```
* 说明: app,source,keyword分别代表应用.日志来源(stdout,stderr,文件路径),关键字三个都是必填.

### 更新一条告警规则
```
PUT http://192.168.59.3:5098/v1/log/rules -d
'{
{
	"id":4,
	"app":"nginx",
	"source":"stderr",
	"keyword":"ABBBAcccd"
}
}'
```

### 删除一条告警规则
```
DELETE http://192.168.59.3:5098/v1/log/rules/:id
```

### 获取指定的告警规则
```
GET http://192.168.59.3:5098/v1/log/rules/:id
```

return
```
{
  "code": 0,
  "data": {
    "id": 1,
    "app": "nginx",
    "keyword": "ABBBA",
    "source": "stderr",
    "createdAt": "2017-03-23T14:54:57+08:00",
    "updatedAt": "2017-03-23T15:46:45+08:00"
  }
}
```

### 获取告警规则列表
```
GET http://192.168.59.3:5098/v1/log/rules?page=1&size=50
```

return
```
{
  "code": 0,
  "data": {
    "count": 2,
    "rules": [
      {
        "id": 1,
        "app": "nginx",
        "keyword": "ABBBA",
        "source": "stderr",
        "createdAt": "2017-03-23T14:54:57+08:00",
        "updatedAt": "2017-03-23T15:46:45+08:00"
      },
      {
        "id": 5,
        "app": "nginx",
        "keyword": "ABBBAcccd",
        "source": "stderr",
        "createdAt": "2017-03-23T17:36:03+08:00",
        "updatedAt": "2017-03-23T17:36:03+08:00"
      }
    ]
  }
}
```

```
{"code":0,"data":{"count":4,"results":[{"id":"AVkHBbEQIIGpJqE63UXA","appid":"dsgsdg","keyword":"gdsgsdg","path":"gsdsdgsd","createtime":"2016-12-16T17:45:30.639081077+08:00"},{"id":"AVkV4YfXIIGpJqE63U2b","appid":"6566y6","keyword":"y65y56","path":"y56y6","createtime":"2016-12-19T15:00:19.024077007+08:00"},{"id":"AVkHBV59IIGpJqE63UW_","appid":"sxacsacs111","keyword":"scsacsac","path":"csacasc","createtime":"2016-12-19T17:42:25.092717794+08:00"},{"id":"AVkAeS_LIIGpJqE63UH7","appid":"work-nginxefef","keyword":"GET11","path":"stdout","createtime":"2016-12-19T17:43:11.950277076+08:00"}]}}

```

POST /v1/search/keyword -d '{"period":1,"appid":"test","keyword":"keyword","condition":1,"enable":true}'



GET /v1/log/prometheus

return

```
{"code":0,"data":{"count":3301,"results":[{"alertname":"DatamanServiceDown","annotations":{"description":"mesos-master of node srymaster1 has been down for more than 1 minutes.","summary":"DatamanService mesos-master down"},"createtime":"2016-12-20T11:14:38.095324531+08:00","endsAt":"2016-12-20T11:13:52.807+08:00","generatorURL":"http://srymaster2:9090/graph#%5B%7B%22expr%22%3A%22consul_catalog_service_node_healthy%7Bservice%21%3D%5C%22alertmanager-vip%5C%22%2Cservice%21%3D%5C%22mysql-vip%5C%22%7D%20%3D%3D%200%22%2C%22tab%22%3A0%7D%5D","id":"AVkaOUWQIIGpJqE63VOl","labels":"{\"alertname\":\"DatamanServiceDown\",\"instance\":\"192.168.1.92:9107\",\"job\":\"consul\",\"node\":\"srymaster1\",\"service\":\"mesos-master\",\"severity\":\"Warning\"}","startsAt":"2016-12-20T11:11:22.807+08:00","status":"resolved"}]}}

```

GET /v1/log/prometheus/:id

return

```
{"code":0,"data":{"alertname":"DatamanServiceDown","annotations":{"description":"mesos-master of node srymaster1 has been down for more than 1 minutes.","summary":"DatamanService mesos-master down"},"createtime":"2016-12-20T11:14:38.095324531+08:00","endsAt":"2016-12-20T11:13:52.807+08:00","generatorURL":"http://srymaster2:9090/graph#%5B%7B%22expr%22%3A%22consul_catalog_service_node_healthy%7Bservice%21%3D%5C%22alertmanager-vip%5C%22%2Cservice%21%3D%5C%22mysql-vip%5C%22%7D%20%3D%3D%200%22%2C%22tab%22%3A0%7D%5D","labels":"{\"alertname\":\"DatamanServiceDown\",\"instance\":\"192.168.1.92:9107\",\"job\":\"consul\",\"node\":\"srymaster1\",\"service\":\"mesos-master\",\"severity\":\"Warning\"}","startsAt":"2016-12-20T11:11:22.807+08:00","status":"resolved"}}

```

GET /v1/log/mointor

GET /v1/monitor/alerts/groups

GET /v1/monitor/silences

return

```
{"code":0,"data":[{"comment":"121","createdAt":"2016-12-19T17:51:51.513126826+08:00","createdBy":"1@1","endsAt":"2016-12-19T09:56:16Z","id":4,"matchers":[{"isRegex":false,"name":"alertname","value":"DatamanServiceDown"},{"isRegex":false,"name":"instance","value":"192.168.1.91:9107"},{"isRegex":false,"name":"job","value":"consul"},{"isRegex":false,"name":"node","value":"srymaster1"},{"isRegex":false,"name":"service","value":"zookeeper"},{"isRegex":false,"name":"severity","value":"Warning"}],"startsAt":"2016-12-19T07:56:16Z"},{"comment":"qgqw","createdAt":"2016-12-19T16:48:55.764252+08:00","createdBy":"adf@qfwf","endsAt":"2016-12-19T08:40:17Z","id":3,"matchers":[{"isRegex":false,"name":"node","value":"srymaster1"},{"isRegex":false,"name":"service","value":"prometheus"},{"isRegex":false,"name":"severity","value":"Warning"}],"startsAt":"2016-12-19T06:40:17Z"},{"comment":"asdfasdf","createdAt":"2016-12-19T11:08:14.816592362+08:00","createdBy":"test@123.com","endsAt":"2016-12-19T07:12:00Z","id":2,"matchers":[{"isRegex":false,"name":"alertname","value":"LogKeyword"},{"isRegex":false,"name":"appid","value":"work-nginx"},{"isRegex":false,"name":"clusterid","value":"work"},{"isRegex":false,"name":"instance","value":"192.168.1.75:5098"},{"isRegex":false,"name":"job","value":"log-proxy"},{"isRegex":false,"name":"keyword","value":"GET"},{"isRegex":false,"name":"offset","value":"1481781258185649664"},{"isRegex":false,"name":"path","value":"stdout"},{"isRegex":false,"name":"severity","value":"Warning"},{"isRegex":false,"name":"taskid","value":"work-nginx.1f17a9f0-c02b-11e6-9030-024245dc84c8"},{"isRegex":false,"name":"userid","value":"4"}],"startsAt":"2016-12-19T03:12:00Z"}]}

```


GET /v1/monitor/silence/:id

return

```
{"code":0,"data":{"comment":"121","createdAt":"2016-12-19T17:51:51.513126826+08:00","createdBy":"1@1","endsAt":"2016-12-19T09:56:16Z","id":4,"matchers":[{"isRegex":false,"name":"alertname","value":"DatamanServiceDown"},{"isRegex":false,"name":"instance","value":"192.168.1.91:9107"},{"isRegex":false,"name":"job","value":"consul"},{"isRegex":false,"name":"node","value":"srymaster1"},{"isRegex":false,"name":"service","value":"zookeeper"},{"isRegex":false,"name":"severity","value":"Warning"}],"startsAt":"2016-12-19T07:56:16Z"}}

```

DELETE /v1/monitor/silence/:id

return

```
{"code":0,"data":"delete success"}
```

PUT /v1/monitor/silence/:id

return

```
{"code":0,"data":"update success"}
```


POST /v1/monitor/silences -d '{"matchers":[{"name":"alertname","value":"cpu_usage","isRegex":false},{"name":"container_label_APP_ID","value":"work-web","isRegex":false},{"name":"container_label_VCLUSTER","value":"work","isRegex":false},{"name":"cpu","value":"cpu01","isRegex":false},{"name":"id","value":"/docker/e4a59106d9f763626c70fc5e6bfa8c46a14a471ad66f102618a142f1e2f3e33f","isRegex":false},{"name":"image","value":"192.168.1.75/library/nginx-stress:1.10","isRegex":false},{"name":"instance","value":"192.168.1.102:5014","isRegex":false},{"name":"job","value":"cadvisor","isRegex":false},{"name":"name","value":"mesos-f011c830-7bb0-4edf-8c25-f9e64fa2246a-S0.3836a388-864b-48dc-95dd-25119b9ee3fa","isRegex":false},{"name":"severity","value":"critical","isRegex":false}],"startsAt":"2016-12-05T11:08:00.000Z","endsAt":"2016-12-05T15:08:00.000Z","createdBy":"yqguo@dataman-inc.com","comment":"this is a test"}'
