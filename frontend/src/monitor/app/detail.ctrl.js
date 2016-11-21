(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorAppDetailCtrl', MonitorAppDetailCtrl);
    /* @ngInject */
    function MonitorAppDetailCtrl(monitorBackend, $stateParams, $q, $timeout, $scope, monitorChart) {
        var timeoutResult;

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

        activate();

        function activate() {

            monitorBackend.monitor({
                metric: 'all',
                appid: $stateParams.appId,
                from: '2016-11-09%2000:01:00',
                to: '2016-11-09%2000:01:30'
            }).get(function (data) {
                self.chartOptions.pushData(data.data, self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
                self.chartOptions.flushCharts(self.cpuApi, self.memApi, self.networkApi, self.fileSysApi);
            });

            tick();
        }

        function tick() {
            $q.all([monitorBackend.monitor({metric: 'all', appid: $stateParams.appId}).get().$promise,
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

                    timeoutResult = $timeout(tick, 5000);
                })

        };

        $scope.$on('$destroy', function () {
            $timeout.cancel(timeoutResult);
        });

    }
})();
