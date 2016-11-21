(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorAppDetailCtrl', MonitorAppDetailCtrl);
    /* @ngInject */
    function MonitorAppDetailCtrl(monitorBackend, $stateParams, $q, $timeout, $scope, monitorChart, moment) {
        var timeoutResult;
        var reloadInterval = 5000;

        var self = this;
        self.instances = {};
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
        self.fromTimestamp = moment().subtract(1, 'hours').unix();

        activate();

        function activate() {

            monitorBackend.monitor({
                metric: 'all',
                appid: $stateParams.appId,
                from: self.fromTimestamp,
                to: self.curTimestamp
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
                to: moment().unix(),
                from: moment().subtract(5, 'seconds').unix()
            }).get().$promise,
                monitorBackend.monitor({metric: 'all', appid: $stateParams.appId, type: 'app'}).get().$promise,
                monitorBackend.listInstance({appid: $stateParams.appId}).get().$promise])
                .then(function (result) {
                    self.chartOptions.pushData(result[0].data, self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
                    self.chartOptions.flushCharts(self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);

                    self.realTimeData.cpuPercent = parseInt(result[1].data.cpu.usage[0].values[0][1]);
                    self.realTimeData.memUsage = parseInt(result[1].data.memory.usage[0].values[0][1]);
                    self.realTimeData.networkReceive = parseInt(result[1].data.network.receive[0].values[0][1]);
                    self.realTimeData.networkTransmit = parseInt(result[1].data.network.transmit[0].values[0][1]);
                    self.realTimeData.fileSysRead = parseInt(result[1].data.filesystem.read[0].values[0][1]);
                    self.realTimeData.fileSysWrite = parseInt(result[1].data.filesystem.write[0].values[0][1]);

                    self.instances = result[2].data.app;

                    timeoutResult = $timeout(tick, reloadInterval);
                })

        }

        $scope.$on('$destroy', function () {
            $timeout.cancel(timeoutResult);
        });

    }
})();
