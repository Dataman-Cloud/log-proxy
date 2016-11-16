(function () {
    'use strict';
    angular.module('app')
        .controller('ListAppCtrl', ListAppCtrl);
    /* @ngInject */
    function ListAppCtrl(appMonitorBackend, timing, $scope, $q) {
        var self = this;
        self.apps = [];

        timing.start($scope, reloadClusters, 5000);
        activate();

        function activate() {
            appMonitorBackend.listApp().query()
        }

        function reloadClusters() {
            var deferred = $q.defer();
            self.apps = appMonitorBackend.listApp().query(function () {
                deferred.resolve();
            });

            return deferred.promise;
        }
    }
})();
