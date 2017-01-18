(function () {
    'use strict';
    angular.module('app')
        .controller('DashboardListAppCtrl', DashboardListAppCtrl);
    /* @ngInject */
    function DashboardListAppCtrl(info, $stateParams) {
        var self = this;

        self.clusterInfo = info.data.clusters[$stateParams.cluster];

        activate();

        function activate() {

        }

    }
})();
