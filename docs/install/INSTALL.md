# 配置监控日志报警服务

*临时部署文档，仅供实验环境部署*

## 1. 修改alertmanager的配置
在运行alertmanager的服务器上修改配置文件`/data/config/alertmanager/alertmanager.yml`
```
receivers:
  - name: 'alert_webhook'
    webhook_configs:
    - send_resolved: true
      url: 'http://172.16.0.20:8080/api/webhooks/prometheus'
```
在后面新增内容
```
    - send_resolved: true
      url: 'http://192.168.1.75:5098/v1/receive/prometheus'
```
重启alertmanager
```
# docker restart dataman-monitor-alertmanager
```

## 2. 修改prometheus的配置
在prometheus的服务器上修改配置文件`/data/config/prometheus/prometheus.yml`,
增加target mola，为关键字报警采集数据用。
```
- job_name: 'mola'
  # Override the global default and scrape targets from this job every 5 seconds.
  scrape_interval: 10s
  scrape_timeout: 5s
  static_configs:
    - targets: ['192.168.1.75:5098']
```
修改配置文件`/data/config/prometheus/alert.rules`, 增加报警规则。
```
ALERT LogKeyword
  IF ceil(increase(log_keyword{id=~'.*'}[2m])) > 3
  FOR 1m
  LABELS { severity = "Warning" }
  ANNOTATIONS {
    summary = "Application {{ $labels.appid }} taskid {{ $labels.taskid }} Log keyword filter {{ $labels.keyword }} trigger times {{ $value }}",
    description = "Application {{ $labels.appid }} taskid {{ $labels.taskid }} Log keyword filter {{ $labels.keyword }} trigger times {{ $value }}",
  }
```
规则描述：
* ALERT LogKeyword 是规则名
* IF increase(log_keyword{id=~'.*'}[2m]) > 3 告警触发条件, 2m是时间，两分钟，3是触发次数阈值。
* FOR 1m 报警延迟1分钟发送
* LABELS { severity = "Warning" } 附加的标签，增加标签severity，表示告警等级。
* ANNOTATIONS 报警消息的内容设定。

重启prometheus
```
# docker restart dataman-monitor-prometheus
```

## 3. 启动mola
在独立的主机上运行mola，也可以和prometheus运行在同一个主机上，

* mola是无状态应用，可以部署多个实例

* 支持接多个elasticsearch服务

* 暂不支持接多个prometheus, alertmanager的地址

进入目录mola,根据实际情况,  修改[mola/docker-compose.yml](mola/docker-compose.yml)中的环境变量

* ES_URL: Elasticsearch地址
```
   - ES_URL=http://192.168.1.75:9200
```
* PROMETHEUS_URL：Prometheus地址
```
   - PROMETHEUS_URL=http://192.168.1.75:9090
```
* ALERTMANAGER_URL：ALERTMANAGER_URL地址
```
   - ALERTMANAGER_URL=http://192.168.1.75:9093
```

运行mola
```
docker-compose -p mola up -d
```

在浏览器打开服务地址：http://IP:5098

## 4. 安装log-agent（新版本borgsphere 3.2的线下包里已经集成了log-agent）
在工作节点或者要收集容器日志的主机上运行log-agent,

修改文件[log-agent/docker-compose.yml](log-agent/docker-compose.yml), 根据实际情况，替换变量。
* OUTPUT: Logstash的地址
```
   - OUTPUT=tcp://192.168.1.74:5044
```

运行log-agent
```
docker-compose -p log up -d
```

## 5. 更新应用（选作）
为了从监控数据中取到集群信息，需要在3.0更新应用，手动支持新加的lable VCLUSTER。
```
"parameters": [
{
"key": "label",
"value": "APP_ID=work-nginx-stress"
},
{
"key": "label",
"value": "VCLUSTER=work"
}
```

## 5. 启动es以后执行命令建立mapping 文件在docs/json_sample里面
```
curl -XPUT localhost:9200/_template/event -d @event.json
curl -XPUT localhost:9200/_template/dataman -d @log.json
curl -XPUT localhost:9200/_template/alert -d @alert.json
curl -XPUT localhost:9200/_template/keyword -d @alert-keyword.json
curl -XPUT localhost:9200/_template/prometheus -d @prometheus.json
```
