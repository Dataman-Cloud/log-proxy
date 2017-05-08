# Logging/Monitoring Proxy API Guide

## Monitoring

### Metric(监控数据)

#### Get the apps

```
GET /v1/monitor/apps
```
- path: /v1/monitor/apps
- HTTP Method: GET
- URL Params: Null
- Query Params: Null
For example:
```
http://127.0.0.1:5098/v1/monitor/apps
```
return
```
{
  "code": 0,
  "data": [
    "web-zdou-datamanmesos"
  ]
}
```

#### Get the tasks of app

```
GET /v1/monitor/apps/:appid/tasks
```
- path: /v1/monitor/apps/:appid/tasks
- HTTP Method: GET
- URL Params: Null
- Query Params: Null
For example:
```
http://127.0.0.1:5098/v1/monitor/apps/web-zdou-datamanmesos/tasks
```
return
```
{
  "code": 0,
  "data": [
    "0"
  ]
}
```

#### Get the metric values

```
GET /v1/monitor/query/items
```
- path: /v1/monitor/query
- HTTP Method: GET
- URL Params: Null
- Query Params: Null

For example:
```
http://127.0.0.1:5098/v1/monitor/query/items
```
return
```
{
  "code": 0,
  "data": [
    "CPU使用率",
    "内存使用字节数",
    "内存分配字节数"
  ]
}
```


```
GET /v1/monitor/query
```
- path: /v1/monitor/query
- HTTP Method: GET
- URL Params: Null
- Query Params: metric, app, task, start, end, step, expr
  - metric=[从API /v1/monitor/query/items拿到的列表]: the metric string.
  - app=<string>: the name of application.
  - task=<number>: the id string of the app slot ID
  - start=<1493708502>: the start time of the query range.
  - end=<1493708502>: the end time of the query range.
  - step=<duration>: Query resolution step width.
  - expr=<string>: Prometheus expression query string. it is conflict with metric.

For example:

Get the metrics by Expr, Refer to "https://prometheus.io/docs/querying/api/"

```
http://127.0.0.1:5098/v1/monitor/query?expr=avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID='work-web',id=~'/docker/.*',name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)&start=1493708502&end=1493708502&step=30s
```

Get the metrics by URL
```
http://127.0.0.1:5098/v1/monitor/query?start=1483942403&end=1483942403&step=30s&metric=CPU使用率&app=web-zdou-datamanmesos&task=0
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
  "data": {
    "CPU使用百分比": "%",
    "Tomcat线程数": "",
    "内存使用百分比": "%"
  }
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
    "group": "dev",
    "app": "web-zdou-datamanmesos",
    "indicator": "CPU使用百分比",
    "severity": "warning",
    "pending": "2m",
    "aggregation": "max",
    "comparison": ">",
    "threshold": 60  
}  
return
{
  "code": 0,
  "data": {
    "ID": 1,
    "CreatedAt": "2017-05-05T10:40:46+08:00",
    "UpdatedAt": "2017-05-05T10:40:45.662149083+08:00",
    "name": "web_zdou_datamanmesos_cpu_usage_warning",
    "group": "dev",
    "app": "web-zdou-datamanmesos",
    "severity": "warning",
    "indicator": "CPU使用百分比",
    "status": "Enabled",
    "pending": "2m",
    "duration": "5m",
    "aggregation": "max",
    "comparison": ">",
    "threshold": 60,
    "unit": "%"
  }
}
```

#### 获取报警规则列表
``` GET /v1/alert/rules```
- path: /v1/alert/rules
- HTTP Method: GET
- URL Params:
- Query Params: group, app

For example:

