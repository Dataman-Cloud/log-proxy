# Logging/Monitoring Proxy API Guide

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

### Get the metric data of nodes

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
            1479444371,
            "0.030511715233357488"
          ]
        },
        "memory": {
          "usage": [
            1479444371,
            "0.5278851770465356"
          ]
        },
        "network": {
          "receive": [
            1479444371,
            "1853086509"
          ],
          "transmit": [
            1479444371,
            "1241003139"
          ]
        }
      },
      "192.168.1.102": {
        "cpu": {
          "usage": [
            1479444371,
            "0.018358025783337933"
          ]
        },
        "memory": {
          "usage": [
            1479444371,
            "0.12986786032589293"
          ]
        },
        "network": {
          "receive": [
            1479444371,
            "483338529"
          ],
          "transmit": [
            1479444371,
            "860170615"
          ]
        }
      },
      "192.168.1.91": {
        "cpu": {
          "usage": [
            1479444371,
            "0.05481752030003312"
          ]
        },
        "memory": {
          "usage": [
            1479444371,
            "0.7841073807249992"
          ]
        },
        "network": {
          "receive": [
            1479444371,
            "1880722606"
          ],
          "transmit": [
            1479444371,
            "2680864446"
          ]
        }
      },
      "192.168.1.92": {
        "cpu": {
          "usage": [
            1479444371,
            "0.04335599093331742"
          ]
        },
        "memory": {
          "usage": [
            1479444371,
            "0.7260611826624831"
          ]
        },
        "network": {
          "receive": [
            1479444371,
            "1502950018"
          ],
          "transmit": [
            1479444371,
            "1279601302"
          ]
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
