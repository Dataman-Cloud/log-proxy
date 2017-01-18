(function () {
    'use strict';
    angular.module('app')
        .controller('LogsCtrl', LogsCtrl);
    /* @ngInject */
    function LogsCtrl(logBackend, $stateParams) {
        var tempLogQuery = {
            cluster: $stateParams.cluster || '',
            user: $stateParams.user || '',
            from: $stateParams.from,
            to: $stateParams.to,
            app: $stateParams.app,
            task: $stateParams.task,
            path: $stateParams.path,
            keyword: $stateParams.keyword,
            page: 1,
            size: 50
        };
        var self = this;

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
                cluster: tempLogQuery.cluster || '',
                user: tempLogQuery.user || '',
                from: tempLogQuery.from,
                to: tempLogQuery.to,
                app: tempLogQuery.app,
                task: tempLogQuery.task,
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
            self.loadingFlag = true;

            logBackend.searchLogs(tempLogQuery).get(function (data) {
                self.loadingFlag = false;
                self.logs = data.data.results;
                self.count = data.data.count;
            })
        }
    }
})();
