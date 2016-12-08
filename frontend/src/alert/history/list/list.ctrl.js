(function () {
    'use strict';
    angular.module('app')
        .controller('AlertHistoriesCtrl', AlertHistoriesCtrl);
    /* @ngInject */
    function AlertHistoriesCtrl(alertBackend, $stateParams) {
        var self = this;
        self.historyDisplaySet = {};
        self.count = 0;
        self.histories = [];

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

        self.listHistory = listHistory;
        self.onPaginate = onPaginate;

        activate();

        function activate() {
            listHistory()
        }

        function listHistory() {
            alertBackend.histories({size: 50, page: 1}).get(function (data) {
                self.histories = data.data.results;
                self.count = data.data.count;
                angular.forEach(self.histories, function (history, index) {
                    history.labels = angular.fromJson(history.labels)
                });
            })
        }

        function onPaginate(page, limit) {
            self.historyDisplaySet = {};

            alertBackend.histories({size: limit, page: page}).get(function (data) {
                self.histories = data.data.results;
                self.count = data.data.count;
                angular.forEach(self.histories, function (history, index) {
                    history.labels = angular.fromJson(history.labels)
                });
            })
        }
    }
})();