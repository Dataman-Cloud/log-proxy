(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorListAppCtrl', MonitorListAppCtrl);
    /* @ngInject */
    function MonitorListAppCtrl(apps) {
        var self = this;

        self.apps = apps.data.apps;

        activate();

        function activate() {

        }

    }
})();
