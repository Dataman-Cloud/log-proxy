# Logging/Monitoring Proxy API Guide

## Monitoring

### Metric(监控数据)

#### Get the metric items

```GET /v1/monitor/query/metrics```
- path: /v1/monitor/query/metrics
- HTTP Method: GET
- URL Params: Null
- Query Params: Null

For example:

Get the metrics list
```
"http://127.0.0.1:5098/v1/monitor/query/metrics"
```

#### Get the clusters list

```GET /v1/monitor/clusters```
- path: /v1/monitor/clusters
- HTTP Method: GET
- URL Params: Null
- Query Params: Null

For example:

Get the metrics list
```
"http://127.0.0.1:5098/v1/monitor/clusters"
```

#### Get the app list by cluster

```GET /v1/monitor/clusters/:clusterid/apps```
- path: /v1/monitor/clusters/:clusterid/apps
- HTTP Method: GET
- URL Params:
  - clusterid:=<string>: vcluster ID
- Query Params: Null

For example:

Get the metrics list
```
"http://127.0.0.1:5098/v1/monitor/clusters/work/apps"
```

#### Get the tasks list by app

```GET /v1/monitor/clusters/:clusterid/apps/:appid/tasks```
- path: /v1/monitor/cluster/:clusterid/apps/:appid/tasks
- HTTP Method: GET
- URL Params:
  - clusterid:=<string>: vcluster ID
  - appid:=<string>: app ID
- Query Params: Null

For example:

Get the metrics list
```
"http://127.0.0.1:5098/v1/monitor/clusters/work/apps/work-web/tasks"
```

#### Get the metric values

```GET /v1/monitor/query```
- path: /v1/monitor/query
- HTTP Method: GET
- URL Params: Null
- Query Params: metric, clusterid, appid, taskid, start, end, step, expr
  - metric=<string>: the list from the api /v1/monitor/query/metrics
  - cluster=<string>: the name of vcluster.
  - app=<string>: the name of application.
  - task=<string>: the id string of the docker instance.
  - start=<2016-12-02T00:00:01.781Z>: the start time of the query range.
  - end=<2016-12-02T00:00:01.781Z>: the end time of the query range.
  - step=<duration>: Query resolution step width.
  - expr=<string>: Prometheus expression query string. it is conflict with metric.

For example:

Get the metrics by Expr, Refer to "https://prometheus.io/docs/querying/api/"
```
http://127.0.0.1:5098/v1/monitor/query?expr=avg(irate(container_cpu_usage_seconds_total{container_label_VCLUSTER='work', container_label_APP_ID='work-web', container_env_mesos_task_id='work-web.dfb9bf02-0f03-11e7-8697-02421d201838', id=~'/docker/.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID, container_env_mesos_task_id) keep_common
```

Get the metrics by URL
```
http://127.0.0.1:5098/v1/monitor/query?metric=cpu_usage&app=work-web&cluster=work&task=work-web.dfb9bf02-0f03-11e7-8697-02421d201838
```

## Alert

### 报警规则

#### 获取报警规则中报警指标

```GET /v1/alert/indicators```
- path: /v1/alert/indicators
- HTTP Method: GET
- URL Params: Null
- Query Params: Null

For example:

```
http://127.0.0.1:5098/v1/alert/indicators
```
return
```
{
  "code": 0,
  "data": [
    "cpu_usage",
    "mem_usage"
  ]
}
```

#### 新建报警规则

```POST /v1/alert/rules```
- path: /v1/alert/rules
- HTTP Method: POST
- URL Params: Null
- Query Params: Null

For example:

创建规则，创建的时候不要加status这个字段的值
```
curl -X POST "http://127.0.0.1:5098/v1/alert/rules" -d '{
       "class": "mola",
       "name": "cpuUsage",
       "cluster": "work",
       "app": "work-nginx",
       "pending": "1m",
       "indicator": "mem_usage",
       "severity": "warning",
       "aggregation": "max",
       "comparison": ">",
       "Threshold": 60,
  }  
  return
  {
  "code": 0,
  "data": {
    "ID": 1,
    "CreatedAt": "2017-03-27T10:56:42+08:00",
    "UpdatedAt": "2017-03-27T10:56:41.725919919+08:00",
    "class": "mola",
    "name": "cpuUsage",
    "status": "Enabled",
    "cluster": "work",
    "app": "work-nginx",
    "user": "",
    "user_group": "",
    "pending": "1m",
    "duration": "5m",
    "indicator": "mem_usage",
    "severity": "warning",
    "aggregation": "max",
    "comparison": ">",
    "threshold": 60
    }
  }
```
#### 获取报警规则列表
``` GET /v1/alert/rules```
- path: /v1/alert/rules
- HTTP Method: GET
- URL Params:
- Query Params: class, cluster, app

For example:

- 获取所有规则：
```
http://127.0.0.1:5098/v1/alert/rules
```
- 获取某个分类的规则
```
http://127.0.0.1:5098/v1/alert/rules?class=mola
```
- 获取某个分类的某个集群规则
```
http://127.0.0.1:5098/v1/alert/rules?class=mola&cluster=work
```
- 获取某个分类的某个集群的某个应用规则
```
http://127.0.0.1:5098/v1/alert/rules?class=mola&cluster=work&app=work-nginx
```
return
```
{
  "code": 0,
  "data": {
    "count": 1,
    "rules": [
      {
        "ID": 1,
        "CreatedAt": "2017-03-27T10:56:42+08:00",
        "UpdatedAt": "2017-03-27T10:56:42+08:00",
        "class": "mola",
        "name": "cpuUsage",
        "status": "Enabled",
        "cluster": "work",
        "app": "work-nginx",
        "user": "",
        "user_group": "",
        "pending": "1m",
        "duration": "5m",
        "indicator": "mem_usage",
        "severity": "warning",
        "aggregation": "max",
        "comparison": ">",
        "threshold": 60
      }
    ]
  }
}
```

