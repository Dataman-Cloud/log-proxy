(function () {
    'use strict';
    angular.module('app')
        .factory('statsChart', statsChart);


    /* @ngInject */
    function statsChart($filter, $rootScope, chartUtil) {

        var optionsCls = createOptionsCls();

        return {
            Options: createOptions
        };

        function createOptions(serialKey) {
            return new optionsCls(serialKey);
        }

        function createOptionsCls() {
            function Options(serialKey) {
                this.serialKey = serialKey;
                this.memOptions = this._createMemOptions();
                this.cpuOptions = this._createCpuOptions();
                this.networkOptions = this._createNetworkOptions();
                this.diskOptions = this._createDiskOptions();
                this.memData = [];
                this.cpuData = [];
                this.networkData = [];
                this.diskData = [];
                this.dataChanged = false;
            }

            Options.prototype._createDiskOptions = function () {
                var options = chartUtil.createDefaultOptions();
                options.chart.yAxis.axisLabel = '速率';
                options.title.text = '磁盘读/写速率';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('netRate')(d);
                };
                options.chart.margin.left = 80;
                return options
            };

            Options.prototype._createNetworkOptions = function () {
                var options = chartUtil.createDefaultOptions();
                options.chart.yAxis.axisLabel = '速率';
                options.title.text = '网络接收/发送速率';
                options.chart.yAxis.tickFormat = function (d) {
                    return $filter('netRate')(d);
                };
                options.chart.margin.left = 80;
                return options
            };

            Options.prototype._createCpuOptions = function () {
                var options = chartUtil.createDefaultOptions();
                options.chart.yAxis.axisLabel = 'CPU使用率';
                options.title.text = 'CPU 使用率';
                return options
            };

            Options.prototype._createMemOptions = function () {
                var options = chartUtil.createDefaultOptions();
                options.chart.yAxis.axisLabel = '内存使用率';
                options.title.text = '内存使用率';
                return options;
            };

            Options.prototype.pushData = function (data, cpuApi, memApi, networkApi, diskApi) {
                var stats = angular.fromJson(data);
                var data = stats.Stats;
                var serialName = stats.TaskID;
                var x = new Date(data.read).getTime();
                this._pushCpuData(serialName, x, data, cpuApi);
                this._pushMemData(serialName, x, data, memApi);
                this._pushNetworkData(serialName, x, stats, networkApi);
                this._pushDiskData(serialName, x, stats, diskApi);
                this.dataChanged = true;
            };

            Options.prototype._pushCpuData = function (serialName, x, data, cpuApi) {
                var serialKey = serialName;
                chartUtil.pushData(this.cpuData, serialKey, {
                    x: x,
                    y: this._getCpuUsageRate(data)
                }, $rootScope.STATS_POINT_NUM);
            };

            Options.prototype._pushMemData = function (serialName, x, data, memApi) {
                var serialKey = serialName;
                chartUtil.pushData(this.memData, serialKey, {
                    x: x, y: data.memory_stats.usage / data.memory_stats.limit * 100,
                    total: data.memory_stats.limit, use: data.memory_stats.usage
                }, $rootScope.STATS_POINT_NUM);
            };

            Options.prototype._pushNetworkData = function (serialName, x, data, networkApi) {

                chartUtil.pushData(this.networkData, serialName + '接收', {
                        x: x,
                        y: data.ReceiveRate
                    },
                    $rootScope.STATS_POINT_NUM);
                chartUtil.pushData(this.networkData, serialName + '发送', {
                        x: x,
                        y: data.SendRate
                    },
                    $rootScope.STATS_POINT_NUM);
            };

            Options.prototype._pushDiskData = function (serialName, x, data, diskApi) {
                chartUtil.pushData(this.diskData, serialName + '读取', {
                        x: x,
                        y: data.BlkIOReadRate
                    },
                    $rootScope.STATS_POINT_NUM);
                chartUtil.pushData(this.diskData, serialName + '写入', {
                        x: x,
                        y: data.BlkIOWriteRate
                    },
                    $rootScope.STATS_POINT_NUM);
            };

            Options.prototype.flushCharts = function (cpuApi, memApi, networkApi, diskApi) {
                if(this.dataChanged) {
                    if (chartUtil.updateForceY(this.cpuOptions.chart, this.cpuData, 100, 1.2, 1, 100)) {
                        cpuApi.refresh();
                    } else {
                        cpuApi.update();
                    }

                    if (chartUtil.updateForceY(this.memOptions.chart, this.memData, 100, 1.2, 1, 100)) {
                        memApi.refresh();
                    } else {
                        memApi.update();
                    }

                    networkApi.update();
                    diskApi.update();
                }
            };

            Options.prototype._getCpuUsageRate = function (data) {
                var cpuPercent = 0;
                var cpuDelta = data.cpu_stats.cpu_usage.total_usage - data.precpu_stats.cpu_usage.total_usage;
                var systemDelta = data.cpu_stats.system_cpu_usage - data.precpu_stats.system_cpu_usage;

                if (systemDelta > 0 && cpuDelta > 0) {
                    cpuPercent = (cpuDelta / systemDelta) * data.cpu_stats.cpu_usage.percpu_usage.length * 100;
                }
                return cpuPercent
            };

            return Options;
        }
    }
})();
