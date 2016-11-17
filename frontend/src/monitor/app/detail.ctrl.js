(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorAppDetailCtrl', MonitorAppDetailCtrl);
    /* @ngInject */
    function MonitorAppDetailCtrl(monitorBackend, $stateParams, $q, $timeout, $scope) {
        var timeoutResult;

        var self = this;
        self.realTimeData = {
            cpuPercent: 0,
            memUsage: 0,
            networkReceive: 0,
            networkTransmit: 0,
            fileSysRead: 0,
            fileSysWrite: 0
        };

        activate();

        function activate() {
            (function tick() {
                $q.all([monitorBackend.monitor({metric: 'all', appid: $stateParams.appId}).get().$promise,
                    monitorBackend.monitor({metric: 'all', appid: $stateParams.appId, type: 'app'}).get().$promise])
                    .then(function (result) {
                        self.realTimeData.cpuPercent = parseInt(result[1].data.cpu.usage[0].values[0][1]);
                        self.realTimeData.memUsage = parseInt(result[1].data.memory.usage[0].values[0][1]);
                        self.realTimeData.networkReceive = parseInt(result[1].data.network.receive[0].values[0][1]);
                        self.realTimeData.networkTransmit = parseInt(result[1].data.network.transmit[0].values[0][1]);
                        self.realTimeData.fileSysRead = parseInt(result[1].data.filesystem.read[0].values[0][1]);
                        self.realTimeData.fileSysWrite = parseInt(result[1].data.filesystem.write[0].values[0][1]);


                        timeoutResult = $timeout(tick, 5000);
                    })

            })();
        }

        $scope.$on('$destroy', function () {
            $timeout.cancel(timeoutResult);
        });

    }
})();
