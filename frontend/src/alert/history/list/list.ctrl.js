(function () {
    'use strict';
    angular.module('app')
        .controller('AlertHistoriesCtrl', AlertHistoriesCtrl);
    /* @ngInject */
    function AlertHistoriesCtrl(alertBackend) {
        var self = this;

        self.monitors = [];

        activate();

        function activate() {
            listHistory()
        }

        function listHistory() {
            alertBackend.history().get(function (data) {
                self.histories = data.data.results;
            })
        }
    }
})();