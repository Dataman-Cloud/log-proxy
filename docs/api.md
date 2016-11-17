# Logging/Monitoring Proxy API Guide

<a href="#itm1">GET /v1/search/applications</a>

## Monitoring

### Get the metric values (CPU/Memory/Network/Filesystem)

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

### Get the list of applications

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
`GET /v1/search/paths/:appid/:taskid

<span id="itm1">For example:</span>
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
