# Logging/Monitoring Proxy API Guide

## Monitoring

### Metric(监控数据)

#### Get the metric values (CPU/Memory/Network/Filesystem)

```GET /v1/monitor/```
- path: /v1/monitor
- HTTP Method: GET
- URL Params: Null
- Query Params: type, metric, appid, instanceid, from, to, step
  - type=[app]: the type string
  - metric=[all/cpu/memory/network_rx/network_tx/fs_read/fs_write]: the metric string.
  - appid=<string>: the name of application.
  - taskid=<string>: the id string of the docker instance.
  - from=<2006-01-02 15:04:05>: the start time of the query range.
  - to=<2006-01-02 15:04:05>: the end time of the query range.
  - step=<duration>: Query resolution step width.

For example:

Get the metrics of the instances in one application.
```
curl http://127.0.0.1:5098/v1/monitor?metric=all&appid=work-nginx&from=2016-11-17%2000:01:00&to=2016-11-17%2000:01:00&step=10s
{
  "code": 0,
  "data": {
    "cpu": {
      "usage": [
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.fb702a72-000d-4244-a644-f85d42fcee08"
          },
          "values": [
            [
              1479340860,
              "0"
            ]
          ]
        },
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/e58177e632ca1a6ef5404a76ae129f047e295cd4aab1c262eaec4811d24f9b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.07346023-b17b-4c97-8b29-f28063e8aa05"
          },
          "values": [
            [
              1479340860,
              "0"
            ]
          ]
        }
      ],
      "count": 2
    },
    "memory": {
      "usage": [
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.fb702a72-000d-4244-a644-f85d42fcee08"
          },
          "values": [
            [
              1479340860,
              "0.0213623046875"
            ]
          ]
        },
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/e58177e632ca1a6ef5404a76ae129f047e295cd4aab1c262eaec4811d24f9b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.07346023-b17b-4c97-8b29-f28063e8aa05"
          },
          "values": [
            [
              1479340860,
              "0.02130126953125"
            ]
          ]
        }
      ],
      "count": 2
    },
    "network": {
      "receive": [
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.fb702a72-000d-4244-a644-f85d42fcee08"
          },
          "values": [
            [
              1479340860,
              "1854"
            ]
          ]
        },
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/e58177e632ca1a6ef5404a76ae129f047e295cd4aab1c262eaec4811d24f9b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.07346023-b17b-4c97-8b29-f28063e8aa05"
          },
          "values": [
            [
              1479340860,
              "1944"
            ]
          ]
        }
      ],
      "transmit": [
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.fb702a72-000d-4244-a644-f85d42fcee08"
          },
          "values": [
            [
              1479340860,
              "648"
            ]
          ]
        },
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/e58177e632ca1a6ef5404a76ae129f047e295cd4aab1c262eaec4811d24f9b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.07346023-b17b-4c97-8b29-f28063e8aa05"
          },
          "values": [
            [
              1479340860,
              "648"
            ]
          ]
        }
      ],
      "count": 2
    },
    "filesystem": {
      "read": [
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.fb702a72-000d-4244-a644-f85d42fcee08"
          },
          "values": [
            [
              1479340860,
              "0"
            ]
          ]
        },
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/e58177e632ca1a6ef5404a76ae129f047e295cd4aab1c262eaec4811d24f9b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.07346023-b17b-4c97-8b29-f28063e8aa05"
          },
          "values": [
            [
              1479340860,
              "0"
            ]
          ]
        }
      ],
      "write": [
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.fb702a72-000d-4244-a644-f85d42fcee08"
          },
          "values": [
            [
              1479340860,
              "0"
            ]
          ]
        },
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/e58177e632ca1a6ef5404a76ae129f047e295cd4aab1c262eaec4811d24f9b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.07346023-b17b-4c97-8b29-f28063e8aa05"
          },
          "values": [
            [
              1479340860,
              "0"
            ]
          ]
        }
      ],
      "count": 2
    }
  }
}
```
Get the metrics of one instance in one application.
```
http://127.0.0.1:5098/v1/monitor?metric=cpu&appid=work-nginx&taskid=d2f1c53
{
  "code": 0,
  "data": {
    "cpu": {
      "usage": [
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "/docker/d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f",
            "image": "nginx:latest",
            "instance": "192.168.1.102:5014",
            "job": "cadvisor",
            "name": "mesos-2f4c9ba3-a8ca-4df1-a4d4-cbec4343a64c-S1.fb702a72-000d-4244-a644-f85d42fcee08"
          },
          "values": [
            [
              1479361498,
              "0"
            ]
          ]
        }
      ],
      "count": 1
    },
    "memory": {
      "usage": null,
      "count": 0
    },
    "network": {
      "receive": null,
      "transmit": null,
      "count": 0
    },
    "filesystem": {
      "read": null,
      "write": null,
      "count": 0
    }
  }
}
```
Get the aggregation of the metrics data of one applications
```
http://127.0.0.1:5098/v1/monitor?type=app&metric=cpu&appid=work-nginx&taskid=d2f1c53
{
  "code": 0,
  "data": {
    "cpu": {
      "usage": [
        {
          "metric": {
            "container_label_APP_ID": "work-nginx",
            "group": "",
            "id": "",
            "image": "",
            "instance": "",
            "job": "",
            "name": ""
          },
          "values": [
            [
              1479361556,
              "0"
            ]
          ]
        }
      ],
      "count": 1
    },
    "memory": {
      "usage": null,
      "count": 0
    },
    "network": {
      "receive": null,
      "transmit": null,
      "count": 0
    },
    "filesystem": {
      "read": null,
      "write": null,
      "count": 0
    }
  }
}
```

