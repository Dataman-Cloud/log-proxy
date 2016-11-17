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
curl http://127.0.0.1:5098/v1/monitor?metric=memory&appid=nginx-stress&from=2016-11-09%2000:01:00&to=2016-11-09%2000:01:30&step=10s

{
  "code": 0,
  "data": {
    "cpu": {
      "usage": [
        {
          "metric": {
            "container_label_APP_ID": "nginx-stress",
            "group": "cadvisor",
            "id": "/docker/063d8a98a5df330c5f22800e0ec212167a816daf4fe7064614db7fdd3927f12a",
            "image": "192.168.1.58/library/nginx-stress:1.10",
            "instance": "192.168.1.137:5014",
            "job": "dataman",
            "name": "mesos-3e79858b-7dea-470c-8038-6ed01fd69485-S0.7d54119e-543c-46cf-8625-4e52cab6ed4b"
          },
          "values": [
            [
              1478649660,
              "0"
            ]
          ]
        }
      ]
    },
    "memory": {
      "usage": [
        {
          "metric": {
            "container_label_APP_ID": "nginx-stress",
            "group": "cadvisor",
            "id": "/docker/063d8a98a5df330c5f22800e0ec212167a816daf4fe7064614db7fdd3927f12a",
            "image": "192.168.1.58/library/nginx-stress:1.10",
            "instance": "192.168.1.137:5014",
            "job": "dataman",
            "name": "mesos-3e79858b-7dea-470c-8038-6ed01fd69485-S0.7d54119e-543c-46cf-8625-4e52cab6ed4b"
          },
          "values": [
            [
              1478649660,
              "0.02130126953125"
            ]
          ]
        }
      ]
    },
    "network": {
      "receive": [
        {
          "metric": {
            "container_label_APP_ID": "nginx-stress",
            "group": "cadvisor",
            "id": "/docker/063d8a98a5df330c5f22800e0ec212167a816daf4fe7064614db7fdd3927f12a",
            "image": "192.168.1.58/library/nginx-stress:1.10",
            "instance": "192.168.1.137:5014",
            "job": "dataman",
            "name": "mesos-3e79858b-7dea-470c-8038-6ed01fd69485-S0.7d54119e-543c-46cf-8625-4e52cab6ed4b"
          },
          "values": [
            [
              1478649660,
              "20552716"
            ]
          ]
        }
      ],
      "transmit": [
        {
          "metric": {
            "container_label_APP_ID": "nginx-stress",
            "group": "cadvisor",
            "id": "/docker/063d8a98a5df330c5f22800e0ec212167a816daf4fe7064614db7fdd3927f12a",
            "image": "192.168.1.58/library/nginx-stress:1.10",
            "instance": "192.168.1.137:5014",
            "job": "dataman",
            "name": "mesos-3e79858b-7dea-470c-8038-6ed01fd69485-S0.7d54119e-543c-46cf-8625-4e52cab6ed4b"
          },
          "values": [
            [
              1478649660,
              "33266850"
            ]
          ]
        }
      ]
    },
    "filesystem": {
      "receive": [
        {
          "metric": {
            "container_label_APP_ID": "nginx-stress",
            "group": "cadvisor",
            "id": "/docker/063d8a98a5df330c5f22800e0ec212167a816daf4fe7064614db7fdd3927f12a",
            "image": "192.168.1.58/library/nginx-stress:1.10",
            "instance": "192.168.1.137:5014",
            "job": "dataman",
            "name": "mesos-3e79858b-7dea-470c-8038-6ed01fd69485-S0.7d54119e-543c-46cf-8625-4e52cab6ed4b"
          },
          "values": [
            [
              1478649660,
              "0"
            ]
          ]
        }
      ],
      "transmit": [
        {
          "metric": {
            "container_label_APP_ID": "nginx-stress",
            "group": "cadvisor",
            "id": "/docker/063d8a98a5df330c5f22800e0ec212167a816daf4fe7064614db7fdd3927f12a",
            "image": "192.168.1.58/library/nginx-stress:1.10",
            "instance": "192.168.1.137:5014",
            "job": "dataman",
            "name": "mesos-3e79858b-7dea-470c-8038-6ed01fd69485-S0.7d54119e-543c-46cf-8625-4e52cab6ed4b"
          },
          "values": [
            [
              1478649660,
              "0"
            ]
          ]
        }
      ]
    }
  }
}
```
Get the metrics of one instance in one application.
```
http://127.0.0.1:5098/v1/monitor?metric=cpu&appid=nginx-stress&taskid=063d8a98a5df330c

{
  "code": 0,
  "data": {
    "cpu": {
      "usage": [
        {
          "metric": {
            "container_label_APP_ID": "nginx-stress",
            "group": "cadvisor",
            "id": "/docker/063d8a98a5df330c5f22800e0ec212167a816daf4fe7064614db7fdd3927f12a",
            "image": "192.168.1.58/library/nginx-stress:1.10",
            "instance": "192.168.1.137:5014",
            "job": "dataman",
            "name": "mesos-3e79858b-7dea-470c-8038-6ed01fd69485-S0.7d54119e-543c-46cf-8625-4e52cab6ed4b"
          },
          "values": [
            [
              1479102521,
              "0"
            ]
          ]
        }
      ]
    },
    "memory": {
      "usage": null
    },
    "network": {
      "receive": null,
      "transmit": null
    },
    "filesystem": {
      "read": null,
      "write": null
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

