(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorNodeCtrl', MonitorNodeCtrl);
    /* @ngInject */
    function MonitorNodeCtrl(nodes, $scope, $timeout, monitorBackend) {
        var timeoutResult;
        var reloadInterval = 5000;

        var self = this;
        self.nodes = nodes.data.nodes;

        activate();

        function activate() {
            tick()
        }

        function tick() {
            monitorBackend.listNode().get().$promise
                .then(function (data) {
                    self.nodes = data.data.nodes;
                    timeoutResult = $timeout(tick, reloadInterval);
                });
        }

        $scope.$on('$destroy', function () {
            $timeout.cancel(timeoutResult);
        });


    }
})();