(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorCtrl', MonitorCtrl);
    /* @ngInject */
    function MonitorCtrl($state, $stateParams, monitorBackend) {
        var self = this;

        self.query = {
            metric: $stateParams.metric,
            appid: $stateParams.appid,
            taskid: $stateParams.taskid,
            start: $stateParams.start,
            end: $stateParams.end,
            step: $stateParams.step,
            expr: $stateParams.expr
        };

        activate();

        function activate() {
            monitorBackend.monitor(self.query).get(function (data) {

            })

        }
    }
})();