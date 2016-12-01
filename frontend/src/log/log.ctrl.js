(function () {
    'use strict';
    angular.module('app')
        .controller('LogCtrl', LogCtrl);
    /* @ngInject */
    function LogCtrl(logBackend, moment, logcurd, logChart) {
        var self = this;
        var tempLogQuery = {};

        self.chartOptions = logChart.Options();

        //md-table parameter
        self.options = {
            rowSelection: true,
            decapitate: false,
            boundaryLinks: false,
            limitSelect: true,
            pageSelect: true
        };

        self.defaultLimit = 50;
        self.limitOptions = [50, 100, 200];
        self.query = {
            limit: 50,
            page: 1
        };

        self.timePeriod = 120;
        self.curTimestamp = moment().unix() * 1000;
        self.fromTimestamp = moment().subtract(self.timePeriod, 'minutes').unix() * 1000;
        self.selectedTabIndex = 0;
        self.logDisplaySet = {};
        self.count = 0;
        self.logs = [];
        self.hasLogContext = false;

        self.onPaginate = onPaginate;
        self.periodChange = periodChange;
        self.selectAppChange = selectAppChange;
        self.selectTasksChange = selectTasksChange;
        self.loadApps = loadApps;
        self.loadTasks = loadTasks;
        self.loadPaths = loadPaths;
        self.searchLog = searchLog;
        self.logContext = logContext;

        activate();

        function activate() {

        }

        function onPaginate(page, limit) {
            self.logDisplaySet = {};
            logBackend.searchLogs({
                from: tempLogQuery.from,
                to: tempLogQuery.to,
                appid: tempLogQuery.appid,
                taskid: tempLogQuery.taskid,
                path: tempLogQuery.path,
                keyword: tempLogQuery.keyword,
                page: page,
                size: limit
            }).get(function (data) {
                self.logs = data.data.results;
                self.count = data.data.count;
            })
        }

        function periodChange(period) {
            self.selectApp = '';
            self.curTimestamp = moment().unix() * 1000;
            self.fromTimestamp = moment().subtract(period, 'minutes').unix() * 1000;
        }

        function selectAppChange(app) {
            self.selectTasks = '';
        }

        function selectTasksChange(tasks) {
            self.selectPaths = [];
        }

        function loadApps() {
            checkTimeRange();
            logBackend.listApp({from: self.fromTimestamp, to: self.curTimestamp}).get(function (data) {
                self.apps = data.data;
            })
        }

        function loadTasks() {
            checkTimeRange();
            logBackend.listTask({
                from: self.fromTimestamp,
                to: self.curTimestamp,
                appid: self.selectApp
            }).get(function (data) {
                self.tasks = data.data;
            })
        }

        function loadPaths() {
            checkTimeRange();
            logBackend.listPath({
                from: self.fromTimestamp,
                to: self.curTimestamp,
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
            self.logDisplaySet = {};
            self.query.page = 1;
            self.query.limit = 50;

            tempLogQuery = {
                from: self.fromTimestamp,
                to: self.curTimestamp,
                appid: self.selectApp,
                taskid: self.selectTasks,
                path: self.selectPaths,
                keyword: self.keyword,
                page: 1,
                size: 50
            };

            checkTimeRange();

            self.hasLogContext = !!self.keyword;

            logBackend.searchLogs({
                from: self.fromTimestamp,
                to: self.curTimestamp,
                appid: self.selectApp,
                taskid: self.selectTasks,
                path: self.selectPaths,
                keyword: self.keyword,
                page: 1,
                size: 50
            }).get(function (data) {
                self.logs = data.data.results;
                self.count = data.data.count;

                self.chartOptions.pushData(data.data.history);
            })
        }

        function logContext(ev, log) {
            logcurd.openLogContext(ev, log);
        }
    }
})();