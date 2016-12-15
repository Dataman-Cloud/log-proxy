(function () {
    'use strict';
    angular.module('app')
        .controller('LogBaseCtrl', LogBaseCtrl);
    /* @ngInject */
    function LogBaseCtrl(moment, $state, logBackend, $stateParams) {
        var self = this;

        self.form = {
            appid: $stateParams.appid || '',
            taskid: $stateParams.taskid || '',
            path: $stateParams.path || '',
            keyword: $stateParams.keyword || ''
        };

        self.apps = [];
        self.tasks = [];
        self.paths = [];

        self.startTime = new Date(parseInt($stateParams.from));
        self.endTime = new Date(parseInt($stateParams.to));

        self.timePeriod = 120;
        self.selectedTabIndex = ($stateParams.from && $stateParams.to) ? 1 : 0;

        self.selectAppChange = selectAppChange;
        self.loadApps = loadApps;
        self.loadTasks = loadTasks;
        self.loadPaths = loadPaths;

        self.searchLog = searchLog;

        activate();

        function activate() {
            loadPaths();
            loadApps();
            loadTasks();
        }

        function selectAppChange(app) {
            self.form.taskid = '';
            self.form.path = '';
        }

        function loadApps() {
            logBackend.listApp().get(function (data) {
                self.apps = data.data;
            })
        }

        function loadTasks() {
            if (self.form.appid) {
                logBackend.listTask({
                    appid: self.form.appid
                }).get(function (data) {
                    self.tasks = data.data;
                })
            }
        }

        function loadPaths() {
            if (self.form.appid) {
                logBackend.listPath({
                    appid: self.form.appid,
                    taskid: self.form.taskid
                }).get(function (data) {
                    self.paths = data.data;
                })
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

            if (self.form.keyword) {
                $state.go('home.logbase.logs', self.form);
            } else {
                $state.go('home.logbase.logWithoutKey', self.form);
            }
        }
    }
})();