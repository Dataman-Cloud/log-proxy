(function () {
    'use strict';
    angular.module('app')
        .controller('LogBaseCtrl', LogBaseCtrl);
    /* @ngInject */
    function LogBaseCtrl(moment, $state, logBackend) {
        var self = this;

        self.timePeriod = 120;
        self.selectedTabIndex = 0;

        self.periodChange = periodChange;
        self.selectAppChange = selectAppChange;
        self.loadApps = loadApps;
        self.loadTasks = loadTasks;
        self.loadPaths = loadPaths;
        self.searchLog = searchLog;

        activate();

        function activate() {

        }

        function periodChange(period) {
            self.curTimestamp = moment().unix() * 1000;
            self.fromTimestamp = moment().subtract(period, 'minutes').unix() * 1000;
        }

        function selectAppChange(app) {
            self.selectTasks = '';
            self.selectPaths = '';
        }

        function loadApps() {
            logBackend.listApp().get(function (data) {
                self.apps = data.data;
            })
        }

        function loadTasks() {
            logBackend.listTask({
                appid: self.selectApp
            }).get(function (data) {
                self.tasks = data.data;
            })
        }

        function loadPaths() {
            logBackend.listPath({
                appid: self.selectApp,
                taskid: self.selectTasks
            }).get(function (data) {
                self.paths = data.data;
            })
        }

        function checkTimeRange() {
            if (self.selectedTabIndex === 0) {
                self.curTimestamp = moment().unix() * 1000;
                self.fromTimestamp = moment().subtract(self.timePeriod, 'minutes').unix() * 1000;
            } else if (self.selectedTabIndex === 1) {
                if (angular.isDate(self.endTime) && angular.isDate(self.startTime)) {
                    self.curTimestamp = self.endTime.getTime();
                    self.fromTimestamp = self.startTime.getTime();
                }
            }
        }

        function searchLog() {
            checkTimeRange();

            if (self.keyword) {
                $state.go('home.logbase.logs', {
                    from: self.fromTimestamp,
                    to: self.curTimestamp,
                    appid: self.selectApp,
                    taskid: self.selectTasks,
                    path: self.selectPaths,
                    keyword: self.keyword
                });
            } else {
                $state.go('home.logbase.logWithoutKey', {
                    from: self.fromTimestamp,
                    to: self.curTimestamp,
                    appid: self.selectApp,
                    taskid: self.selectTasks,
                    path: self.selectPaths
                });
            }
        }
    }
})();