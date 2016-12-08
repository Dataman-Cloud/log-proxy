(function () {
    'use strict';
    angular.module('app')
        .controller('DashboardListAppCtrl', DashboardListAppCtrl);
    /* @ngInject */
    function DashboardListAppCtrl(info) {
        var self = this;
        self.clusterInfo = info.data;

        activate();

        function activate() {

        }

    }
})();
