(function () {
    'use strict';
    angular.module('app')
        .factory('monitorSingleChart', monitorSingleChart);


    /* @ngInject */
    function monitorSingleChart($filter, chartUtil) {

        var optionsCls = createOptionsCls();
        var POINT_NUM = 180;

        return {
            Options: createOptions
        };

        function createOptions(type) {
            return new optionsCls(type);
        }

        function formatCpuOrMem(usage) {
            angular.forEach(usage, function (item, index, usage) {
                angular.forEach(item.values, function (value, valueIndex, values) {
                    value[1] = (parseFloat(value[1]) * 100).toFixed(2)
                })
            });
        }

        function createOptionsCls(type) {
            function Options(type) {
                this.type = type;
                this.memOptions = this._createMemOptions();
                this.memUsageOptions = this._createMemUsageOptions();
                this.memTotalOptions = this._createMemTotalOptions();
                this.cpuOptions = this._createCpuOptions();
                this.networkRxOptions = this._createNetworkRxOptions();
                this.networkTxOptions = this._createNetworkTxOptions();
                this.fileSysReadOptions = this._createFileSysReadOptions();
                this.fileSysWriteOptions = this._createFileSysWriteOptions();
                this.exprOptions = this._createExprOptions();
                this.memData = [];
                this.memUsageData = [];
                this.memTotalData = [];
                this.cpuData = [];
                this.networkRxData = [];
                this.networkTxData = [];
                this.fileSysReadData = [];
                this.fileSysWriteData = [];
                this.exprData = [];
            }

            Options.prototype._createExprOptions = function () {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = 'expr';
                options.chart.yAxis.tickFormat = function (d) {
                    return d;
                };
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype._createFileSysWriteOptions = function () {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = '磁盘写入速率';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('netRate')(d);
                };
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype._createFileSysReadOptions = function () {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = '磁盘读取速率';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('netRate')(d);
                };
                options.chart.margin.left = 100;
                return options
            };


            Options.prototype._createNetworkTxOptions = function () {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = '网络读取速率';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('netRate')(d);
                };
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype._createNetworkRxOptions = function () {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = '网络接收速率';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('netRate')(d);
                };
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype._createCpuOptions = function () {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = 'CPU 使用率';
                options.chart.yAxis.tickFormat = function (d) {
                    return d + '%';
                };
                options.chart.forceY = [0, 100];
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype._createMemTotalOptions = function () {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = '内存总量';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('size')(d);
                };
                options.chart.forceY = [0, 100];
                options.chart.margin.left = 100;
                return options;
            };

            Options.prototype._createMemUsageOptions = function () {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = '内存使用量';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('size')(d);
                };
                options.chart.forceY = [0, 100];
                options.chart.margin.left = 100;
                return options;
            };

            Options.prototype._createMemOptions = function () {
                var options = chartUtil.createDefaultOptions({height: 500, showLegend: true, hideGuideline: true});
                options.title.text = '内存使用率';
                options.chart.yAxis.tickFormat = function (d) {
                    return d + '%';
                };
                options.chart.forceY = [0, 100];
                options.chart.margin.left = 100;
                return options;
            };

            Options.prototype.pushData = function (data, api) {
                switch (this.type) {
                    case 'cpu':
                        this._pushCpuData(data, api);
                        break;
                    case 'memory':
                        this._pushMemData(data, api);
                        break;
                    case 'memory_usage':
                        this._pushMemUsageData(data, api);
                        break;
                    case 'memory_total':
                        this._pushMemTotalData(data, api);
                        break;
                    case 'network_rx':
                        this._pushNetworRxkData(data, api);
                        break;
                    case 'network_tx':
                        this._pushNetworTxkData(data, api);
                        break;
                    case 'fs_read':
                        this._pushFileSysReadData(data, api);
                        break;
                    case 'fs_write':
                        this._pushFileSysWriteData(data, api);
                        break;
                    case 'expr':
                        this._pushExprData(data, api);
                        break;
                    default:
                }
            };

            Options.prototype._pushCpuData = function (data, cpuApi) {
                formatCpuOrMem(data.cpu.usage);
                this._pushData(data.cpu.usage, cpuApi, this.cpuData);
            };

            Options.prototype._pushMemTotalData = function (data, memApi) {
                this._pushData(data.memory.total_bytes, memApi, this.memTotalData);
            };

            Options.prototype._pushMemUsageData = function (data, memApi) {
                this._pushData(data.memory.usage_bytes, memApi, this.memUsageData);
            };

            Options.prototype._pushMemData = function (data, memApi) {
                formatCpuOrMem(data.memory.usage);
                this._pushData(data.memory.usage, memApi, this.memData);
            };

            Options.prototype._pushNetworRxkData = function (data, networkApi) {
                this._pushData(data.network.receive, networkApi, this.networkRxData, function (insId) {
                    return insId + '接收';
                });
            };

            Options.prototype._pushNetworTxkData = function (data, networkApi) {
                this._pushData(data.network.transmit, networkApi, this.networkTxData, function (insId) {
                    return insId + '发送';
                });

            };

            Options.prototype._pushFileSysReadData = function (data, fileSysApi) {
                this._pushData(data.filesystem.read, fileSysApi, this.fileSysReadData, function (insId) {
                    return insId + '读取';
                });
            };

            Options.prototype._pushFileSysWriteData = function (data, fileSysApi) {
                this._pushData(data.filesystem.write, fileSysApi, this.fileSysWriteData, function (insId) {
                    return insId + '写入';
                });
            };

            Options.prototype._pushExprData = function (data, exprApi) {
                this._pushData(data.result, exprApi, this.exprData, function (insId) {
                    return insId;
                });
            };

            Options.prototype._pushData = function (data, api, target, serialKeyBuilder) {
                angular.forEach(data, function (item, index) {
                    var serialKey;
                    if (serialKeyBuilder) {
                        if (this.type === 'expr') {
                            serialKey = '';
                            angular.forEach(item.metric, function (value, key) {
                                serialKey += value + ' ';
                            });
                            serialKey = serialKeyBuilder(serialKey);
                        } else {
                            serialKey = serialKeyBuilder(item.metric.container_label_SLOT + '-' + item.metric.id);
                        }
                    } else {
                        serialKey = item.metric.container_label_SLOT + '-' + item.metric.id;
                    }
                    angular.forEach(item.values, function (value) {
                        chartUtil.pushData(target, serialKey, {
                            x: value[0] * 1000,
                            y: parseFloat(value[1])
                        }, POINT_NUM);
                    })

                }.bind(this))
            };

            Options.prototype.flushCharts = function (api) {
                api.update();
            };


            return Options;
        }
    }
})();
