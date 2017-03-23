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

#### Get the cluster list

```GET /v1/monitor/cluster```
- path: /v1/monitor/cluster
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

```GET /v1/monitor/clusters/:clusterid/apps/:appid/task```
- path: /v1/monitor/cluster/:clusterid/apps/:appid/task
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

#### 新建报警规则

```POST /v1/alert/rules```
- path: /v1/alert/rules
- HTTP Method: POST
- URL Params: Null
- Query Params: Null

For example:

Create the alert rule
```
curl -X POST "http://127.0.0.1:5098/v1/alert/rules" -d '{
  "name" : "admin",
  "alert" : "cpu",
  "if": "irate(container_cpu_usage_seconds_total{id=~\"/docker/.*\",name=~\"mesos.*\"}[5m]) * 100 > 80",
  "for" : "1m",
  "labels" : "{ service=\"nginx-app\",severity=\"Warning\" }",
  "description" : "High CPU usage on {{ $labels.name }}",
  "summary" :  "The value is {{ $value }}"
}'
```

#### 获取报警规则

```GET /v1/alert/rules/:id```
- path: /v1/alert/rules
- HTTP Method: GET
- URL Params:
  - id=<string>: the ID of the rule
- Query Params: Null

For example:

Get the alert rules
```
"http://127.0.0.1:5098/v1/alert/rules"
```

Get one alert rule by id
```
"http://127.0.0.1:5098/v1/alert/rules/1"
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
  "name" : "admin",
}'
```

#### 更新报警规则

```PUT /v1/alert/rules/```
- path: /v1/alert/rules/
- HTTP Method: DELETE
- URL Params: Null
- Query Params: Null

For example:

Delete the alert rule
```
curl -X PUT "http://127.0.0.1:5098/v1/alert/rules" -d '{
  "name" : "admin",
  "alert" : "cpu",
  "if": "irate(container_cpu_usage_seconds_total{id=~\"/docker/.*\",name=~\"mesos.*\"}[5m]) * 100 > 80",
  "for" : "1m",
  "labels" : "{ service=\"nginx-app\",severity=\"Warning\" }",
  "description" : "High CPU usage on {{ $labels.name }}",
  "summary" :  "The value is {{ $value }}"
}'
```

### 报警事件

#### 获取报警事件

```GET /v1/alert/events```
- path: /v1/alert/events
- HTTP Method: GET
- URL Params: Null
- Query Params: Null

For example:

Get the alert events
```
"http://127.0.0.1:5098/v1/alert/events"
```

### 设置报警事件ACK

```PUT /v1/alert/events/:id```
- path: /v1/alert/events/
- HTTP Method: PUT
- URL Params:
  - id=<string>, event id
- Query Params: Null

For example:

Delete the alert rule
```
curl -X PUT "http://127.0.0.1:5098/v1/alert/events/1" -d '{
  "action": "ack"
}'
```
