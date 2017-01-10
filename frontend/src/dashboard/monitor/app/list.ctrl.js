(function () {
    'use strict';
    angular.module('app')
        .controller('DashboardListAppCtrl', DashboardListAppCtrl);
    /* @ngInject */
    function DashboardListAppCtrl(info, $stateParams) {
        var self = this;

        self.clusterInfo = info.data.clusters[$stateParams.clusterId];

        activate();

        function activate() {

        }

    }
})();
