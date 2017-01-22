(function () {
    'use strict';
    angular.module('app')
        .controller('LogBaseCtrl', LogBaseCtrl);
    /* @ngInject */
    function LogBaseCtrl(moment, $state, logBackend, $stateParams) {
        var self = this;

        self.form = {
            cluster: $stateParams.cluster || '',
            user: $stateParams.user || '',
            app: $stateParams.app || '',
            task: $stateParams.task || '',
            path: $stateParams.path || '',
            keyword: $stateParams.keyword || ''
        };

        self.paths = [];

        self.startTime = new Date(parseInt($stateParams.from));
        self.endTime = new Date(parseInt($stateParams.to));

        self.timePeriod = 60;
        self.selectedTabIndex = ($stateParams.from && $stateParams.to) ? 1 : 0;

        self.loadPaths = loadPaths;

        self.searchLog = searchLog;

        activate();

        function activate() {
            loadPaths();
        }

        function loadPaths() {
            if (self.form.app) {
                logBackend.listPath({
                    cluster: self.form.cluster,
                    user: self.form.user,
                    app: self.form.app,
                    task: self.form.task
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

            if (self.form.keyword || (!self.form.keyword && !self.form.task)) {
                $state.go('home.logbase.logs', self.form);
            } else {
                $state.go('home.logbase.logWithoutKey', self.form);
            }
        }
    }
})();
