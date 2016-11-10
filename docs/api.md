# Logging/Monitoring Proxy API Guide

`GET /v1/search/applications`(#itm1 "itm1")


## Monitoring

### Get the metric values (CPU/Memory)

```GET /v1/monitor/query```

- path: /v1/monitor/query
- HTTP Method: GET
- URL Params: Null
- Query Params: metric, appid, from, to, step
  - metric=[all/cpu/memory]: the metric string.
  - appid=<string>: the name of application.
  - from=<2006-01-02 15:04:05>: the start time of the query range.
  - to=<2006-01-02 15:04:05>: the end time of the query range.
  - step=<duration>: Query resolution step width.

For example:
```
curl http://127.0.0.1:5098/v1/monitor/query?metric=memory&appid=nginx-stress&from=2016-11-09%2000:01:00&to=2016-11-09%2000:01:30&step=10s
```
return
```
{
    "code":0,
    "data":{
        "cpu":null,
        "memory":[
            {
                "metric":{
                    "container_label_APP_ID":"nginx-stress",
                    "group":"cadvisor",
                    "id":"/docker/063d8a98a5df330c5f22800e0ec212167a816daf4fe7064614db7fdd3927f12a",
                    "image":"192.168.1.58/library/nginx-stress:1.10",
                    "instance":"192.168.1.137:5014",
                    "job":"dataman",
                    "name":"mesos-3e79858b-7dea-470c-8038-6ed01fd69485-S0.7d54119e-543c-46cf-8625-4e52cab6ed4b"
                },
                "values":[
                    [
                        1478649660,
                        "0.02130126953125"
                    ],
                    [
                        1478649670,
                        "0.02130126953125"
                    ],
                    [
                        1478649680,
                        "0.02130126953125"
                    ],
                    [
                        1478649690,
                        "0.02130126953125"
                    ]
                ]
            }
        ]
    }
}
```

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
`GET /v1/search/paths/:appid/:taskid`

For example:
```
curl -XGET http://localhost:5098/v1/search/paths/appid/taskid
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
