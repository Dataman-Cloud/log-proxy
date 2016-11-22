(function () {
    'use strict';
    angular.module('app')
        .controller('LogCtrl', LogCtrl);
    /* @ngInject */
    function LogCtrl(logBackend, moment) {
        var self = this;

        self.timePeriod = 120;
        self.curTimestamp = moment().unix() * 1000;
        self.fromTimestamp = moment().subtract(self.timePeriod, 'minutes').unix() * 1000;

        self.periodChange = periodChange;
        self.selectAppChange = selectAppChange;
        self.selectTasksChange = selectTasksChange;
        self.loadApps = loadApps;
        self.loadTasks = loadTasks;
        self.loadPaths = loadPaths;

        self.searchLog = searchLog;

        activate();

        function activate() {

        }

        function periodChange(period) {
            self.selectApp = '';
            self.curTimestamp = moment().unix() * 1000;
            self.fromTimestamp = moment().subtract(period, 'minutes').unix() * 1000;
        }

        function selectAppChange(app) {
            self.selectTasks = [];
        }

        function selectTasksChange(tasks) {
            self.selectPaths = [];
        }

        function loadApps() {
            logBackend.listApp({from: self.fromTimestamp, to: self.curTimestamp}).get(function (data) {
                self.apps = data.data;
            })
        }

        function loadTasks() {
            logBackend.listTask({
                from: self.fromTimestamp,
                to: self.curTimestamp,
                appid: self.selectApp
            }).get(function (data) {
                self.tasks = data.data;
            })
        }

        function loadPaths() {
            logBackend.listTask({
                from: self.fromTimestamp,
                to: self.curTimestamp,
                appid: self.selectApp,
                taskid: self.selectTasks
            }).get(function (data) {
                self.paths = data.data;
            })
        }

        function searchLog() {
            console.log(self.selectTasks)
        }
    }
})();