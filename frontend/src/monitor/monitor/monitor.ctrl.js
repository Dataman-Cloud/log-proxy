(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorCtrl', MonitorCtrl);
    /* @ngInject */
    function MonitorCtrl($stateParams, monitorBackend, monitorSingleChart) {
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

        if ($stateParams.metric) {
            self.chartOptions = monitorSingleChart.Options($stateParams.metric);
        } else if ($stateParams.expr) {
            self.chartOptions = monitorSingleChart.Options('expr');
        }

        activate();

        function activate() {
            if (self.chartOptions) {
                monitorBackend.monitor(self.query).get(function (data) {
                    self.chartOptions.pushData(data.data, self.api);
                    setTimeout(function () {
                        self.chartOptions.flushCharts(self.api)
                    }, 0);
                })
            }
        }
    }
})();