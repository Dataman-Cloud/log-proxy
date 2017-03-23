(function () {
    'use strict';
    angular.module('app')
        .controller('LogBaseCtrl', LogBaseCtrl);
    /* @ngInject */
    function LogBaseCtrl(moment, $state, logBackend, $stateParams) {
        var self = this;

        self.timePeriod = 60;
        self.selectedTabIndex = ($stateParams.from && $stateParams.to) ? 1 : 0;

        self.form = {
            cluster: $stateParams.cluster || '',
            app: $stateParams.app || '',
            task: $stateParams.task || '',
            path: $stateParams.path || '',
            keyword: $stateParams.keyword || ''
        };

        self.clusters = {};
        self.apps = {};
        self.tasks = [];
        self.paths = [];

        self.startTime = new Date(parseInt($stateParams.from));
        self.endTime = new Date(parseInt($stateParams.to));

        self.loadClusters = loadClusters;
        self.loadApps = loadApps;
        self.loadTasks = loadTasks;
        self.loadPaths = loadPaths;
        self.checkTimeVali = checkTimeVali;

        self.searchLog = searchLog;

        activate();

        function activate() {
            loadClusters();
            loadApps();
            loadTasks();
            // loadPaths();
        }

        function loadClusters() {
            checkTimeRange();

            logBackend.clusters({
                from: self.selectedTabIndex ? self.startTime.getTime() : self.form.from,
                to: self.selectedTabIndex ? self.endTime.getTime() : self.form.to
            }).get(function (data) {
                self.clusters = data.data;
            })
        }

        function loadApps() {
            if (self.form.cluster) {
                logBackend.apps(self.form.cluster, {
                    from: self.selectedTabIndex ? self.startTime.getTime() : self.form.from,
                    to: self.selectedTabIndex ? self.endTime.getTime() : self.form.to
                }).get(function (data) {
                    self.apps = data.data;
                })
            }
        }

        function loadTasks() {
            if (self.form.cluster && self.form.app) {
                logBackend.tasks(self.form.cluster, self.form.app, {
                    from: self.selectedTabIndex ? self.startTime.getTime() : self.form.from,
                    to: self.selectedTabIndex ? self.endTime.getTime() : self.form.to
                }).get(function (data) {
                    if (self.isTaskRunning && angular.isArray(data.data)) {
                        self.tasks = data.data.filter(function (task) {
                            return task.status === 'running'
                        })
                    } else {
                        self.tasks = data.data;
                    }
                })
            }
        }

        function loadPaths() {
            if (self.form.app && self.form.cluster) {
                logBackend.paths({
                    cluster: self.form.cluster,
                    app: self.form.app,
                    task: self.form.task,
                    from: self.selectedTabIndex ? self.startTime.getTime() : self.form.from,
                    to: self.selectedTabIndex ? self.endTime.getTime() : self.form.to
                }).get(function (data) {
                    self.paths = data.data;
                })
            }
        }

        function checkTimeVali(tabIndex) {
            if (tabIndex === 1) {
                return !(self.startTime && self.endTime)
            } else {
                return false
            }
        }

        function checkTimeRange() {
            if (self.selectedTabIndex === 0) {
                self.form.to = moment().unix() * 1000;
                self.form.from = moment().subtract(self.timePeriod, 'minutes').unix() * 1000;
            } else if (self.selectedTabIndex === 1) {
                if (angular.isDate(self.endTime) && angular.isDate(self.startTime)) {
                    self.form.to = self.endTime.getTime();
                    self.form.from = self.startTime.getTime();
                }
            }
        }

        function searchLog() {
            checkTimeRange();

            if (!self.form.keyword) {
                $state.go('home.logbase.logdetail', self.form);
            } else {
                $state.go('home.logbase.logs', self.form);
            }
        }
    }
})();
