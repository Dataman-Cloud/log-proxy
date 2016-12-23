# Logging/Monitoring Proxy API Guide

## Monitoring

### Metric(监控数据)


#### Get the metric values (CPU/Memory/Network/Filesystem)

```GET /v1/monitor/query```
- path: /v1/monitor/query
- HTTP Method: GET
- URL Params: Null
- Query Params: metric, clusterid, appid, taskid, start, end, step, expr
  - metric=[cpu/memory/memory_usage/memory_total/network_rx/network_tx/fs_read/fs_write]: the metric string.
  - appid=<string>: the name of application.
  - taskid=<string>: the id string of the docker instance.
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
http://127.0.0.1:5098/v1/monitor/query?start=2016-12-02T00:00:01.781Z&end=2016-12-02T00:01:00.781Z&step=30s&metric=memory_total&appid=work-web
```

#### Get the info of clusters, cluster, application

```GET /v1/monitor/info```

- path: /v1/monitor/applications
- HTTP Method: GET
- URL Params: Null
- Query Params: clusterid, appid
  - clusterid=<string>: the name of cluster.
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
http://127.0.0.1:5098/v1/monitor/info?appid=work-web
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

#### 获取配置的报警规则
`GET /v1/monitor/alerts/rules`
```
curl -XGET http://127.0.0.1:5098/v1/monitor/alerts/rules
{
    "data":[
        "ALERT cpu_usage IF irate(container_cpu_usage_seconds_total{id=~"/docker/.*",name=~"mesos.*"}[5m]) * 100 > 80 FOR 1m LABELS {severity="critical"} ANNOTATIONS {description="High CPU usage on {{ $labels.name }} of App {{ $labels.container_label_APP_ID }}", summary="CPU Usage on {{ $labels.name }} of App {{ $labels.container_label_APP_ID }}"}",
        "ALERT mem_usage IF container_memory_usage_bytes{id=~"/docker/.*",name=~"mesos.*"} / container_spec_memory_limit_bytes{id=~"/docker.*",name=~"mesos.*"} * 100 > 80 FOR 1m LABELS {severity="critical"} ANNOTATIONS {description="High Mem usage on {{ $labels.name }} of App {{ $labels.container_label_APP_ID }}", summary="Mem Usage on {{ $labels.name }} of App {{ $labels.container_label_APP_ID }}"}",
        "ALERT InstanceDown IF up == 0 FOR 5m LABELS {severity="page"} ANNOTATIONS {description="{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes.", summary="Instance {{ $labels.instance }} down"}"
    ],
    "status":"success"
}
```

## 日志

### 获取所有应用
`GET /v1/search/applications`

For example:
```
curl -XGET http://localhost:5098/v1/search/applications
```

- Query Params:from, to
 - from=1478769333000
 - to=1478769333000

return

```
{
    "code": 0,
    "data": {
        "cluster1-maliao": 79273,
        "cluster1-proxytest": 78595,
        "cluster1-test": 88599
    }
}
```

### 根据应用获取所有实例
`GET /v1/search/tasks/:appid`

For example:

```
curl -XGET http://localhost:5098/v1/search/tasks/test
```

- URL Params: appid
 - appid=test
- Query Params:from, to
 - from=1478769333000
 - to=1478769333000

 return

```
{
    "code": 0,
    "data": {
        "cluster1-maliao.e43f85e6-9f7f-11e6-8313-02421fb0085c": 79273
    }
}
```

### 根据应用实例获取所有日志来源
`GET /v1/search/paths/:appid

<span id="itm1">For example:</span>
```
curl -XGET http://localhost:5098/v1/search/paths/appid
```

- URL Params: appid
 - appid=test
 - taskid=taskid
- Query Params:from, to
 - from=1478769333000
 - to=1478769333000

 return

```
{
    "code": 0,
    "data": {
        "stderr": 5,
        "stdout": 79268
    }
}
```

### 日志搜索

`GET /v1/search/index`

```
http://192.168.1.46:5098/v1/search/index?appid=cluster1-maliao&from=now-7d&taskid=cluster1-maliao.e43f85e6-9f7f-11e6-8313-02421fb0085c,cluster1-maliao.e43f85e6-9f7f-11e6-8313-02421fb0085c&path=stdout,stderr&keyword=container
```

- Query Params: appid,taskid,path,keyword,from,to
 - appid=test
 - taskid=tasktest
 - path=stdout
 - keyword=test
 - from=now-7d
 - to=now

return

```
{
"code": 0,
"data": [
{
"@timestamp": "2016-11-08T10:27:56.759Z",
"@version": "1",
"appid": "cluster1-maliao",
"clusterid": "cluster1",
"groupid": "9",
"host": "192.168.1.71",
"id": "75145e5517b7f2038e39012bc471db59bd8c7dab1b5779075603fdb452fbac27",
"message": "--container=\"mesos-c62c27ef-c144-4a38-b9fb-684794919bc7-S5.6b08bb27-409c-49fa-9606-7ed56e9d8366\" --docker=\"docker\" --docker_socket=\"/var/run/docker.sock\" --help=\"false\" --initialize_driver_logging=\"true\" --launcher_dir=\"/usr/libexec/mesos\" --logbufsecs=\"0\" --logging_level=\"INFO\" --mapped_directory=\"/mnt/mesos/sandbox\" --quiet=\"false\" --sandbox_directory=\"/data/mesos/slaves/c62c27ef-c144-4a38-b9fb-684794919bc7-S5/frameworks/c62c27ef-c144-4a38-b9fb-684794919bc7-0000/executors/cluster1-maliao.e43f85e6-9f7f-11e6-8313-02421fb0085c/runs/6b08bb27-409c-49fa-9606-7ed56e9d8366\" --stop_timeout=\"0ns\"\n",
"offset": 1,
"path": "stdout",
"port": 39426,
"taskid": "cluster1-maliao.e43f85e6-9f7f-11e6-8313-02421fb0085c",
"time": "2016-11-09T02:16:09.405026732+08:00",
"userid": "23"
}
]
}
```

`GET /v1/search/context?appid=x&taskid=x&path=x&offset=x`

return

```
{
"code": 0,
"data": [
{
"@timestamp": "2016-11-08T10:27:56.759Z",
"@version": "1",
"appid": "cluster1-maliao",
"clusterid": "cluster1",
"groupid": "9",
"host": "192.168.1.71",
"id": "75145e5517b7f2038e39012bc471db59bd8c7dab1b5779075603fdb452fbac27",
"message": "--container=\"mesos-c62c27ef-c144-4a38-b9fb-684794919bc7-S5.6b08bb27-409c-49fa-9606-7ed56e9d8366\" --docker=\"docker\" --docker_socket=\"/var/run/docker.sock\" --help=\"false\" --initialize_driver_logging=\"true\" --launcher_dir=\"/usr/libexec/mesos\" --logbufsecs=\"0\" --logging_level=\"INFO\" --mapped_directory=\"/mnt/mesos/sandbox\" --quiet=\"false\" --sandbox_directory=\"/data/mesos/slaves/c62c27ef-c144-4a38-b9fb-684794919bc7-S5/frameworks/c62c27ef-c144-4a38-b9fb-684794919bc7-0000/executors/cluster1-maliao.e43f85e6-9f7f-11e6-8313-02421fb0085c/runs/6b08bb27-409c-49fa-9606-7ed56e9d8366\" --stop_timeout=\"0ns\"\n",
"offset": 1,
"path": "stdout",
"port": 39426,
"taskid": "cluster1-maliao.e43f85e6-9f7f-11e6-8313-02421fb0085c",
"time": "2016-11-09T02:16:09.405026732+08:00",
"userid": "23"
}
]
}
```


GET /v1/search/keyword

return

```
{"code":0,"data":{"count":4,"results":[{"id":"AVkHBbEQIIGpJqE63UXA","appid":"dsgsdg","keyword":"gdsgsdg","path":"gsdsdgsd","createtime":"2016-12-16T17:45:30.639081077+08:00"},{"id":"AVkV4YfXIIGpJqE63U2b","appid":"6566y6","keyword":"y65y56","path":"y56y6","createtime":"2016-12-19T15:00:19.024077007+08:00"},{"id":"AVkHBV59IIGpJqE63UW_","appid":"sxacsacs111","keyword":"scsacsac","path":"csacasc","createtime":"2016-12-19T17:42:25.092717794+08:00"},{"id":"AVkAeS_LIIGpJqE63UH7","appid":"work-nginxefef","keyword":"GET11","path":"stdout","createtime":"2016-12-19T17:43:11.950277076+08:00"}]}}

```

POST /v1/search/keyword -d '{"period":1,"appid":"test","keyword":"keyword","condition":1,"enable":true}'

return 

```
{"code":0,"data":"create success"}
```


PUT /v1/search/keyword -d '{"id":"x","period":1,"appid":"test","keyword":"keyword","condition":1,"enable":true}'

return
 
```
{"code":0,"data":"update success"}
```


DELETE /v1/search/keyword/:id

return

```
{"code":0,"data":"delete success"}
```

GET /v1/search/keyword/:id

return

```
{"code":0,"data":{"id":"AVkAeS_LIIGpJqE63UH7","appid":"work-nginxefef","keyword":"GET11","path":"stdout","createtime":"2016-12-19T17:43:11.950277076+08:00"}}

```


GET /v1/search/prometheus

return

```
{"code":0,"data":{"count":3301,"results":[{"alertname":"DatamanServiceDown","annotations":{"description":"mesos-master of node srymaster1 has been down for more than 1 minutes.","summary":"DatamanService mesos-master down"},"createtime":"2016-12-20T11:14:38.095324531+08:00","endsAt":"2016-12-20T11:13:52.807+08:00","generatorURL":"http://srymaster2:9090/graph#%5B%7B%22expr%22%3A%22consul_catalog_service_node_healthy%7Bservice%21%3D%5C%22alertmanager-vip%5C%22%2Cservice%21%3D%5C%22mysql-vip%5C%22%7D%20%3D%3D%200%22%2C%22tab%22%3A0%7D%5D","id":"AVkaOUWQIIGpJqE63VOl","labels":"{\"alertname\":\"DatamanServiceDown\",\"instance\":\"192.168.1.92:9107\",\"job\":\"consul\",\"node\":\"srymaster1\",\"service\":\"mesos-master\",\"severity\":\"Warning\"}","startsAt":"2016-12-20T11:11:22.807+08:00","status":"resolved"}]}}

```

GET /v1/search/prometheus/:id

return

```
{"code":0,"data":{"alertname":"DatamanServiceDown","annotations":{"description":"mesos-master of node srymaster1 has been down for more than 1 minutes.","summary":"DatamanService mesos-master down"},"createtime":"2016-12-20T11:14:38.095324531+08:00","endsAt":"2016-12-20T11:13:52.807+08:00","generatorURL":"http://srymaster2:9090/graph#%5B%7B%22expr%22%3A%22consul_catalog_service_node_healthy%7Bservice%21%3D%5C%22alertmanager-vip%5C%22%2Cservice%21%3D%5C%22mysql-vip%5C%22%7D%20%3D%3D%200%22%2C%22tab%22%3A0%7D%5D","labels":"{\"alertname\":\"DatamanServiceDown\",\"instance\":\"192.168.1.92:9107\",\"job\":\"consul\",\"node\":\"srymaster1\",\"service\":\"mesos-master\",\"severity\":\"Warning\"}","startsAt":"2016-12-20T11:11:22.807+08:00","status":"resolved"}}

```

GET /v1/search/mointor

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
