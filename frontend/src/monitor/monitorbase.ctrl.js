(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorBaseCtrl', MonitorBaseCtrl);
    /* @ngInject */
    function MonitorBaseCtrl($state, $stateParams, moment) {
        var self = this;
        self.timePeriod = 120;
        self.selectedTabIndex = ($stateParams.start || $stateParams.end) ? 1 : 0;
        self.queryTabIndex = $stateParams.expr ? 1 : 0;

        self.form = {
            metric: $stateParams.metric,
            appid: $stateParams.appid,
            taskid: $stateParams.taskid,
            step: $stateParams.step,
            expr: $stateParams.expr
        };

        self.startTime = new Date(parseInt($stateParams.start) * 1000);
        self.endTime = new Date(parseInt($stateParams.end) * 1000);

        self.search = search;

        activate();

        function activate() {

        }

        function checkTimeRange() {
            if (!self.selectedTabIndex) {
                self.form.end = moment().unix();
                self.form.start = moment().subtract(self.timePeriod, 'minutes').unix();
            } else {
                if (angular.isDate(self.endTime) && angular.isDate(self.startTime)) {
                    self.form.end = self.endTime.getTime() / 1000;
                    self.form.start = self.startTime.getTime() / 1000;
                }
            }
            // 180 is the max number of points in Chart
            self.form.step = Math.ceil((self.form.end - self.form.start) / 180);
        }

        function search() {
            checkTimeRange();
            var form = angular.copy(self.form);

            if (self.queryTabIndex) {
                form = {
                    expr: self.form.expr,
                    start: self.form.start,
                    end: self.form.end,
                    step: self.form.step
                };
            } else {
                form = {
                    metric: self.form.metric,
                    appid: self.form.appid,
                    taskid: self.form.taskid,
                    start: self.form.start,
                    end: self.form.end,
                    step: self.form.step
                };
            }
            $state.go('home.monitor.chart', form, {inherit: false})

        }
    }
})();
