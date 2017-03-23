(function () {
    'use strict';
    angular.module('app')
        .controller('MonitorBaseCtrl', MonitorBaseCtrl);
    /* @ngInject */
    function MonitorBaseCtrl($state, $stateParams, moment, monitorBackend) {
        var self = this;

        self.timePeriod = 60;
        self.selectedTabIndex = ($stateParams.start || $stateParams.end) ? 1 : 0;
        self.queryTabIndex = $stateParams.expr ? 1 : 0;

        self.form = {
            cluster: $stateParams.cluster,
            metric: $stateParams.metric,
            app: $stateParams.app,
            task: $stateParams.task,
            step: $stateParams.step,
            expr: $stateParams.expr
        };

        self.startTime = new Date(parseInt($stateParams.start) * 1000);
        self.endTime = new Date(parseInt($stateParams.end) * 1000);

        self.clusters = [];
        self.apps = [];
        self.tasks = [];
        self.metrics = [];

        self.loadClusters = loadClusters;
        self.loadApps = loadApps;
        self.loadTasks = loadTasks;
        self.loadMetric = loadMetric;
        self.checkTimeVali = checkTimeVali;

        self.search = search;

        activate();

        function activate() {
            loadClusters();
            loadApps();
            loadTasks();
            loadMetric();
        }

        function loadClusters() {
            checkTimeRange();

            monitorBackend.clusters({
                start: self.selectedTabIndex ? self.startTime.getTime() / 1000: self.form.start,
                end: self.selectedTabIndex ? self.endTime.getTime() / 1000: self.form.end,
                step: self.form.step
            }).get(function (data) {
                self.clusters = data.data;
            })
        }

        function loadApps() {
            if (self.form.cluster) {
                monitorBackend.apps(self.form.cluster, {
                    start: self.selectedTabIndex ? self.startTime.getTime() / 1000: self.form.start,
                    end: self.selectedTabIndex ? self.endTime.getTime() / 1000: self.form.end,
                    step: self.form.step
                }).get(function (data) {
                    self.apps = data.data;
                })
            }
        }

        function loadTasks() {
            if (self.form.cluster && self.form.app) {
                monitorBackend.tasks(self.form.cluster, self.form.app, {
                    start: self.selectedTabIndex ? self.startTime.getTime() / 1000: self.form.start,
                    end: self.selectedTabIndex ? self.endTime.getTime() / 1000: self.form.end,
                    step: self.form.step
                }).get(function (data) {
                    self.tasks = data.data;
                })
            }
        }

        function loadMetric() {
            monitorBackend.metrics()
                .get(function (data) {
                    self.metrics = data.data;
                })
        }

        function checkTimeVali(tabIndex) {
            if (tabIndex === 1) {
                return !(self.startTime && self.endTime)
            } else {
                return false
            }
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
                    cluster: self.form.cluster,
                    metric: self.form.metric,
                    app: self.form.app,
                    task: self.form.task,
                    start: self.form.start,
                    end: self.form.end,
                    step: self.form.step
                };
            }
            $state.go('home.monitor.chart', form, {inherit: false})

        }
    }
})();
