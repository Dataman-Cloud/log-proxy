(function () {
    'use strict';
    angular.module('app')
        .controller('DashboardAppDetailCtrl', DashboardAppDetailCtrl);
    /* @ngInject */
    function DashboardAppDetailCtrl(dashboardBackend, monitorBackend, $q, $timeout, $scope, $stateParams, moment, monitorChart) {
        var timeoutResult;
        var reloadInterval = 15000;
        var monitorData = {
            network: {receive: [], transmit: []},
            filesystem: {read: [], write: []}
        };

        var self = this;
        self.tasks = {};
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
                cluster: $stateParams.cluster,
                user: $stateParams.user,
                app: $stateParams.app,
                metric: 'cpu',
                start: self.fromTimestamp,
                end: self.curTimestamp
            }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'memory',
                    start: self.fromTimestamp,
                    end: self.curTimestamp
                }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'network_rx',
                    start: self.fromTimestamp,
                    end: self.curTimestamp
                }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'network_tx',
                    start: self.fromTimestamp,
                    end: self.curTimestamp
                }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'fs_read',
                    start: self.fromTimestamp,
                    end: self.curTimestamp
                }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'fs_write',
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
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'cpu',
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'memory',
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'network_rx',
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'network_tx',
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'fs_read',
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise,
                monitorBackend.monitor({
                    cluster: $stateParams.cluster,
                    user: $stateParams.user,
                    app: $stateParams.app,
                    metric: 'fs_write',
                    start: moment().subtract(10, 'seconds').unix(),
                    end: moment().unix()
                }).get().$promise
            ]).then(function (result) {
                monitorData = mergeMonitorData(result);
                self.chartOptions.pushData(monitorData, self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
                self.chartOptions.flushCharts(self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);

                timeoutResult = $timeout(tick, reloadInterval);
            });

            dashboardBackend.info({cluster: $stateParams.cluster, user: $stateParams.user ,app: $stateParams.app}).get(function (data) {
                self.appInfo = data.data.clusters[$stateParams.cluster].users[$stateParams.user].applications[$stateParams.app];

                self.realTimeData.cpuPercent = parseFloat(self.appInfo.cpu.usage[1]) * 100;
                self.realTimeData.memUsage = parseFloat(self.appInfo.memory.usage[1]) * 100;
                self.realTimeData.networkReceive = parseInt(self.appInfo.network.receive[1]);
                self.realTimeData.networkTransmit = parseInt(self.appInfo.network.transmit[1]);
                self.realTimeData.fileSysRead = parseInt(self.appInfo.filesystem.read[1]);
                self.realTimeData.fileSysWrite = parseInt(self.appInfo.filesystem.write[1]);

                self.tasks = self.appInfo.tasks;
            })

        }

        $scope.$on('$destroy', function () {
            $timeout.cancel(timeoutResult);
        });

    }
})();