#### 获取单个报警规则

```GET /v1/alert/rules/:id```
- path: /v1/alert/rules/:id
- HTTP Method: GET
- URL Params:
  - id=<string>: the ID of the rule
- Query Params: Null

For example:

Get one alert rule by id
```
"http://127.0.0.1:5098/v1/alert/rules/1"
```
return
```
{
  "code": 0,
  "data": {
    "ID": 1,
    "CreatedAt": "2017-03-27T10:56:42+08:00",
    "UpdatedAt": "2017-03-27T10:56:42+08:00",
    "class": "mola",
    "name": "cpuUsage",
    "status": "Enabled",
    "cluster": "work",
    "app": "work-nginx",
    "user": "",
    "user_group": "",
    "pending": "1m",
    "duration": "5m",
    "indicator": "mem_usage",
    "severity": "warning",
    "aggregation": "max",
    "comparison": ">",
    "threshold": 60
  }
}
```

#### 更新报警规则

```PUT /v1/alert/rules/:id```
- path: /v1/alert/rules/:id
- HTTP Method: DELETE
- URL Params:
  - id=<string>: the ID of the rule
- Query Params: Null

For example:

更新报警规则
```
curl -X PUT "http://127.0.0.1:5098/v1/alert/rules/1" -d '{
       "class": "mola",
       "name": "cpuUsage",
       "cluster": "work",
       "app": "work-nginx",
       "pending": "1m",
       "indicator": "mem_usage",
       "severity": "warning",
       "aggregation": "max",
       "comparison": ">",
       "Threshold": 60
  }'
```
更新报警规则状态，报警规则分三种状态:
- Uninit 未初始化
- Enabled 活跃
- Disabled 暂停

```
curl -X PUT "http://127.0.0.1:5098/v1/alert/rules/1" -d '{
       "class": "mola",
       "name": "cpuUsage",
       "cluster": "work",
       "app": "work-nginx",
       "status": "Disabled"
  }'
```
return
```
{
  "code": 0,
  "data": {
    "ID": 1,
    "CreatedAt": "2017-03-27T10:56:42+08:00",
    "UpdatedAt": "2017-03-27T11:17:00+08:00",
    "class": "mola",
    "name": "cpuUsage",
    "status": "Disabled",
    "cluster": "work",
    "app": "work-nginx",
    "user": "",
    "user_group": "",
    "pending": "1m",
    "duration": "5m",
    "indicator": "mem_usage",
    "severity": "warning",
    "aggregation": "max",
    "comparison": ">",
    "threshold": 60
  }
}
```

#### 删除报警规则

```DELETE /v1/alert/rules/:id```
- path: /v1/alert/rules/:id
- HTTP Method: DELETE
- URL Params: Null
- Query Params: Null

For example:

Delete the alert rule

```
curl -X DELETE "http://127.0.0.1:5098/v1/alert/rules/1" -d '{
  "class" : "mola",
}'
```
return
```
{
  "code": 0,
  "data": "success"
}
```

### 报警事件

#### 获取报警事件

```GET /v1/alert/events```
- path: /v1/alert/events
- HTTP Method: GET
- URL Params: Null
- Query Params: Null

For example:

- 获取所有事件
```
"http://127.0.0.1:5098/v1/alert/events"
```
- 获取所有ACK事件
```
"http://127.0.0.1:5098/v1/alert/events?ack=true"
```
- 获取集群的所有事件
```
http://127.0.0.1:5098/v1/alert/events?cluster=work
```
- 获取集群的应用的所有事件
```
http://127.0.0.1:5098/v1/alert/events?cluster=work&app=work-nginx
```
return
```
{
  "code": 0,
  "data": {
    "count": 1,
    "events": [
      {
        "ID": 1,
        "CreatedAt": "2017-03-27T11:21:38+08:00",
        "UpdatedAt": "2017-03-27T11:21:38+08:00",
        "DeletedAt": null,
        "count": 1,
        "severity": "warning",
        "cluster": "work",
        "app": "work-nginx",
        "task": "work-nginx.cb7844e5-1172-11e7-b3da-0242a9f0feba",
        "user_name": "",
        "group_name": "",
        "container_id": "/docker/bce3a8f0dcb9496fc169be4499e1451d7d04a385eb99b924dc9abe99f5849607",
        "container_name": "mesos-b069b31c-cc93-4c12-992e-026c51913d06-S0.169a52a6-f007-457a-a0cd-ce9b7155e81e",
        "alert_name": "mola_cpuUsage_work_nginx",
        "ack": false,
        "value": "98.91357421875",
        "description": "",
        "summary": ""
      }
    ]
  }
}
```

### 设置报警事件ACK

```PUT /v1/alert/events/:id```
- path: /v1/alert/events/
- HTTP Method: PUT
- URL Params:
  - id=<string>, event id
- Query Params: Null

For example:

Ack the event
```
curl -X PUT "http://127.0.0.1:5098/v1/alert/events/1" -d '{
  "action":"ack",
	"cluster": "work",
	"app": "work-nginx"}'
```
return
```
{
  "code": 0,
  "data": {
    "status": "success"
  }
}
```
