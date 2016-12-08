(function () {
    'use strict';
    angular.module('app')
        .controller('DashboardCtrl', DashboardCtrl);
    /* @ngInject */
    function DashboardCtrl(dashboardBackend, $stateParams, $scope, $timeout, $q) {
        var timeoutResult;
        var reloadInterval = 10000;

        var self = this;
        self.nodes = [];
        self.info = {
            clusters: {},
            applications: [],
            tasks: [],
            nodes: []
        };


        activate();

        function activate() {
            tick();
        }

        function tick() {
            $q.all([dashboardBackend.listNode().get().$promise,
                dashboardBackend.info().get().$promise])
                .then(function (result) {
                    self.nodes = result[0].data.nodes;

                    self.info = result[1].data;

                    timeoutResult = $timeout(tick, reloadInterval);
                })

        }

        $scope.$on('$destroy', function () {
            $timeout.cancel(timeoutResult);
        });

    }
})();