(function () {
    'use strict';
    angular.module('app')
        .controller('AlertHistoriesCtrl', AlertHistoriesCtrl);
    /* @ngInject */
    function AlertHistoriesCtrl(alertBackend, $stateParams, moment, logBackend) {
        var self = this;

        self.timePeriod = 60;

        self.count = 0;
        self.histories = [];

        self.clusters = {};
        self.apps = {};
        self.sources = {};

        self.form = {
            cluster: $stateParams.cluster || '',
            app: $stateParams.app || '',
            start: moment().subtract(self.timePeriod, 'minutes').unix() * 1000,
            end: moment().unix() * 1000,
            page: 1,
            size: 100
        };

        //md-table parameter
        self.options = {
            rowSelection: true,
            decapitate: false,
            boundaryLinks: false,
            limitSelect: true,
            pageSelect: true
        };

        self.defaultLimit = 100;
        self.limitOptions = [100, 200, 500];
        self.query = {
            limit: 100,
            page: 1
        };

        self.loadClusters = loadClusters;
        self.loadApps = loadApps;
        self.loadSources = loadSources;
        self.onPaginate = onPaginate;
        self.searchHistory = searchHistory;

        activate();

        function activate() {
            loadClusters();
            loadApps();
            loadSources();
            listHistory();
        }

        function loadClusters() {
            logBackend.clusters().get(function (data) {
                self.clusters = data.data;
            })
        }

        function loadApps() {
            if (self.form.cluster) {
                logBackend.apps(self.form.cluster).get(function (data) {
                    self.apps = data.data;
                })
            }
        }

        function loadSources() {
            if (self.form.app && self.form.cluster) {
                logBackend.sources({
                    cluster: self.form.cluster,
                    app: self.form.app
                }).get(function (data) {
                    self.sources = data.data;
                })
            }
        }

        function fetchHistory(data) {
            alertBackend.histories(data).get(function (data) {
                if (data.data) {
                    self.histories = data.data.events;
                    self.count = data.data.count;

                }
            })
        }

        function searchHistory() {
            if (!self.timePeriod) {
                fetchHistory({page: 1, size: 100});
            } else {
                self.form.size = 100;
                self.form.page = 1;

                self.form.to = moment().unix() * 1000;
                self.form.from = moment().subtract(self.timePeriod, 'minutes').unix() * 1000;

                fetchHistory(self.form);
            }

        }

        function listHistory() {
            if ($stateParams.cluster && $stateParams.app) {
                self.form.size = 100;
                self.form.page = 1;

                fetchHistory(self.form)
            }
        }

        function onPaginate(page, limit) {
            self.form.size = limit;
            self.form.page = page;

            fetchHistory(self.form)
        }
    }
})();
