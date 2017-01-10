(function () {
    'use strict';
    angular.module('app')
        .controller('DashboardInstanceCtrl', DashboardInstanceCtrl);
    /* @ngInject */
    function DashboardInstanceCtrl(monitorBackend, $stateParams, $q, $timeout, $scope, monitorChart, moment) {
        var timeoutResult;
        var reloadInterval = 10000;
        var monitorData = {
            network: {receive: [], transmit: []},
            filesystem: {read: [], write: []}
        };

        var self = this;
        self.realTimeData = {
            cpuPercent: 0,
            memUsage: 0,
            networkReceive: 0,
            networkTransmit: 0,
            fileSysRead: 0,
            fileSysWrite: 0
        };

        self.chartOptions = monitorChart.Options();
        self.curTimestamp = moment().unix();
        self.fromTimestamp = moment().subtract(2, 'hours').unix();

        activate();

        function activate() {
            $q.all([monitorBackend.monitor({
                clusterid: $stateParams.clusterId,
                userid: $stateParams.userId,
                appid: $stateParams.appId,
                metric: 'cpu',
                taskid: $stateParams.taskId,
                start: self.fromTimestamp,
                end: self.curTimestamp
            }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'memory',
                    taskid: $stateParams.taskId,
                    start: self.fromTimestamp,
                    end: self.curTimestamp
                }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'network_rx',
                    taskid: $stateParams.taskId,
                    start: self.fromTimestamp,
                    end: self.curTimestamp
                }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'network_tx',
                    taskid: $stateParams.taskId,
                    start: self.fromTimestamp,
                    end: self.curTimestamp
                }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'fs_read',
                    taskid: $stateParams.taskId,
                    start: self.fromTimestamp,
                    end: self.curTimestamp
                }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'fs_write',
                    taskid: $stateParams.taskId,
                    start: self.fromTimestamp,
                    end: self.curTimestamp
                }).get().$promise
            ]).then(function (results) {
                monitorData = mergeMonitorData(results);

                self.chartOptions.pushData(monitorData, self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
                self.chartOptions.flushCharts(self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
                tick();
            })

        }

        function mergeMonitorData(results) {
            angular.forEach(results, function (result, index) {
                if (result.data.cpu.usage) {
                    monitorData.cpu = result.data.cpu;
                }

                if (result.data.memory.usage) {
                    monitorData.memory = result.data.memory;
                }

                if (result.data.network.receive) {
                    monitorData.network.receive = result.data.network.receive;
                }

                if (result.data.network.transmit) {
                    monitorData.network.transmit = result.data.network.transmit;
                }

                if (result.data.filesystem.read) {
                    monitorData.filesystem.read = result.data.filesystem.read;
                }

                if (result.data.filesystem.write) {
                    monitorData.filesystem.write = result.data.filesystem.write;
                }
            });

            return monitorData;
        }

        function tick() {
            $q.all([
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'cpu',
                    taskid: $stateParams.taskId,
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'memory',
                    taskid: $stateParams.taskId,
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'network_rx',
                    taskid: $stateParams.taskId,
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'network_tx',
                    taskid: $stateParams.taskId,
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'fs_read',
                    taskid: $stateParams.taskId,
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    clusterid: $stateParams.clusterId,
                    userid: $stateParams.userId,
                    appid: $stateParams.appId,
                    metric: 'fs_write',
                    taskid: $stateParams.taskId,
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise
            ]).then(function (result) {
                monitorData = mergeMonitorData(result);
                self.chartOptions.pushData(monitorData, self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
                self.chartOptions.flushCharts(self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);

                self.realTimeData.cpuPercent = parseFloat(result[0].data.cpu.usage[0].values[0][1]);
                self.realTimeData.memUsage = parseFloat(result[1].data.memory.usage[0].values[0][1]);
                self.realTimeData.networkReceive = parseInt(result[2].data.network.receive[0].values[0][1]);
                self.realTimeData.networkTransmit = parseInt(result[3].data.network.transmit[0].values[0][1]);
                self.realTimeData.fileSysRead = parseInt(result[4].data.filesystem.read[0].values[0][1]);
                self.realTimeData.fileSysWrite = parseInt(result[5].data.filesystem.write[0].values[0][1]);

                timeoutResult = $timeout(tick, reloadInterval);
            });
        }

        $scope.$on('$destroy', function () {
            $timeout.cancel(timeoutResult);
        });

    }
})();