#### Get the list of applications

```GET /v1/monitor/applications```

- path: /v1/monitor/applications
- HTTP Method: GET
- URL Params: Null
- Query Params: appid
  - appid=<string>: the name of application.

For example:
Get the instances of all applications
```
curl http://127.0.0.1:5098/v1/monitor/applications
{
  "code": 0,
  "data": {
    "apps": {
      "work-nginx": [
        "/docker/d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f",
        "/docker/e58177e632ca1a6ef5404a76ae129f047e295cd4aab1c262eaec4811d24f9b6f"
      ],
      "work-nginx2": [
        "/docker/130c1a4d5ec79315e98deafb1445bb79a94125297f27d20847efb9ba961ebd7f",
        "/docker/5b8d5c0c1dea877c21e225ab9dfc6b7b0b5255f4e593b80a335613a5439ab425"
      ]
    }
  }
}
```

Get the instances of one Applications
```
http://127.0.0.1:5098/v1/monitor/applications?appid=work-nginx
{
  "code": 0,
  "data": {
    "apps": {
      "work-nginx": [
        "/docker/e58177e632ca1a6ef5404a76ae129f047e295cd4aab1c262eaec4811d24f9b6f",
        "/docker/d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f"
      ]
    }
  }
}
```

#### Get the metric data of nodes

```GET /v1/monitor/nodes```

- path: /v1/monitor/nodes
- HTTP Method: GET
- URL Params: Null
- Query Params: node
  - node=<string>: the IP address of node.

For example:
Get the metric data of all nodes
```
http://127.0.0.1:5098/v1/monitor/nodes
{
  "code": 0,
  "data": {
    "nodes": {
      "192.168.1.101": {
        "cpu": {
          "usage": [
            1480409373,
            "0.03590153046667789"
          ]
        },
        "memory": {
          "usage_bytes": [
            1480409373,
            "1189441536"
          ],
          "total_bytes": [
            1480409373,
            "3975892992"
          ]
        },
        "network": {
          "eno16777984": {
            "receive": [
              1480409373,
              "5120.2"
            ],
            "transmit": [
              1480409373,
              "4275.733333333334"
            ]
          },
          "eno33557248": {
            "receive": [
              1480409373,
              "281.6"
            ],
            "transmit": [
              1480409373,
              "0"
            ]
          }
        },
        "filesystem": {
          "/dev/mapper/centos-root": {
            "usage_bytes": [
              1480409373,
              "7745433600"
            ],
            "total_bytes": [
              1480409373,
              "28406726656"
            ]
          },
          "/dev/sda1": {
            "usage_bytes": [
              1480409373,
              "172580864"
            ],
            "total_bytes": [
              1480409373,
              "520794112"
            ]
          }
        }
      }
    }
  }
}
```
Get the metric data of one node
```
http://127.0.0.1:5098/v1/monitor/nodes?node=192.168.1.101
```


#### Get the metric data of application

```GET /v1/monitor/application```

- path: /v1/monitor/application
- HTTP Method: GET
- URL Params: Null
- Query Params: appid
  - appid=<string>: the name of application.

For example:

Get the metric data of the tasks in one application
```
http://127.0.0.1:5098/v1/monitor/application?appid=work-nginx
{
  "code": 0,
  "data": {
    "app": {
      "d2f1c5324cced328d766fb858055bc0a5f9fa04402343379077e89ce6c9c0b6f": {
        "cpu": {
          "usage": [
            1479452528,
            "0"
          ]
        },
        "memory": {
          "usage": [
            1479452528,
            "0.0213623046875"
          ]
        },
        "network": {
          "receive": [
            1479452528,
            "1854"
          ],
          "transmit": [
            1479452528,
            "648"
          ]
        },
        "filesystem": {
          "read": [
            1479452528,
            "0"
          ],
          "write": [
            1479452528,
            "0"
          ]
        }
      },
      "e58177e632ca1a6ef5404a76ae129f047e295cd4aab1c262eaec4811d24f9b6f": {
        "cpu": {
          "usage": [
            1479452528,
            "0"
          ]
        },
        "memory": {
          "usage": [
            1479452528,
            "0.02130126953125"
          ]
        },
        "network": {
          "receive": [
            1479452528,
            "1944"
          ],
          "transmit": [
            1479452528,
            "648"
          ]
        },
        "filesystem": {
          "read": [
            1479452528,
            "0"
          ],
          "write": [
            1479452528,
            "0"
          ]
        }
      }
    }
  }
}
```

### PromQL API

Refer to "https://prometheus.io/docs/querying/api/"

#### 单时间点查询
`GET /v1/monitor/promql/query`

https://prometheus.io/docs/querying/api/#instant-queries

#### 时间范围查询
`GET /v1/monitor/promql/query_range`

https://prometheus.io/docs/querying/api/#range-queries

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
