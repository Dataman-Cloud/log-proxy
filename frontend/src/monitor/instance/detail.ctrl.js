(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorInstanceCtrl', MonitorInstanceCtrl);
    /* @ngInject */
    function MonitorInstanceCtrl(monitorBackend, $stateParams, $q, $timeout, $scope, monitorChart, moment) {
        var timeoutResult;
        var reloadInterval = 5000;

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
            monitorBackend.monitor({
                metric: 'all',
                appid: $stateParams.appId,
                taskid: $stateParams.taskId,
                to: self.curTimestamp,
                from: self.fromTimestamp
            }).get(function (data) {
                self.chartOptions.pushData(data.data, self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
                self.chartOptions.flushCharts(self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
            }).$promise.then(function () {
                tick();
            });
        }

        function tick() {
            $q.all([monitorBackend.monitor({
                metric: 'all',
                appid: $stateParams.appId,
                taskid: $stateParams.taskId,
                to: moment().unix(),
                from: moment().subtract(5, 'seconds').unix()
            }).get().$promise,
                monitorBackend.monitor({
                    metric: 'all',
                    appid: $stateParams.appId,
                    taskid: $stateParams.taskId
                }).get().$promise])
                .then(function (result) {
                    self.chartOptions.pushData(result[0].data, self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
                    self.chartOptions.flushCharts(self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);

                    self.realTimeData.cpuPercent = parseFloat(result[1].data.cpu.usage[0].values[0][1]) * 100;
                    self.realTimeData.memUsage = parseFloat(result[1].data.memory.usage[0].values[0][1]) * 100;
                    self.realTimeData.networkReceive = parseInt(result[1].data.network.receive[0].values[0][1]);
                    self.realTimeData.networkTransmit = parseInt(result[1].data.network.transmit[0].values[0][1]);
                    self.realTimeData.fileSysRead = parseInt(result[1].data.filesystem.read[0].values[0][1]);
                    self.realTimeData.fileSysWrite = parseInt(result[1].data.filesystem.write[0].values[0][1]);

                    timeoutResult = $timeout(tick, reloadInterval);
                })
        }

        $scope.$on('$destroy', function () {
            $timeout.cancel(timeoutResult);
        });
    }
})();