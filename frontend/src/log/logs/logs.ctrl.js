(function () {
    'use strict';
    angular.module('app')
        .controller('LogsCtrl', LogsCtrl);
    /* @ngInject */
    function LogsCtrl(logBackend, $stateParams) {
        var tempLogQuery = {
            from: $stateParams.from,
            to: $stateParams.to,
            appid: $stateParams.appid,
            taskid: $stateParams.taskid,
            path: $stateParams.path,
            keyword: $stateParams.keyword,
            page: 1,
            size: 50
        };
        var self = this;

        self.count = null; //log total

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

        self.logDisplaySet = {};
        self.count = 0;
        self.logs = [];

        self.onPaginate = onPaginate;

        activate();

        function activate() {
            getLogs();
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

        function getLogs() {
            logBackend.searchLogs(tempLogQuery).get(function (data) {
                self.logs = data.data.results;
                self.count = data.data.count;
            })
        }
    }
})();