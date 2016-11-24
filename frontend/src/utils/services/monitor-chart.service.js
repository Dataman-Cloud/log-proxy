(function () {
    'use strict';
    angular.module('app')
        .factory('monitorChart', monitorChart);


    /* @ngInject */
    function monitorChart($filter, chartUtil) {

        var optionsCls = createOptionsCls();
        var POINT_NUM = 180;

        return {
            Options: createOptions
        };

        function createOptions() {
            return new optionsCls();
        }

        function formatCpuOrMem(usage) {
            angular.forEach(usage, function (item, index, usage) {
                angular.forEach(item.values, function (value, valueIndex, values) {
                    value[1] = (parseFloat(value[1]) * 100).toFixed(2)
                })
            });
        }

        function createOptionsCls() {
            function Options() {
                this.memOptions = this._createMemOptions();
                this.cpuOptions = this._createCpuOptions();
                this.networkOptions = this._createNetworkOptions();
                this.fileSysOptions = this._createFileSysOptions();
                this.memData = [];
                this.cpuData = [];
                this.networkData = [];
                this.fileSysData = [];
                this.dataChanged = false;
            }

            Options.prototype._createFileSysOptions = function () {
                var options = chartUtil.createDefaultOptions();
                options.title.text = '磁盘读/写速率';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('netRate')(d);
                };
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype._createNetworkOptions = function () {
                var options = chartUtil.createDefaultOptions();
                options.title.text = '网络接收/发送速率';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('netRate')(d);
                };
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype._createCpuOptions = function () {
                var options = chartUtil.createDefaultOptions();
                options.title.text = 'CPU 使用率';
                options.chart.yAxis.tickFormat = function (d) {
                    return d + '%';
                };
                options.chart.forceY = [0, 100];
                options.chart.margin.left = 100;
                return options
            };

            Options.prototype._createMemOptions = function () {
                var options = chartUtil.createDefaultOptions();
                options.title.text = '内存使用率';
                options.chart.yAxis.tickFormat = function (d) {
                    return d + '%';
                };
                options.chart.forceY = [0, 100];
                options.chart.margin.left = 100;
                return options;
            };

            Options.prototype.pushData = function (data, cpuApi, memApi, networkApi, fileSysApi) {
                this._pushCpuData(data, cpuApi);
                this._pushMemData(data, memApi);
                this._pushNetworkData(data, networkApi);
                this._pushFileSysData(data, fileSysApi);
                this.dataChanged = true;
            };

            Options.prototype._pushCpuData = function (data, cpuApi) {
                formatCpuOrMem(data.cpu.usage);
                this._pushData(data.cpu.usage, cpuApi, this.cpuData);
            };

            Options.prototype._pushMemData = function (data, memApi) {
                formatCpuOrMem(data.memory.usage);
                this._pushData(data.memory.usage, memApi, this.memData);
            };

            Options.prototype._pushNetworkData = function (data, networkApi) {
                this._pushData(data.network.receive, networkApi, this.networkData, function (insId) {
                    return insId + '接收';
                });
                this._pushData(data.network.transmit, networkApi, this.networkData, function (insId) {
                    return insId + '发送';
                });

            };

            Options.prototype._pushFileSysData = function (data, fileSysApi) {
                this._pushData(data.filesystem.read, fileSysApi, this.fileSysData, function (insId) {
                    return insId + '读取';
                });
                this._pushData(data.filesystem.write, fileSysApi, this.fileSysData, function (insId) {
                    return insId + '写入';
                });
            };

            Options.prototype._pushData = function (data, api, target, serialKeyBuilder) {
                angular.forEach(data, function (item) {
                    var serialKey;
                    if (serialKeyBuilder) {
                        serialKey = serialKeyBuilder(item.metric.id);
                    } else {
                        serialKey = item.metric.id;
                    }
                    angular.forEach(item.values, function (value) {
                        chartUtil.pushData(target, serialKey, {
                            x: value[0] * 1000,
                            y: parseFloat(value[1])
                        }, POINT_NUM);
                    })

                })
            };

            Options.prototype.flushCharts = function (cpuApi, memApi, networkApi, diskApi) {
                if (this.dataChanged) {
                    cpuApi.update();
                    memApi.update();
                    networkApi.update();
                    diskApi.update();
                }
            };


            return Options;
        }
    }
})();
