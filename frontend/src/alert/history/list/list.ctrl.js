(function () {
    'use strict';
    angular.module('app')
        .controller('AlertHistoriesCtrl', AlertHistoriesCtrl);
    /* @ngInject */
    function AlertHistoriesCtrl(alertBackend,$stateParams) {
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

        self.search = search;
        self.listHistory = listHistory;
        // self.listOperations = listOperations;
        // self.periodChange = periodChange;
        self.onPaginate = onPaginate;

        self.monitors = [];

        activate();


        function activate() {
            listHistory()
        }

        function search() {
            checkTimeRange();
            alertBackend.history()
                .then(function (data) {
                    self.history = data.data.results;
                    self.count = data.count;
                })
        }

        function listHistory() {
            self.loadingFlag = true;
            alertBackend.history().get(function (data) {
                self.histories = data.data.results;
                $.each(self.histories, function(index, value){
                    value.labels = JSON.parse(value.labels);
                });
                self.loadingFlag = false;
            })
        }

        function onPaginate(page, limit) {
            alertBackend.searchLog({
                operator: self.operator,
                operation: self.operationTypes,
                size: limit,
                page: page
            }).then(function (data) {
                self.auditLogs = data.audits;
                self.count = data.count;
            })
        }

    }
})();