- 获取所有规则：
```
http://127.0.0.1:5098/v1/alert/rules
```
- 获取某个组的规则
```
http://127.0.0.1:5098/v1/alert/rules?group=dev
```
- 获取某个应用的规则
```
http://127.0.0.1:5098/v1/alert/rules?group=dev&app=web-zdou-datamanmesos
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
        "CreatedAt": "2017-05-05T12:21:29+08:00",
        "UpdatedAt": "2017-05-05T12:21:29+08:00",
        "name": "web_zdou_datamanmesos_cpu_usage_warning",
        "group": "dev",
        "app": "web-zdou-datamanmesos",
        "severity": "warning",
        "indicator": "CPU使用百分比",
        "status": "Enabled",
        "pending": "2m",
        "duration": "5m",
        "aggregation": "max",
        "comparison": ">",
        "threshold": 60,
        "unit": "%"
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
http://127.0.0.1:5098/v1/alert/rules/1
```
return
```
{
  "code": 0,
  "data": {
    "ID": 1,
    "CreatedAt": "2017-05-05T12:21:29+08:00",
    "UpdatedAt": "2017-05-05T12:21:29+08:00",
    "name": "web_zdou_datamanmesos_cpu_usage_warning",
    "group": "dev",
    "app": "web-zdou-datamanmesos",
    "severity": "warning",
    "indicator": "CPU使用百分比",
    "status": "Enabled",
    "pending": "2m",
    "duration": "5m",
    "aggregation": "max",
    "comparison": ">",
    "threshold": 60,
    "unit": "%"
  }
}
```

#### 更新报警规则

```PUT /v1/alert/rules/:id```
- path: /v1/alert/rules/:id
- HTTP Method: PUT
- URL Params:
  - id=<string>: the ID of the rule
- Query Params: Null

For example:

更新报警规则
```
curl -X PUT "http://127.0.0.1:5098/v1/alert/rules/1" -d '{
    "group": "dev",
    "pending": "1m",
    "aggregation": "max",
    "comparison": ">",
    "threshold": 60,
    "status": "Enabled"  
    }'
```
更新报警规则状态，报警规则分三种状态:
- Uninitialized 未初始化
- Enabled 活跃
- Disabled 暂停

```
curl -X PUT "http://127.0.0.1:5098/v1/alert/rules/1" -d '{
    "group": "Undefine",
    "pending": "1m",
    "aggregation": "max",
    "comparison": ">",
    "threshold": 60,
    "status": "Enabled"  
    }'
```
return
```
{
  "code": 0,
  "data": {
    "ID": 1,
    "CreatedAt": "2017-05-04T14:01:36+08:00",
    "UpdatedAt": "2017-05-05T10:13:59+08:00",
    "name": "work_nginx_mem_usage_warning",
    "group": "Undefine",
    "app": "work-nginx",
    "severity": "warning",
    "indicator": "内存使用百分比",
    "status": "Enabled",
    "pending": "1m",
    "duration": "5m",
    "aggregation": "max",
    "comparison": ">",
    "threshold": 60,
    "unit": "%"
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
curl -X DELETE "http://127.0.0.1:5098/v1/alert/rules/1?group=dev"

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
- 获取用户组的报警事件
```
http://127.0.0.1:5098/v1/alert/events?group=dev
```
- 获取应用的所有事件
```
http://127.0.0.1:5098/v1/alert/events?group=dev&app=web-zdou-datamanmesos&ack=true
```
- 按时间获取事件
```
http://127.0.0.1:5098/v1/alert/events?start=1490660541&end=1490660542
```
return
```
{
  "code": 0,
  "data": {
    "count": 1,
    "events": [
      {
        "ID": 13,
        "CreatedAt": "2017-05-06T21:10:02+08:00",
        "UpdatedAt": "2017-05-07T10:55:44+08:00",
        "DeletedAt": null,
        "alert_name": "web_zdou_datamanmesos_cpu_usage_warning",
        "group": "dev",
        "app": "web-zdou-datamanmesos",
        "task": "0",
        "severity": "warning",
        "indicator": "cpu_usage",
        "judgement": "max > 60%",
        "container_id": "/docker/4828b086a265a03ce4b34ac3662c6d67d5cb1a5230e1eb77a253b3fd382c1e19",
        "container_name": "mesos-3fca175f-7593-46ad-bcc6-2875e089f8b5-S0.c5a378af-7ea7-4cc7-9f80-de748cf0bbe9",
        "ack": true,
        "value": "0.01727926502595149",
        "description": "",
        "summary": "",
        "count": 823
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
	"group": "dev",
	"app": "web-zdou-datamanmesos"
  }'
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
## 日志

### 获取所有应用
`GET /v2/log/apps`

For example:
```
curl -XGET http://localhost:5098/v2/log/apps
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
`GET /v2/log/apps/:app/slots`

For example:

```
curl -XGET http://localhost:5098/v2/log/apps/test/slots
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
        "0": 79273
    }
}
```

### 根据应用和实例索引获取所有的实例ID
`GET /v2/log/apps/:app/slots/:slot/tasks`

For example

```
curl -XGET http://localhost:5098/v2/log/apps/aaaaaaaaaaa-yaoyun-datamanmesos/slots/0/tasks
```

return 

```
{
    "code": 0,
    "data": {
        "0-aaaaaaaaaaa-yaoyun-datamanmesos-1a4cde96aef34c03a8ed9f07dae5ec63": 79273
    }
}
```

### 根据应用实例获取所有日志来源
`GET /v2/log/apps/:app/sources`

for example

```
curl -XGET http://localhost:5098/v2/log/apps/aaaaaaaaaaa-yaoyun-datamanmesos/sources?slot=0&task=test
```

- URL Params: app
 - app=test
- Query Params: slot, task, from, to
 - slot=test
 - task=test
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

`GET /v2/log/apps/:app/search"`

```
curl -XGET http://localhost:5098/v2/log/apps/aaaaaaaaaaa-yaoyun-datamanmesos/search\?slot\=0\&keyword\=GET
```
- URL Params: app
 - app=test
- Query Params: slot,task,source,keyword,from,to
 - slot=test
 - task=tasktest
 - source=stdout
 - keyword=test
 - from=now-7d
 - to=now

return

```
{
    "code":0,
    "data":{
        "count":11,
        "results":[
            {
                "@timestamp":"2017-05-03T02:17:11.874Z",
                "@version":"1",
                "DM_APP_ID":"aaaaaaaaaaa-yaoyun-datamanmesos",
                "DM_APP_NAME":"aaaaaaaaaaa",
                "DM_CLUSTER":"datamanmesos",
                "DM_GROUP_NAME":"yaoyun",
                "DM_SLOT_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos",
                "DM_SLOT_INDEX":"0",
                "DM_TASK_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos-1a4cde96aef34c03a8ed9f07dae5ec63",
                "DM_USER":"yaoyun",
                "DM_USER_NAME":"yaoyun",
                "DM_VCLUSTER":"ymola",
                "containerid":"2101477298d0622d2fb81156ec6ad328c778035dfc528c475f8e1d6908ece819",
                "host":"192.168.1.102",
                "logtime":"2017-05-03T02:17:11.83244967Z",
                "message":"192.168.1.146 - - [03/May/2017:02:17:08 +0000] "<mark>GET</mark> / HTTP/1.1" 304 0 "http://localhost:5013/ui/app",
                "offset":1493777831832449800,
                "path":"stdout",
                "port":44940
            },
            {
                "@timestamp":"2017-05-03T02:20:06.855Z",
                "@version":"1",
                "DM_APP_ID":"aaaaaaaaaaa-yaoyun-datamanmesos",
                "DM_APP_NAME":"aaaaaaaaaaa",
                "DM_CLUSTER":"datamanmesos",
                "DM_GROUP_NAME":"yaoyun",
                "DM_SLOT_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos",
                "DM_SLOT_INDEX":"0",
                "DM_TASK_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos-1a4cde96aef34c03a8ed9f07dae5ec63",
                "DM_USER":"yaoyun",
                "DM_USER_NAME":"yaoyun",
                "DM_VCLUSTER":"ymola",
                "containerid":"2101477298d0622d2fb81156ec6ad328c778035dfc528c475f8e1d6908ece819",
                "host":"192.168.1.102",
                "logtime":"2017-05-03T02:20:06.850765691Z",
                "message":"192.168.1.146 - - [03/May/2017:02:20:05 +0000] "<mark>GET</mark> / HTTP/1.1" 304 0 "http://localhost:5013/ui/app",
                "offset":1493778006850765600,
                "path":"stdout",
                "port":45096
            }
        ]
    }
}
```

* 注: keyword不为空是关键字搜索, 关键字为空是普通日志搜索

### 日志上下文搜索

`GET /v2/log/apps/:app/context`

For example
```
curl -XGET http://localhost:5098/v2/log/apps/aaaaaaaaaaa-yaoyun-datamanmesos/context\?offset\=1493801584377272300
```

return

```
{
    "code":0,
    "data":[
        {
            "@timestamp":"2017-05-03T02:17:11.874Z",
            "@version":"1",
            "DM_APP_ID":"aaaaaaaaaaa-yaoyun-datamanmesos",
            "DM_APP_NAME":"aaaaaaaaaaa",
            "DM_CLUSTER":"datamanmesos",
            "DM_GROUP_NAME":"yaoyun",
            "DM_SLOT_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos",
            "DM_SLOT_INDEX":"0",
            "DM_TASK_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos-1a4cde96aef34c03a8ed9f07dae5ec63",
            "DM_USER":"yaoyun",
            "DM_USER_NAME":"yaoyun",
            "DM_VCLUSTER":"ymola",
            "containerid":"2101477298d0622d2fb81156ec6ad328c778035dfc528c475f8e1d6908ece819",
            "host":"192.168.1.102",
            "logtime":"2017-05-03T02:17:11.83244967Z",
            "message":"192.168.1.146 - - [03/May/2017:02:17:08 +0000] "GET / HTTP/1.1" 304 0 "http://localhost:5013/ui/app/aaaaaaaaaaa-yaoyun-datamanmesos/instance" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36" "-" ",
            "offset":1493777831832449800,
            "path":"stdout",
            "port":44940
        },
        {
            "@timestamp":"2017-05-03T02:20:06.855Z",
            "@version":"1",
            "DM_APP_ID":"aaaaaaaaaaa-yaoyun-datamanmesos",
            "DM_APP_NAME":"aaaaaaaaaaa",
            "DM_CLUSTER":"datamanmesos",
            "DM_GROUP_NAME":"yaoyun",
            "DM_SLOT_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos",
            "DM_SLOT_INDEX":"0",
            "DM_TASK_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos-1a4cde96aef34c03a8ed9f07dae5ec63",
            "DM_USER":"yaoyun",
            "DM_USER_NAME":"yaoyun",
            "DM_VCLUSTER":"ymola",
            "containerid":"2101477298d0622d2fb81156ec6ad328c778035dfc528c475f8e1d6908ece819",
            "host":"192.168.1.102",
            "logtime":"2017-05-03T02:20:06.850765691Z",
            "message":"192.168.1.146 - - [03/May/2017:02:20:05 +0000] "GET / HTTP/1.1" 304 0 "http://localhost:5013/ui/app/aaaaaaaaaaa-yaoyun-datamanmesos/instance" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36" "-" ",
            "offset":1493778006850765600,
            "path":"stdout",
            "port":45096
        }
    ]
}
```

### 日志告警

#### 创建日志告警规则
`POST /log/alert/rules`

For example

```
curl -XPOST http://localhost:5098/v2/log/alert/rules -d
'
{
	"user": "yaoyun",
	"group":"yaoyun",
	"cluster":"ymola",
	"app":"aaaaaaaaaaa-yaoyun-datamanmesos-3",
	"keyword":"GET",
	"source":"stdout"
}
'
```

* 注: user为当前用户, group为当前所在用户组

return
```
{
	"code": 0,
	"data": {
		"app": "aaaaaaaaaaa-yaoyun-datamanmesos-3",
		"cluster": "ymola",
		"createdAt": "2017-05-04T11:41:33.889755527+08:00",
		"description": "",
		"group": "yaoyun",
		"id": 5,
		"keyword": "GET",
		"source": "stdout",
		"status": "Enabled",
		"updatedAt": "2017-05-04T11:41:33.889755527+08:00",
		"user": "yaoyun"
	}
}
```

#### 更新单挑日志告警规则
`PUT /v2/log/alert/rules/:id`
For example

```
curl -XPUT http://localhost:5098/v2/log/alert/rules/1 -d
'
{
	"user": "yaoyun",
	"group":"yaoyun",
	"cluster":"ymola",
	"app":"aaaaaaaaaaa-yaoyun-datamanmesos-3",
	"keyword":"GET",
	"source":"stdout"
}
'
```

return
```
{
	"code": 0,
	"data": {
		"app": "aaaaaaaaaaa-yaoyun-datamanmesos-3",
		"cluster": "ymola",
		"createdAt": "2017-05-04T11:41:33.889755527+08:00",
		"description": "",
		"group": "yaoyun",
		"id": 5,
		"keyword": "GET",
		"source": "stdout",
		"status": "Enabled",
		"updatedAt": "2017-05-04T11:41:33.889755527+08:00",
		"user": "yaoyun"
	}
}
```

#### 查询单条报警规则
`GET /v2/log/alert/rules/:id`

For example
```
curl -XPUT http://localhost:5098/v2/log/alert/rules/1
```

return
```
{
	"code": 0,
	"data": {
		"app": "yaoyun-nginx",
		"cluster": "yaoyun",
		"createdAt": "2017-04-19T11:50:51+08:00",
		"description": "sfvsfv威武是否vava",
		"group": "",
		"id": 1,
		"keyword": "werer",
		"source": "stderr",
		"status": "Enabled",
		"updatedAt": "2017-04-19T11:51:06+08:00",
		"user": ""
	}
}
```

#### 查询报警规则列表
`GET /v2/log/alert/rules`

For example
```
curl -XPUT http://localhost:5098/v2/log/alert/rules?group=test&from=111&to=2222
```

* group 为用户当前所在的组, from 和 to代表其实和结束时间戳, 单位为秒,分页参数使用之前默认的size和page

return
```
{
	"code": 0,
	"data": {
		"count": 5,
		"rules": [
			{
				"app": "yaoyun-nginx",
				"cluster": "yaoyun",
				"createdAt": "2017-04-19T11:50:51+08:00",
				"description": "sfvsfv威武是否vava",
				"group": "",
				"id": 1,
				"keyword": "werer",
				"source": "stderr",
				"status": "Enabled",
				"updatedAt": "2017-04-19T11:51:06+08:00",
				"user": ""
			},
			{
				"app": "nginx",
				"cluster": "mola",
				"createdAt": "2017-04-28T15:24:19+08:00",
				"description": "",
				"group": "yaoyun",
				"id": 2,
				"keyword": "xxx",
				"source": "GET",
				"status": "Enabled",
				"updatedAt": "2017-04-28T15:24:19+08:00",
				"user": "yaoyun"
			}
		]
	}
}
```

 #### 删除指定的报警规则
 `DELETE /v2/log/alert/rules/:id`


#### 获取报警事件的过滤条件 apps
`GET /v2/log/alert/apps"`

For example
```
curl -XGET http://localhost:5098/v2/log/alert/apps
```

return 
```
{
    "code":0,
    "data":[
        {
            "app":"aaaaaaaaaaa-yaoyun-datamanmesos"
        }
    ]
}
```

#### 查询报警事件
`GET /v2/log/alert/events`

For example
```
curl -XGET http://localhost:5098/v2/log/alert/events\?app\=aaaaaaaaaaa-yaoyun-datamanmesos\&group\=yaoyun
```
* 注: group 为用户当前所在的组

return

```
{
    "code":0,
    "data":{
        "count":1,
        "events":[
            {
                "id":1,
                "containerid":"2101477298d0622d2fb81156ec6ad328c778035dfc528c475f8e1d6908ece819",
                "message":"192.168.1.146 - - [03/May/2017:08:53:02 +0000] "GET / HTTP/1.1" 304 0 "http://localhost:5013/ui/app/aaaaaaaaaaa-yaoyun-datamanmesos/instance" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36" "-" ",
                "logtime":"2017-05-03T16:53:04+08:00",
                "path":"stdout",
                "offset":1493801584377272300,
                "DM_SLOT_INDEX":"0",
                "DM_APP_ID":"aaaaaaaaaaa-yaoyun-datamanmesos",
                "DM_USER":"yaoyun",
                "DM_TASK_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos-1a4cde96aef34c03a8ed9f07dae5ec63",
                "DM_GROUP_NAME":"yaoyun",
                "DM_VCLUSTER":"ymola",
                "keyword":"GET",
                "ack":false,
                "description":""
            }
        ]
    }
}
```

#### 设置报警事件为一睹
`PATCH /v2/log/alert/events/:id`

For exmaple
```
curl -XPATCH http://localhost:5098/v2/log/alert/events/1 -d '{"action":"ack"}'
```




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
